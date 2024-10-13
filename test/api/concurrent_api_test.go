package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/payload/response"
	"github.com/google/uuid"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

// TestConcurrentCreateAndSignTransactionHandler tests concurrent requests for CreateSignatureDeviceHandler and SignTransactionHandler
func TestConcurrentCreateAndSignTransactionHandler(t *testing.T) {
	// Initialize the server
	server := setup()

	// Number of concurrent requests
	concurrency := 10000

	// WaitGroup to wait for all goroutines to finish
	var createWg sync.WaitGroup
	var signWg sync.WaitGroup

	// Add the number of create and sign requests
	createWg.Add(concurrency)
	signWg.Add(concurrency)

	// Create channel for any errors
	errCh := make(chan error, concurrency*2)

	// Store device IDs to ensure we use them for signing transactions
	deviceIDs := make([]string, concurrency)

	// First, create devices
	for i := 0; i < concurrency; i++ {
		go func(i int) {
			defer createWg.Done()

			deviceID := uuid.New().String() // Generate a new UUID for each device
			deviceIDs[i] = deviceID         // Store the device ID for later use in signing

			createReqBody := `{
				"id": "` + deviceID + `",
				"algorithm": "RSA",
				"label": "test-device-` + deviceID + `"
			}`
			req := httptest.NewRequest("POST", "/api/v0/create-signature-device", bytes.NewBufferString(createReqBody))
			req.Header.Set("Content-Type", "application/json")
			recorder := httptest.NewRecorder()

			handler := http.HandlerFunc(server.CreateSignatureDeviceHandler)
			handler.ServeHTTP(recorder, req)

			// Check if the request was successful
			if recorder.Code != http.StatusOK {
				errCh <- errors.New("CreateSignatureDeviceHandler returned wrong status code: got " + http.StatusText(recorder.Code))
				return
			}
		}(i)
	}

	// Wait for all device creation goroutines to complete
	createWg.Wait()

	// Now sign transactions using the created devices
	for i := 0; i < concurrency; i++ {
		go func(i int) {
			defer signWg.Done()

			// Use the stored device ID for signing the transaction
			signReqBody := `{
				"deviceId": "` + deviceIDs[i] + `",
				"data": "sample-transaction-data-` + deviceIDs[i] + `"
			}`
			req := httptest.NewRequest("POST", "/api/v0/sign-transaction", bytes.NewBufferString(signReqBody))
			req.Header.Set("Content-Type", "application/json")
			recorder := httptest.NewRecorder()

			handler := http.HandlerFunc(server.SignTransactionHandler)
			handler.ServeHTTP(recorder, req)

			// Check if the request was successful
			if recorder.Code != http.StatusOK {
				errCh <- errors.New("SignTransactionHandler returned wrong status code: got " + http.StatusText(recorder.Code))
				return
			}

			// Validate the response body
			var res response.SignTransactionResponse
			err := json.Unmarshal(recorder.Body.Bytes(), &res)
			if err != nil {
				errCh <- errors.New("unexpected error in sign response unmarshalling: " + err.Error())
				return
			}
			if res.SignedData == "" {
				errCh <- errors.New("expected signed data but got empty")
				return
			}
		}(i)
	}

	// Wait for all signing operations to complete
	signWg.Wait()

	// Close the error channel
	close(errCh)

	// Check if any errors occurred during concurrent requests
	for err := range errCh {
		if err != nil {
			t.Error(err)
		}
	}
}

