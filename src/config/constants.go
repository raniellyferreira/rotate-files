package config

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

const (
	HOURLY_FLAG  = "hourly"
	DAILY_FLAG   = "daily"
	WEEKLY_FLAG  = "weekly"
	MONTHLY_FLAG = "monthly"
	YEARLY_FLAG  = "yearly"
	DRYRUN_FLAG  = "dry-run"
)

const (
	HOURLY_SHORT_FLAG  = "h"
	DAILY_SHORT_FLAG   = "d"
	WEEKLY_SHORT_FLAG  = "w"
	MONTHLY_SHORT_FLAG = "m"
	YEARLY_SHORT_FLAG  = "y"
	DRYRUN_SHORT_FLAG  = "D"
)

const (
	DEFAULT_HOURLY  = 24
	DEFAULT_DAILY   = 7
	DEFAULT_WEEKLY  = 14
	DEFAULT_MONTHLY = 12
	DEFAULT_YEARLY  = -1
)
