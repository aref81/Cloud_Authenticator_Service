package datasource

import (
	"bytes"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/labstack/gommon/log"
	"io"
	"os"
)

var (
	accessKey  = os.Getenv("S3_ACCESS_KEY")
	secretKey  = os.Getenv("S3_SECRET_KEY")
	s3Region   = os.Getenv("S3_REGION")
	s3Bucket   = os.Getenv("S3_BUCKET")
	s3Endpoint = os.Getenv("S3_ENDPOINT")
)

func DownloadPic(fileName string) (*bytes.Buffer, error) {
	curSession, err := session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
		Region:      aws.String(s3Region),
		Endpoint:    aws.String(s3Endpoint),
	})

	var buffer []byte
	buf := bytes.NewBuffer(buffer)
	s3Client := s3.New(curSession)

	obj, err := s3Client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(s3Bucket),
		Key:    aws.String(fileName),
	})

	_, err = io.Copy(buf, obj.Body)
	if err != nil {
		log.Printf("cant copy file")
	}

	if err != nil {
		log.Warnf("Unable to download item %q, %v", "hep", err)
	}

	return buf, nil
}
