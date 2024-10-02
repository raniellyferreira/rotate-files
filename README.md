# Rotate Files

[![Go Report Card](https://goreportcard.com/badge/github.com/raniellyferreira/rotate-files)](https://goreportcard.com/report/github.com/raniellyferreira/rotate-files) ![License](https://img.shields.io/github/license/raniellyferreira/rotate-files) [![GoDoc](https://img.shields.io/static/v1?label=godoc&message=reference&color=blue)](https://pkg.go.dev/github.com/raniellyferreira/rotate-files)

## Description

Rotate Files is a project that allows you to rotate files locally or in a compatible third-party storage with S3 API, such as Amazon S3, MinIO, DigitalOcean Spaces, Wasabi, Backblaze B2, Azure Blob Storage, and Google Cloud Storage, based on a custom backup rotation scheme.

## Requirements

To run this project, you need to have the following installed:

- bash
- curl
- openssl

## Installation

To install Rotate Files, run the following command in the terminal:

```bash
curl -SsL https://github.com/raniellyferreira/rotate-files/raw/master/environment/scripts/get | bash
```

## Configuration

You can configure Rotate Files using the following options:

- `-d, --daily`: number of daily files to preserve (default: 7)
- `-D, --dry-run`: simulate deletion process (default: false)
- `-h, --help`: displays usage information of the application or a command
- `-h, --hourly`: number of hourly files to preserve (default: 24)
- `-m, --monthly`: number of monthly files to preserve (default: 12)
- `-v, --version`: displays version number
- `-w, --weekly`: number of weekly files to preserve (default: 14)
- `-y, --yearly`: number of yearly files to preserve, set to 0 for no preservation (default is -1 for preserve always)

## Environment Vars

### Amazon S3

- `AWS_ACCESS_KEY_ID`: The access key ID for your AWS account.
- `AWS_SECRET_ACCESS_KEY`: The secret access key for your AWS account.
- `AWS_REGION`: The AWS region where your S3 bucket is located.
- `AWS_ENDPOINT_OVERRIDE`: The endpoint override for your third-party storage with S3 API based.

Supported third-party storage providers:

- CloudFlare R2
- MinIO
- DigitalOcean Spaces
- Wasabi
- Backblaze B2
- and more...

### Google Cloud Storage

- `GOOGLE_APPLICATION_CREDENTIALS`: The path to your Google Cloud service account key file.

### Azure Blob Storage

- `AZURE_STORAGE_CONNECTION_STRING`: The connection string for your Azure storage account.
- `AZURE_CLIENT_ID`: The client ID for your Azure AD application.
- `AZURE_CLIENT_SECRET`: The client secret for your Azure AD application.
- `AZURE_TENANT_ID`: The tenant ID for your Azure AD application.
- `AZURE_SUBSCRIPTION_ID`: The subscription ID for your Azure account.

**Note: If the `AZURE_STORAGE_CONNECTION_STRING` environment variable is set, all other variables related to Azure Blob Storage will be ignored.**

## Usage

To use Rotate Files, run the following command in the terminal:

```bash
rotate help
```

This will display information about how to use the application and its commands.

## Examples

Here is an example of how to use Rotate Files:

```bash
# Amazon S3
rotate s3://example-bucket/backups/ --hourly 24 --daily 7 --weekly 10 --monthly 12

# Google Cloud Storage
rotate gs://example-bucket/backups/ --hourly 24 --daily 7 --weekly 10 --monthly 12

# Azure Blob Storage
rotate blob://example-storage-account-name/example-container/backups/ --hourly 24 --daily 7 --weekly 10 --monthly 12
```

This command will rotate the files in the specified S3 bucket, preserving 24 hourly files, 7 daily files, 10 weekly files, and 12 monthly files.

## Contribution

If you want to contribute to this project, feel free to submit issues, request features, or submit pull requests.

## License

This project is licensed under the [Apache-2.0 License](https://www.apache.org/licenses/LICENSE-2.0).

## Support

If you have any questions or need help with this project, please contact us at contato@awesomeapi.com.br
