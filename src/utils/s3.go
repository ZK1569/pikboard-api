package util

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	errs "github.com/zk1569/pikboard-api/src/errors"
)

type S3Connection struct {
	Client *s3.Client
}

var singleS3Intance *S3Connection

func GetS3Instance() *S3Connection {
	if singleS3Intance == nil {
		lock.Lock()
		defer lock.Unlock()

		if singleS3Intance == nil {
			creds := aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(
				EnvVariable.S3.AccessKey,
				EnvVariable.S3.SecretAccessKey,
				"",
			))

			cfg, err := awsConfig.LoadDefaultConfig(context.TODO(),
				awsConfig.WithCredentialsProvider(creds),
				awsConfig.WithRegion(EnvVariable.S3.Region),
			)

			if err != nil {
				log.Fatalf("%s - ❌ Unable to load SDK config: %v", errs.S3Error, err)
			}

			client := s3.NewFromConfig(cfg)

			singleS3Intance = &S3Connection{
				Client: client,
			}
		}
	}

	return singleS3Intance
}
