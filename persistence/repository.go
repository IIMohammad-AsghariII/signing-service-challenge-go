package persistence

import "github.com/fiskaly/coding-challenges/signing-service-challenge/domain"

// DeviceRepository defines the interface for storage backends, allowing flexibility for future implementations.
type DeviceRepository interface {
	AddDevice(id, label string, algorithm domain.AlgorithmType, publicKey, privateKey, lastSignature string) (*domain.SignatureDevice, error)
	GetDevice(id string) (*domain.SignatureDevice, error)
	ListDevices() ([]*domain.SignatureDevice, error)
	IncrementSignatureCount(id string) error
	UpdateLastSignature(id string, lastSignature string) error
}
