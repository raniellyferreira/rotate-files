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

package utils_test

import (
	"testing"

	"github.com/raniellyferreira/rotate-files/pkg/utils"
)

func TestGetAccountContainerAndPath(t *testing.T) {
	tests := []struct {
		input             string
		expectedAccount   string
		expectedContainer string
		expectedPath      string
	}{
		{"blob://account/container/path/to/file", "account", "container", "path/to/file"},
		{"blob://account/container/", "account", "container", ""},
		{"blob://account/container", "account", "container", ""},
		{"blob://account/", "", "", ""},
		{"blob://account", "", "", ""},
		{"invalid", "", "", ""},
		{"", "", "", ""},
	}

	for _, test := range tests {
		account, container, path := utils.GetAccountContainerAndPath(test.input)
		if account != test.expectedAccount || container != test.expectedContainer || path != test.expectedPath {
			t.Errorf("Input: %s - Expected: (%s, %s, %s), Got: (%s, %s, %s)",
				test.input, test.expectedAccount, test.expectedContainer, test.expectedPath, account, container, path)
		}
	}
}

func TestGetBucketAndKey(t *testing.T) {
	tests := []struct {
		input          string
		expectedBucket string
		expectedKey    string
	}{
		{"s3://my-bucket/my-key/path", "my-bucket", "my-key/path"},
		{"s3://my-bucket/", "my-bucket", ""},
		{"s3://my-bucket", "my-bucket", ""},
		{"gs://another-bucket/another-key", "another-bucket", "another-key"},
		{"invalid-path", "invalid-path", ""},
		{"", "", ""},
	}

	for _, test := range tests {
		bucket, key := utils.GetBucketAndKey(test.input)
		if bucket != test.expectedBucket || key != test.expectedKey {
			t.Errorf("Input: %s - Expected: (%s, %s), Got: (%s, %s)",
				test.input, test.expectedBucket, test.expectedKey, bucket, key)
		}
	}
}
