package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/peithosecure/peitho-backend/internal/auth/keycloak"
	"github.com/peithosecure/peitho-backend/internal/config"
	corestub "github.com/peithosecure/peitho-backend/internal/corestub"
	"github.com/peithosecure/peitho-backend/internal/db/sqlite"
	"github.com/peithosecure/peitho-backend/internal/utils"
)

var appConfig *config.Config

func InjectConfig(cfg *config.Config) {
	appConfig = cfg
}

// PasswordResetRequest defines the input for requesting a reset link
type PasswordResetRequest struct {
	Email string `json:"email" example:"user@example.com"`
}

// PasswordResetConfirm defines the input for setting a new password
type PasswordResetConfirm struct {
	Token    string `json:"token" example:"abc123token"`
	Password string `json:"password" example:"SuperSecurePassword123!"`
}

// RequestPasswordResetHandler godoc
// @Summary Request password reset link
// @Description Sends a password reset email to the user with a one-time token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body PasswordResetRequest true "Email to receive reset link"
// @Success 200 {object} GenericMessageResponse
// @Failure 400 {object} map[string]string "Missing or invalid email"
// @Failure 500 {object} map[string]string "Internal error or email send failed"
// @Router /api/v1/auth/request-password-reset [post]
func RequestPasswordResetHandler(w http.ResponseWriter, r *http.Request) {
	var body map[string]string
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body["email"] == "" {
		corestub.RespondWithTraceError(w, "missing_email", http.StatusBadRequest)
		return
	}

	email := strings.ToLower(strings.TrimSpace(body["email"]))
	token := utils.GenerateSecureToken(32)

	if err := sqlite.InsertEmailToken(email, token, "password_reset"); err != nil {
		corestub.RespondWithTraceError(w, "reset_token_insert_failed", http.StatusInternalServerError)
		return
	}

	if err := SendPasswordResetEmail(email, token); err != nil {
		corestub.RespondWithTraceError(w, "reset_email_failed", http.StatusInternalServerError)
		return
	}

	fmt.Printf("🔑 Password reset token sent to %s\n", email)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(GenericMessageResponse{
		Message: "Password reset email sent",
	})
}

// ResetPasswordHandler godoc
// @Summary Reset password using token
// @Description Accepts a one-time token and new password, and resets the account password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body PasswordResetConfirm true "Token and new password"
// @Success 200 {object} GenericMessageResponse
// @Failure 400 {object} map[string]string "Invalid or missing token/password"
// @Failure 500 {object} map[string]string "Internal error or reset failure"
// @Router /api/v1/auth/reset-password [post]
func ResetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	var body map[string]string
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		corestub.RespondWithTraceError(w, "reset_body_invalid", http.StatusBadRequest)
		return
	}
	token := body["token"]
	newPass := body["password"]
	if token == "" || newPass == "" {
		corestub.RespondWithTraceError(w, "reset_missing_token_or_pass", http.StatusBadRequest)
		return
	}

	email, err := sqlite.GetUsernameByTokenAndType(token, "password_reset")
	if err != nil || email == "" {
		corestub.RespondWithTraceError(w, "reset_token_invalid", http.StatusBadRequest)
		return
	}

	user, err := sqlite.GetUserByEmail(email)
	if err != nil || user == nil {
		corestub.RespondWithTraceError(w, "user_not_found", http.StatusInternalServerError)
		return
	}

	if appConfig == nil {
		corestub.RespondWithTraceError(w, "missing_config", http.StatusInternalServerError)
		return
	}

	if err := keycloak.ResetPassword(appConfig, user.Username, newPass); err != nil {
		corestub.RespondWithTraceError(w, "kc_reset_failed", http.StatusInternalServerError)
		return
	}

	_ = sqlite.DeleteVerificationToken(token)

	fmt.Printf("✅ Password reset for user: %s at %v\n", user.Username, time.Now())
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(GenericMessageResponse{
		Message: "Password updated successfully",
	})
}
