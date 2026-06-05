package lib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type FaceRegisterRequest struct {
	EmployeeID   string `json:"employee_id"`
	FaceImageURL string `json:"face_image_url"`
}

type FaceRecognizeRequest struct {
	EmployeeID   string `json:"employee_id"`
	FaceImageURL string `json:"face_image_url"`
}

type FaceRecognizeResponse struct {
	Match      bool    `json:"match"`
	FakeFace   bool    `json:"fake_face"`
	Confidence float64 `json:"confidence"`
}

type FaceVerifyRequest struct {
	EmployeeID   string `json:"employee_id"`
	FaceImageURL string `json:"face_image_url"`
}

type FaceVerifyResponse struct {
	IsVerified bool    `json:"is_verified"`
	Confidence float64 `json:"confidence"`
}

type FaceCheckExistResponse struct {
	Exists bool `json:"exists"`
}

func VerifyFace(FaceVerifyURL, employeeID, faceImageURL string) (*FaceVerifyResponse, error) {
	payload, _ := json.Marshal(FaceVerifyRequest{
		EmployeeID:   employeeID,
		FaceImageURL: faceImageURL,
	})

	resp, err := http.Post(FaceVerifyURL, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result FaceVerifyResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func RegisterFace(FaceRegisterURL, employeeID, faceImageURL string) error {
	payload, _ := json.Marshal(FaceRegisterRequest{
		EmployeeID:   employeeID,
		FaceImageURL: faceImageURL,
	})

	resp, err := http.Post(FaceRegisterURL, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("face register failed with status: %d", resp.StatusCode)
	}

	return nil
}

func CheckFaceExistence(FaceCheckExistURL, employeeID string) *FaceCheckExistResponse {
	resp, err := http.Get(FaceCheckExistURL + "/" + employeeID)
	if err != nil {
		return &FaceCheckExistResponse{Exists: false}
	}
	defer resp.Body.Close()

	var result FaceCheckExistResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return &FaceCheckExistResponse{Exists: false}
	}

	return &result
}

func DeleteFace(FaceDeleteURL, employeeID string) error {
	req, err := http.NewRequest(http.MethodDelete, FaceDeleteURL+"/"+employeeID, nil)
	if err != nil {
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("face delete failed with status: %d", resp.StatusCode)
	}

	return nil
}

func RecognizeFace(FaceRecognizeURL, employeeID, faceImageURL string) (*FaceRecognizeResponse, error) {
	payload, _ := json.Marshal(FaceRecognizeRequest{
		EmployeeID:   employeeID,
		FaceImageURL: faceImageURL,
	})

	resp, err := http.Post(FaceRecognizeURL, "application/json", bytes.NewBuffer(payload))
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
