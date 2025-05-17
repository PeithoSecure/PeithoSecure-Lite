package handlers

import (
	"encoding/json"
	"net/http"
)

// StatusResponse defines branding and health status metadata
type StatusResponse struct {
	Brand         string `json:"brand" example:"PeithoSecure Lite"`
	LicenseStatus string `json:"license_status" example:"Valid"`
	Copyright     string `json:"copyright" example:"© 2025 Peitho"`
	Mood          string `json:"mood" example:"Stable"`
}

// StatusHandler godoc
// @Summary Service status and branding info
// @Description Returns branding, license status, copyright, and system mood
// @Tags health
// @Produce json
// @Success 200 {object} StatusResponse
// @Router /status [get]
func StatusHandler(w http.ResponseWriter, r *http.Request) {
	resp := StatusResponse{
		Brand:         "Peitho",
		LicenseStatus: "Valid",
		Copyright:     "© 2025 Peitho",
		Mood:          "Stable",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
