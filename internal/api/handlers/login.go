package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/peithosecure/peitho-backend/internal/auth/keycloak"
	corestub "github.com/peithosecure/peitho-backend/internal/corestub"
	"github.com/peithosecure/peitho-backend/internal/db/models"
	"github.com/peithosecure/peitho-backend/internal/db/sqlite"
	"github.com/peithosecure/peitho-backend/internal/metrics"
	"github.com/peithosecure/peitho-backend/internal/middleware"
)

// LoginHandler godoc
// @Summary Authenticate user credentials
// @Description Logs in the user and returns a Keycloak-issued JWT token pair along with email verification status
// @Tags auth
// @Accept json
// @Produce json
// @Param loginRequest body models.LoginRequest true "Username and password payload"
// @Success 200 {object} models.LoginResponse "Authentication successful"
// @Failure 400 {object} map[string]string "Malformed request or JSON parsing failed"
// @Failure 401 {object} map[string]string "Invalid credentials"
// @Failure 429 {object} map[string]string "Too many failed login attempts"
// @Router /api/v1/auth/login [post]
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var loginReq models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		corestub.RespondWithTraceError(w, "login_parse_fail", http.StatusBadRequest)
		return
	}

	if middleware.IsUserLocked(loginReq.Username) {
		retryAfter := middleware.GetRetryAfterSeconds(loginReq.Username)
		w.Header().Set("Retry-After", retryAfter)
		corestub.RespondWithTraceError(w, "user_locked", http.StatusTooManyRequests)
		return
	}

	tokenResp, err := keycloak.LoginWithPassword(GlobalConfig, loginReq.Username, loginReq.Password)
	if err != nil {
		middleware.IncrementLoginFailure(loginReq.Username)
		corestub.RespondWithTraceError(w, "auth_failed", http.StatusUnauthorized)
		return
	}

	user, err := sqlite.GetUserByUsername(loginReq.Username)
	if err != nil || user == nil {
		corestub.RespondWithTraceError(w, "user_lookup_failed", http.StatusInternalServerError)
		return
	}

	middleware.ClearLoginAttempts(loginReq.Username)

	metrics.IncIssued()
	_ = sqlite.LogAuditEvent(user.Username, "login", r.RemoteAddr, r.UserAgent())

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.LoginResponse{
		AccessToken:   tokenResp.AccessToken,
		RefreshToken:  tokenResp.RefreshToken,
		TokenType:     tokenResp.TokenType,
		ExpiresIn:     tokenResp.ExpiresIn,
		EmailVerified: user.EmailVerified != 0,
	})
}
