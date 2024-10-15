package api

import (
	"encoding/json"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/payload/request"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/persistence"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"regexp"
)

// The store variable for interacting with the data layer (DeviceRepositoryInterface)
var store persistence.DeviceRepository

// Create deviceService to use the device service
var deviceService DeviceServiceInterface

// Initialize the repository and service
func init() {
	// Regex to match the current working directory
	re := regexp.MustCompile(`^(.*` + "signing-service-challenge-go" + `)`)
	cwd, _ := os.Getwd()
	rootPath := re.Find([]byte(cwd))

	// Load environment variables from .env file
	if err := godotenv.Load(string(rootPath) + `/.env`); err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Get the DATA_STORE environment variable
	dataStore := os.Getenv("DATA_STORE")

	var err error
	// Initialize the store based on the DATA_STORE value
	switch dataStore {
	case "db":
		dataSourceName := "devices.db"
		store, err = persistence.NewSQLiteDeviceRepository(dataSourceName)
		if err != nil {
			log.Fatalf("failed to create SQLite repository: %v", err)
		}
	case "memory":
		store = persistence.NewInMemoryDeviceRepository()
	default:
		log.Fatalf("Invalid DATA_STORE value: %v. Use 'memory' or 'db'", dataStore)
	}

	// Initialize the device service with the store
	deviceService = NewDeviceService(store)
}

// CreateSignatureDeviceHandler API handler for creating a signature device
// @Summary Create a new signature device
// @Description Create a new signature device with a specified ID, label, and algorithm
// @Tags devices
// @Accept json
// @Produce json
// @Param device body DeviceRequest true "Device information"
// @Success 200 {object} DeviceResponse "Successful response"
// @Failure 400 {object} ErrorResponse "Invalid input"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /api/v0/create-signature-device [post]
func (s *Server) CreateSignatureDeviceHandler(w http.ResponseWriter, r *http.Request) {
	// Ensure the request method is POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req request.DeviceRequest
	// Decode the incoming request body into req
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Create the signature device using the device service
	deviceResponse, err := deviceService.CreateSignatureDevice(&req)
	if err != nil {
		if err.Error() == "device with this ID already exists" {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Set the response header and encode the response to JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	//	json.NewEncoder(w).Encode(deviceResponse)
	if err := json.NewEncoder(w).Encode(deviceResponse); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// SignTransactionHandler API handler for signing a transaction
// @Summary Sign a transaction
// @Description Sign the transaction data with the specified device
// @Tags transactions
// @Accept json
// @Produce json
// @Param transaction body SignTransactionRequest true "Transaction data"
// @Success 200 {object} SignTransactionResponse "Successful response"
// @Failure 400 {object} ErrorResponse "Invalid input"
// @Failure 404 {object} ErrorResponse "Device not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /api/v0/sign-transaction [post]
func (s *Server) SignTransactionHandler(w http.ResponseWriter, r *http.Request) {
	// Ensure the request method is POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req request.SignTransactionRequest
	// Decode the incoming request body into req
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Sign the transaction using the device service
	signResponse, err := deviceService.SignTransaction(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Set the response header and encode the response to JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	//json.NewEncoder(w).Encode(signResponse)
	if err := json.NewEncoder(w).Encode(signResponse); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// ListSignatureDevicesHandler API handler for listing signature devices
// @Summary List all signature devices
// @Description Retrieve a list of all signature devices
// @Tags devices
// @Accept json
// @Produce json
// @Success 200 {array} DeviceResponse "Successful response"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /api/v0/devices [get]
func (s *Server) ListSignatureDevicesHandler(w http.ResponseWriter, r *http.Request) {
	// Ensure the request method is GET
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Retrieve the list of devices from the device service
	devices, err := deviceService.ListSignatureDevices()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Set the response header and encode the response to JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	//json.NewEncoder(w).Encode(devices)
	if err := json.NewEncoder(w).Encode(devices); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// GetSignatureDeviceByIdHandler API handler for retrieving information about a specific device by ID
// @Summary Get a signature device by ID
// @Description Retrieve information about a specific signature device using its ID
// @Tags devices
// @Accept json
// @Produce json
// @Param id query string true "Device ID"
// @Success 200 {object} DeviceResponse "Successful response"
// @Failure 400 {object} ErrorResponse "Device ID is required"
// @Failure 404 {object} ErrorResponse "Device not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /api/v0/device [get]
func (s *Server) GetSignatureDeviceByIdHandler(w http.ResponseWriter, r *http.Request) {
	// Ensure the request method is GET
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Get the device ID from the query parameters
	deviceID := r.URL.Query().Get("id")
	//deviceID := r.URL.Path[len("/api/v0/devices/"):]
	if deviceID == "" {
		http.Error(w, "Device ID is required", http.StatusBadRequest)
		return
	}
	// Retrieve the device information using the device service
	deviceResponse, err := deviceService.GetSignatureDeviceById(deviceID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	// Set the response header and encode the response to JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	//json.NewEncoder(w).Encode(deviceResponse)
	if err := json.NewEncoder(w).Encode(deviceResponse); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
