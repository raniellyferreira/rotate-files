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

package files

import (
	"os"
	"path/filepath"

	"github.com/golang-module/carbon"
	"github.com/raniellyferreira/rotate-files/pkg/providers"
)

type LocalProvider struct{}

// NewLocalProvider creates a new instance of LocalProvider for managing local file operations.
func NewLocalProvider() *LocalProvider {
	return &LocalProvider{}
}

// Delete removes a file from the local filesystem using the specified full path.
func (l *LocalProvider) Delete(fullPath string) error {
	return os.Remove(fullPath)
}

// ListFiles traverses the local directory specified by fullPath and returns a list of files.
func (l *LocalProvider) ListFiles(fullPath string) ([]*providers.FileInfo, error) {
	var files []*providers.FileInfo

	err := filepath.Walk(fullPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			files = append(files, &providers.FileInfo{
				Path:      path,
				Size:      info.Size(),
				Timestamp: carbon.FromStdTime(info.ModTime()),
			})
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return files, nil
}
