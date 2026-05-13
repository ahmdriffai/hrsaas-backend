package http

import (
	"hr-sas/internal/model"
	"hr-sas/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type UploadController struct {
	UseCase *usecase.UploadUseCase
	Log     *logrus.Logger
}

func NewUploadController(useCase *usecase.UploadUseCase, log *logrus.Logger) *UploadController {
	return &UploadController{
		UseCase: useCase,
		Log:     log,
	}
}

func (c *UploadController) Upload(ctx *fiber.Ctx) error {
	file, err := ctx.FormFile("file")
	if err != nil {
		c.Log.WithError(err).Error("failed to read uploaded file")
		return fiber.NewError(fiber.StatusBadRequest, "file is required")
	}

	request := &model.UploadRequest{File: file}

	response, err := c.UseCase.Upload(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("failed to upload file")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.UploadResponse]{Data: response})
}

func (c *UploadController) Uploads(ctx *fiber.Ctx) error {
	form, err := ctx.MultipartForm()
	if err != nil {
		c.Log.WithError(err).Error("failed to parse multipart form")
		return fiber.NewError(fiber.StatusBadRequest, "multipart form is required")
	}

	files := form.File["file"]
	if len(files) == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "at least one file is required")
	}

	request := &model.UploadsRequest{Files: files}

	response, err := c.UseCase.Uploads(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("failed to upload files")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.UploadResponses]{Data: response})
}
