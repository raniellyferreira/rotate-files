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
	"context"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/golang-module/carbon"

	"github.com/raniellyferreira/rotate-files/v1/pkg/rotation"
)

var ClientS3 *s3.Client

func loadConfig() {
	if ClientS3 != nil {
		return
	}

	profile := GetEnv("AWS_PROFILE", "")

	var err error
	var cfg aws.Config

	if profile == "" {
		cfg, err = config.LoadDefaultConfig(context.TODO())
	} else {
		cfg, err = config.LoadDefaultConfig(context.TODO(), config.WithSharedConfigProfile(profile))
	}

	if err != nil {
		log.Fatalf("failed to load aws configuration, %v", err)
	}

	ClientS3 = s3.NewFromConfig(cfg)
}

func DeleteS3File(bucket, path string) error {
	_, err := ClientS3.DeleteObject(context.Background(), &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(path),
	})
	return err
}

func GetS3FilesList(bucket, prefix string) *rotation.BackupFiles {
	loadConfig()

	// TODO fazer loop para pegar todos os itens
	result, err := ClientS3.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket:  aws.String(bucket),
		Prefix:  aws.String(prefix),
		MaxKeys: 1e5,
	})

	if err != nil {
		log.Fatal(err)
	}

	var backups = rotation.BackupFiles{}
	for _, obj := range result.Contents {
		backups = append(backups, rotation.Backup{
			Bucket:    bucket,
			Path:      *obj.Key,
			Timestamp: carbon.FromStdTime(*obj.LastModified),
		})
	}

	return &backups
}

func GetBucketAndPrefix(fullPath string) (string, string) {
	path := strings.SplitN(strings.TrimPrefix(fullPath, "s3://"), "/", 2)
	prefix := ""
	if len(path) > 1 && strings.TrimSpace(path[1]) != "" {
		prefix = path[1]
	}
	return path[0], prefix
}