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

	"github.com/raniellyferreira/rotate-files/pkg/aws"
	"github.com/raniellyferreira/rotate-files/pkg/azure"
	"github.com/raniellyferreira/rotate-files/pkg/files"
	"github.com/raniellyferreira/rotate-files/pkg/google"
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

	switch {
	case
		strings.HasPrefix(path, "s3://"):
		performRotateOnS3(path, rotationScheme)
	case
		strings.HasPrefix(path, "gc://"),
		strings.HasPrefix(path, "gs://"),
		strings.HasPrefix(path, "gcs://"):
		performRotateOnGCS(path, rotationScheme)
	case
		strings.HasPrefix(path, "blob://"),
		strings.HasPrefix(path, "azure://"):
		performRotateOnAzureBlobStorage(path, rotationScheme)
	default:
		performRotateLocally(path, rotationScheme)
	}
}

func performRotateLocally(path string, rotationScheme *rotate.BackupRotationScheme) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		log.Println("directory does not exist")
		os.Exit(1)
		return
	}

	filesList, err := files.ListDir(path)
	if err != nil {
		log.Println("Error on listing directory ", err.Error())
		os.Exit(1)
		return
	}

	if filesList.Len() == 0 {
		log.Println("No files found to rotate")
		os.Exit(0)
		return
	}

	if filesList.Len() == 1 {
		log.Println("One file is not eligible for rotate")
		os.Exit(0)
		return
	}

	summary := filesList.Rotate(rotationScheme)

	if len(summary.ForDelete) == 0 {
		log.Println("No files eligible for deletion")
	} else {
		if !rotationScheme.DryRun {
			for _, v := range summary.ForDelete {
				log.Println("Deleting file...", v.Path)
				if err := files.DeleteLocalFile(v.Path); err != nil {
					log.Println("Error on delete local file: ", v.Path, err)
					continue
				}
			}
		} else {
			for _, v := range summary.ForDelete {
				log.Println("DRYRUN: simulate file delete...", v.Path)
			}
		}
	}

	summary.Print()
}

func performRotateOnS3(path string, rotationScheme *rotate.BackupRotationScheme) {
	bucket, prefix := rotate.GetBucketAndPrefix(path)
	s3Files := aws.GetS3FilesList(bucket, prefix)

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

	if len(summary.ForDelete) == 0 {
		log.Println("No files eligible for deletion")
	} else {
		if !rotationScheme.DryRun {
			for _, v := range summary.ForDelete {
				log.Println("Deleting file...", v.Path)
				if err := aws.DeleteS3File(v.Bucket, v.Path); err != nil {
					log.Println("Error on delete object from S3: ", v.Bucket, v.Path, err)
				}
			}
		} else {
			for _, v := range summary.ForDelete {
				log.Println("DRYRUN: simulate file delete...", v.Path)
			}
		}
	}

	summary.Print()
}

func performRotateOnGCS(path string, rotationScheme *rotate.BackupRotationScheme) {
	bucket, prefix := rotate.GetBucketAndPrefix(path)
	gcsFiles := google.GetGCSFilesList(bucket, prefix)

	if gcsFiles.Len() == 0 {
		log.Println("No files found to rotate")
		os.Exit(0)
		return
	}

	if gcsFiles.Len() == 1 {
		log.Println("One file is not eligible for rotate")
		os.Exit(0)
		return
	}

	summary := gcsFiles.Rotate(rotationScheme)

	if len(summary.ForDelete) == 0 {
		log.Println("No files eligible for deletion")
	} else {
		if !rotationScheme.DryRun {
			for _, v := range summary.ForDelete {
				log.Println("Deleting file...", v.Path)
				err := google.DeleteGCSFile(v.Bucket, v.Path)
				if err != nil {
					log.Println("Error on delete object from GCS: ", v.Bucket, v.Path, err)
				}
			}
		} else {
			for _, v := range summary.ForDelete {
				log.Println("DRYRUN: simulate file delete...", v.Path)
			}
		}
	}

	summary.Print()
}

func performRotateOnAzureBlobStorage(path string, rotationScheme *rotate.BackupRotationScheme) {
	accountName, containerName, containerPath := azure.GetAccountContainerAndPath(path)
	if accountName == "" || containerName == "" {
		log.Fatal("Error: You must specify the storage account and container. Example: blob://storage-account-name/container-name")
	}

	files, err := azure.GetAllFiles(accountName, containerName, containerPath)
	if err != nil {
		log.Fatalf("error on get files from blob storage: %v", err)
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

	if len(summary.ForDelete) == 0 {
		log.Println("No files eligible for deletion")
	} else {
		if !rotationScheme.DryRun {
			for _, v := range summary.ForDelete {
				log.Println("Deleting file...", v.Path)
				err := azure.DeleteFile(&v)
				if err != nil {
					log.Println("Error on delete object from blob storage: ", v.Bucket, v.Path, err)
				}
			}
		} else {
			for _, v := range summary.ForDelete {
				log.Println("DRYRUN: simulate file delete...", v.Path)
			}
		}
	}

	summary.Print()
}
