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
	"sort"

	"github.com/golang-module/carbon"
	"github.com/raniellyferreira/rotate-files/pkg/providers"
)

type RotationManager struct {
	provider providers.Provider
}

func NewRotationManager(provider providers.Provider) *RotationManager {
	return &RotationManager{provider: provider}
}

func (r *RotationManager) ListBackups(path string) ([]*Backup, error) {
	infos, err := r.provider.ListFiles(path)
	if err != nil {
		return nil, err
	}

	backups := make([]*Backup, len(infos))
	for i, info := range infos {
		backups[i] = &Backup{
			Path:      info.Path,
			Size:      info.Size,
			Timestamp: info.Timestamp,
		}
	}
	return backups, nil
}

func (r *RotationManager) RemoveBackup(fullPath string) error {
	return r.provider.Delete(fullPath)
}

func (r *RotationManager) RotateBackups(backups []*Backup, scheme *BackupRotationScheme) *BackupSummary {
	return RotateBackupsOf(backups, scheme, carbon.Now())
}

func RotateBackupsOf(backups []*Backup, scheme *BackupRotationScheme, current carbon.Carbon) *BackupSummary {
	sort.Sort(BackupFiles(backups))

	var hourly, daily, weekly, monthly, yearly, forDelete BackupFiles
	var prevYearlyBackup, prevMonthlyBackup, prevWeeklyBackup, prevDailyBackup, prevHourlyBackup *carbon.Carbon
	var totalSizeHourly, totalSizeDaily, totalSizeWeekly, totalSizeMonthly, totalSizeYearly, totalSizeForDelete int64

	for _, backup := range backups {
		addedToCategory := false

		if backup.IsHourlyOf(current, prevHourlyBackup) && len(hourly) < scheme.Hourly {
			hourly = append(hourly, backup)
			prevHourlyBackup = &backup.Timestamp
			totalSizeHourly += backup.Size
			addedToCategory = true
		}

		if backup.IsDailyOf(current, prevDailyBackup) && len(daily) < scheme.Daily {
			daily = append(daily, backup)
			prevDailyBackup = &backup.Timestamp
			totalSizeDaily += backup.Size
			addedToCategory = true
		}

		if backup.IsWeeklyOf(current, prevWeeklyBackup, scheme.Weekly) && len(weekly) < scheme.Weekly {
			weekly = append(weekly, backup)
			prevWeeklyBackup = &backup.Timestamp
			totalSizeWeekly += backup.Size
			addedToCategory = true
		}

		if backup.IsMonthlyOf(current, prevMonthlyBackup) && len(monthly) < scheme.Monthly {
			monthly = append(monthly, backup)
			prevMonthlyBackup = &backup.Timestamp
			totalSizeMonthly += backup.Size
			addedToCategory = true
		}

		if backup.IsYearlyOf(current, prevYearlyBackup) {
			if scheme.Yearly == -1 || len(yearly) < scheme.Yearly {
				yearly = append(yearly, backup)
				prevYearlyBackup = &backup.Timestamp
				totalSizeYearly += backup.Size
				addedToCategory = true
			}
		}

		if !addedToCategory {
			forDelete = append(forDelete, backup)
			totalSizeForDelete += backup.Size
		}
	}

	return &BackupSummary{
		Hourly:             hourly,
		Daily:              daily,
		Weekly:             weekly,
		Monthly:            monthly,
		Yearly:             yearly,
		ForDelete:          forDelete,
		SizeTotalHourly:    totalSizeHourly,
		SizeTotalDaily:     totalSizeDaily,
		SizeTotalWeekly:    totalSizeWeekly,
		SizeTotalMonthly:   totalSizeMonthly,
		SizeTotalYearly:    totalSizeYearly,
		SizeTotalForDelete: totalSizeForDelete,
	}
}
