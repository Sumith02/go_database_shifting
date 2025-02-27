package api

import (
	"bytes"
	"db_shifter/model"
	"encoding/json"
	"fmt"
	"net/http"
)

const newAPIURL = "https://biometric.http.vsensetech.in/user/create/student"

type APIClient struct{}

func NewAPIClient() *APIClient {
	return &APIClient{}
}

func (api *APIClient) SendStudentData(student model.StudentData) error {
	jsonData, err := json.Marshal(student)
	if err != nil {
		return err
	}

	resp, err := http.Post(newAPIURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send data, status code: %d", resp.StatusCode)
	}

	return nil
}
