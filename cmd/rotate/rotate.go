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
	"strings"

	"github.com/joho/godotenv"
	"github.com/raniellyferreira/rotate-files/v1/internal/version"
	"github.com/thatisuday/commando"
)

func main() {
	_ = godotenv.Load()

	commando.
		SetExecutableName("rotate").
		SetVersion(version.GetVersion()).
		SetDescription("Rotate files locally or in S3 bucket based on custom backup rotation scheme")

	commando.
		Register(nil).
		AddArgument(
			"path",
			"local directory path or s3:// path",
			"").
		AddFlag(
			strings.Join([]string{HOURLY_FLAG, HOURLY_SHORT_FLAG}, ","),
			"number of hourly backups to preserve",
			commando.Int,
			DEFAULT_HOURLY).
		AddFlag(
			strings.Join([]string{DAILY_FLAG, DAILY_SHORT_FLAG}, ","),
			"number of daily backups to preserve",
			commando.Int,
			DEFAULT_DAILY).
		AddFlag(
			strings.Join([]string{WEEKLY_FLAG, WEEKLY_SHORT_FLAG}, ","),
			"number of weekly backups to preserve",
			commando.Int,
			DEFAULT_WEEKLY).
		AddFlag(
			strings.Join([]string{MONTHLY_FLAG, MONTHLY_SHORT_FLAG}, ","),
			"number of monthly backups to preserve",
			commando.Int,
			DEFAULT_MONTHLY).
		AddFlag(
			strings.Join([]string{YEARLY_FLAG, YEARLY_SHORT_FLAG}, ","),
			"number of yearly backups to preserve, set 0 to preserver always",
			commando.Int,
			DEFAULT_YEARLY).
		AddFlag(
			strings.Join([]string{DRYRUN_FLAG, DRYRUN_SHORT_FLAG}, ","),
			"simulate deletion process",
			commando.Bool,
			false).
		SetAction(HandlerRotate)

	commando.Parse(nil)
}
