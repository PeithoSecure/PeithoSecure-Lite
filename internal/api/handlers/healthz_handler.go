package handlers

import (
	"encoding/json"
	"net/http"
)

// HealthzResponse defines the structure for infra probe responses
type HealthzResponse struct {
	Status  string `json:"status" example:"ok"`
	Service string `json:"service" example:"PeithoSecure Lite"`
	Version string `json:"version" example:"1.0.0"`
}

// HealthzHandler godoc
// @Summary Infra health probe
// @Description Lightweight endpoint for readiness/liveness probes (Kubernetes, CI/CD, etc.)
// @Tags Health
// @Produce json
// @Success 200 {object} HealthzResponse
// @Router /healthz [get]
func HealthzHandler(w http.ResponseWriter, r *http.Request) {
	response := HealthzResponse{
		Status:  "ok",
		Service: "PeithoSecure Lite",
		Version: "1.0.0",
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}
