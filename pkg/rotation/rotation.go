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

package rotation

import (
	"fmt"
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
	Hourly    []Backup
	Daily     []Backup
	Weekly    []Backup
	Monthly   []Backup
	Yearly    []Backup
	ForDelete []Backup
}

func (s BackupSummary) GetTotalCategorized() int {
	return len(s.Hourly) +
		len(s.Daily) +
		len(s.Weekly) +
		len(s.Monthly) +
		len(s.Yearly) +
		len(s.ForDelete)
}

type Backup struct {
	Bucket    string
	Path      string
	Timestamp carbon.Carbon
	DeleteMe  bool
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

func (b Backup) IsWeekly(limit int64) bool {
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

func (b Backup) IsWeeklyOf(date carbon.Carbon, limit int64) bool {
	return b.Timestamp.DiffInWeeks(date) <= limit && b.Timestamp.IsSunday()
}

func (b Backup) IsMonthlyOf(date carbon.Carbon) bool {
	return int(b.Timestamp.DiffInMonths(date)) <= 13 && int(b.Timestamp.DiffInWeeks(date)) >= 4
}

func (b Backup) IsYearlyOf(date carbon.Carbon) bool {
	return int(b.Timestamp.DiffInYears(date)) >= 1 && !b.Timestamp.IsSameYear(date)
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

	// Define os limites de backups a serem mantidos em cada grupo
	hourlyLimit := rotationScheme.Hourly
	dailyLimit := rotationScheme.Daily
	weeklyLimit := rotationScheme.Weekly
	monthlyLimit := rotationScheme.Monthly
	yearlyLimit := rotationScheme.Yearly

	var hourly, daily, weekly, monthly, yearly, forDelete BackupFiles
	var prevHourly, prevDaily, prevWeekly, prevMonthly, prevYearly Backup

	sort.Sort(backups)

	for _, backup := range backups {

		if backup.IsHourlyOf(date) && !backup.IsSameHour(prevHourly) {
			if len(hourly) < hourlyLimit {
				hourly = append(hourly, backup)
				prevHourly = backup
				continue
			}
		}

		if backup.IsDailyOf(date) && !backup.IsSameDay(prevDaily) {
			if len(daily) < dailyLimit {
				daily = append(daily, backup)
				prevDaily = backup
				continue
			}
		}

		if backup.IsWeeklyOf(date, int64(weeklyLimit)) && !backup.IsSameWeek(prevWeekly) {
			if len(weekly) < weeklyLimit {
				weekly = append(weekly, backup)
				prevWeekly = backup
				continue
			}
		}

		if backup.IsMonthlyOf(date) && !backup.IsSameMonth(prevMonthly) {
			if len(monthly) < monthlyLimit {
				monthly = append(monthly, backup)
				prevMonthly = backup
				continue
			}
		}

		if backup.IsYearlyOf(date) && !backup.IsSameYear(prevYearly) {
			if yearlyLimit < 1 || len(yearly) < yearlyLimit {
				yearly = append(yearly, backup)
				prevYearly = backup
				continue
			}
		}

		forDelete = append(forDelete, backup)
	}

	return BackupSummary{
		Hourly:    hourly,
		Daily:     daily,
		Weekly:    weekly,
		Monthly:   monthly,
		Yearly:    yearly,
		ForDelete: forDelete,
	}
}
