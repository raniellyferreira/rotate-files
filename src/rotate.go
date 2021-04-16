package main

import (
	"log"
	"rotate/src/config"
	"rotate/src/handler"
	"strings"

	"github.com/joho/godotenv"
	"github.com/thatisuday/commando"
)

// interface func(map[string]ArgValue, map[string]FlagValue);

func main() {
	if envErr := godotenv.Load(); envErr != nil {
		log.Fatal(envErr)
	}

	commando.
		SetExecutableName("rotate").
		SetVersion("0.0.1").
		SetDescription("Rotate backups locally or in S3 bucket based on backup rotation scheme")

	// configure the root command
	commando.
		Register(nil).
		AddArgument(
			"path",
			"local directory path or s3:// path",
			"./").
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
			"number of yearly backups to preserve, set -1 to preserver always",
			commando.Int,
			config.DEFAULT_YEARLY).
		AddFlag(
			strings.Join([]string{config.DRYRUN_FLAG, config.DRYRUN_SHORT_FLAG}, ","),
			"simulate deletion process",
			commando.Bool,
			false).
		SetAction(handler.HanderRotate)

	commando.Parse(nil)
}
