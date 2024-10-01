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
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/golang-module/carbon"
	"github.com/raniellyferreira/rotate-files/internal/environment"
	"github.com/raniellyferreira/rotate-files/pkg/providers"
	"github.com/raniellyferreira/rotate-files/pkg/utils"
)

type AWSProvider struct {
	client *s3.Client
}

func NewAWSProvider() (*AWSProvider, error) {
	region := environment.GetEnv("AWS_REGION", "us-east-1")
	endpoint := environment.GetEnv("AWS_ENDPOINT_OVERRIDE", "")

	var cfg aws.Config
	var err error

	if endpoint != "" {
		customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			return aws.Endpoint{
				URL:               endpoint,
				HostnameImmutable: true,
			}, nil
		})
		cfg, err = config.LoadDefaultConfig(context.Background(),
			config.WithRegion(region),
			config.WithEndpointResolverWithOptions(customResolver),
		)
	} else {
		cfg, err = config.LoadDefaultConfig(context.Background(), config.WithRegion(region))
	}

	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(cfg)
	return &AWSProvider{client: client}, nil
}

func (a *AWSProvider) Delete(path string) error {
	bucket, key := utils.GetBucketAndKey(path)
	_, err := a.client.DeleteObject(context.Background(), &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	return err
}

func (a *AWSProvider) ListFiles(fullPath string) ([]*providers.BackupInfo, error) {
	var continuationToken *string
	var files []*providers.BackupInfo

	bucket, path := utils.GetBucketAndKey(fullPath)

	for {
		resp, err := a.client.ListObjectsV2(context.Background(), &s3.ListObjectsV2Input{
			Bucket:            aws.String(bucket),
			Prefix:            aws.String(path),
			ContinuationToken: continuationToken,
		})
		if err != nil {
			return nil, err
		}

		for _, obj := range resp.Contents {
			files = append(files, &providers.BackupInfo{
				Path:      fmt.Sprintf("s3://%s/%s", bucket, aws.ToString(obj.Key)),
				Size:      aws.ToInt64(obj.Size),
				Timestamp: carbon.FromStdTime(aws.ToTime(obj.LastModified)),
			})
		}

		if aws.ToBool(resp.IsTruncated) {
			continuationToken = resp.NextContinuationToken
		} else {
			break
		}
	}

	return files, nil
}
