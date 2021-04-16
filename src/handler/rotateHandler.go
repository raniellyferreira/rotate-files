package handler

import (
	"log"
	"rotate/src/config"
	"rotate/src/helper"
	"strings"

	"github.com/thatisuday/commando"
)

func HanderRotate(args map[string]commando.ArgValue, flags map[string]commando.FlagValue) {

	path := args["path"].Value

	log.Println("Starting rotation on", path)

	hourlyInt, _ := flags[config.HOURLY_FLAG].GetInt()
	dailyInt, _ := flags[config.DAILY_FLAG].GetInt()
	weeklyInt, _ := flags[config.WEEKLY_FLAG].GetInt()
	monthlyInt, _ := flags[config.MONTHLY_FLAG].GetInt()
	yearlyInt, _ := flags[config.YEARLY_FLAG].GetInt()
	dryRunBool, _ := flags[config.DRYRUN_FLAG].GetBool()

	rotationScheme := &helper.RotationScheme{
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

func performRotateOnS3(path string, rotationScheme *helper.RotationScheme) {
	bucket, prefix := helper.GetBucketAndPrefix(path)
	s3Files := helper.GetS3FilesList(bucket, prefix)
	helper.RotateObjects(*s3Files, rotationScheme)

	for _, v := range *s3Files {
		log.Println(*v.Path, v.Status)
	}
}

func performRotateLocally(path string, rotationScheme *helper.RotationScheme) {}

func validateS3EnvironmentsVars() {
	if !helper.EnvExists("AWS_ACCESS_KEY_ID") || !helper.EnvExists("AWS_SECRET_ACCESS_KEY") {
		log.Fatal("To use S3 is required defined AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY environment vars ou set ~/.aws/credentials file")
	}
}
