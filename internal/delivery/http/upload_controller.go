package http

import (
	"context"
	"hr-sas/internal/model"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/sirupsen/logrus"
)

type UploadController struct {
	Log *logrus.Logger
}

func NewUploadController(log *logrus.Logger) *UploadController {
	return &UploadController{Log: log}
}

// TODO: Replace local filesystem storage with object storage (S3/IDCloud) when ready.
// TODO: IDCloud flow (S3-compatible):
// 1) init client with access_key/secret_key/endpoint
//    - IDCLOUD_ACCESS_KEY=ganti dengan access key
//    - IDCLOUD_SECRET_KEY=ganti dengan secret key
//    - IDCLOUD_ENDPOINT=https://<endpoint>
//    - IDCLOUD_BUCKET=<bucket-name>
//    - IDCLOUD_PUBLIC_BASE_URL=https://<public-base>
// 2) upload file bytes to bucket with key = pathPrefix + filename
// 3) return file_url = public_base_url + "/" + key
func (c *UploadController) Upload(ctx *fiber.Ctx) error {
	file, err := ctx.FormFile("file")
	if err != nil {
		c.Log.WithError(err).Error("failed to read uploaded file")
		return fiber.NewError(fiber.StatusBadRequest, "file is required")
	}

	pathPrefix := strings.TrimSpace(ctx.FormValue("path", ""))
	cleanPrefix := filepath.Clean(pathPrefix)
	if strings.Contains(cleanPrefix, "..") {
		return fiber.NewError(fiber.StatusBadRequest, "invalid path")
	}
	if cleanPrefix == "." {
		cleanPrefix = ""
	}

	baseDir := "uploads"
	fileName := filepath.Base(file.Filename)
	timestamp := time.Now().UnixMilli()
	storedName := strings.Join([]string{strconv.FormatInt(timestamp, 10), fileName}, "_")
	relativePath := filepath.Join(cleanPrefix, storedName)
	fullDir := filepath.Join(baseDir, cleanPrefix)
	fullPath := filepath.Join(baseDir, relativePath)

	// Try IDCloud (S3-compatible) upload when config is present.
	if fileURL, ok, err := c.uploadToIDCloud(file, cleanPrefix, storedName); err != nil {
		c.Log.WithError(err).Error("failed to upload to IDCloud")
		return fiber.ErrInternalServerError
	} else if ok {
		return ctx.JSON(model.WebResponse[map[string]string]{
			Data: map[string]string{
				"file_url": fileURL,
			},
		})
	}

	if err := os.MkdirAll(fullDir, 0o755); err != nil {
		c.Log.WithError(err).Error("failed to create upload directory")
		return fiber.ErrInternalServerError
	}

	if err := ctx.SaveFile(file, fullPath); err != nil {
		c.Log.WithError(err).Error("failed to save uploaded file")
		return fiber.ErrInternalServerError
	}

	return ctx.JSON(model.WebResponse[map[string]string]{
		Data: map[string]string{
			"file_url": ctx.Protocol() + "://" + ctx.Hostname() + "/" + filepath.ToSlash(fullPath),
		},
	})
}

func (c *UploadController) uploadToIDCloud(file *fiber.FileHeader, cleanPrefix, storedName string) (string, bool, error) {
	accessKey := strings.TrimSpace(os.Getenv("IDCLOUD_ACCESS_KEY"))
	secretKey := strings.TrimSpace(os.Getenv("IDCLOUD_SECRET_KEY"))
	endpoint := strings.TrimSpace(os.Getenv("IDCLOUD_ENDPOINT"))
	bucket := strings.TrimSpace(os.Getenv("IDCLOUD_BUCKET"))
	publicBaseURL := strings.TrimSpace(os.Getenv("IDCLOUD_PUBLIC_BASE_URL"))

	if accessKey == "" || secretKey == "" || endpoint == "" || bucket == "" || publicBaseURL == "" {
		return "", false, nil
	}

	secure := true
	if strings.HasPrefix(endpoint, "http://") || strings.HasPrefix(endpoint, "https://") {
		parsed, err := url.Parse(endpoint)
		if err != nil {
			return "", false, err
		}
		secure = parsed.Scheme == "https"
		endpoint = parsed.Host
	}

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: secure,
	})
	if err != nil {
		return "", false, err
	}

	key := storedName
	if cleanPrefix != "" {
		key = path.Join(cleanPrefix, storedName)
	}

	src, err := file.Open()
	if err != nil {
		return "", false, err
	}
	defer src.Close()

	contentType := file.Header.Get("Content-Type")
	_, err = client.PutObject(context.Background(), bucket, key, src, file.Size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", false, err
	}

	publicBaseURL = strings.TrimRight(publicBaseURL, "/")
	return publicBaseURL + "/" + key, true, nil
}
