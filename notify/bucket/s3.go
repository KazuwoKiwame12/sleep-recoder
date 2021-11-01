package bucket

import (
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type ImageUploader struct {
	uploader *s3.S3
	bucket   *string
	key      *string
}

func NewImageUploader(sess *session.Session, bucket, key string) *ImageUploader {
	return &ImageUploader{
		uploader: s3.New(sess),
		bucket:   aws.String(bucket),
		key:      aws.String(key),
	}
}

func (i *ImageUploader) UploadImage(imagePath string) (string, error) {
	image, err := os.Open(imagePath)
	if err != nil {
		return "", err
	}
	defer image.Close()
	req, _ := i.uploader.PutObjectRequest(&s3.PutObjectInput{
		Bucket: i.bucket,
		Key:    i.key,
		Body:   image,
	})

	url, err := req.Presign(24 * 5 * time.Hour)
	if err != nil {
		return "", err
	}

	log.Printf("url: %s\n", url)

	return url, nil
}
