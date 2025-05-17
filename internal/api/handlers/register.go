package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	corestub "github.com/peithosecure/peitho-backend/internal/corestub"
	"github.com/peithosecure/peitho-backend/internal/db/sqlite"
)

// RegisterRequest defines payload for new account registration
type RegisterRequest struct {
	Username string `json:"username" example:"johndoe"`
	Email    string `json:"email" example:"john@example.com"`
}

// RegisterResponse defines response after registration
type RegisterResponse struct {
	Message string `json:"message" example:"user registered successfully, please verify email"`
}

// RegisterHandler godoc
// @Summary Register a new user (deferred setup)
// @Description Creates a local user, sends verification email, and defers password/keycloak setup
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Email and username payload"
// @Success 201 {object} RegisterResponse
// @Failure 400 {object} map[string]string "Malformed request or missing fields"
// @Failure 409 {object} map[string]string "Email already registered and verified"
// @Failure 500 {object} map[string]string "Database or email error"
// @Router /api/v1/auth/register [post]
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		corestub.RespondWithTraceError(w, "invalid_request", http.StatusBadRequest)
		return
	}

	if req.Email == "" || !isValidEmail(req.Email) {
		corestub.RespondWithTraceError(w, "email_format_invalid", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(req.Username) == "" {
		corestub.RespondWithTraceError(w, "username_missing", http.StatusBadRequest)
		return
	}

	existingUser, _ := sqlite.GetUserByEmail(req.Email)
	if existingUser != nil && existingUser.EmailVerified == 1 {
		corestub.RespondWithTraceError(w, "email_already_verified", http.StatusConflict)
		return
	}

	db := sqlite.GetDB()
	_, err := db.Exec(`
		INSERT INTO users (username, email, role, email_verified)
		VALUES (?, ?, ?, ?)
		ON CONFLICT(email) DO UPDATE SET username=excluded.username, email_verified=0
	`, req.Username, req.Email, "user", 0)
	if err != nil {
		corestub.RespondWithTraceError(w, "user_save_fail", http.StatusInternalServerError)
		return
	}

	_ = sqlite.LogAuditEvent(req.Username, "user_registered", r.RemoteAddr, r.UserAgent())

	go func() {
		if err := SendVerificationEmail(req.Username, req.Email); err != nil {
			fmt.Printf("[!] Failed to send email to %s: %v\n", req.Email, err)
		}
	}()

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(RegisterResponse{
		Message: "user registered successfully, please verify email",
	})
}
