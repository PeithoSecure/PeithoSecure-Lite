package handlers

import (
	"encoding/json"
	"net/http"

	corestub "github.com/peithosecure/peitho-backend/internal/corestub"
	"github.com/peithosecure/peitho-backend/pkg/prowler"
)

// ProwlerScanResponse is the full response payload for /security-scan
type ProwlerScanResponse struct {
	Status  string             `json:"status" example:"completed"`
	Summary prowler.ScanResult `json:"summary"`
}

// ProwlerScanHandler godoc
// @Summary Trigger Prowler security audit
// @Description Runs a limited Peitho Prowler scan to check for signs of tampering, token leakage, or license gaps
// @Tags security
// @Produce json
// @Success 200 {object} ProwlerScanResponse
// @Failure 403 {object} map[string]string "License lock or tamper guard triggered"
// @Router /api/v1/security-scan [get]
func ProwlerScanHandler(w http.ResponseWriter, r *http.Request) {
	corestub.TrackEvent("prowler_scan_triggered")

	result := prowler.RunSecurityScan()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ProwlerScanResponse{
		Status:  "completed",
		Summary: result,
	})
}
