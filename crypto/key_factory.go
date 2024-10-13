package crypto

import (
	"errors"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
)

// KeyPairGenerator defines an interface for key generation and private key unmarshaling.
type KeyPairGenerator interface {
	GenerateKeyPair() ([]byte, []byte, error)   // Generates the public and private keys.
	UnmarshalPrivateKey([]byte) (Signer, error) // Converts the private key to a Signer.
}

// KeyPairFactory is a factory for generating key pairs and signers based on the algorithm.
type KeyPairFactory struct{}

// NewKeyPairFactory creates a new KeyPairFactory.
func NewKeyPairFactory() *KeyPairFactory {
	return &KeyPairFactory{}
}

// GetKeyPair returns a key pair generator for the given algorithm.
func (f *KeyPairFactory) GetKeyPair(algorithm domain.AlgorithmType) (KeyPairGenerator, error) {
	switch algorithm {
	case domain.RSA:
		return &RSAKeyPairGenerator{}, nil
	case domain.ECC:
		return &ECCKeyPairGenerator{}, nil
	default:
		return nil, errors.New("unsupported algorithm")
	}
}

// RSAKeyPairGenerator handles RSA key generation and signing.
type RSAKeyPairGenerator struct{}

// GenerateKeyPair generates an RSA key pair and returns the public and private keys.
func (g *RSAKeyPairGenerator) GenerateKeyPair() ([]byte, []byte, error) {
	rsaGen := RSAGenerator{}
	keyPair, err := rsaGen.Generate()
	if err != nil {
		return nil, nil, err
	}

	marshaler := RSAMarshaler{}
	publicKey, privateKey, err := marshaler.Marshal(*keyPair)
	if err != nil {
		return nil, nil, err
	}
	return publicKey, privateKey, nil
}

// UnmarshalPrivateKey converts the private key bytes to a Signer for RSA.
func (g *RSAKeyPairGenerator) UnmarshalPrivateKey(privateKeyBytes []byte) (Signer, error) {
	rsaKeyPair, err := (&RSAMarshaler{}).Unmarshal(privateKeyBytes)
	if err != nil {
		return nil, err
	}
	return NewRSASigner(*rsaKeyPair), nil
}

// ECCKeyPairGenerator handles ECC key generation and signing.
type ECCKeyPairGenerator struct{}

// GenerateKeyPair generates an ECC key pair and returns the public and private keys.
func (g *ECCKeyPairGenerator) GenerateKeyPair() ([]byte, []byte, error) {
	eccGen := ECCGenerator{}
	keyPair, err := eccGen.Generate()
	if err != nil {
		return nil, nil, err
	}

	marshaler := ECCMarshaler{}
	publicKey, privateKey, err := marshaler.Encode(*keyPair)
	if err != nil {
		return nil, nil, err
	}
	return publicKey, privateKey, nil
}

// UnmarshalPrivateKey converts the private key bytes to a Signer for ECC.
func (g *ECCKeyPairGenerator) UnmarshalPrivateKey(privateKeyBytes []byte) (Signer, error) {
	ecKeyPair, err := (&ECCMarshaler{}).Decode(privateKeyBytes)
	if err != nil {
		return nil, err
	}
	return NewECDSASigner(*ecKeyPair), nil
}
