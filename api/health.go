package api

import "net/http"

// HealthResponse defines the response for health check
// swagger:model
type HealthResponse struct {
	Status  string `json:"status"`
	Version string `json:"version"`
}

// Health evaluates the health of the service and writes a standardized response.
// swagger:route GET /health health checkHealth
//
// Evaluates the health of the service.
//
// Responses:
//
//	200: HealthResponse
//	405: ErrorResponse
func (s *Server) Health(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		WriteErrorResponse(response, http.StatusMethodNotAllowed, []string{
			http.StatusText(http.StatusMethodNotAllowed),
		})
		return
	}

	health := HealthResponse{
		Status:  "pass",
		Version: "v0",
	}

	WriteAPIResponse(response, http.StatusOK, health)
}
