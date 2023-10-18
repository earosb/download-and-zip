package storage

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type DOStorage struct {
	client *manager.Downloader
}

func newDOStorage() ClientManager {
	ctx := context.Background()

	spacesKey := os.Getenv("DO_SPACES_KEY")
	spacesSecret := os.Getenv("DO_SPACES_SECRET")

	creds := credentials.NewStaticCredentialsProvider(spacesKey, spacesSecret, "")

	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL: os.Getenv("DO_SPACES_ENDPOINT"),
		}, nil
	})
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion("us-east-1"),
		config.WithCredentialsProvider(creds),
		config.WithEndpointResolverWithOptions(customResolver))
	if err != nil {
		log.Panic(err)
	}
	// Create an Amazon S3 service client
	awsS3Client := s3.NewFromConfig(cfg)

	downloader := manager.NewDownloader(awsS3Client)

	return &DOStorage{
		client: downloader,
	}
}

func (d DOStorage) Download(destFolder string, filepath string) {
	ctx := context.Background()

	fullURLFileParts := strings.Split(filepath, "/")
	filename := fullURLFileParts[len(fullURLFileParts)-1]

	input := &s3.GetObjectInput{
		Bucket: aws.String(os.Getenv("DO_SPACES_BUCKET")),
		Key:    aws.String(filepath),
	}

	newFile, err := os.Create(destFolder + "/" + filename)
	if err != nil {
		log.Panic(err)
	}
	defer newFile.Close()

	_, err = d.client.Download(ctx, newFile, input)
	if err != nil {
		log.Panic(err)
	}

}
