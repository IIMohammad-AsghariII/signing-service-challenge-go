package api

import (
	"github.com/fiskaly/coding-challenges/signing-service-challenge/payload/request"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/payload/response"
)

// DeviceServiceInterface defines the interface for device services
type DeviceServiceInterface interface {
	// CreateSignatureDevice creates a new signature device.
	CreateSignatureDevice(req *request.DeviceRequest) (*response.DeviceResponse, error)
	// SignTransaction signs a transaction with a specified signature device.
	SignTransaction(req *request.SignTransactionRequest) (*response.SignTransactionResponse, error)
	// ListSignatureDevices retrieves a list of all available signature devices.
	ListSignatureDevices() ([]*response.DeviceResponse, error)
	// GetSignatureDeviceById retrieves a specific signature device by its ID.
	GetSignatureDeviceById(deviceID string) (*response.DeviceResponse, error)
}
