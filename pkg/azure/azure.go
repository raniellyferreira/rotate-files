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

package azure

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/golang-module/carbon"
	"github.com/raniellyferreira/rotate-files/internal/environment"
	"github.com/raniellyferreira/rotate-files/pkg/rotate"
)

var (
	Client *azblob.Client
)

func newBlobStorageClient(storageAccountName string) (*azblob.Client, error) {
	connectionString := environment.GetEnv("AZURE_STORAGE_CONNECTION_STRING", "")

	var err error
	var client *azblob.Client

	if connectionString == "" {
		credentials, err := azidentity.NewDefaultAzureCredential(nil)
		if err != nil {
			return nil, err
		}
		client, err = azblob.NewClient(fmt.Sprintf("https://%s.blob.core.windows.net", storageAccountName), credentials, nil)
		if err != nil {
			return nil, err
		}
	} else {
		client, err = azblob.NewClientFromConnectionString(connectionString, nil)
		if err != nil {
			return nil, err
		}
	}

	return client, nil
}

func GetAllFiles(accountName, container, prefix string) (*rotate.BackupFiles, error) {
	if Client == nil {
		var err error
		Client, err = newBlobStorageClient(accountName)
		if err != nil {
			return nil, err
		}
	}

	files := make(rotate.BackupFiles, 0)
	pager := Client.NewListBlobsFlatPager(container, &azblob.ListBlobsFlatOptions{
		Prefix: &prefix,
	})
	for pager.More() {
		resp, err := pager.NextPage(context.TODO())
		if err != nil {
			log.Fatal(err)
		}

		for _, blob := range resp.Segment.BlobItems {
			files = append(files, rotate.Backup{
				Bucket:    container,
				Path:      *blob.Name,
				Size:      *blob.Properties.ContentLength,
				Timestamp: carbon.FromStdTime(*blob.Properties.CreationTime),
			})
		}
	}

	return &files, nil
}

func DeleteFile(file *rotate.Backup) error {
	_, err := Client.DeleteBlob(context.Background(), file.Bucket, file.Path, nil)
	return err
}

func GetAccountContainerAndPath(fullPath string) (string, string, string) {
	parts := strings.SplitN(fullPath, "://", 2)
	if len(parts) < 2 {
		return "", "", ""
	}

	pathParts := strings.SplitN(parts[1], "/", 3)
	if len(parts) < 2 {
		return "", "", ""
	}

	path := ""
	if len(pathParts) > 2 && strings.TrimSpace(pathParts[2]) != "" {
		path = pathParts[2]
	}

	return pathParts[0], pathParts[1], path
}
