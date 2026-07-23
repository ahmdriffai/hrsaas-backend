package usecase

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"hr-sas/internal/model"
	"hr-sas/internal/pkg"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	smithymiddleware "github.com/aws/smithy-go/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type UploadUseCase struct {
	Log       *logrus.Logger
	Validator *validator.Validate
	S3Client  *pkg.S3Client
	Viper     *viper.Viper
}

func NewUploadUseCase(log *logrus.Logger, validator *validator.Validate, s3Client *pkg.S3Client, viper *viper.Viper) *UploadUseCase {
	return &UploadUseCase{
		Log:       log,
		Validator: validator,
		Viper:     viper,
		S3Client:  s3Client,
	}
}

func (u *UploadUseCase) Upload(ctx context.Context, request *model.UploadRequest) (*model.UploadResponse, error) {
	if err := u.Validator.Struct(request); err != nil {
		u.Log.WithError(err).Error("Failed to validate request body")
		return nil, fiber.ErrBadRequest
	}

	if request.File.Size > 3*1024*1024 {
		u.Log.Error("file too large")
		return nil, fiber.ErrBadRequest
	}

	var (
		url *string
		err error
	)

	if u.Viper.GetString("s3.bucket") != "" {
		url, err = u.SaveToS3(ctx, request.File)
	} else {
		baseUrl := u.Viper.GetString("app.url")
		url, err = u.Save(request.File, baseUrl)
	}

	if err != nil {
		u.Log.WithError(err).Error("failed to save file")
		return nil, fiber.ErrBadRequest
	}

	return &model.UploadResponse{Url: *url}, nil
}

func (u *UploadUseCase) Uploads(ctx context.Context, request *model.UploadsRequest) (*model.UploadResponses, error) {
	if err := u.Validator.Struct(request); err != nil {
		u.Log.WithError(err).Error("Failed to validate request body")
		return nil, fiber.ErrBadRequest
	}

	useS3 := u.Viper.GetString("s3.bucket") != ""
	baseUrl := u.Viper.GetString("app.url")

	var results []string

	for _, file := range request.Files {
		if file.Size > 10*1024*1024 {
			u.Log.Warn("skip file: too large")
			continue
		}

		var (
			url *string
			err error
		)

		if useS3 {
			url, err = u.SaveToS3(ctx, file)
		} else {
			url, err = u.Save(file, baseUrl)
		}

		if err != nil {
			u.Log.WithError(err).Warn("skip file: failed to save")
			continue
		}

		results = append(results, *url)
	}

	if len(results) == 0 {
		return nil, fiber.ErrBadRequest
	}

	return &model.UploadResponses{Urls: results}, nil
}

// Save menyimpan file ke disk lokal dan mengembalikan URL-nya.
func (u *UploadUseCase) Save(file *multipart.FileHeader, baseURL string) (*string, error) {
	ext := filepath.Ext(file.Filename)
	folder := time.Now().Format("2006/01")
	dir := fmt.Sprintf("./uploads/%s", folder)

	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return nil, err
	}

	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	fullPath := fmt.Sprintf("%s/%s", dir, filename)

	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	dst, err := os.Create(fullPath)
	if err != nil {
		return nil, err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/uploads/%s/%s", baseURL, folder, filename)
	return &url, nil
}

// SaveToS3 mengupload file ke S3-compatible object storage dan mengembalikan public URL-nya.
// Mendukung AWS S3, MinIO, dan Cloudflare R2 (atur s3.endpoint di config).
func (u *UploadUseCase) SaveToS3(ctx context.Context, file *multipart.FileHeader) (*string, error) {
	accessKey := u.Viper.GetString("s3.access_key")
	secretKey := u.Viper.GetString("s3.secret_key")
	region := u.Viper.GetString("s3.region")
	bucket := u.Viper.GetString("s3.bucket")
	endpoint := u.Viper.GetString("s3.endpoint")
	publicURL := u.Viper.GetString("s3.public_url")

	opts := []func(*awsconfig.LoadOptions) error{
		awsconfig.WithRegion(region),
		awsconfig.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(accessKey, secretKey, ""),
		),
	}

	cfg, err := awsconfig.LoadDefaultConfig(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("s3 config: %w", err)
	}

	clientOpts := []func(*s3.Options){}
	if endpoint != "" {
		clientOpts = append(clientOpts, func(o *s3.Options) {
			o.BaseEndpoint = aws.String(endpoint)
			o.UsePathStyle = true // diperlukan untuk MinIO
			// Disable strict payload checksum — fixes XAmzContentSHA256Mismatch
			o.RequestChecksumCalculation = aws.RequestChecksumCalculationWhenRequired
			o.ResponseChecksumValidation = aws.ResponseChecksumValidationWhenRequired
		})
	}

	client := s3.NewFromConfig(cfg, clientOpts...)

	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}
	data, err := io.ReadAll(src)
	src.Close()
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}

	ext := filepath.Ext(file.Filename)
	key := fmt.Sprintf("uploads/%s/%d%s", time.Now().Format("2006/01"), time.Now().UnixNano(), ext)

	_, err = client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:            aws.String(bucket),
		Key:               aws.String(key),
		Body:              bytes.NewReader(data),
		ContentLength:     aws.Int64(int64(len(data))),
		ContentType:       aws.String(file.Header.Get("Content-Type")),
		ACL:               types.ObjectCannedACLPublicRead,
		ChecksumAlgorithm: "", // leave empty
	}, func(o *s3.Options) {
		// Ceph (id cloudhost) requires actual SHA256, not UNSIGNED-PAYLOAD.
		// The default dynamic middleware sends UNSIGNED-PAYLOAD for HTTPS, so we swap it.
		o.APIOptions = append(o.APIOptions, func(stack *smithymiddleware.Stack) error {
			_, err := stack.Finalize.Swap("ComputePayloadHash", &v4.ComputePayloadSHA256{})
			return err
		})
	})
	if err != nil {
		return nil, fmt.Errorf("s3 put object: %w", err)
	}

	var url string
	if publicURL != "" {
		url = fmt.Sprintf("%s/%s", publicURL, key)
	} else {
		url = fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", bucket, region, key)
	}

	return &url, nil
}

var allowedMimeTypes = map[string]bool{
	"image/jpeg":      true,
	"image/png":       true,
	"image/webp":      true,
	"application/pdf": true,
}

func (u *UploadUseCase) GenerateUploadURL(ctx context.Context, request *model.PresignRequest) (*model.PresignResponse, error) {
	if !allowedMimeTypes[request.MimeType] {
		return nil, errors.New("unsupported mime type")
	}

	uploadUrl, filename, err := u.S3Client.GeneratPresignURL(request.MimeType)
	if err != nil {
		return nil, errors.New("error generate url")
	}

	return &model.PresignResponse{
		UploadURL: uploadUrl,
		ObjectKey: filename,
	}, nil
}
