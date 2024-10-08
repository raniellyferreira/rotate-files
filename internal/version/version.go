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

# The script is based off of the Apache-licensed script from Helm,
# the package manager for Go: https://github.com/helm/helm/blob/main/internal/version/version.go
*/

package version // import "github.com/raniellyferreira/rotate-files/internal/version"

import (
	"flag"
	"runtime"
	"strings"
)

var (
	// version is the current version of Helm.
	// Update this whenever making a new release.
	// The version is of the format Major.Minor.Patch[-Prerelease][+BuildMetadata]
	//
	// Increment major number for new feature additions and behavioral changes.
	// Increment minor number for medium feature changes.
	// Increment patch number for bug fixes and performance enhancements.
	version = "1.3.1"

	// metadata is extra build time data
	metadata = ""
	// gitCommit is the git sha1
	gitCommit = ""
	// gitTreeState is the state of the git tree
	gitTreeState = ""
)

// BuildInfo describes the compile time information.
type BuildInfo struct {
	// Version is the current semver.
	Version string `json:"version,omitempty"`
	// GitCommit is the git sha1.
	GitCommit string `json:"git_commit,omitempty"`
	// GitTreeState is the state of the git tree.
	GitTreeState string `json:"git_tree_state,omitempty"`
	// GoVersion is the version of the Go compiler used.
	GoVersion string `json:"go_version,omitempty"`
}

// GetVersion returns the semver string of the version
func GetVersion() string {
	if metadata == "" {
		return version
	}
	return version + "+" + metadata
}

// GetUserAgent returns a user agent for user with an HTTP client
func GetUserAgent() string {
	return "Rotate/" + strings.TrimPrefix(GetVersion(), "v")
}

// Get returns build info
func Get() BuildInfo {
	v := BuildInfo{
		Version:      GetVersion(),
		GitCommit:    gitCommit,
		GitTreeState: gitTreeState,
		GoVersion:    runtime.Version(),
	}

	// HACK(bacongobbler): strip out GoVersion during a test run for consistent test output
	if flag.Lookup("test.v") != nil {
		v.GoVersion = ""
	}
	return v
}
