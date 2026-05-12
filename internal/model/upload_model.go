package model

import "mime/multipart"

type UploadRequest struct {
	File *multipart.FileHeader `form:"file" validate:"required"`
}

type UploadsRequest struct {
	Files []*multipart.FileHeader `form:"file" validate:"required"`
}

type UploadResponse struct {
	Url string `json:"url"`
}

type UploadResponses struct {
	Urls []string `json:"urls"`
}
