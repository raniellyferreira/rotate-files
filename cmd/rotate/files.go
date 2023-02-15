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

package main

import (
	"os"
	"path/filepath"
	"time"

	"github.com/golang-module/carbon"

	"github.com/raniellyferreira/rotate-files/v1/pkg/rotation"
)

func GetFileInfo(path string) (string, time.Time, error) {
	file, err := os.Stat(path)

	if err != nil {
		return "", time.Time{}, err
	}

	return path, file.ModTime(), nil
}

func ListDir(path string) (rotation.BackupFiles, error) {
	var files rotation.BackupFiles

	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			file, date, err := GetFileInfo(path)
			if err != nil {
				return err
			}
			files = append(files, rotation.Backup{Path: file, Timestamp: carbon.FromStdTime(date)})
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return files, nil
}

func DeleteLocalFile(path string) error {
	return os.Remove(path)
}
