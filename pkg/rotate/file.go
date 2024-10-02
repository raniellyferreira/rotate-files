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

	"github.com/golang-module/carbon"
)

// File represents a backup file with its path, size, and timestamp.
type File struct {
	Path      string
	Size      int64
	Timestamp carbon.Carbon
}

// String returns the string representation of the File, including path and timestamp.
func (b File) String() string {
	return fmt.Sprintf("Path: %s, Timestamp: %s", b.Path, b.Timestamp)
}

// IsHourlyOf checks if the file is an hourly backup based on the provided date.
func (b File) IsHourlyOf(date carbon.Carbon, prev *carbon.Carbon) bool {
	if b.IsSameHour(prev) {
		return false
	}
	return b.Timestamp.DiffInHours(date) <= carbon.HoursPerDay
}

// IsDailyOf checks if the file is a daily backup based on the provided date.
func (b File) IsDailyOf(date carbon.Carbon, prev *carbon.Carbon) bool {
	if b.IsSameDay(prev) {
		return false
	}
	diff := int(b.Timestamp.DiffInDays(date))
	return diff >= 1 && diff <= carbon.DaysPerWeek
}

// IsWeeklyOf checks if the file is a weekly backup based on the provided date and limit.
func (b File) IsWeeklyOf(date carbon.Carbon, prev *carbon.Carbon, limit int) bool {
	if b.IsSameWeek(prev) {
		return false
	}
	return b.Timestamp.DiffInWeeks(date) <= int64(limit) && b.Timestamp.IsSunday()
}

// IsMonthlyOf checks if the file is a monthly backup based on the provided date.
func (b File) IsMonthlyOf(date carbon.Carbon, prev *carbon.Carbon) bool {
	if b.IsSameMonth(prev) {
		return false
	}
	return b.Timestamp.DiffInMonths(date) <= 13 && b.Timestamp.DiffInWeeks(date) >= 4
}

// IsYearlyOf checks if the file is a yearly backup based on the provided date.
func (b File) IsYearlyOf(date carbon.Carbon, prevBackup *carbon.Carbon) bool {
	// Se houver um backup anterior no mesmo ano, não o consideramos anual
	if b.IsSameYear(prevBackup) {
		return false
	}

	monthsDiff := b.Timestamp.DiffInMonths(date)

	// Verificamos se o backup tem pelo menos 12 meses ou mais de 6 meses e é de um ano diferente
	return monthsDiff >= 12 || (monthsDiff > 6 && !b.IsSameYear(&date))
}

// IsSameHour checks if the file has the same hour as the provided date.
func (b File) IsSameHour(compare *carbon.Carbon) bool {
	if compare == nil {
		return false
	}
	return b.Timestamp.IsSameHour(*compare)
}

// IsSameDay checks if the file has the same day as the provided date.
func (b File) IsSameDay(compare *carbon.Carbon) bool {
	if compare == nil {
		return false
	}
	return b.Timestamp.IsSameDay(*compare)
}

// IsSameWeek checks if the file has the same week as the provided date.
func (b File) IsSameWeek(compare *carbon.Carbon) bool {
	if compare == nil {
		return false
	}
	return b.Timestamp.Between(compare.StartOfWeek(), compare.EndOfWeek())
}

// IsSameMonth checks if the file has the same month as the provided date.
func (b File) IsSameMonth(compare *carbon.Carbon) bool {
	if compare == nil {
		return false
	}
	return b.Timestamp.IsSameMonth(*compare)
}

// IsSameYear checks if the file has the same year as the provided date.
func (b File) IsSameYear(compare *carbon.Carbon) bool {
	if compare == nil {
		return false
	}
	return b.Timestamp.IsSameYear(*compare)
}

// Files is a slice of File pointers.
type Files []*File

// Less compares the elements at the given indexes.
func (b Files) Less(i, j int) bool {
	return b[i].Timestamp.Gt(b[j].Timestamp)
}

// Swap swaps the elements at the given indexes.
func (b Files) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

// Len returns the length of the Files slice.
func (b Files) Len() int {
	return len(b)
}
