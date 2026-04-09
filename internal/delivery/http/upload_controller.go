package http

import (
	"hr-sas/internal/model"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type UploadController struct {
	Log *logrus.Logger
}

func NewUploadController(log *logrus.Logger) *UploadController {
	return &UploadController{Log: log}
}

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

	if err := os.MkdirAll(fullDir, 0o755); err != nil {
		c.Log.WithError(err).Error("failed to create upload directory")
		return fiber.ErrInternalServerError
	}

	if err := ctx.SaveFile(file, fullPath); err != nil {
		c.Log.WithError(err).Error("failed to save uploaded file")
		return fiber.ErrInternalServerError
	}

	fileURL := ctx.Protocol() + "://" + ctx.Hostname() + "/" + filepath.ToSlash(fullPath)

	return ctx.JSON(model.WebResponse[map[string]string]{
		Data: map[string]string{
			"file_url": fileURL,
		},
	})
}
