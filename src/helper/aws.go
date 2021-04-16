package helper

import (
	"context"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var ClientS3 *s3.Client

func loadConfig() {
	if ClientS3 != nil {
		return
	}

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("failed to load aws configuration, %v", err)
	}

	ClientS3 = s3.NewFromConfig(cfg)
}

// TODO desenvolver
func DeleteS3File(bucket, prefix string) {}

func GetS3FilesList(bucket, prefix string) *[]RotationObject {
	loadConfig()

	output, err := ClientS3.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
		Prefix: aws.String(prefix),
	})

	if err != nil {
		log.Fatal(err)
	}

	rotationFiles := []RotationObject{}

	for _, v := range output.Contents {
		rotationFiles = append(rotationFiles, RotationObject{
			Bucket:    &bucket,
			Path:      v.Key,
			CreatedAt: v.LastModified,
			Status:    OBJECT_WAITING,
		})
	}

	return &rotationFiles
}

func GetBucketAndPrefix(fullPath string) (string, string) {
	path := strings.SplitN(strings.TrimPrefix(fullPath, "s3://"), "/", 2)
	prefix := ""
	if len(path) > 1 && strings.TrimSpace(path[1]) != "" {
		prefix = path[1]
	}
	return path[0], prefix
}
