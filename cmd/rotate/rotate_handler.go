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

// HandlerRotate is the command handler for rotating files.
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

	manager := rotate.NewRotationManager(
		provider,
		rotationScheme,
		path,
	)

	var summary *rotate.Summary
	summary, err = manager.RotateFiles()

	if err != nil {
		switch err {
		case rotate.ErrEmptyFileList:
			log.Println("No files to rotate")
			return
		case rotate.ErrSingleFile:
			log.Println("Only one file to rotate, ignoring rotation")
			return
		default:
			log.Fatal("Unknown error:", err)
		}
	}

	handleFileDeletion(manager, summary, rotationScheme)

	summary.Print()
}

// initializeProvider initializes the provider based on the path.
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
		return azure.NewAzureProvider(path)
	default:
		return files.NewLocalProvider(), nil
	}
}

// handleFileDeletion deletes the files from the file provider.
func handleFileDeletion(manager *rotate.RotationManager, summary *rotate.Summary, scheme *rotate.RotationScheme) {
	if len(summary.ForDelete) == 0 {
		log.Println("No files eligible for deletion")
		return
	}

	if scheme.DryRun {
		simulateDeletion(summary)
	} else {
		executeDeletion(manager, summary)
	}
}

// simulateDeletion prints the files that would be deleted in a dry run.
func simulateDeletion(summary *rotate.Summary) {
	for _, backup := range summary.ForDelete {
		log.Println("DRYRUN: simulate file delete...", backup.Path)
	}
}

// executeDeletion deletes the files from the file provider.
func executeDeletion(manager *rotate.RotationManager, summary *rotate.Summary) {
	for _, backup := range summary.ForDelete {
		log.Println("Deleting file...", backup.Path)
		if err := manager.RemoveFile(backup.Path); err != nil {
			log.Println("Error deleting file:", err)
		}
	}
}
