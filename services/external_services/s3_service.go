package external_services

import (
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/aws/session"
	"os"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"fmt"
	"strings"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"io"
)

type S3Service struct {

	region string

	bucket string

	basePath string

	s3Svc *s3.S3

	sessionS3 *session.Session

	uploader *s3manager.Uploader
}

func NewS3Service(basePath string) S3Service {

	s3Service := S3Service{}

	var err error

	s3Service.region = os.Getenv("AWS_REGION")

	s3Service.bucket = os.Getenv("AWS_BUCKET")

	s3Service.basePath = os.Getenv("AWS_BASE_PATH")

	s3Service.sessionS3, err = session.NewSession(&aws.Config{
		Region:      aws.String(s3Service.region),
		Credentials: credentials.NewStaticCredentials(
			os.Getenv("AWS_ACCESS_KEY_ID"),
			os.Getenv("AWS_SECRET_ACCESS_KEY"),
			"",
		),
	})

	s3Service.s3Svc = s3.New(s3Service.sessionS3)

	s3Service.uploader = s3manager.NewUploader(s3Service.sessionS3)

	if err != nil {
		fmt.Errorf(err.Error())
		os.Exit(1)
	}


	return s3Service
}

func (s *S3Service) ListFiles(directory string) ([]*s3.Object, error){

	resp, err := s.s3Svc.ListObjects(&s3.ListObjectsInput{
		Bucket: aws.String(s.bucket),
		Prefix: aws.String(directory),
	})

	if err != nil {
		fmt.Printf("Unable to list items in bucket %q \n %v", s.bucket, err)
		return nil, err
	}

	return resp.Contents, nil

}

func (s *S3Service) UploadFile(filename string, file io.Reader) error{

	_, err := s.uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(s.bucket),
		Key: aws.String(filename),
		Body: file,
		CacheControl: aws.String("max-age=86400"),
		ACL: aws.String("public-read"),
	})

	return err

}

func (s *S3Service) DeleteFile(filePath string) error{

	_, err := s.s3Svc.GetObject(&s3.GetObjectInput{Bucket: aws.String(s.bucket), Key: aws.String(filePath)})

	if err != nil {
		return err
	}
	_, err = s.s3Svc.DeleteObject(&s3.DeleteObjectInput{Bucket: aws.String(s.bucket), Key: aws.String(filePath)})

	if err != nil {
		return err
	}

	err = s.s3Svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(filePath),
	})

	return err

}

func (s *S3Service) GetFileUrl(fileKey string) string{

	return fmt.Sprintf(
		"https://s3.amazonaws.com/%s/%s",
		s.bucket,
		strings.Replace(fileKey, " ", "+", -1),
	)

}