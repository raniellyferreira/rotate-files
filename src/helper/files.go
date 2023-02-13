package helper

import (
	"os"
	"path/filepath"
	"time"

	"github.com/golang-module/carbon"
)

func GetFileInfo(path string) (string, time.Time, error) {
	file, err := os.Stat(path)

	if err != nil {
		return "", time.Time{}, err
	}

	return path, file.ModTime(), nil
}

func ListDir(path string) (BackupFiles, error) {
	var files BackupFiles

	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			file, date, err := GetFileInfo(path)
			if err != nil {
				return err
			}
			files = append(files, Backup{Path: file, Timestamp: carbon.FromStdTime(date)})
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return files, nil
}
