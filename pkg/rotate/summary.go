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
)

// Summary represents the categorized backup files and their sizes.
type Summary struct {
	Hourly             []*File
	Daily              []*File
	Weekly             []*File
	Monthly            []*File
	Yearly             []*File
	ForDelete          []*File
	SizeTotalHourly    int64
	SizeTotalDaily     int64
	SizeTotalWeekly    int64
	SizeTotalMonthly   int64
	SizeTotalYearly    int64
	SizeTotalForDelete int64
}

// GetTotalCategorized returns the total number of categorized files in the summary.
func (s Summary) GetTotalCategorized() int {
	total := 0
	total += len(s.Hourly)
	total += len(s.Daily)
	total += len(s.Weekly)
	total += len(s.Monthly)
	total += len(s.Yearly)
	total += len(s.ForDelete)
	return total
}

// Print displays the categorized backup files and their sizes.
func (s Summary) Print() {
	log.Println("")
	s.printBackups("Delete", s.ForDelete, s.SizeTotalForDelete)
	s.printBackups("Yearly", s.Yearly, s.SizeTotalYearly)
	s.printBackups("Monthly", s.Monthly, s.SizeTotalMonthly)
	s.printBackups("Weekly", s.Weekly, s.SizeTotalWeekly)
	s.printBackups("Daily", s.Daily, s.SizeTotalDaily)
	s.printBackups("Hourly", s.Hourly, s.SizeTotalHourly)
}

// printBackups displays the backup files in the specified category.
func (s Summary) printBackups(category string, backups []*File, sizeTotal int64) {
	formattedSize := s.formatSize(sizeTotal)
	log.Printf("%s matched [%d]:", category, len(backups))
	if len(backups) == 0 {
		log.Println("  No files")
	} else {
		for _, v := range backups {
			log.Println(" ", v.Path, s.formatSize(v.Size), v.Timestamp)
		}
		log.Printf("  Total Size: %s", formattedSize)
	}
	log.Println("")
}

// formatSize converts the size in bytes to a human-readable format.
func (s Summary) formatSize(size int64) string {
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
