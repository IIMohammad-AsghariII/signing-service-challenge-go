package api

import (
	"bytes"
	"encoding/json"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/api"
	_ "github.com/fiskaly/coding-challenges/signing-service-challenge/api"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/payload/response"
	_ "github.com/fiskaly/coding-challenges/signing-service-challenge/persistence"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Setup a new server for testing
func setup() *api.Server {
	return api.NewServer(":8080")
}

// TestCreateSignatureDeviceHandler tests the CreateSignatureDeviceHandler function
func TestCreateSignatureDeviceHandler(t *testing.T) {
	// Initialize the server and test recorder
	server := setup()
	reqBody := `{
		"id": "123e4567-e89b-12d3-a456-426614174000",
		"algorithm": "RSA",
		"label": "test-device"
	}`

	req := httptest.NewRequest("POST", "/api/v0/create-signature-device", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	// Call the handler
	handler := http.HandlerFunc(server.CreateSignatureDeviceHandler)
	handler.ServeHTTP(recorder, req)

	// Validate the response
	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check response body
	var res response.DeviceResponse
	err := json.Unmarshal(recorder.Body.Bytes(), &res)
	if err != nil {
		t.Errorf("unexpected error in response unmarshalling: %v", err)
	}
	if res.ID != "123e4567-e89b-12d3-a456-426614174000" {
		t.Errorf("unexpected device ID: got %v want %v", res.ID, "123e4567-e89b-12d3-a456-426614174000")
	}
}

// TestSignTransactionHandler tests the SignTransactionHandler function
func TestSignTransactionHandler(t *testing.T) {
	// Initialize the server and test recorder
	server := setup()

	// First, create a device
	createReqBody := `{
		"id": "123e4567-e89b-12d3-a456-426614174000",
		"algorithm": "RSA",
		"label": "test-device"
	}`
	createReq := httptest.NewRequest("POST", "/api/v0/create-signature-device", bytes.NewBufferString(createReqBody))
	createReq.Header.Set("Content-Type", "application/json")
	createRecorder := httptest.NewRecorder()

	createHandler := http.HandlerFunc(server.CreateSignatureDeviceHandler)
	createHandler.ServeHTTP(createRecorder, createReq)

	// Now, sign a transaction with that device
	signReqBody := `{
		"deviceId": "123e4567-e89b-12d3-a456-426614174000",
		"data": "sample-transaction-data"
	}`
	signReq := httptest.NewRequest("POST", "/api/v0/sign-transaction", bytes.NewBufferString(signReqBody))
	signReq.Header.Set("Content-Type", "application/json")
	signRecorder := httptest.NewRecorder()

	signHandler := http.HandlerFunc(server.SignTransactionHandler)
	signHandler.ServeHTTP(signRecorder, signReq)

	// Validate the response
	if status := signRecorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check response body
	var resSign response.SignTransactionResponse
	err := json.Unmarshal(signRecorder.Body.Bytes(), &resSign)
	if err != nil {
		t.Errorf("unexpected error in response unmarshalling: %v", err)
	}
	if resSign.SignedData == "" {
		t.Errorf("expected signed data but got empty")
	}
}

// TestListSignatureDevicesHandler tests the ListSignatureDevicesHandler function
func TestListSignatureDevicesHandler(t *testing.T) {
	// Initialize the server and test recorder
	server := setup()

	// Create a device
	createReqBody := `{
		"id": "123e4567-e89b-12d3-a456-426614174000",
		"algorithm": "RSA",
		"label": "test-device"
	}`
	createReq := httptest.NewRequest("POST", "/api/v0/create-signature-device", bytes.NewBufferString(createReqBody))
	createReq.Header.Set("Content-Type", "application/json")
	createRecorder := httptest.NewRecorder()

	createHandler := http.HandlerFunc(server.CreateSignatureDeviceHandler)
	createHandler.ServeHTTP(createRecorder, createReq)

	// List devices
	listReq := httptest.NewRequest("GET", "/api/v0/devices", nil)
	listRecorder := httptest.NewRecorder()

	listHandler := http.HandlerFunc(server.ListSignatureDevicesHandler)
	listHandler.ServeHTTP(listRecorder, listReq)

	// Validate the response
	if status := listRecorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check response body
	var res []response.DeviceResponse
	err := json.Unmarshal(listRecorder.Body.Bytes(), &res)
	if err != nil {
		t.Errorf("unexpected error in response unmarshalling: %v", err)
	}
	if len(res) == 0 {
		t.Errorf("expected at least one device but got none")
	}
}

// TestGetSignatureDeviceByIdHandler tests the GetSignatureDeviceByIdHandler function
func TestGetSignatureDeviceByIdHandler(t *testing.T) {
	// Initialize the server and test recorder
	server := setup()

	// Create a device
	createReqBody := `{
		"id": "123e4567-e89b-12d3-a456-426614174000",
		"algorithm": "RSA",
		"label": "test-device"
	}`
	createReq := httptest.NewRequest("POST", "/api/v0/create-signature-device", bytes.NewBufferString(createReqBody))
	createReq.Header.Set("Content-Type", "application/json")
	createRecorder := httptest.NewRecorder()

	createHandler := http.HandlerFunc(server.CreateSignatureDeviceHandler)
	createHandler.ServeHTTP(createRecorder, createReq)

	// Get device by ID
	getReq := httptest.NewRequest("GET", "/api/v0/device?id=123e4567-e89b-12d3-a456-426614174000", nil)
	getRecorder := httptest.NewRecorder()

	getHandler := http.HandlerFunc(server.GetSignatureDeviceByIdHandler)
	getHandler.ServeHTTP(getRecorder, getReq)

	// Validate the response
	if status := getRecorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check response body
	var res response.DeviceResponse
	err := json.Unmarshal(getRecorder.Body.Bytes(), &res)
	if err != nil {
		t.Errorf("unexpected error in response unmarshalling: %v", err)
	}
	if res.ID != "123e4567-e89b-12d3-a456-426614174000" {
		t.Errorf("unexpected device ID: got %v want %v", res.ID, "123e4567-e89b-12d3-a456-426614174000")
	}
}
