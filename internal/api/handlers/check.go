package handlers

import (
	"encoding/json"
	"net/http"

	corestub "github.com/peithosecure/peitho-backend/internal/corestub"
	"github.com/peithosecure/peitho-backend/internal/db/sqlite"
	"github.com/peithosecure/peitho-backend/internal/middleware"
)

// CheckEmailHandler godoc
// @Summary Check if email exists and is verified
// @Description Determines if the email exists in the system and whether it has been verified
// @Tags Auth
// @Produce json
// @Param email query string true "Email address to check"
// @Success 200 {object} CheckEmailResponse
// @Failure 400 {object} map[string]string "Invalid email"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/auth/check [get]
func CheckEmailHandler(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	if email == "" || !isValidEmail(email) {
		corestub.RespondWithTraceError(w, "email_invalid", http.StatusBadRequest)
		return
	}

	user, err := sqlite.GetUserByEmail(email)
	if err != nil {
		corestub.RespondWithTraceError(w, "check_email_failed", http.StatusInternalServerError)
		return
	}

	exists := false
	verified := false

	if user != nil {
		verified = user.EmailVerified != 0
		exists = verified // Only "exists" if verified
	}

	resp := CheckEmailResponse{
		Exists:   exists,
		Verified: verified,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		corestub.RespondWithTraceError(w, "encode_fail", http.StatusInternalServerError)
	}
}

// CheckEmailResponse defines the structure returned by CheckEmailHandler
// @Description Email presence and verification status
type CheckEmailResponse struct {
	Exists   bool `json:"exists" example:"true"`
	Verified bool `json:"verified" example:"true"`
}

// SecureSampleHandler godoc
// @Summary Sample protected route
// @Description Returns a success message for authenticated users
// @Tags Auth
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 401 {object} map[string]string "Unauthorized"
// @Security BearerAuth
// @Router /api/v1/auth/secure-sample [get]
func SecureSampleHandler(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(middleware.UserContextKey)
	if user == nil {
		corestub.RespondWithTraceError(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Secure data accessed successfully.",
	})
}
