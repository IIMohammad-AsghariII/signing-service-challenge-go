package request

// DeviceRequest request for creating a device
type DeviceRequest struct {
	ID        string `json:"id"`        // JSON label for ID
	Algorithm string `json:"algorithm"` // JSON label for Algorithm
	Label     string `json:"label"`     // JSON label for Label (optional)
}
