package lib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/url"
	"time"
)

const faceServiceTimeout = 30 * time.Second

var faceHTTPClient = &http.Client{Timeout: faceServiceTimeout}

type FaceRecognizeResponse struct {
	Match   bool   `json:"match"`
	Message string `json:"message"`
}

type FaceRegisterResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
}

type FaceCheckExistResponse struct {
	EmployeeID string `json:"employee_id"`
	Registered bool   `json:"registered"`
}

func postImageMultipart(serviceURL, employeeID, filename string, image []byte) (*http.Response, error) {
	if filename == "" {
		filename = "image.jpg"
	}

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	part, err := writer.CreateFormFile("image", filename)
	if err != nil {
		return nil, err
	}
	if _, err := part.Write(image); err != nil {
		return nil, err
	}
	if err := writer.WriteField("employee_id", employeeID); err != nil {
		return nil, err
	}
	if err := writer.Close(); err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, serviceURL, &body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	return faceHTTPClient.Do(req)
}

func RegisterFace(faceRegisterURL, employeeID, filename string, image []byte) error {
	resp, err := postImageMultipart(faceRegisterURL, employeeID, filename, image)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var result FaceRegisterResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("face register failed with status: %d", resp.StatusCode)
	}

	if result.Error || (resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated) {
		return fmt.Errorf("face register failed: %s", result.Message)
	}

	return nil
}

func RecognizeFace(faceRecognizeURL, employeeID, filename string, image []byte) (*FaceRecognizeResponse, error) {
	resp, err := postImageMultipart(faceRecognizeURL, employeeID, filename, image)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result FaceRecognizeResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func CheckFaceExistence(faceCheckExistURL, employeeID string) *FaceCheckExistResponse {
	reqURL := fmt.Sprintf("%s?employee_id=%s", faceCheckExistURL, url.QueryEscape(employeeID))

	resp, err := faceHTTPClient.Get(reqURL)
	if err != nil {
		return &FaceCheckExistResponse{EmployeeID: employeeID, Registered: false}
	}
	defer resp.Body.Close()

	var result FaceCheckExistResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return &FaceCheckExistResponse{EmployeeID: employeeID, Registered: false}
	}

	return &result
}

func DeleteFace(faceDeleteURL, employeeID string) error {
	reqURL := fmt.Sprintf("%s?employee_id=%s", faceDeleteURL, url.QueryEscape(employeeID))

	req, err := http.NewRequest(http.MethodDelete, reqURL, nil)
	if err != nil {
		return err
	}

	resp, err := faceHTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("face delete failed with status: %d", resp.StatusCode)
	}

	return nil
}
