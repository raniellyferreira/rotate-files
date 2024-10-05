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
	"errors"
	"testing"

	"github.com/golang-module/carbon"
	"github.com/raniellyferreira/rotate-files/pkg/providers"
	"github.com/raniellyferreira/rotate-files/pkg/rotate"
)

// DummyProvider is a mock implementation of the Provider interface for testing.
type DummyProvider struct {
	files []*providers.FileInfo
	err   error
}

func (d *DummyProvider) ListFiles(path string) ([]*providers.FileInfo, error) {
	return d.files, d.err
}

func (d *DummyProvider) Delete(path string) error {
	return d.err
}

func TestRotationManager_RotateFiles(t *testing.T) {
	// Setup a dummy provider with test files
	files := []*providers.FileInfo{
		{Path: "file1", Size: 100, Timestamp: carbon.Now().SubHours(1)},
		{Path: "file2", Size: 200, Timestamp: carbon.Now().SubDays(1)},
	}
	provider := &DummyProvider{files: files, err: nil}

	// Define a rotation scheme
	scheme := &rotate.RotationScheme{
		Hourly:  1,
		Daily:   1,
		Weekly:  1,
		Monthly: 1,
		Yearly:  1,
	}

	manager := rotate.NewRotationManager(provider, scheme, "dummy/path")
	summary, err := manager.RotateFiles()

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(summary.Hourly) != 1 || len(summary.Daily) != 1 {
		t.Errorf("unexpected number of rotated files")
	}
}

func TestRotationManager_EmptyFileList(t *testing.T) {
	provider := &DummyProvider{files: []*providers.FileInfo{}, err: nil}
	scheme := &rotate.RotationScheme{}
	manager := rotate.NewRotationManager(provider, scheme, "dummy/path")

	_, err := manager.RotateFiles()
	if err == nil || !errors.Is(err, rotate.ErrEmptyFileList) {
		t.Errorf("expected ErrEmptyFileList, got %v", err)
	}
}

func TestRotationManager_ListFilesError(t *testing.T) {
	provider := &DummyProvider{files: nil, err: errors.New("list error")}
	scheme := &rotate.RotationScheme{}
	manager := rotate.NewRotationManager(provider, scheme, "dummy/path")

	_, err := manager.RotateFiles()
	if err == nil || err.Error() != "list error" {
		t.Errorf("expected 'list error', got %v", err)
	}
}

func TestRotationManager_SingleFileError(t *testing.T) {
	files := []*providers.FileInfo{
		{Path: "file1", Size: 100, Timestamp: carbon.Now().SubHours(1)},
	}
	provider := &DummyProvider{files: files, err: nil}
	scheme := &rotate.RotationScheme{}
	manager := rotate.NewRotationManager(provider, scheme, "dummy/path")

	_, err := manager.RotateFiles()
	if err == nil || !errors.Is(err, rotate.ErrSingleFile) {
		t.Errorf("expected ErrSingleFile, got %v", err)
	}
}

func TestRotationManager_NilRotationScheme(t *testing.T) {
	files := []*providers.FileInfo{
		{Path: "file1", Size: 100, Timestamp: carbon.Now().SubHours(1)},
		{Path: "file2", Size: 200, Timestamp: carbon.Now().SubDays(1)},
	}
	provider := &DummyProvider{files: files, err: nil}
	manager := rotate.NewRotationManager(provider, nil, "dummy/path")

	_, err := manager.RotateFiles()
	if err == nil || !errors.Is(err, rotate.ErrNilRotationScheme) {
		t.Errorf("expected ErrNilRotationScheme, got %v", err)
	}
}

func TestRotationManager_RemoveFile(t *testing.T) {
	// Setup a dummy provider with no initial error
	provider := &DummyProvider{
		files: []*providers.FileInfo{
			{Path: "file1", Size: 100, Timestamp: carbon.Now().SubHours(1)},
		},
		err: nil,
	}

	// Define a rotation scheme
	scheme := &rotate.RotationScheme{}
	manager := rotate.NewRotationManager(provider, scheme, "dummy/path")

	// Attempt to remove a file
	err := manager.RemoveFile("file1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

// Test error handling in Delete
func TestRotationManager_RemoveFileError(t *testing.T) {
	// Setup a dummy provider with an error on delete
	provider := &DummyProvider{
		files: []*providers.FileInfo{
			{Path: "file1", Size: 100, Timestamp: carbon.Now().SubHours(1)},
		},
		err: errors.New("delete error"),
	}

	// Define a rotation scheme
	scheme := &rotate.RotationScheme{}
	manager := rotate.NewRotationManager(provider, scheme, "dummy/path")

	// Attempt to remove a file
	err := manager.RemoveFile("file1")
	if err == nil || err.Error() != "delete error" {
		t.Errorf("expected 'delete error', got %v", err)
	}
}
