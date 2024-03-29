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

package aws

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/golang-module/carbon"
	"github.com/raniellyferreira/rotate-files/internal/environment"
	"github.com/raniellyferreira/rotate-files/pkg/rotate"
)

var ClientS3 *s3.Client

func loadConfig() {
	if ClientS3 != nil {
		return
	}

	var err error
	var cfg aws.Config

	region := environment.GetEnv("AWS_REGION", "us-east-1")
	endpointOverride := environment.GetEnv("AWS_ENDPOINT_OVERRIDE", "")

	if endpointOverride == "" {
		// Using the SDK's default configuration, loading additional config
		// and credentials values from the environment variables, shared
		// credentials, and shared configuration files
		cfg, err = config.LoadDefaultConfig(
			context.TODO(),
			config.WithRegion(region),
		)
	} else {
		customResolver := aws.EndpointResolverWithOptionsFunc(func(service, signingRegion string, options ...interface{}) (aws.Endpoint, error) {
			return aws.Endpoint{
				URL:               endpointOverride,
				HostnameImmutable: true,
				SigningName:       service,
				SigningRegion:     signingRegion,
			}, nil
		})
		cfg, err = config.LoadDefaultConfig(
			context.TODO(),
			config.WithRegion(region),
			config.WithEndpointResolverWithOptions(customResolver),
		)
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

func GetAllS3Files(bucket, prefix string) ([]types.Object, error) {
	var continuationToken *string
	files := make([]types.Object, 0)
	for {
		resp, err := ClientS3.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
			Bucket:            aws.String(bucket),
			Prefix:            aws.String(prefix),
			MaxKeys:           1e4, //10.000
			ContinuationToken: continuationToken,
		})
		if err != nil {
			return nil, err
		}

		files = append(files, resp.Contents...)

		// Verifica se houve paginação e atualiza o continuationToken
		if !resp.IsTruncated {
			break
		}
		continuationToken = resp.NextContinuationToken
	}
	return files, nil
}

func GetS3FilesList(bucket, prefix string) *rotate.BackupFiles {
	loadConfig()

	result, err := GetAllS3Files(bucket, prefix)
	if err != nil {
		log.Fatal(err)
	}

	backups := rotate.BackupFiles{}

	for _, obj := range result {
		backups = append(backups, rotate.Backup{
			Bucket:    bucket,
			Path:      *obj.Key,
			Size:      obj.Size,
			Timestamp: carbon.FromStdTime(*obj.LastModified),
		})
	}

	return &backups
}
