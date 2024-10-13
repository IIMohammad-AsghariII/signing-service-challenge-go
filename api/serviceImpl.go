package api

import (
	"errors"
	"fmt"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/crypto"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/payload/request"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/payload/response"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/persistence"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/utils"
	"github.com/google/uuid"
)

// DeviceService implements the service
type DeviceService struct {
	store persistence.DeviceRepository
}

// NewDeviceService function to create a new service
func NewDeviceService(store persistence.DeviceRepository) DeviceServiceInterface {
	return &DeviceService{store: store}
}

// ValidateDeviceRequest validates the DeviceRequest
func (s *DeviceService) ValidateDeviceRequest(r *request.DeviceRequest) error {
	if r.ID == "" {
		return errors.New("ID is required")
	}
	if r.Algorithm == "" {
		return errors.New("algorithm is required")
	}
	if _, err := uuid.Parse(r.ID); err != nil {
		return errors.New("invalid UUID")
	}
	return nil
}

// CreateSignatureDevice creates a new signature device
func (s *DeviceService) CreateSignatureDevice(req *request.DeviceRequest) (*response.DeviceResponse, error) {
	var publicKey, privateKey []byte
	var err error

	// Validate the request
	if err = s.ValidateDeviceRequest(req); err != nil {
		return nil, err
	}
	// Check for duplicate ID
	existingDevice, err := s.store.GetDevice(req.ID)
	if err != nil {
		// Assume if the device does not exist, a specific error like ErrNotFound is returned
		if err.Error() != "device not found" { // If another error occurs, return it
			return nil, errors.New("failed to check existing device")
		}
	}

	// If the device exists, return a duplicate ID error
	if existingDevice != nil {
		return nil, errors.New("device with this ID already exists")
	}

	// Generate key pair based on the algorithm using the factory.
	factory := crypto.NewKeyPairFactory()
	keyGenerator, err := factory.GetKeyPair(domain.AlgorithmType(req.Algorithm))
	if err != nil {
		return nil, errors.New("invalid algorithm")
	}

	publicKey, privateKey, err = keyGenerator.GenerateKeyPair()
	if err != nil {
		return nil, errors.New("key generation failed")
	}

	// Create and store the device
	device, err := s.store.AddDevice(req.ID, req.Label, domain.AlgorithmType(req.Algorithm), string(publicKey), string(privateKey), "")
	if err != nil {
		return nil, errors.New("failed to add device")
	}

	// Return response
	return &response.DeviceResponse{
		ID:             device.GetID(),
		PublicKey:      device.GetPublicKey(),
		Label:          device.GetLabel(),
		SignatureCount: device.GetSignatureCount(),
	}, nil
}

// ValidateSignTransactionRequest validates the SignTransactionRequest
func (s *DeviceService) ValidateSignTransactionRequest(req *request.SignTransactionRequest) error {
	if req.DeviceID == "" {
		return errors.New("DeviceID is required")
	}
	if req.Data == "" {
		return errors.New("data is required")
	}
	if _, err := uuid.Parse(req.DeviceID); err != nil {
		return errors.New("invalid UUID for DeviceID")
	}
	return nil
}

// SignTransaction signs the transaction data with the specified device
func (s *DeviceService) SignTransaction(req *request.SignTransactionRequest) (*response.SignTransactionResponse, error) {

	// Validate the request
	if err := s.ValidateSignTransactionRequest(req); err != nil {
		return nil, err
	}

	// Retrieve the signature device
	device, err := s.store.GetDevice(req.DeviceID)
	if err != nil {
		return nil, errors.New("device not found")
	}

	// Variable to hold the signed data
	var signedData string

	// Prepare signed data using either the last signature or the device ID
	if device.GetSignatureCount() == 0 {
		// First transaction - use the device ID instead of the last signature
		signedData = fmt.Sprintf("%d_%s_%s", device.GetSignatureCount(), req.Data, utils.Base64Encode(device.GetID()))
	} else {
		// Use the last saved signature
		signedData = fmt.Sprintf("%d_%s_%s", device.GetSignatureCount(), req.Data, device.GetLastSignature())
	}

	// Choose the signing algorithm based on the device's private key using the factory.
	factory := crypto.NewKeyPairFactory()
	keyGenerator, err := factory.GetKeyPair(device.GetAlgorithm())
	if err != nil {
		return nil, errors.New("invalid algorithm")
	}

	signer, err := keyGenerator.UnmarshalPrivateKey([]byte(device.GetPrivateKey()))
	if err != nil {
		return nil, errors.New("failed to unmarshal private key")
	}

	signature, err := signer.Sign([]byte(signedData))
	if err != nil {
		return nil, errors.New("signing failed")
	}

	// Update the last signature with the new signature
	err = s.store.UpdateLastSignature(req.DeviceID, utils.Base64Encode(string(signature)))
	if err != nil {
		return nil, errors.New("failed to update last signature")
	}

	// Increment the signature count
	err = s.store.IncrementSignatureCount(req.DeviceID)
	if err != nil {
		return nil, errors.New("failed to increment signature count")
	}

	return &response.SignTransactionResponse{
		Signature:  utils.Base64Encode(string(signature)),
		SignedData: signedData,
	}, nil
}

// ListSignatureDevices method to list signature devices
func (s *DeviceService) ListSignatureDevices() ([]*response.DeviceResponse, error) {
	devices, _ := s.store.ListDevices()

	if len(devices) == 0 {
		return []*response.DeviceResponse{}, nil
	}

	var deviceResponses []*response.DeviceResponse
	for _, device := range devices {
		deviceResponses = append(deviceResponses, &response.DeviceResponse{
			ID:             device.GetID(),
			PublicKey:      device.GetPublicKey(),
			Label:          device.GetLabel(),
			SignatureCount: device.GetSignatureCount(),
		})
	}
	return deviceResponses, nil
}

// GetSignatureDeviceById method to retrieve the specified device's information by ID
func (s *DeviceService) GetSignatureDeviceById(deviceID string) (*response.DeviceResponse, error) {
	device, err := s.store.GetDevice(deviceID)
	if err != nil {
		return nil, errors.New("device not found")
	}

	return &response.DeviceResponse{
		ID:             device.GetID(),
		PublicKey:      device.GetPublicKey(),
		Label:          device.GetLabel(),
		SignatureCount: device.GetSignatureCount(),
	}, nil
}
