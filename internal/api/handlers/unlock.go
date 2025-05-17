package handlers

import (
	"encoding/json"
	"net/http"
	"os"

	corestub "github.com/peithosecure/peitho-backend/internal/corestub"
)

// UnlockValidateHandler godoc
// @Summary Validate unlock license
// @Description Accepts a signed PQC license block, writes it to disk, and validates it
// @Tags license
// @Accept json
// @Produce json
// @Param unlockRequest body map[string]string true "License block payload"
// @Success 200 {object} UnlockSuccessResponse "Unlock succeeded"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 403 {object} map[string]string "Unlock validation failed"
// @Failure 500 {object} map[string]string "Write failed"
// @Router /api/v1/auth/unlock/validate [post]
func UnlockValidateHandler(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Block string `json:"block"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		corestub.RespondWithTraceError(w, "unlock_block_parse_fail", http.StatusBadRequest)
		return
	}

	unlockPath := os.Getenv("UNLOCK_PATH")
	if unlockPath == "" {
		unlockPath = "/app/peitho-core/unlock.lic"
	}

	// Optional: backup existing license
	if _, err := os.Stat(unlockPath); err == nil {
		_ = os.Rename(unlockPath, unlockPath+".bak")
	}

	if err := os.WriteFile(unlockPath, []byte(body.Block), 0600); err != nil {
		corestub.RespondWithTraceError(w, "unlock_write_failed", http.StatusInternalServerError)
		return
	}

	if _, err := corestub.ValidateUnlock(); err != nil {
		corestub.RespondWithTraceError(w, "unlock_invalid", http.StatusForbidden)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(UnlockSuccessResponse{
		Message: "unlocked",
	})
}
