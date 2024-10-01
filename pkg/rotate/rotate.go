/*
Copyright The Rotate Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package rotate

import (
	"fmt"
	"log"

	"github.com/golang-module/carbon"
)

type BackupRotationScheme struct {
	Hourly  int
	Daily   int
	Weekly  int
	Monthly int
	Yearly  int
	DryRun  bool
}

type BackupSummary struct {
	Hourly             []*Backup
	Daily              []*Backup
	Weekly             []*Backup
	Monthly            []*Backup
	Yearly             []*Backup
	ForDelete          []*Backup
	SizeTotalHourly    int64
	SizeTotalDaily     int64
	SizeTotalWeekly    int64
	SizeTotalMonthly   int64
	SizeTotalYearly    int64
	SizeTotalForDelete int64
}

func (s BackupSummary) GetTotalCategorized() int {
	total := 0
	total += len(s.Hourly)
	total += len(s.Daily)
	total += len(s.Weekly)
	total += len(s.Monthly)
	total += len(s.Yearly)
	total += len(s.ForDelete)
	return total
}

func (summary BackupSummary) Print() {
	log.Println("")
	summary.printBackups("Yearly", summary.Yearly, summary.SizeTotalYearly)
	summary.printBackups("Monthly", summary.Monthly, summary.SizeTotalMonthly)
	summary.printBackups("Weekly", summary.Weekly, summary.SizeTotalWeekly)
	summary.printBackups("Daily", summary.Daily, summary.SizeTotalDaily)
	summary.printBackups("Hourly", summary.Hourly, summary.SizeTotalHourly)
	summary.printBackups("Delete", summary.ForDelete, summary.SizeTotalForDelete)
}

func (summary BackupSummary) printBackups(category string, backups []*Backup, sizeTotal int64) {
	formattedSize := summary.formatSize(sizeTotal)
	log.Printf("%s matched [%d]:", category, len(backups))
	if len(backups) == 0 {
		log.Println("  No files")
	} else {
		for _, v := range backups {
			log.Println(" ", v.Path, summary.formatSize(v.Size), v.Timestamp)
		}
		log.Printf("  Total Size: %s", formattedSize)
	}
	log.Println("")
}

func (summary BackupSummary) formatSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f%cB", float64(size)/float64(div), "KMGTPE"[exp])
}

type Backup struct {
	Bucket    string
	Path      string
	Size      int64
	Timestamp carbon.Carbon
}

func (b Backup) String() string {
	return fmt.Sprintf("Path: %s, Timestamp: %s", b.Path, b.Timestamp)
}

func (b Backup) IsHourlyOf(date carbon.Carbon, prev *carbon.Carbon) bool {
	if b.IsSameHour(prev) {
		return false
	}
	return b.Timestamp.DiffInHours(date) <= carbon.HoursPerDay
}

func (b Backup) IsDailyOf(date carbon.Carbon, prev *carbon.Carbon) bool {
	if b.IsSameDay(prev) {
		return false
	}
	diff := int(b.Timestamp.DiffInDays(date))
	return diff >= 1 && diff <= carbon.DaysPerWeek
}

func (b Backup) IsWeeklyOf(date carbon.Carbon, prev *carbon.Carbon, limit int) bool {
	if b.IsSameWeek(prev) {
		return false
	}
	return b.Timestamp.DiffInWeeks(date) <= int64(limit) && b.Timestamp.IsSunday()
}

func (b Backup) IsMonthlyOf(date carbon.Carbon, prev *carbon.Carbon) bool {
	if b.IsSameMonth(prev) {
		return false
	}
	return b.Timestamp.DiffInMonths(date) <= 13 && b.Timestamp.DiffInWeeks(date) >= 4
}

func (b Backup) IsYearlyOf(date carbon.Carbon, prevBackup *carbon.Carbon) bool {
	// Se houver um backup anterior no mesmo ano, não o consideramos anual
	if b.IsSameYear(prevBackup) {
		return false
	}

	monthsDiff := b.Timestamp.DiffInMonths(date)

	// Verificamos se o backup tem pelo menos 12 meses ou mais de 6 meses e é de um ano diferente
	return monthsDiff >= 12 || (monthsDiff > 6 && !b.IsSameYear(&date))
}

func (b Backup) IsSameHour(compare *carbon.Carbon) bool {
	if compare == nil {
		return false
	}
	return b.Timestamp.IsSameHour(*compare)
}

func (b Backup) IsSameDay(compare *carbon.Carbon) bool {
	if compare == nil {
		return false
	}
	return b.Timestamp.IsSameDay(*compare)
}

func (b Backup) IsSameWeek(compare *carbon.Carbon) bool {
	if compare == nil {
		return false
	}
	return b.Timestamp.Between(compare.StartOfWeek(), compare.EndOfWeek())
}

func (b Backup) IsSameMonth(compare *carbon.Carbon) bool {
	if compare == nil {
		return false
	}
	return b.Timestamp.IsSameMonth(*compare)
}

func (b Backup) IsSameYear(compare *carbon.Carbon) bool {
	if compare == nil {
		return false
	}
	return b.Timestamp.IsSameYear(*compare)
}

type BackupFiles []*Backup

func (b BackupFiles) Less(i, j int) bool {
	return b[i].Timestamp.Gt(b[j].Timestamp)
}

func (b BackupFiles) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

func (b BackupFiles) Len() int {
	return len(b)
}
