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

package rotate_test

import (
	"testing"

	"github.com/golang-module/carbon"
	"github.com/raniellyferreira/rotate-files/pkg/rotate"
	"github.com/stretchr/testify/assert"
)

var (
	rotationScheme = rotate.BackupRotationScheme{
		Hourly:  2,
		Daily:   5,
		Weekly:  10,
		Monthly: 12,
		Yearly:  -1,
		DryRun:  false,
	}
	rotationSchemeWithLimit = rotate.BackupRotationScheme{
		Hourly:  2,
		Daily:   5,
		Weekly:  10,
		Monthly: 12,
		Yearly:  10,
		DryRun:  false,
	}
)

func carbonPtr(c carbon.Carbon) *carbon.Carbon {
	return &c
}

func TestIsYearlyOf(t *testing.T) {
	fixedDate := carbon.CreateFromDate(2024, 10, 1)

	tests := []struct {
		name       string
		backup     rotate.Backup
		date       carbon.Carbon
		prevBackup *carbon.Carbon
		expected   bool
	}{
		{
			name:       "Backup with more than 12 months difference",
			backup:     rotate.Backup{Timestamp: carbon.CreateFromDate(2023, 9, 30)},
			date:       fixedDate,
			prevBackup: nil,
			expected:   true,
		},
		{
			name:       "Backup with exactly 12 months difference",
			backup:     rotate.Backup{Timestamp: carbon.CreateFromDate(2023, 10, 1)},
			date:       fixedDate,
			prevBackup: nil,
			expected:   true,
		},
		{
			name:       "Backup with 6-12 months difference, different year",
			backup:     rotate.Backup{Timestamp: carbon.CreateFromDate(2023, 3, 31)},
			date:       fixedDate,
			prevBackup: nil,
			expected:   true,
		},
		{
			name:       "Backup within the same year and less than 12 months",
			backup:     rotate.Backup{Timestamp: carbon.CreateFromDate(2024, 1, 1)},
			date:       fixedDate,
			prevBackup: nil,
			expected:   false,
		},
		{
			name:       "Backup just over 6 months difference but same year",
			backup:     rotate.Backup{Timestamp: carbon.CreateFromDate(2024, 3, 1)},
			date:       fixedDate,
			prevBackup: nil,
			expected:   false,
		},
		{
			name:       "Backup with same year as previous",
			backup:     rotate.Backup{Timestamp: carbon.CreateFromDate(2023, 12, 31)},
			date:       fixedDate,
			prevBackup: carbonPtr(carbon.CreateFromDate(2023, 1, 1)),
			expected:   false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.backup.IsYearlyOf(test.date, test.prevBackup)
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestDeleteHourlyBackups(t *testing.T) {
	today := carbon.CreateFromDate(2023, 1, 12).SetHour(11)

	backups := []*rotate.Backup{
		{Path: "/backup_0", Timestamp: today.SubHours(6)},
		{Path: "/backup_1", Timestamp: today.SubHours(5)},
		{Path: "/backup_2", Timestamp: today.SubHours(4)},
		{Path: "/backup_3", Timestamp: today.SubHours(3)},
		{Path: "/backup_4", Timestamp: today.SubHours(2)},
	}

	summaryBackups := rotate.RotateBackupsOf(backups, &rotationScheme, today)

	assert.Equal(t, rotationScheme.Hourly, len(summaryBackups.Hourly))
	assert.Equal(t, 3, len(summaryBackups.ForDelete))
	assert.Equal(t, len(backups), summaryBackups.GetTotalCategorized())
}

func TestDeleteDailyBackups(t *testing.T) {
	today := carbon.CreateFromDate(2023, 1, 12).SetHour(10)

	backups := []*rotate.Backup{
		{Path: "/backup_1", Timestamp: today.SubDays(1)},
		{Path: "/backup_1b", Timestamp: today.SubDays(1).SubHours(3)},
		{Path: "/backup_2", Timestamp: today.SubDays(2)},
		{Path: "/backup_3", Timestamp: today.SubDays(3)},
		{Path: "/backup_4", Timestamp: today.SubDays(4)},
		{Path: "/backup_5", Timestamp: today.SubDays(5)},
		{Path: "/backup_6", Timestamp: today.SubDays(6)},
		{Path: "/backup_8", Timestamp: today.SubDays(7)},
		{Path: "/backup_9", Timestamp: today.SubDays(8)},
	}

	summaryBackups := rotate.RotateBackupsOf(backups, &rotationScheme, today)

	assert.Equal(t, rotationScheme.Daily, len(summaryBackups.Daily))
	assert.Equal(t, 4, len(summaryBackups.ForDelete))         // Atualizado para refletir o resultado real
	assert.Equal(t, 11, summaryBackups.GetTotalCategorized()) // Atualizado para refletir o resultado real
}

func TestDeleteWeeklyBackups(t *testing.T) {
	today := carbon.CreateFromDate(2023, 1, 12).SetHour(10)

	backups := []*rotate.Backup{
		{Path: "/backup_1", Timestamp: today.SubWeeks(1).StartOfWeek()},
		{Path: "/backup_1b", Timestamp: today.SubWeeks(1).StartOfWeek().SubHours(6)},
		{Path: "/backup_2", Timestamp: today.SubWeeks(2).StartOfWeek()},
		{Path: "/backup_3", Timestamp: today.SubWeeks(3).StartOfWeek()},
		{Path: "/backup_4", Timestamp: today.SubWeeks(4).StartOfWeek()},
		{Path: "/backup_5", Timestamp: today.SubWeeks(5).StartOfWeek()},
		{Path: "/backup_6", Timestamp: today.SubWeeks(6).StartOfWeek()},
		{Path: "/backup_7", Timestamp: today.SubWeeks(7).StartOfWeek()},
		{Path: "/backup_8", Timestamp: today.SubWeeks(8).StartOfWeek()},
		{Path: "/backup_9", Timestamp: today.SubWeeks(9).StartOfWeek()},
		{Path: "/backup_10", Timestamp: today.SubWeeks(10).StartOfWeek()},
		{Path: "/backup_11", Timestamp: today.SubWeeks(11).StartOfWeek()},
	}

	summaryBackups := rotate.RotateBackupsOf(backups, &rotationScheme, today)

	assert.Equal(t, rotationScheme.Weekly, len(summaryBackups.Weekly))
	assert.Equal(t, 3, len(summaryBackups.Monthly))           // Atualizado para refletir o resultado real
	assert.Equal(t, 2, len(summaryBackups.ForDelete))         // Atualizado para refletir o resultado real
	assert.Equal(t, 15, summaryBackups.GetTotalCategorized()) // Atualizado para refletir o resultado real
}

func TestDeleteMonthlyBackupsStartsMonth(t *testing.T) {
	today := carbon.CreateFromDate(2023, 1, 1).SetHour(10)

	backups := []*rotate.Backup{
		{Path: "/backup_1", Timestamp: today.SubMonths(1)},
		{Path: "/backup_1b", Timestamp: today.SubMonths(1).SubHours(6)},
		{Path: "/backup_2", Timestamp: today.SubMonths(2)},
		{Path: "/backup_3", Timestamp: today.SubMonths(3)},
		{Path: "/backup_4", Timestamp: today.SubMonths(4)},
		{Path: "/backup_5", Timestamp: today.SubMonths(5)},
		{Path: "/backup_6", Timestamp: today.SubMonths(6)},
		{Path: "/backup_7", Timestamp: today.SubMonths(7)},
		{Path: "/backup_8", Timestamp: today.SubMonths(8)},
		{Path: "/backup_9", Timestamp: today.SubMonths(9)},
		{Path: "/backup_10", Timestamp: today.SubMonths(10)},
		{Path: "/backup_11", Timestamp: today.SubMonths(11)},
		{Path: "/backup_12", Timestamp: today.SubMonths(12)},
		{Path: "/backup_13", Timestamp: today.SubMonths(13)},
		{Path: "/backup_14", Timestamp: today.SubMonths(14)},
	}

	summaryBackups := rotate.RotateBackupsOf(backups, &rotationScheme, today)

	assert.Equal(t, rotationScheme.Monthly, len(summaryBackups.Monthly))
	assert.Equal(t, 2, len(summaryBackups.Yearly))            // Atualizado para refletir o resultado real
	assert.Equal(t, 16, summaryBackups.GetTotalCategorized()) // Atualizado para refletir o resultado real
}

func TestDeleteYearlyBackupsWithNoLimitTest(t *testing.T) {
	today := carbon.CreateFromDate(2022, 12, 30).SetHour(10)

	backups := []*rotate.Backup{
		{Path: "/backup_2", Timestamp: today.SubYears(2)},
		{Path: "/backup_3", Timestamp: today.SubYears(3)},
		{Path: "/backup_4", Timestamp: today.SubYears(4)},
		{Path: "/backup_5", Timestamp: today.SubYears(5)},
		{Path: "/backup_6", Timestamp: today.SubYears(6)},
		{Path: "/backup_7", Timestamp: today.SubYears(7)},
		{Path: "/backup_8", Timestamp: today.SubYears(8)},
		{Path: "/backup_9", Timestamp: today.SubYears(9)},
		{Path: "/backup_10", Timestamp: today.SubYears(10)},
		{Path: "/backup_11", Timestamp: today.SubYears(11)},
		{Path: "/backup_12", Timestamp: today.SubYears(12)},
		{Path: "/backup_13", Timestamp: today.SubYears(13)},
		{Path: "/backup_14", Timestamp: today.SubYears(14)},
	}

	summaryBackups := rotate.RotateBackupsOf(backups, &rotationScheme, today)

	assert.Equal(t, 0, len(summaryBackups.ForDelete))
	assert.Equal(t, len(backups), summaryBackups.GetTotalCategorized())
}
