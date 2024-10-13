package crypto

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"fmt"
)

// Signer defines a contract for different types of signing implementations.
type Signer interface {
	Sign(dataToBeSigned []byte) ([]byte, error)
}

// TODO: implement RSA and ECDSA signing ...

// RSASigner implements the Signer interface for RSA signing.
type RSASigner struct {
	keyPair RSAKeyPair
}

// NewRSASigner creates a new RSASigner using an RSAKeyPair.
func NewRSASigner(keyPair RSAKeyPair) *RSASigner {
	return &RSASigner{keyPair: keyPair}
}

// Sign signs the given data using the RSA private key.
func (s *RSASigner) Sign(dataToBeSigned []byte) ([]byte, error) {
	// Hash the data using SHA-256
	hashed := sha256.Sum256(dataToBeSigned)

	// Sign the data using the RSA private key and specify SHA-256 as the hash function
	signature, err := rsa.SignPKCS1v15(rand.Reader, s.keyPair.Private, crypto.SHA256, hashed[:])
	if err != nil {
		return nil, fmt.Errorf("RSA signing failed: %w", err)
	}

	return signature, nil
}

// ECDSASigner implements the Signer interface for ECDSA signing.
type ECDSASigner struct {
	keyPair ECCKeyPair
}

// NewECDSASigner creates a new ECDSASigner using an ECCKeyPair.
func NewECDSASigner(keyPair ECCKeyPair) *ECDSASigner {
	return &ECDSASigner{keyPair: keyPair}
}

// Sign signs the given data using the ECDSA private key.
func (s *ECDSASigner) Sign(dataToBeSigned []byte) ([]byte, error) {
	// Hash the data using SHA-256
	hashed := sha256.Sum256(dataToBeSigned)

	// Sign the data using the ECDSA private key
	r, sigS, err := ecdsa.Sign(rand.Reader, s.keyPair.Private, hashed[:])
	if err != nil {
		return nil, fmt.Errorf("ECDSA signing failed: %w", err)
	}

	// Convert r and s values to byte slices
	rBytes := r.Bytes()
	sBytes := sigS.Bytes()

	// Combine r and s into a single byte slice
	signature := append(rBytes, sBytes...)

	return signature, nil
}
