package utils

import (
	"bytes"
	"context"
	"log"

	cred "github.com/ashiqsabith123/user-details-svc/pkg/config"
	interfaces "github.com/ashiqsabith123/user-details-svc/pkg/utils/interface"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
)

type S3Client struct {
	client *s3.Client
}

func NewS3Client(awsConfig cred.Config) interfaces.Utils {

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(awsConfig.AWS.Region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(awsConfig.AWS.AccessKey, awsConfig.AWS.SecretKey, "")),
	)

	if err != nil {
		log.Fatal(err)
	}

	client := s3.NewFromConfig(cfg)

	return &S3Client{client: client}

}

func (S3 *S3Client) UploadPhotos(key string, image []byte) error {

	uploader := manager.NewUploader(S3.client)

	_, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(cred.GetAWSBucket()),
		Key:    aws.String(key),
		Body:   bytes.NewReader(image),
		ACL:    "public-read",
	})

	if err != nil {
		return err
	}

	return nil

}
