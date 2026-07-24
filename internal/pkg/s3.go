package pkg

import (
	"context"
	"fmt"
	"io"
	"mime"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/spf13/viper"
)

type S3Client struct {
	Client *s3.Client
	Config *viper.Viper
}

func NewS3Client(config *viper.Viper) *S3Client {
	cfg := aws.Config{
		Region: config.GetString("s3.region"),
		Credentials: credentials.NewStaticCredentialsProvider(
			config.GetString("s3.access_key"),
			config.GetString("s3.secret_key"),
			"",
		),
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(config.GetString("s3.endpoint"))
		o.UsePathStyle = true

		// Untuk S3 Compatible
		o.RequestChecksumCalculation = aws.RequestChecksumCalculationWhenRequired
		o.ResponseChecksumValidation = aws.ResponseChecksumValidationWhenRequired
	})

	return &S3Client{
		Client: client,
		Config: config,
	}
}

func (s *S3Client) GenerateDownloadURL(presignClient *s3.PresignClient, objectKey string, isPublic ...bool) (string, error) {
	bucket := s.Config.GetString("s3.bucket")
	if len(isPublic) > 0 && isPublic[0] {
		bucket = s.Config.GetString("s3.public_bucket")
	}

	req, err := presignClient.PresignGetObject(
		context.Background(),
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(objectKey),
		},
		s3.WithPresignExpires(30*time.Minute),
	)
	if err != nil {
		return "", err
	}

	return req.URL, nil
}

func (s *S3Client) GeneratPresignURL(mimeType string, isPublic ...bool) (string, string, error) {

	bucket := s.Config.GetString("s3.bucket")
	if len(isPublic) > 0 && isPublic[0] {
		bucket = s.Config.GetString("s3.public_bucket")
	}

	exts, _ := mime.ExtensionsByType(mimeType)
	ext := ""
	if len(exts) > 0 {
		ext = exts[0]
	}

	env := s.Config.GetString("app.env") // development, staging, production

	folder := time.Now().Format("2006/01")

	baseDir := "uploads"

	if env == "development" {
		baseDir = "development/uploads"
	}

	dir := fmt.Sprintf("%s/%s", baseDir, folder)

	filename := fmt.Sprintf(
		"%s/%d%s",
		dir,
		time.Now().UnixMilli(),
		ext,
	)

	presignClient := s3.NewPresignClient(s.Client)

	req, err := presignClient.PresignPutObject(
		context.Background(),
		&s3.PutObjectInput{
			Bucket:      aws.String(bucket),
			Key:         aws.String(filename),
			ContentType: aws.String(mimeType),
		},
		s3.WithPresignExpires(5*time.Minute),
	)

	if err != nil {
		return "", "", nil
	}

	return req.URL, filename, nil
}

func (s *S3Client) GetObjectBytes(objectKey string, isPublic ...bool) ([]byte, error) {
	bucket := s.Config.GetString("s3.bucket")
	if len(isPublic) > 0 && isPublic[0] {
		bucket = s.Config.GetString("s3.public_bucket")
	}

	obj, err := s.Client.GetObject(
		context.Background(),
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(objectKey),
		},
	)
	if err != nil {
		return nil, err
	}
	defer obj.Body.Close()

	return io.ReadAll(obj.Body)
}

func (s *S3Client) GetPublicURL(objectKey string) string {
	return strings.TrimRight(s.Config.GetString("s3.public_url"), "/") + "/" + strings.TrimLeft(objectKey, "/")
}
