package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/peithosecure/peitho-backend/internal/auth/keycloak"
	corestub "github.com/peithosecure/peitho-backend/internal/corestub"
	"github.com/peithosecure/peitho-backend/internal/db/sqlite"
)

// VerifyEmailHandler godoc
// @Summary Verify user email
// @Description Verifies email via token, marks user as verified, registers in Keycloak, and issues PQC license
// @Tags Email
// @Produce json
// @Param token query string true "Verification token"
// @Success 200 {object} EmailVerificationResponse
// @Failure 400 {object} map[string]string "Invalid token"
// @Failure 409 {object} map[string]string "Already verified"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/auth/verify-email [get]
func VerifyEmailHandler(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		corestub.RespondWithTraceError(w, "missing_token", http.StatusBadRequest)
		return
	}

	fmt.Printf("🔍 [VerifyEmailHandler] Incoming token: %s\n", token)

	email, err := sqlite.GetUsernameByToken(token)
	if err != nil || email == "" {
		fmt.Printf("❌ [VerifyEmailHandler] Token lookup failed: %v\n", err)
		corestub.RespondWithTraceError(w, "invalid_token", http.StatusBadRequest)
		return
	}

	user, err := sqlite.GetUserByEmail(email)
	if err != nil || user == nil {
		fmt.Printf("❌ [VerifyEmailHandler] User not found for email: %s\n", email)
		corestub.RespondWithTraceError(w, "user_not_found", http.StatusInternalServerError)
		return
	}

	if user.EmailVerified == 1 {
		fmt.Printf("⚠️ [VerifyEmailHandler] Email already verified: %s\n", user.Username)
		corestub.RespondWithTraceError(w, "already_verified", http.StatusConflict)
		return
	}

	if err := sqlite.MarkEmailVerified(user.Username); err != nil {
		corestub.RespondWithTraceError(w, "verify_fail", http.StatusInternalServerError)
		return
	}

	fmt.Printf("✅ [VerifyEmailHandler] Email verified for user: %s at %v\n", user.Username, time.Now())

	err = keycloak.RegisterNewUserWithEmail(GlobalConfig, user.Username, "", user.Email)
	if err != nil {
		corestub.RespondWithTraceError(w, "keycloak_user_create_failed", http.StatusInternalServerError)
		return
	}

	go func() {
		deviceID := r.Header.Get("Device-ID")
		if deviceID == "" {
			deviceID = "web-default"
		}

		block, err := corestub.GenerateSignedLicense(user.Username, deviceID, true)
		if err != nil {
			fmt.Printf("[!] Failed to generate license for %s: %v\n", user.Username, err)
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
			fmt.Printf("[!] Failed to write unlock.lic for %s: %v\n", user.Username, err)
		} else {
			fmt.Printf("[+] PQC license written for user: %s → %s\n", user.Username, unlockPath)
		}
	}()

	_ = sqlite.LogAuditEvent(user.Username, "email_verified", r.RemoteAddr, r.UserAgent())

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(EmailVerificationResponse{
		Verified: true,
		Username: user.Username,
		Message:  "Email verified and Keycloak account created.",
	})
}
