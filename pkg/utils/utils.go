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

package utils

import "strings"

// GetBucketAndKey returns the bucket and key from a full path. (for AWS and Google)
func GetBucketAndKey(fullPath string) (string, string) {
	parts := strings.SplitN(strings.TrimRight(fullPath, "/"), "://", 2)
	if len(parts) < 2 {
		return parts[0], ""
	}
	pathParts := strings.SplitN(parts[1], "/", 2)
	bucket := pathParts[0]
	prefix := ""
	if len(pathParts) > 1 && strings.TrimSpace(pathParts[1]) != "" {
		prefix = pathParts[1]
	}
	return bucket, prefix
}

// GetAccountContainerAndPath returns the account, container and path from a full path. (for Azure)
func GetAccountContainerAndPath(fullPath string) (string, string, string) {
	parts := strings.SplitN(strings.TrimRight(fullPath, "/"), "://", 2)
	if len(parts) < 2 {
		return "", "", ""
	}

	pathParts := strings.SplitN(parts[1], "/", 3)
	if len(pathParts) < 2 {
		return "", "", ""
	}

	path := ""
	if len(pathParts) > 2 && strings.TrimSpace(pathParts[2]) != "" && pathParts[2] != "/" {
		path = pathParts[2]
	}

	return pathParts[0], pathParts[1], path
}