// TestConcurrentCreateAndSignTransactionMixedHandler tests concurrent creation and signing transactions
func TestConcurrentCreateAndSignTransactionMixedHandler(t *testing.T) {
	// Initialize the server
	server := setup()

	// Number of initial devices
	initialDeviceCount := 1000

	// Number of new devices created during signing of initial ones
	newDeviceCount := 500

	// Number of sign requests per device
	signRequestsPerDevice := 100

	// WaitGroup for creating devices
	var createWg sync.WaitGroup
	// WaitGroup for signing transactions
	var signWg sync.WaitGroup

	// Create channel for any errors
	errCh := make(chan error, (initialDeviceCount+newDeviceCount)*signRequestsPerDevice)

	// Store device IDs to ensure we use them for signing transactions
	deviceIDs := make([]string, initialDeviceCount+newDeviceCount)

	// First, create initial 1000 devices
	createWg.Add(initialDeviceCount)
	for i := 0; i < initialDeviceCount; i++ {
		go func(i int) {
			defer createWg.Done()

			deviceID := uuid.New().String() // Generate a new UUID for each device
			deviceIDs[i] = deviceID         // Store the device ID for later use in signing

			createReqBody := `{
				"id": "` + deviceID + `",
				"algorithm": "RSA",
				"label": "test-device-` + deviceID + `"
			}`
			req := httptest.NewRequest("POST", "/api/v0/create-signature-device", bytes.NewBufferString(createReqBody))
			req.Header.Set("Content-Type", "application/json")
			recorder := httptest.NewRecorder()

			handler := http.HandlerFunc(server.CreateSignatureDeviceHandler)
			handler.ServeHTTP(recorder, req)

			// Check if the request was successful
			if recorder.Code != http.StatusOK {
				errCh <- errors.New("CreateSignatureDeviceHandler returned wrong status code: got " + http.StatusText(recorder.Code))
				return
			}
		}(i)
	}

	// Wait for all initial device creation goroutines to complete
	createWg.Wait()

	// Now concurrently sign transactions using the created devices
	for i := 0; i < initialDeviceCount; i++ {
		signWg.Add(signRequestsPerDevice)
		for j := 0; j < signRequestsPerDevice; j++ {
			go func(i int) {
				defer signWg.Done()

				signReqBody := `{
					"deviceId": "` + deviceIDs[i] + `",
					"data": "sample-transaction-data-` + deviceIDs[i] + `"
				}`
				req := httptest.NewRequest("POST", "/api/v0/sign-transaction", bytes.NewBufferString(signReqBody))
				req.Header.Set("Content-Type", "application/json")
				recorder := httptest.NewRecorder()

				handler := http.HandlerFunc(server.SignTransactionHandler)
				handler.ServeHTTP(recorder, req)

				// Check if the request was successful
				if recorder.Code != http.StatusOK {
					errCh <- errors.New("SignTransactionHandler returned wrong status code: got " + http.StatusText(recorder.Code))
					return
				}

				// Validate the response body
				var res response.SignTransactionResponse
				err := json.Unmarshal(recorder.Body.Bytes(), &res)
				if err != nil {
					errCh <- errors.New("unexpected error in sign response unmarshalling: " + err.Error())
					return
				}
				if res.SignedData == "" {
					errCh <- errors.New("expected signed data but got empty")
					return
				}
			}(i)
		}
	}

	// Create new devices while signing transactions
	createWg.Add(newDeviceCount)
	for i := initialDeviceCount; i < initialDeviceCount+newDeviceCount; i++ {
		go func(i int) {
			defer createWg.Done()

			deviceID := uuid.New().String()
			deviceIDs[i] = deviceID

			createReqBody := `{
				"id": "` + deviceID + `",
				"algorithm": "RSA",
				"label": "test-device-` + deviceID + `"
			}`
			req := httptest.NewRequest("POST", "/api/v0/create-signature-device", bytes.NewBufferString(createReqBody))
			req.Header.Set("Content-Type", "application/json")
			recorder := httptest.NewRecorder()

			handler := http.HandlerFunc(server.CreateSignatureDeviceHandler)
			handler.ServeHTTP(recorder, req)

			// Check if the request was successful
			if recorder.Code != http.StatusOK {
				errCh <- errors.New("CreateSignatureDeviceHandler returned wrong status code: got " + http.StatusText(recorder.Code))
				return
			}
		}(i)
	}

	// Wait for all new device creation to complete
	createWg.Wait()

	// Concurrently sign transactions for the new 500 devices
	for i := initialDeviceCount; i < initialDeviceCount+newDeviceCount; i++ {
		signWg.Add(signRequestsPerDevice)
		for j := 0; j < signRequestsPerDevice; j++ {
			go func(i int) {
				defer signWg.Done()

				signReqBody := `{
					"deviceId": "` + deviceIDs[i] + `",
					"data": "sample-transaction-data-` + deviceIDs[i] + `"
				}`
				req := httptest.NewRequest("POST", "/api/v0/sign-transaction", bytes.NewBufferString(signReqBody))
				req.Header.Set("Content-Type", "application/json")
				recorder := httptest.NewRecorder()

				handler := http.HandlerFunc(server.SignTransactionHandler)
				handler.ServeHTTP(recorder, req)

				// Check if the request was successful
				if recorder.Code != http.StatusOK {
					errCh <- errors.New("SignTransactionHandler returned wrong status code: got " + http.StatusText(recorder.Code))
					return
				}

				// Validate the response body
				var res response.SignTransactionResponse
				err := json.Unmarshal(recorder.Body.Bytes(), &res)
				if err != nil {
					errCh <- errors.New("unexpected error in sign response unmarshalling: " + err.Error())
					return
				}
				if res.SignedData == "" {
					errCh <- errors.New("expected signed data but got empty")
					return
				}
			}(i)
		}
	}

	// Wait for all signing operations to complete
	signWg.Wait()

	// Close the error channel
	close(errCh)

	// Check if any errors occurred during concurrent requests
	for err := range errCh {
		if err != nil {
			t.Error(err)
		}
	}
}
