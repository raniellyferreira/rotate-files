package handler

import (
	"log"
	"os"
	"rotate/src/config"
	"rotate/src/helper"
	"strings"

	"github.com/thatisuday/commando"
)

func HandlerRotate(args map[string]commando.ArgValue, flags map[string]commando.FlagValue) {

	path := args["path"].Value

	log.Println("Starting rotation on", path)

	hourlyInt, _ := flags[config.HOURLY_FLAG].GetInt()
	dailyInt, _ := flags[config.DAILY_FLAG].GetInt()
	weeklyInt, _ := flags[config.WEEKLY_FLAG].GetInt()
	monthlyInt, _ := flags[config.MONTHLY_FLAG].GetInt()
	yearlyInt, _ := flags[config.YEARLY_FLAG].GetInt()
	dryRunBool, _ := flags[config.DRYRUN_FLAG].GetBool()

	rotationScheme := &helper.BackupRotationScheme{
		Hourly:  hourlyInt,
		Daily:   dailyInt,
		Weekly:  weeklyInt,
		Monthly: monthlyInt,
		Yearly:  yearlyInt,
		DryRun:  dryRunBool,
	}

	if isS3BucketPath := strings.HasPrefix(path, "s3://"); isS3BucketPath {
		validateS3EnvironmentsVars()
		performRotateOnS3(path, rotationScheme)
		return
	}

	performRotateLocally(path, rotationScheme)
}

func performRotateOnS3(path string, rotationScheme *helper.BackupRotationScheme) {
	bucket, prefix := helper.GetBucketAndPrefix(path)
	s3Files := helper.GetS3FilesList(bucket, prefix)
	s3FilesCount := len(*s3Files)

	if s3FilesCount == 0 {
		log.Println("No files found to rotate")
		os.Exit(0)
		return
	}

	if s3FilesCount == 1 {
		log.Println("One file is not eligible for rotate")
		os.Exit(0)
		return
	}

	summary := s3Files.Rotate(rotationScheme)

	log.Println("Yearly matched:")
	for _, v := range summary.Yearly {
		log.Println(v.Path, v.Timestamp)
	}

	log.Println("Monthly matched:")
	for _, v := range summary.Monthly {
		log.Println(v.Path, v.Timestamp)
	}

	log.Println("Weekly matched:")
	for _, v := range summary.Weekly {
		log.Println(v.Path, v.Timestamp)
	}

	log.Println("Daily matched:")
	for _, v := range summary.Daily {
		log.Println(v.Path, v.Timestamp)
	}

	log.Println("Hourly matched:")
	for _, v := range summary.Hourly {
		log.Println(v.Path, v.Timestamp)
	}

	log.Println("For Delete matched:")
	for _, v := range summary.ForDelete {
		log.Println(v.Path, v.Timestamp)
	}
}

func performRotateLocally(path string, rotationScheme *helper.BackupRotationScheme) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		log.Println("directory does not exist")
		os.Exit(1)
		return
	}

	files, err := helper.ListDir(path)
	if err != nil {
		log.Println("Error on listing directory ", err.Error())
		os.Exit(1)
		return
	}

	summary := files.Rotate(rotationScheme)

	log.Println("Yearly matched:")
	for _, v := range summary.Yearly {
		log.Println(v.Path, v.Timestamp)
	}

	log.Println("Monthly matched:")
	for _, v := range summary.Monthly {
		log.Println(v.Path, v.Timestamp)
	}

	log.Println("Weekly matched:")
	for _, v := range summary.Weekly {
		log.Println(v.Path, v.Timestamp)
	}

	log.Println("Daily matched:")
	for _, v := range summary.Daily {
		log.Println(v.Path, v.Timestamp)
	}

	log.Println("Hourly matched:")
	for _, v := range summary.Hourly {
		log.Println(v.Path, v.Timestamp)
	}

	log.Println("For Delete matched:")
	for _, v := range summary.ForDelete {
		log.Println(v.Path, v.Timestamp)
	}
}

func validateS3EnvironmentsVars() {
	if !helper.EnvExists("AWS_ACCESS_KEY_ID") || !helper.EnvExists("AWS_SECRET_ACCESS_KEY") {
		log.Fatal("To use S3 is required defined AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY environment vars ou set ~/.aws/credentials file")
	}
}
