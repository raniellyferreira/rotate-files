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
	provider       providers.Provider
	rotationScheme *RotationScheme
	path           string
}

// NewRotationManager creates a new RotationManager with a rotation scheme.
func NewRotationManager(provider providers.Provider, scheme *RotationScheme, path string) *RotationManager {
	return &RotationManager{
		provider:       provider,
		rotationScheme: scheme,
		path:           path,
	}
}

// Validate checks if the rotation manager is ready to rotate files.
func (r *RotationManager) Validate(fileList []*File) error {
	if r.rotationScheme == nil {
		return ErrNilRotationScheme
	}

	if r.provider == nil {
		return ErrNilProvider
	}

	if len(fileList) == 0 {
		return ErrEmptyFileList
	}

	if len(fileList) == 1 {
		return ErrSingleFile
	}

	return nil
}

// ListFiles retrieves a list of files from the specified path.
func (r *RotationManager) ListFiles(path string) ([]*File, error) {
	infos, err := r.provider.ListFiles(path)
	if err != nil {
		return nil, err
	}

	fileList := make([]*File, len(infos))
	for i, info := range infos {
		fileList[i] = &File{
			Path:      info.Path,
			Size:      info.Size,
			Timestamp: info.Timestamp,
		}
	}
	return fileList, nil
}

// RemoveFile deletes a file from the filesystem using the specified full path.
func (r *RotationManager) RemoveFile(fullPath string) error {
	return r.provider.Delete(fullPath)
}

// RotateFiles retrieves the files and categorizes them based on the rotation scheme and the current time.
func (r *RotationManager) RotateFiles() (*Summary, error) {
	fileList, err := r.ListFiles(r.path)
	if err != nil {
		return nil, err
	}

	if err := r.Validate(fileList); err != nil {
		return nil, err
	}

	return RotateFilesOf(fileList, r.rotationScheme, carbon.Now()), nil
}

// RotateFilesOf categorizes the files based on the rotation scheme and the current time.
func RotateFilesOf(files []*File, scheme *RotationScheme, current carbon.Carbon) *Summary {
	sort.Sort(Files(files))

	var hourly, daily, weekly, monthly, yearly, forDelete Files
	var prevYearly, prevMonthly, prevWeekly, prevDaily, prevHourly *carbon.Carbon
	var totalSizeHourly, totalSizeDaily, totalSizeWeekly, totalSizeMonthly, totalSizeYearly, totalSizeForDelete int64

	for _, file := range files {
		addedToCategory := false

		if file.IsHourlyOf(current, prevHourly) && len(hourly) < scheme.Hourly {
			hourly = append(hourly, file)
			prevHourly = &file.Timestamp
			totalSizeHourly += file.Size
			addedToCategory = true
		}

		if file.IsDailyOf(current, prevDaily) && len(daily) < scheme.Daily {
			daily = append(daily, file)
			prevDaily = &file.Timestamp
			totalSizeDaily += file.Size
			addedToCategory = true
		}

		if file.IsWeeklyOf(current, prevWeekly, scheme.Weekly) && len(weekly) < scheme.Weekly {
			weekly = append(weekly, file)
			prevWeekly = &file.Timestamp
			totalSizeWeekly += file.Size
			addedToCategory = true
		}

		if file.IsMonthlyOf(current, prevMonthly) && len(monthly) < scheme.Monthly {
			monthly = append(monthly, file)
			prevMonthly = &file.Timestamp
			totalSizeMonthly += file.Size
			addedToCategory = true
		}

		if file.IsYearlyOf(current, prevYearly) {
			if scheme.Yearly == -1 || len(yearly) < scheme.Yearly {
				yearly = append(yearly, file)
				prevYearly = &file.Timestamp
				totalSizeYearly += file.Size
				addedToCategory = true
			}
		}

		if !addedToCategory {
			forDelete = append(forDelete, file)
			totalSizeForDelete += file.Size
		}
	}

	return &Summary{
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
