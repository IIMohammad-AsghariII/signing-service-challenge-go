package request

// SignTransactionRequest request for signing a transaction
type SignTransactionRequest struct {
	DeviceID string `json:"deviceId"` // JSON label for DeviceID
	Data     string `json:"data"`     // JSON label for Data
}
