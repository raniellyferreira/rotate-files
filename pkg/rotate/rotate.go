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
	"sort"

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
	Hourly             []Backup
	Daily              []Backup
	Weekly             []Backup
	Monthly            []Backup
	Yearly             []Backup
	ForDelete          []Backup
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

func (summary BackupSummary) printBackups(category string, backups []Backup, sizeTotal int64) {
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

func (b Backup) IsHourly() bool {
	return b.IsHourlyOf(carbon.Now())
}

func (b Backup) IsDaily() bool {
	return b.IsDailyOf(carbon.Now())
}

func (b Backup) IsWeekly(limit int) bool {
	return b.IsWeeklyOf(carbon.Now(), limit)
}

func (b Backup) IsMonthly() bool {
	return b.IsMonthlyOf(carbon.Now())
}

func (b Backup) IsYearly() bool {
	return b.IsYearlyOf(carbon.Now())
}

func (b Backup) IsHourlyOf(date carbon.Carbon) bool {
	return b.Timestamp.DiffInHours(date) <= carbon.HoursPerDay
}

func (b Backup) IsDailyOf(date carbon.Carbon) bool {
	diff := int(b.Timestamp.DiffInDays(date))
	return diff >= 1 && diff <= carbon.DaysPerWeek
}

func (b Backup) IsWeeklyOf(date carbon.Carbon, limit int) bool {
	return b.Timestamp.DiffInWeeks(date) <= int64(limit) && b.Timestamp.IsSunday()
}

func (b Backup) IsMonthlyOf(date carbon.Carbon) bool {
	return b.Timestamp.DiffInMonths(date) <= 13 && b.Timestamp.DiffInWeeks(date) >= 4
}

func (b Backup) IsYearlyOf(date carbon.Carbon) bool {
	return (b.Timestamp.DiffInMonths(date) > 6 && !b.Timestamp.IsSameYear(date)) || b.Timestamp.DiffInMonths(date) >= 12
}

func (b Backup) IsSameHour(compare Backup) bool {
	if compare == (Backup{}) {
		return false
	}
	return b.Timestamp.IsSameHour(compare.Timestamp)
}

func (b Backup) IsSameDay(compare Backup) bool {
	if compare == (Backup{}) {
		return false
	}
	return b.Timestamp.IsSameDay(compare.Timestamp)
}

func (b Backup) IsSameWeek(compare Backup) bool {
	if compare == (Backup{}) {
		return false
	}
	return b.Timestamp.Between(compare.Timestamp.StartOfWeek(), compare.Timestamp.EndOfWeek())
}

func (b Backup) IsSameMonth(compare Backup) bool {
	if compare == (Backup{}) {
		return false
	}
	return b.Timestamp.IsSameMonth(compare.Timestamp)
}

func (b Backup) IsSameYear(compare Backup) bool {
	if compare == (Backup{}) {
		return false
	}
	return b.Timestamp.IsSameYear(compare.Timestamp)
}

type BackupFiles []Backup

func (b BackupFiles) Less(i, j int) bool {
	return b[i].Timestamp.Gt(b[j].Timestamp)
}

func (b BackupFiles) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

func (b BackupFiles) Len() int {
	return len(b)
}

func (b BackupFiles) Rotate(rotationScheme *BackupRotationScheme) BackupSummary {
	return b.RotateOf(rotationScheme, carbon.Now())
}

func (backups BackupFiles) RotateOf(rotationScheme *BackupRotationScheme, date carbon.Carbon) BackupSummary {
	var hourly, daily, weekly, monthly, yearly, forDelete BackupFiles
	var prevHourly, prevDaily, prevWeekly, prevMonthly, prevYearly Backup
	var totalSizeHourly, totalSizeDaily, totalSizeWeekly, totalSizeMonthly, totalSizeYearly, totalSizeForDelete int64

	sort.Sort(backups)

	for _, backup := range backups {
		switch {
		case backup.IsHourlyOf(date) && !backup.IsSameHour(prevHourly) && len(hourly) < rotationScheme.Hourly:
			hourly = append(hourly, backup)
			totalSizeHourly += backup.Size
			prevHourly = backup

		case backup.IsDailyOf(date) && !backup.IsSameDay(prevDaily) && len(daily) < rotationScheme.Daily:
			daily = append(daily, backup)
			totalSizeDaily += backup.Size
			prevDaily = backup

		case backup.IsWeeklyOf(date, rotationScheme.Weekly) && !backup.IsSameWeek(prevWeekly) && len(weekly) < rotationScheme.Weekly:
			weekly = append(weekly, backup)
			totalSizeWeekly += backup.Size
			prevWeekly = backup

		case backup.IsMonthlyOf(date) && !backup.IsSameMonth(prevMonthly) && len(monthly) < rotationScheme.Monthly:
			monthly = append(monthly, backup)
			totalSizeMonthly += backup.Size
			prevMonthly = backup

		case backup.IsYearlyOf(date) && !backup.IsSameYear(prevYearly) && (rotationScheme.Yearly < 1 || len(yearly) < rotationScheme.Yearly):
			yearly = append(yearly, backup)
			totalSizeYearly += backup.Size
			prevYearly = backup

		default:
			forDelete = append(forDelete, backup)
			totalSizeForDelete += backup.Size
		}
	}

	return BackupSummary{
		Hourly:    hourly,
		Daily:     daily,
		Weekly:    weekly,
		Monthly:   monthly,
		Yearly:    yearly,
		ForDelete: forDelete,

		SizeTotalHourly:    totalSizeHourly,
		SizeTotalDaily:     totalSizeDaily,
		SizeTotalWeekly:    totalSizeWeekly,
		SizeTotalMonthly:   totalSizeMonthly,
		SizeTotalYearly:    totalSizeYearly,
		SizeTotalForDelete: totalSizeForDelete,
	}
}
