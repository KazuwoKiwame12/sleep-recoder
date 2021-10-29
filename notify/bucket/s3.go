package bucket

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type ImageUploader struct {
	uploader *s3manager.Uploader
	bucket   *string
	key      *string
}

func NewImageUploader(sess *session.Session, bucket, key string) *ImageUploader {
	return &ImageUploader{
		uploader: s3manager.NewUploader(sess),
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
	output, err := i.uploader.Upload(&s3manager.UploadInput{
		Bucket: i.bucket,
		Key:    i.key,
		Body:   image,
	})
	if err != nil {
		return "", err
	}
	return output.Location, nil
}
