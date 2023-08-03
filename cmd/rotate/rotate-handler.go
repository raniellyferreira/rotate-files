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
	"log"
	"os"
	"strings"

	"github.com/thatisuday/commando"

	"github.com/raniellyferreira/rotate-files/pkg/rotate"
)

func HandlerRotate(args map[string]commando.ArgValue, flags map[string]commando.FlagValue) {

	path := args["path"].Value

	log.Println("Starting rotation on", path)

	hourlyInt, _ := flags[HOURLY_FLAG].GetInt()
	dailyInt, _ := flags[DAILY_FLAG].GetInt()
	weeklyInt, _ := flags[WEEKLY_FLAG].GetInt()
	monthlyInt, _ := flags[MONTHLY_FLAG].GetInt()
	yearlyInt, _ := flags[YEARLY_FLAG].GetInt()
	dryRunBool, _ := flags[DRYRUN_FLAG].GetBool()

	rotationScheme := &rotate.BackupRotationScheme{
		Hourly:  hourlyInt,
		Daily:   dailyInt,
		Weekly:  weeklyInt,
		Monthly: monthlyInt,
		Yearly:  yearlyInt,
		DryRun:  dryRunBool,
	}

	if rotationScheme.DryRun {
		log.Println(" -------- DryRun mode on")
	}

	if isS3BucketPath := strings.HasPrefix(path, "s3://"); isS3BucketPath {
		performRotateOnS3(path, rotationScheme)
		return
	}

	performRotateLocally(path, rotationScheme)
}

func performRotateOnS3(path string, rotationScheme *rotate.BackupRotationScheme) {
	bucket, prefix := GetBucketAndPrefix(path)
	s3Files := GetS3FilesList(bucket, prefix)

	if s3Files.Len() == 0 {
		log.Println("No files found to rotate")
		os.Exit(0)
		return
	}

	if s3Files.Len() == 1 {
		log.Println("One file is not eligible for rotate")
		os.Exit(0)
		return
	}

	summary := s3Files.Rotate(rotationScheme)

	if !rotationScheme.DryRun {
		for _, v := range summary.ForDelete {
			log.Println("Deleting file...", v.Path)
			if err := DeleteS3File(v.Bucket, v.Path); err != nil {
				log.Println("Error on delete object from S3: ", v.Bucket, v.Path, err)
			}
		}
	} else {
		for _, v := range summary.ForDelete {
			log.Println("DRYRUN: simulate file delete...", v.Path)
		}
	}

	summary.Print()
}

func performRotateLocally(path string, rotationScheme *rotate.BackupRotationScheme) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		log.Println("directory does not exist")
		os.Exit(1)
		return
	}

	files, err := ListDir(path)
	if err != nil {
		log.Println("Error on listing directory ", err.Error())
		os.Exit(1)
		return
	}

	if files.Len() == 0 {
		log.Println("No files found to rotate")
		os.Exit(0)
		return
	}

	if files.Len() == 1 {
		log.Println("One file is not eligible for rotate")
		os.Exit(0)
		return
	}

	summary := files.Rotate(rotationScheme)
	for _, v := range summary.ForDelete {
		if !rotationScheme.DryRun {
			if err := DeleteLocalFile(v.Path); err != nil {
				log.Println("Error on delete local file: ", v.Path, err)
				continue
			}
		}
	}

	summary.Print()
}
