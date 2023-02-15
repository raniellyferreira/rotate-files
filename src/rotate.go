package main

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

import (
	"rotate/src/config"
	"rotate/src/handler"
	"strings"

	"github.com/joho/godotenv"
	"github.com/thatisuday/commando"
)

func main() {
	_ = godotenv.Load()

	commando.
		SetExecutableName("rotate").
		SetVersion("1.0.3").
		SetDescription("Rotate files locally or in S3 bucket based on backup rotation scheme")

	// configure the root command
	commando.
		Register(nil).
		AddArgument(
			"path",
			"local directory path or s3:// path",
			"").
		AddFlag(
			strings.Join([]string{config.HOURLY_FLAG, config.HOURLY_SHORT_FLAG}, ","),
			"number of hourly backups to preserve",
			commando.Int,
			config.DEFAULT_HOURLY).
		AddFlag(
			strings.Join([]string{config.DAILY_FLAG, config.DAILY_SHORT_FLAG}, ","),
			"number of daily backups to preserve",
			commando.Int,
			config.DEFAULT_DAILY).
		AddFlag(
			strings.Join([]string{config.WEEKLY_FLAG, config.WEEKLY_SHORT_FLAG}, ","),
			"number of weekly backups to preserve",
			commando.Int,
			config.DEFAULT_WEEKLY).
		AddFlag(
			strings.Join([]string{config.MONTHLY_FLAG, config.MONTHLY_SHORT_FLAG}, ","),
			"number of monthly backups to preserve",
			commando.Int,
			config.DEFAULT_MONTHLY).
		AddFlag(
			strings.Join([]string{config.YEARLY_FLAG, config.YEARLY_SHORT_FLAG}, ","),
			"number of yearly backups to preserve, set 0 to preserver always",
			commando.Int,
			config.DEFAULT_YEARLY).
		AddFlag(
			strings.Join([]string{config.DRYRUN_FLAG, config.DRYRUN_SHORT_FLAG}, ","),
			"simulate deletion process",
			commando.Bool,
			false).
		SetAction(handler.HandlerRotate)

	commando.Parse(nil)
}
