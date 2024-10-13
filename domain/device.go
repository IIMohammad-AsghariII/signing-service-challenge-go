package domain

// TODO: signature device domain model ...

// SignatureDevice represents a signature device with public/private keys
type SignatureDevice struct {
	id             string 
	publicKey      string
	privateKey     string 
	label          string
	signatureCount uint64
	algorithm      AlgorithmType
	lastSignature  string
}

// NewSignatureDevice creates a new signature device with generated keys and an initial label
func NewSignatureDevice(id, label string, algorithm AlgorithmType, publicKey, privateKey, lastSignature string) *SignatureDevice {
	return &SignatureDevice{
		id:             id,
		label:          label,
		algorithm:      algorithm,
		publicKey:      publicKey,
		privateKey:     privateKey,
		signatureCount: 0,
		lastSignature:  lastSignature,
	}
}

// GetID returns the ID of the device
func (device *SignatureDevice) GetID() string {
	return device.id
}

// GetSignatureCount returns the current signature count
func (device *SignatureDevice) GetSignatureCount() uint64 {
	return device.signatureCount
}

// GetLastSignature returns the last signature value
func (device *SignatureDevice) GetLastSignature() string {
	return device.lastSignature
}

// GetLabel returns the label of the device
func (device *SignatureDevice) GetLabel() string {
	return device.label
}

func (device *SignatureDevice) GetAlgorithm() AlgorithmType {
	return device.algorithm
}

// GetPublicKey returns the public key of the device
func (device *SignatureDevice) GetPublicKey() string {
	return device.publicKey
}

// GetPrivateKey returns the private key of the device
func (device *SignatureDevice) GetPrivateKey() string {
	return device.privateKey
}

// SetLabel sets the label of the device
func (device *SignatureDevice) SetLabel(label string) {
	device.label = label
}

// SetPublicKey sets the public key of the device
func (device *SignatureDevice) SetPublicKey(publicKey string) {
	device.publicKey = publicKey
}

// SetPrivateKey sets the private key of the device
func (device *SignatureDevice) SetPrivateKey(privateKey string) {
	device.privateKey = privateKey
}

func (device *SignatureDevice) SetSignatureCount(count uint64) {
	device.signatureCount = count
}

// SetLastSignature sets the last signature of the device
func (device *SignatureDevice) SetLastSignature(lastSignature string) {
	device.lastSignature = lastSignature
}
