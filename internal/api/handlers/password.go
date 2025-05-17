package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/peithosecure/peitho-backend/internal/auth/keycloak"
	corestub "github.com/peithosecure/peitho-backend/internal/corestub"
	"github.com/peithosecure/peitho-backend/internal/db/models"
	"github.com/peithosecure/peitho-backend/internal/db/sqlite"
)

// SetupPasswordRequest is used to bind token and password
type SetupPasswordRequest struct {
	Token    string `json:"token" example:"verify-token-abc123"`
	Password string `json:"password" example:"StrongPassword123!"`
}

// SetupPasswordHandler godoc
// @Summary Finalize user account setup
// @Description Sets initial password after email verification, creates Keycloak user, and writes PQC license
// @Tags auth
// @Accept json
// @Produce json
// @Param request body SetupPasswordRequest true "Token and new password"
// @Success 200 {object} GenericMessageResponse
// @Failure 400 {object} map[string]string "Missing or invalid input"
// @Failure 404 {object} map[string]string "User not found"
// @Failure 409 {object} map[string]string "Already verified or setup attempted twice"
// @Failure 500 {object} map[string]string "Server, Keycloak, or license error"
// @Router /api/v1/auth/setup-password [post]
func SetupPasswordHandler(w http.ResponseWriter, r *http.Request) {
	var body map[string]string
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		corestub.RespondWithTraceError(w, "parse_error", http.StatusBadRequest)
		return
	}

	token := body["token"]
	newPass := body["password"]
	if token == "" || newPass == "" {
		corestub.RespondWithTraceError(w, "missing_token_or_password", http.StatusBadRequest)
		return
	}

	var user *models.User
	var err error

	if token == "pending" {
		email := r.Header.Get("Email")
		if email == "" {
			corestub.RespondWithTraceError(w, "missing_email_header", http.StatusBadRequest)
			return
		}
		user, err = sqlite.GetUserByEmail(email)
		if err != nil || user == nil || user.EmailVerified == 0 {
			corestub.RespondWithTraceError(w, "user_not_verified", http.StatusBadRequest)
			return
		}
	} else {
		email, err := sqlite.GetUsernameByTokenAndType(token, "verify")
		if err != nil || email == "" {
			corestub.RespondWithTraceError(w, "invalid_token", http.StatusBadRequest)
			return
		}
		user, err = sqlite.GetUserByEmail(email)
		if err != nil || user == nil {
			corestub.RespondWithTraceError(w, "user_not_found", http.StatusNotFound)
			return
		}
	}

	exists, err := keycloak.UserExists(GlobalConfig, user.Username)
	if err != nil {
		corestub.RespondWithTraceError(w, "keycloak_check_failed", http.StatusInternalServerError)
		return
	}

	if exists {
		err = keycloak.ResetPassword(GlobalConfig, user.Username, newPass)
		if err != nil {
			corestub.RespondWithTraceError(w, "password_reset_failed", http.StatusInternalServerError)
			return
		}
	} else {
		err = keycloak.RegisterNewUserWithEmail(GlobalConfig, user.Username, newPass, user.Email)
		if err != nil {
			corestub.RespondWithTraceError(w, "user_create_failed", http.StatusInternalServerError)
			return
		}
	}

	deviceID := r.Header.Get("Device-ID")
	if deviceID == "" {
		deviceID = "web-default"
	}

	block, err := corestub.GenerateSignedLicense(user.Username, deviceID, true)
	if err != nil {
		corestub.RespondWithTraceError(w, "license_gen_failed", http.StatusInternalServerError)
		return
	}

	unlockPath := os.Getenv("UNLOCK_PATH")
	if unlockPath == "" {
		unlockPath = "/app/peitho-core/unlock.lic"
	}
	if _, err := os.Stat(unlockPath); err == nil {
		_ = os.Rename(unlockPath, unlockPath+".bak")
	}

	if err := os.WriteFile(unlockPath, []byte(block), 0600); err != nil {
		corestub.RespondWithTraceError(w, "license_write_failed", http.StatusInternalServerError)
		return
	}

	if token != "pending" {
		_ = sqlite.DeleteVerificationToken(token)
	}
	_ = sqlite.LogAuditEvent(user.Username, "password_set", r.RemoteAddr, r.UserAgent())

	fmt.Printf("✅ Account finalized and password SET for: %s\n", user.Username)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(GenericMessageResponse{
		Message: "Password set and account activated.",
	}); err != nil {
		corestub.RespondWithTraceError(w, "response_encode_failed", http.StatusInternalServerError)
	}
}
