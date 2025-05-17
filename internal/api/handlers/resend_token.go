package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	corestub "github.com/peithosecure/peitho-backend/internal/corestub"
	"github.com/peithosecure/peitho-backend/internal/db/sqlite"
)

// ResendVerificationTokenHandler godoc
// @Summary Resend email verification link
// @Description Sends a new verification email if the account exists and is not verified
// @Tags auth
// @Produce json
// @Param email query string true "User's registered email"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string "Missing or already verified"
// @Failure 404 {object} map[string]string "User not found"
// @Failure 500 {object} map[string]string "Email sending or DB error"
// @Router /api/v1/auth/resend-token [get]
func ResendVerificationTokenHandler(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	fmt.Println("💌 Received resend-token request for:", email)

	if email == "" {
		corestub.RespondWithTraceError(w, "missing_email_param", http.StatusBadRequest)
		return
	}

	user, err := sqlite.GetUserByEmail(email)
	if err != nil {
		corestub.RespondWithTraceError(w, "user_lookup_failed", http.StatusInternalServerError)
		return
	}
	if user == nil {
		corestub.RespondWithTraceError(w, "user_not_found", http.StatusNotFound)
		return
	}
	if user.EmailVerified == 1 {
		corestub.RespondWithTraceError(w, "already_verified", http.StatusBadRequest)
		return
	}

	err = SendVerificationEmail(user.Username, user.Email)
	if err != nil {
		corestub.RespondWithTraceError(w, "email_send_failed", http.StatusInternalServerError)
		return
	}

	fmt.Println("✅ Verification email resent to:", email)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "resent",
		"message": "Verification email sent successfully",
	})
}
