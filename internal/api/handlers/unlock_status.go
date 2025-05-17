package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	corestub "github.com/peithosecure/peitho-backend/internal/corestub"
)

// UnlockStatusResponse defines license status payload
type UnlockStatusResponse struct {
	Unlocked          bool      `json:"unlocked" example:"true"`
	UnlockedAt        time.Time `json:"unlocked_at,omitempty" example:"2025-05-16T02:50:19Z"`
	ExpiresAt         time.Time `json:"expires_at,omitempty" example:"2025-08-14T02:50:19Z"`
	ServerTime        time.Time `json:"server_time" example:"2025-05-16T09:30:00Z"`
	SecuredBy         string    `json:"secured_by" example:"Peitho 🔐"`
	TraceEngineActive bool      `json:"trace_engine_active" example:"true"`
	BrandingLocked    bool      `json:"branding_locked" example:"true"`
}

const enforcedBrandingSignature = "Peitho 🔐"

// UnlockStatusHandler godoc
// @Summary License Unlock Status
// @Description Returns license unlock state, branding status, and server clock
// @Tags License
// @Produce json
// @Success 200 {object} UnlockStatusResponse
// @Failure 403 {object} map[string]string
// @Router /api/v1/auth/unlock-status [get]
func UnlockStatusHandler(w http.ResponseWriter, r *http.Request) {
	status, unlockedAt := corestub.UnlockStatus()

	if !status {
		corestub.RespondWithTraceError(w, "core_locked", http.StatusForbidden)
		return
	}

	if enforcedBrandingSignature != "Peitho 🔐" {
		corestub.RespondWithTraceError(w, "branding_override_detected", http.StatusForbidden)
		return
	}

	payload := UnlockStatusResponse{
		Unlocked:          status,
		UnlockedAt:        unlockedAt,
		ExpiresAt:         unlockedAt.Add(90 * 24 * time.Hour),
		ServerTime:        time.Now().UTC(),
		SecuredBy:         enforcedBrandingSignature,
		TraceEngineActive: true,
		BrandingLocked:    corestub.ValidateBrand(),
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(payload)
}
