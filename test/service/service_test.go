package api

import (
	"github.com/fiskaly/coding-challenges/signing-service-challenge/api"
	_ "github.com/fiskaly/coding-challenges/signing-service-challenge/api"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/payload/request"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/persistence"
	_ "github.com/fiskaly/coding-challenges/signing-service-challenge/persistence"
	"testing"
)

// Setup a new DeviceService for testing
func setupService() api.DeviceServiceInterface {
	store := persistence.NewInMemoryDeviceRepository()
	return api.NewDeviceService(store)
}

// TestCreateSignatureDevice tests the CreateSignatureDevice function
func TestCreateSignatureDevice(t *testing.T) {
	service := setupService()

	// Define a valid request
	id := "123e4567-e89b-12d3-a456-426614174000"
	label := "test-device"
	algorithm := domain.RSA

	req := request.DeviceRequest{
		ID:        id,
		Label:     label,
		Algorithm: string(algorithm),
	}

	// Call the service method
	response, err := service.CreateSignatureDevice(&req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Validate the response
	if response.ID != id {
		t.Errorf("expected device ID %v, but got %v", id, response.ID)
	}
	if response.Label != label {
		t.Errorf("expected label %v, but got %v", label, response.Label)
	}
	if response.SignatureCount != 0 {
		t.Errorf("expected signature count 0, but got %v", response.SignatureCount)
	}
}

// TestCreateSignatureDeviceDuplicateID tests the duplicate ID error in CreateSignatureDevice function
func TestCreateSignatureDeviceDuplicateID(t *testing.T) {
	service := setupService()

	// Create the first device
	id := "123e4567-e89b-12d3-a456-426614174000"
	label := "test-device"
	algorithm := domain.RSA

	req := request.DeviceRequest{
		ID:        id,
		Label:     label,
		Algorithm: string(algorithm),
	}

	_, err := service.CreateSignatureDevice(&req)
	if err != nil {
		t.Fatalf("unexpected error during first creation: %v", err)
	}

	// Try to create the device with the same ID again
	_, err = service.CreateSignatureDevice(&req)
	if err == nil {
		t.Errorf("expected error due to duplicate ID, but got none")
	}
}

// TestSignTransaction tests the SignTransaction function
func TestSignTransaction(t *testing.T) {
	service := setupService()

	// First, create a device
	id := "123e4567-e89b-12d3-a456-426614174000"
	label := "test-device"
	algorithm := domain.RSA

	reqCSD := request.DeviceRequest{
		ID:        id,
		Label:     label,
		Algorithm: string(algorithm),
	}

	_, err := service.CreateSignatureDevice(&reqCSD)
	if err != nil {
		t.Fatalf("unexpected error during device creation: %v", err)
	}

	// Sign a transaction
	data := "sample-transaction-data"

	reqST := request.SignTransactionRequest{
		DeviceID: id,
		Data:     data,
	}
	signResponse, err := service.SignTransaction(&reqST)
	if err != nil {
		t.Fatalf("unexpected error during signing: %v", err)
	}

	// Validate the response
	if signResponse.Signature == "" {
		t.Errorf("expected a signature, but got empty")
	}
	if signResponse.SignedData == "" {
		t.Errorf("expected signed data, but got empty")
	}
}

// TestListSignatureDevices tests the ListSignatureDevices function
func TestListSignatureDevices(t *testing.T) {
	service := setupService()

	// Initially, the list should be empty
	devices, err := service.ListSignatureDevices()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(devices) != 0 {
		t.Errorf("expected no devices, but got %d", len(devices))
	}

	// Add a device
	id := "123e4567-e89b-12d3-a456-426614174000"
	label := "test-device"
	algorithm := domain.RSA

	req := request.DeviceRequest{
		ID:        id,
		Label:     label,
		Algorithm: string(algorithm),
	}

	_, err = service.CreateSignatureDevice(&req)
	if err != nil {
		t.Fatalf("unexpected error during device creation: %v", err)
	}

	// List devices again
	devices, err = service.ListSignatureDevices()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(devices) != 1 {
		t.Errorf("expected 1 device, but got %d", len(devices))
	}

	// Validate the details of the listed device
	device := devices[0]
	if device.ID != id {
		t.Errorf("expected device ID %v, but got %v", id, device.ID)
	}
	if device.Label != label {
		t.Errorf("expected label %v, but got %v", label, device.Label)
	}
	if device.SignatureCount != 0 {
		t.Errorf("expected signature count 0, but got %v", device.SignatureCount)
	}
}

// TestGetSignatureDeviceById tests the GetSignatureDeviceById function
func TestGetSignatureDeviceById(t *testing.T) {
	service := setupService()

	// First, create a device
	id := "123e4567-e89b-12d3-a456-426614174000"
	label := "test-device"
	algorithm := domain.RSA

	req := request.DeviceRequest{
		ID:        id,
		Label:     label,
		Algorithm: string(algorithm),
	}

	_, err := service.CreateSignatureDevice(&req)
	if err != nil {
		t.Fatalf("unexpected error during device creation: %v", err)
	}

	// Get the device by ID
	deviceResponse, err := service.GetSignatureDeviceById(id)
	if err != nil {
		t.Fatalf("unexpected error during device retrieval: %v", err)
	}

	// Validate the response
	if deviceResponse.ID != id {
		t.Errorf("expected device ID %v, but got %v", id, deviceResponse.ID)
	}
	if deviceResponse.Label != label {
		t.Errorf("expected label %v, but got %v", label, deviceResponse.Label)
	}
	if deviceResponse.SignatureCount != 0 {
		t.Errorf("expected signature count 0, but got %v", deviceResponse.SignatureCount)
	}
}

// TestGetSignatureDeviceByIdNotFound tests the GetSignatureDeviceById function when the device is not found
func TestGetSignatureDeviceByIdNotFound(t *testing.T) {
	service := setupService()

	// Try to get a device that doesn't exist
	_, err := service.GetSignatureDeviceById("non-existent-id")
	if err == nil {
		t.Errorf("expected error for non-existent device ID, but got none")
	}
}
