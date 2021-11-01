package bucket

import (
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type ImageUploader struct {
	client *s3.S3
	bucket *string
	key    *string
}

func NewImageUploader(sess *session.Session, bucket, key string) *ImageUploader {
	return &ImageUploader{
		client: s3.New(sess),
		bucket: aws.String(bucket),
		key:    aws.String(key),
	}
}

func (i *ImageUploader) UploadImage(imagePath string) (string, error) {
	image, err := os.Open(imagePath)
	if err != nil {
		return "", err
	}
	defer image.Close()
	req, _ := i.client.PutObjectRequest(&s3.PutObjectInput{
		Bucket: i.bucket,
		Key:    i.key,
		Body:   image,
	})
	if err := req.Send(); err != nil {
		return "", nil
	}

	req, _ = i.client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: i.bucket,
		Key:    i.key,
	})
	url, err := req.Presign(7 * 24 * time.Hour)
	if err != nil {
		return "", err
	}

	return url, nil
}
