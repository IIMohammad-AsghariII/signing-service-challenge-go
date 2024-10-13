package utils

import (
	"crypto/sha256"
	"encoding/base64"
)

// Base64Encode is a utility function to base64 encode strings
func Base64Encode(data string) string {
	return base64.StdEncoding.EncodeToString([]byte(data))
}

// Base64Decode is a utility function to base64 encode strings
func Base64Decode(data string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(data)
}

// HashData function to compute the hash of the data
func HashData(data string) []byte {
	hash := sha256.New()
	hash.Write([]byte(data))
	return hash.Sum(nil)
}
