package handlers

import (
	"encoding/json"
	"net/http"
)

// HealthCheckHandler godoc
// @Summary Service health check
// @Description Returns the operational status of the PeithoSecure Lite backend
// @Tags Health
// @Produce json
// @Success 200 {object} HealthCheckResponse "Service is running"
// @Router /health [get]
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	resp := HealthCheckResponse{
		Status:  "ok",
		Service: "PeithoSecure Lite",
		Version: "1.0.0",
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

// HealthCheckResponse defines the structure of the /health response
type HealthCheckResponse struct {
	Status  string `json:"status" example:"ok"`
	Service string `json:"service" example:"PeithoSecure Lite"`
	Version string `json:"version" example:"1.0.0"`
}
