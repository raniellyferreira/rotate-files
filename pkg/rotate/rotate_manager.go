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

	var backups []*Backup
	for _, info := range infos {
		backups = append(backups, &Backup{
			Path:      info.Path,
			Size:      info.Size,
			Timestamp: info.Timestamp,
		})
	}
	return backups, nil
}

func (r *RotationManager) RotateBackups(backups []*Backup, scheme *BackupRotationScheme) *BackupSummary {
	sort.Sort(BackupFiles(backups))

	var hourly, daily, weekly, monthly, yearly, forDelete BackupFiles
	var totalSizeHourly, totalSizeDaily, totalSizeWeekly, totalSizeMonthly, totalSizeYearly, totalSizeForDelete int64

	current := carbon.Now()

	for _, backup := range backups {
		switch {
		case backup.IsHourlyOf(current) && len(hourly) < scheme.Hourly:
			hourly = append(hourly, backup)
			totalSizeHourly += backup.Size

		case backup.IsDailyOf(current) && len(daily) < scheme.Daily:
			daily = append(daily, backup)
			totalSizeDaily += backup.Size

		case backup.IsWeeklyOf(current, scheme.Weekly) && len(weekly) < scheme.Weekly:
			weekly = append(weekly, backup)
			totalSizeWeekly += backup.Size

		case backup.IsMonthlyOf(current) && len(monthly) < scheme.Monthly:
			monthly = append(monthly, backup)
			totalSizeMonthly += backup.Size

		case backup.IsYearlyOf(current) && len(yearly) < scheme.Yearly:
			yearly = append(yearly, backup)
			totalSizeYearly += backup.Size

		default:
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

func (r *RotationManager) RemoveBackup(path string) error {
	return r.provider.Delete(path)
}
