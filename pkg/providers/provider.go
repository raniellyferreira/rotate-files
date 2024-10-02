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

package providers

import (
	"github.com/golang-module/carbon"
)

type FileInfo struct {
	Path      string
	Size      int64
	Timestamp carbon.Carbon
}

// Provider defines the interface for cloud storage operations such as delete and list files.
type Provider interface {
	Delete(fullPath string) error
	ListFiles(fullPath string) ([]*FileInfo, error)
}
