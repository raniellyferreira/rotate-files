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

package google

import (
	"context"
	"log"

	"cloud.google.com/go/storage"
	"github.com/golang-module/carbon"
	"github.com/raniellyferreira/rotate-files/internal/environment"
	"github.com/raniellyferreira/rotate-files/pkg/rotate"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

var (
	GCSClient *storage.Client
)

func loadGCSConfig() error {
	// Carrega as credenciais do arquivo JSON
	credsFile := environment.GetEnv("GOOGLE_APPLICATION_CREDENTIALS", "")

	opts := []option.ClientOption{}
	if credsFile != "" {
		opts = append(opts, option.WithCredentialsFile(credsFile))
	}

	client, err := storage.NewClient(context.Background(), opts...)
	if err != nil {
		return err
	}

	GCSClient = client

	return nil
}

func createGCSClient() (*storage.Client, error) {
	if GCSClient != nil {
		return GCSClient, nil
	}

	err := loadGCSConfig()
	if err != nil {
		return nil, err
	}

	return GCSClient, nil
}

func GetAllGCSFiles(bucket, prefix string) ([]*storage.ObjectAttrs, error) {
	client, err := createGCSClient()
	if err != nil {
		return nil, err
	}

	it := client.Bucket(bucket).Objects(context.Background(), &storage.Query{
		Prefix: prefix,
	})

	var files []*storage.ObjectAttrs
	for {
		obj, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		files = append(files, obj)
	}

	return files, nil
}

func GetGCSFilesList(bucket, prefix string) *rotate.BackupFiles {
	result, err := GetAllGCSFiles(bucket, prefix)
	if err != nil {
		log.Fatal(err)
	}

	backups := rotate.BackupFiles{}
	for _, obj := range result {
		backups = append(backups, rotate.Backup{
			Bucket:    bucket,
			Path:      obj.Name,
			Size:      obj.Size,
			Timestamp: carbon.CreateFromTimestamp(obj.Created.Unix()),
		})
	}

	return &backups
}

func DeleteGCSFile(bucket, path string) error {
	obj := GCSClient.Bucket(bucket).Object(path)
	return obj.Delete(context.Background())
}
