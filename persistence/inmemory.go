package persistence

// TODO: in-memory persistence ...
import (
	"errors"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
	"sync"
)

// InMemoryDeviceRepository implements the DeviceRepositoryInterface
type InMemoryDeviceRepository struct {
	devices map[string]*domain.SignatureDevice
	mu      sync.RWMutex
}

// NewInMemoryDeviceRepository creates a new instance of InMemoryDeviceRepository
func NewInMemoryDeviceRepository() DeviceRepository {
	return &InMemoryDeviceRepository{
		devices: make(map[string]*domain.SignatureDevice),
	}
}

// AddDevice saves a new SignatureDevice to the repository
func (repo *InMemoryDeviceRepository) AddDevice(id, label string, algorithm domain.AlgorithmType, publicKey, privateKey, lastSignature string) (*domain.SignatureDevice, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	device := domain.NewSignatureDevice(id, label, algorithm, publicKey, privateKey, lastSignature)
	repo.devices[device.GetID()] = device
	return device, nil
}

// GetDevice retrieves a SignatureDevice by its ID
func (repo *InMemoryDeviceRepository) GetDevice(id string) (*domain.SignatureDevice, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	device, exists := repo.devices[id]
	if !exists {
		return nil, errors.New("device not found")
	}

	return device, nil
}

// ListDevices returns all SignatureDevices in the repository
func (repo *InMemoryDeviceRepository) ListDevices() ([]*domain.SignatureDevice, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	if len(repo.devices) == 0 {
		return nil, errors.New("no devices found") 
	}

	devices := make([]*domain.SignatureDevice, 0, len(repo.devices))
	for _, device := range repo.devices {
		devices = append(devices, device)
	}

	return devices, nil
}

// IncrementSignatureCount updates the signature count of a device
func (repo *InMemoryDeviceRepository) IncrementSignatureCount(id string) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	device, exists := repo.devices[id]
	if !exists {
		return errors.New("device not found")
	}

	// Increment the signature count directly
	currentCount := device.GetSignatureCount()
	device.SetSignatureCount(currentCount + 1)
	return nil
}

// UpdateLastSignature updates the last signature of a device
func (repo *InMemoryDeviceRepository) UpdateLastSignature(id string, lastSignature string) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	device, exists := repo.devices[id]
	if !exists {
		return errors.New("device not found")
	}

	// Use the setter to update the last signature
	device.SetLastSignature(lastSignature)
	return nil
}
