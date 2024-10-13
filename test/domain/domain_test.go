package domain

import (
	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewSignatureDevice(t *testing.T) {
	// داده‌های نمونه
	id := "device-1"
	label := "Test Device"
	algorithm := domain.AlgorithmType("RSA")
	publicKey := "public-key"
	privateKey := "private-key"
	lastSignature := "signature-1"

	// ایجاد دستگاه جدید
	device := domain.NewSignatureDevice(id, label, algorithm, publicKey, privateKey, lastSignature)

	// بررسی مقادیر اولیه
	assert.Equal(t, id, device.GetID())
	assert.Equal(t, label, device.GetLabel())
	assert.Equal(t, algorithm, device.GetAlgorithm())
	assert.Equal(t, publicKey, device.GetPublicKey())
	assert.Equal(t, privateKey, device.GetPrivateKey())
	assert.Equal(t, uint64(0), device.GetSignatureCount()) // مقدار امضای اولیه 0 است
	assert.Equal(t, lastSignature, device.GetLastSignature())
}

func TestGetID(t *testing.T) {
	device := domain.NewSignatureDevice("device-1", "Test Device", domain.AlgorithmType("RSA"), "public-key", "private-key", "signature-1")
	assert.Equal(t, "device-1", device.GetID())
}

func TestGetSignatureCount(t *testing.T) {
	device := domain.NewSignatureDevice("device-1", "Test Device", domain.AlgorithmType("RSA"), "public-key", "private-key", "signature-1")
	assert.Equal(t, uint64(0), device.GetSignatureCount())
}

func TestGetLastSignature(t *testing.T) {
	device := domain.NewSignatureDevice("device-1", "Test Device", domain.AlgorithmType("RSA"), "public-key", "private-key", "signature-1")
	assert.Equal(t, "signature-1", device.GetLastSignature())
}

func TestGetLabel(t *testing.T) {
	device := domain.NewSignatureDevice("device-1", "Test Device", domain.AlgorithmType("RSA"), "public-key", "private-key", "signature-1")
	assert.Equal(t, "Test Device", device.GetLabel())
}

func TestGetAlgorithm(t *testing.T) {
	device := domain.NewSignatureDevice("device-1", "Test Device", domain.AlgorithmType("RSA"), "public-key", "private-key", "signature-1")
	assert.Equal(t, domain.AlgorithmType("RSA"), device.GetAlgorithm())
}

func TestGetPublicKey(t *testing.T) {
	device := domain.NewSignatureDevice("device-1", "Test Device", domain.AlgorithmType("RSA"), "public-key", "private-key", "signature-1")
	assert.Equal(t, "public-key", device.GetPublicKey())
}

func TestGetPrivateKey(t *testing.T) {
	device := domain.NewSignatureDevice("device-1", "Test Device", domain.AlgorithmType("RSA"), "public-key", "private-key", "signature-1")
	assert.Equal(t, "private-key", device.GetPrivateKey())
}

func TestSetLabel(t *testing.T) {
	device := domain.NewSignatureDevice("device-1", "Test Device", domain.AlgorithmType("RSA"), "public-key", "private-key", "signature-1")
	newLabel := "New Device Label"
	device.SetLabel(newLabel)
	assert.Equal(t, newLabel, device.GetLabel())
}

func TestSetPublicKey(t *testing.T) {
	device := domain.NewSignatureDevice("device-1", "Test Device", domain.AlgorithmType("RSA"), "public-key", "private-key", "signature-1")
	newPublicKey := "new-public-key"
	device.SetPublicKey(newPublicKey)
	assert.Equal(t, newPublicKey, device.GetPublicKey())
}

func TestSetPrivateKey(t *testing.T) {
	device := domain.NewSignatureDevice("device-1", "Test Device", domain.AlgorithmType("RSA"), "public-key", "private-key", "signature-1")
	newPrivateKey := "new-private-key"
	device.SetPrivateKey(newPrivateKey)
	assert.Equal(t, newPrivateKey, device.GetPrivateKey())
}

func TestSetSignatureCount(t *testing.T) {
	device := domain.NewSignatureDevice("device-1", "Test Device", domain.AlgorithmType("RSA"), "public-key", "private-key", "signature-1")
	newCount := uint64(5)
	device.SetSignatureCount(newCount)
	assert.Equal(t, newCount, device.GetSignatureCount())
}

func TestSetLastSignature(t *testing.T) {
	device := domain.NewSignatureDevice("device-1", "Test Device", domain.AlgorithmType("RSA"), "public-key", "private-key", "signature-1")
	newSignature := "signature-2"
	device.SetLastSignature(newSignature)
	assert.Equal(t, newSignature, device.GetLastSignature())
}
