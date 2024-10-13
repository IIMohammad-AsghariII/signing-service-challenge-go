package rsa

import (
	"bytes"
	"crypto"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/api"
	_ "github.com/fiskaly/coding-challenges/signing-service-challenge/api"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/payload/response"
	_ "github.com/fiskaly/coding-challenges/signing-service-challenge/persistence"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/utils"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Setup a new server for testing
func setupRsa() *api.Server {
	return api.NewServer(":8080")
}

// TestVerifySignatureAfterSigning tests the signature verification after signing a transaction
func TestVerifySignatureAfterSigningRsa(t *testing.T) {
	// Initialize the server and test recorder
	server := setupRsa()

	// First, create a device
	createReqBody := `{
		"id": "123e4567-e89b-12d3-a456-426614174000",
		"algorithm": "RSA",
		"label": "test-device"
	}`
	createReq := httptest.NewRequest("POST", "/api/v0/create-signature-device", bytes.NewBufferString(createReqBody))
	createReq.Header.Set("Content-Type", "application/json")
	createRecorder := httptest.NewRecorder()

	createHandler := http.HandlerFunc(server.CreateSignatureDeviceHandler)
	createHandler.ServeHTTP(createRecorder, createReq)

	// Extract public key from response
	var createResponse response.DeviceResponse
	devErr := json.Unmarshal(createRecorder.Body.Bytes(), &createResponse)
	if devErr != nil {
		t.Errorf("unexpected error in response unmarshalling: %v", devErr)
	}
	publicKey := createResponse.PublicKey

	// Now, sign a transaction with that device
	signReqBody := `{
		"deviceId": "123e4567-e89b-12d3-a456-426614174000",
		"data": "sample-transaction-data"
	}`
	signReq := httptest.NewRequest("POST", "/api/v0/sign-transaction", bytes.NewBufferString(signReqBody))
	signReq.Header.Set("Content-Type", "application/json")
	signRecorder := httptest.NewRecorder()

	signHandler := http.HandlerFunc(server.SignTransactionHandler)
	signHandler.ServeHTTP(signRecorder, signReq)

	// Validate the response
	if status := signRecorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check response body for the signature
	var signResponse response.SignTransactionResponse
	err := json.Unmarshal(signRecorder.Body.Bytes(), &signResponse)
	if err != nil {
		t.Errorf("unexpected error in response unmarshalling: %v", err)
	}
	if signResponse.SignedData == "" {
		t.Errorf("expected signed data but got empty")
	}

	// Now verify the signature
	signature := signResponse.Signature
	signedData := signResponse.SignedData

	valid, err := verifySignatureRsa(publicKey, signature, signedData)
	if err != nil {
		t.Errorf("Error verifying signature: %v", err)
	}
	if !valid {
		t.Error("Expected signature to be valid, but it was not")
	}
}

// verifySignature function to verify the validity of a signature
func verifySignatureRsa(publicKeyStr, signatureStr, signedDataStr string) (bool, error) {
	// Load the public key
	block, _ := pem.Decode([]byte(publicKeyStr))
	if block == nil || block.Type != "RSA_PUBLIC_KEY" {
		return false, errors.New("invalid public key")
	}

	// Parse the PKCS#1 public key
	rsaPublicKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return false, err
	}

	// Decode the signature
	signature, err := utils.Base64Decode(signatureStr)
	if err != nil {
		return false, err
	}

	// Parse signedData
	firstUnderscoreIndex := strings.Index(signedDataStr, "_")
	lastUnderscoreIndex := strings.LastIndex(signedDataStr, "_")

	if firstUnderscoreIndex == -1 || lastUnderscoreIndex == -1 || firstUnderscoreIndex == lastUnderscoreIndex {
		return false, errors.New("invalid signed data format")
	}

	// Extract parts
	signatureCounter := signedDataStr[:firstUnderscoreIndex]
	dataToBeSigned := signedDataStr[firstUnderscoreIndex+1 : lastUnderscoreIndex]
	lastSignature := signedDataStr[lastUnderscoreIndex+1:]

	// Display values for verification
	//fmt.Println("Signature:", signatureStr)
	//fmt.Println("Data To Be Signed:", dataToBeSigned)
	//fmt.Println("Last Signature:", lastSignature)
	//fmt.Println("Signature Counter:", signatureCounter)

	// Create secured data to be signed
	securedDataToBeSigned := fmt.Sprintf("%s_%s_%s", signatureCounter, dataToBeSigned, lastSignature)
	//fmt.Println("Secured Data To Be Signed:", securedDataToBeSigned)

	// Hash the data
	hashed := utils.HashData(securedDataToBeSigned)
	//fmt.Printf("Hashed Data: %x\n", hashed)

	// Verify the signature
	err = rsa.VerifyPKCS1v15(rsaPublicKey, crypto.SHA256, hashed, signature)
	if err != nil {
		return false, err
	}

	return true, nil
}
