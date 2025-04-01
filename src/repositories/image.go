package repository

import (
	"bytes"
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	util "github.com/zk1569/pikboard-api/src/utils"
)

type Image struct {
	PresignClient *s3.PresignClient
	s3Client      *util.S3Connection
}

var singleImageInstance ImageInterface

func GetImageInstance() ImageInterface {
	if singleImageInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if singleImageInstance == nil {
			singleImageInstance = &Image{
				s3Client: util.GetS3Instance(),
			}
		}
	}

	return singleImageInstance
}

func (self *Image) UploadImage(fileName string, img []byte, ext string) (string, error) {

	bodyReader := bytes.NewReader(img)

	_, err := self.s3Client.Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(util.EnvVariable.S3.BucketName),
		Key:    aws.String("profile_image/" + fileName + "." + ext),
		Body:   bodyReader,
	})
	if err != nil {
		return "", err
	}

	publicURL := fmt.Sprintf("https://pikboard.s3.eu-west-1.amazonaws.com/profile_image/%s.%s", fileName, ext)
	return publicURL, nil
}

func (self *Image) UploadForChat(name string, img []byte) (string, error) {
	bodyReader := bytes.NewReader(img)

	_, err := self.s3Client.Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(util.EnvVariable.S3.BucketName),
		Key:    aws.String("chat/" + name + ".png"),
		Body:   bodyReader,
	})
	if err != nil {
		return "", err
	}

	publicURL := fmt.Sprintf("https://pikboard.s3.eu-west-1.amazonaws.com/chat/%s.%s", name, "png")
	return publicURL, nil

}
