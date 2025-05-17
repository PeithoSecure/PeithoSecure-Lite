package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/peithosecure/peitho-backend/internal/auth/keycloak"
	corestub "github.com/peithosecure/peitho-backend/internal/corestub"
	"github.com/peithosecure/peitho-backend/internal/db/sqlite"
	"github.com/peithosecure/peitho-backend/internal/metrics"
	"github.com/peithosecure/peitho-backend/internal/middleware"
)

// LogoutRequest represents the incoming payload to logout
type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" example:"eyJhbGc..."`
}

// LogoutResponse represents the standard logout success response
type LogoutResponse struct {
	Message string `json:"message" example:"User logged out successfully."`
}

// LogoutHandler godoc
// @Summary Logout current user
// @Description Invalidates the user's refresh token and clears session state
// @Tags auth
// @Accept json
// @Produce json
// @Param logoutRequest body LogoutRequest true "Refresh token payload"
// @Success 200 {object} LogoutResponse
// @Failure 400 {object} map[string]string "Malformed request or missing token"
// @Failure 401 {object} map[string]string "Invalid or expired token"
// @Router /api/v1/auth/logout [post]
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	var body LogoutRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		corestub.RespondWithTraceError(w, "logout_parse_fail", http.StatusBadRequest)
		return
	}

	if err := keycloak.RevokeRefreshToken(GlobalConfig, body.RefreshToken); err != nil {
		corestub.RespondWithTraceError(w, "logout_failed", http.StatusUnauthorized)
		return
	}

	metrics.IncRevoked()

	ctxUser := r.Context().Value(middleware.UserContextKey)
	if claims, ok := ctxUser.(map[string]interface{}); ok {
		if uname, ok := claims["preferred_username"].(string); ok {
			_ = sqlite.LogAuditEvent(uname, "logout", r.RemoteAddr, r.UserAgent())
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(LogoutResponse{
		Message: "User logged out successfully.",
	})
}
