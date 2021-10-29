package bucket_test

import (
	"notify/bucket"
	"path/filepath"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

func getImageUploader() *bucket.ImageUploader {
	endPoint := "http://localhost:9000"
	sess := session.Must(session.NewSession(&aws.Config{
		Region:           aws.String("ap-northeast-3"),
		Endpoint:         aws.String(endPoint),
		DisableSSL:       aws.Bool(true),
		S3ForcePathStyle: aws.Bool(true),
	}))
	return bucket.NewImageUploader(sess, "static", "sample.png")
}

// https://docs.min.io/docs/how-to-use-aws-sdk-for-go-with-minio-server.html
// https://github.com/reflet/docker-s3
// https://hub.docker.com/r/minio/minio/
// https://docs.aws.amazon.com/ja_jp/cli/latest/userguide/cli-configure-quickstart.html#cli-configure-quickstart-precedence
func Test_UploadImage(t *testing.T) {
	// 環境変数は認証情報ファイルよりも優先順位が高い
	t.Setenv("AWS_ACCESS_KEY_ID", "SAMPLE")
	t.Setenv("AWS_SECRET_ACCESS_KEY", "SAMPLESAMPLESAMPLE")
	t.Setenv("AWS_DEFAULT_REGION", "ap-northeast-3")

	uploader := getImageUploader()
	path, err := filepath.Abs("../testdata/sample.png")
	if err != nil {
		t.Error(err)
	}
	url, err := uploader.UploadImage(path)
	if err != nil {
		t.Error(err)
	}
	if len(url) == 0 {
		t.Error("The url must be more than one character.")
	}
}
