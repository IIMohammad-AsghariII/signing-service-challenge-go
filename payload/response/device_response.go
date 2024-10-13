package response

// DeviceResponse response for creating a device
type DeviceResponse struct {
	ID             string
	PublicKey      string
	Label          string
	SignatureCount uint64
}
