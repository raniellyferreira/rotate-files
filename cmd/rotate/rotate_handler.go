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
	"errors"
	"log"
	"strings"

	"github.com/raniellyferreira/rotate-files/pkg/aws"
	"github.com/raniellyferreira/rotate-files/pkg/azure"
	"github.com/raniellyferreira/rotate-files/pkg/files"
	"github.com/raniellyferreira/rotate-files/pkg/google"
	"github.com/raniellyferreira/rotate-files/pkg/providers"
	"github.com/raniellyferreira/rotate-files/pkg/rotate"
	"github.com/thatisuday/commando"
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

	rotationScheme := &rotate.RotationScheme{
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

	provider, err := initializeProvider(path)
	if err != nil {
		log.Fatal("Failed to initialize provider:", err)
	}

	manager := rotate.NewRotationManager(provider)
	backups, err := manager.ListFiles(path)
	if err != nil {
		log.Fatal("Error listing backups:", err)
	}

	if len(backups) == 0 {
		log.Println("No files found to rotate")
		return
	}

	summary := manager.RotateFiles(backups, rotationScheme)

	if len(summary.ForDelete) > 0 {
		if !rotationScheme.DryRun {
			for _, backup := range summary.ForDelete {
				log.Println("Deleting file...", backup.Path)
				if err := manager.RemoveFile(backup.Path); err != nil {
					log.Println("Error deleting file:", err)
				}
			}
		} else {
			for _, backup := range summary.ForDelete {
				log.Println("DRYRUN: simulate file delete...", backup.Path)
			}
		}
	} else {
		log.Println("No files eligible for deletion")
	}

	summary.Print()
}

func initializeProvider(path string) (providers.Provider, error) {
	prov := strings.SplitN(path, "://", 2)

	if len(prov) == 0 {
		return nil, errors.New("invalid path")
	}

	switch prov[0] {
	case "s3":
		return aws.NewAWSProvider()
	case "gc", "gs", "gcs":
		return google.NewGoogleProvider()
	case "blob", "azure", "azr", "az":
		return azure.NewAzureProvider()
	default:
		return files.NewLocalProvider(), nil
	}
}
