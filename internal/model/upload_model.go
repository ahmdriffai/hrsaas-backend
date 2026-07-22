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

type PresignResponse struct {
	UploadURL string `json:"upload_url"`
	ObjectKey string `json:"object_key"`
}

type PresignRequest struct {
	MimeType string `json:"mime_type" validate:"required"`
}
