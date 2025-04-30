package models

// HealthResponse represents the response from the health check endpoint
type HealthResponse struct {
	Status string `json:"status"`
}
