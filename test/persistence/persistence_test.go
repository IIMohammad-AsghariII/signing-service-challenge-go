package persistence

import (
	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/persistence"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAddDevice(t *testing.T) {
	repo := persistence.NewInMemoryDeviceRepository()
	id := "device-1"
	label := "Test Device"
	algorithm := domain.AlgorithmType("RSA")
	publicKey := "public-key"
	privateKey := "private-key"
	lastSignature := "signature-1"

	device, err := repo.AddDevice(id, label, algorithm, publicKey, privateKey, lastSignature)
	assert.NoError(t, err)
	assert.NotNil(t, device)
	assert.Equal(t, id, device.GetID())
	assert.Equal(t, label, device.GetLabel())
	assert.Equal(t, algorithm, device.GetAlgorithm())
}

func TestGetDevice(t *testing.T) {
	repo := persistence.NewInMemoryDeviceRepository()
	id := "device-1"
	repo.AddDevice(id, "Test Device", domain.AlgorithmType("RSA"), "public-key", "private-key", "signature-1")

	device, err := repo.GetDevice(id)
	assert.NoError(t, err)
	assert.NotNil(t, device)
	assert.Equal(t, id, device.GetID())
}

func TestGetDevice_NonExistent(t *testing.T) {
	repo := persistence.NewInMemoryDeviceRepository()
	_, err := repo.GetDevice("non-existent")
	assert.Error(t, err)
	assert.EqualError(t, err, "device not found")
}

func TestListDevices(t *testing.T) {
	repo := persistence.NewInMemoryDeviceRepository()
	id := "device-1"
	repo.AddDevice(id, "Test Device", domain.AlgorithmType("RSA"), "public-key", "private-key", "signature-1")

	devices, err := repo.ListDevices()
	assert.NoError(t, err)
	assert.Len(t, devices, 1)
	assert.Equal(t, id, devices[0].GetID())
}

func TestListDevices_NoDevices(t *testing.T) {
	repo := persistence.NewInMemoryDeviceRepository()
	devices, err := repo.ListDevices()
	assert.Error(t, err)
	assert.EqualError(t, err, "no devices found")
	assert.Nil(t, devices)
}

func TestIncrementSignatureCount(t *testing.T) {
	repo := persistence.NewInMemoryDeviceRepository()
	id := "device-1"
	repo.AddDevice(id, "Test Device", domain.AlgorithmType("RSA"), "public-key", "private-key", "signature-1")

	err := repo.IncrementSignatureCount(id)
	assert.NoError(t, err)

	device, err := repo.GetDevice(id)
	assert.NoError(t, err)
	assert.Equal(t, uint64(1), device.GetSignatureCount())
}

func TestIncrementSignatureCount_NonExistent(t *testing.T) {
	repo := persistence.NewInMemoryDeviceRepository()
	err := repo.IncrementSignatureCount("non-existent")
	assert.Error(t, err)
	assert.EqualError(t, err, "device not found")
}

func TestUpdateLastSignature(t *testing.T) {
	repo := persistence.NewInMemoryDeviceRepository()
	id := "device-1"
	repo.AddDevice(id, "Test Device", domain.AlgorithmType("RSA"), "public-key", "private-key", "signature-1")

	newSignature := "signature-2"
	err := repo.UpdateLastSignature(id, newSignature)
	assert.NoError(t, err)

	device, err := repo.GetDevice(id)
	assert.NoError(t, err)
	assert.Equal(t, newSignature, device.GetLastSignature())
}

func TestUpdateLastSignature_NonExistent(t *testing.T) {
	repo := persistence.NewInMemoryDeviceRepository()
	err := repo.UpdateLastSignature("non-existent", "new-signature")
	assert.Error(t, err)
	assert.EqualError(t, err, "device not found")
}
