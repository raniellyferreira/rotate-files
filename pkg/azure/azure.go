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

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/golang-module/carbon"
	"github.com/raniellyferreira/rotate-files/internal/environment"
	"github.com/raniellyferreira/rotate-files/internal/utils"
	"github.com/raniellyferreira/rotate-files/pkg/providers"
)

type AzureProvider struct {
	client *azblob.Client
}

// NewAzureProvider initializes a new AzureProvider using the connection string from environment variables.
func NewAzureProvider() (*AzureProvider, error) {
	connectionString := environment.GetEnv("AZURE_STORAGE_CONNECTION_STRING", "")
	client, err := azblob.NewClientFromConnectionString(connectionString, nil)
	if err != nil {
		return nil, err
	}
	return &AzureProvider{client: client}, nil
}

// Delete removes a blob from an Azure container using the specified full path.
func (az *AzureProvider) Delete(fullPath string) error {
	_, container, path := utils.GetAccountContainerAndPath(fullPath)
	_, err := az.client.DeleteBlob(context.Background(), container, path, nil)
	return err
}

// ListFiles retrieves and lists all blobs within an Azure container with the given full path.
func (az *AzureProvider) ListFiles(fullPath string) ([]*providers.FileInfo, error) {
	account, container, prefix := utils.GetAccountContainerAndPath(fullPath)
	pager := az.client.NewListBlobsFlatPager(container, &azblob.ListBlobsFlatOptions{Prefix: &prefix})

	var files []*providers.FileInfo
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		if err != nil {
			return nil, err
		}

		for _, blob := range resp.Segment.BlobItems {
			files = append(files, &providers.FileInfo{
				Path:      fmt.Sprintf("blob://%s/%s/%s", account, container, aws.ToString(blob.Name)),
				Size:      aws.ToInt64(blob.Properties.ContentLength),
				Timestamp: carbon.FromStdTime(aws.ToTime(blob.Properties.CreationTime)),
			})
		}
	}

	return files, nil
}
