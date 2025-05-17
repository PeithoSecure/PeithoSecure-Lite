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

// RefreshRequest is used to request a new access token
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" example:"eyJhbGc..."`
}

// RefreshResponse represents the token response payload
type RefreshResponse struct {
	AccessToken  string `json:"access_token" example:"eyJhbGc..."`
	RefreshToken string `json:"refresh_token" example:"eyJhbGc..."`
	TokenType    string `json:"token_type" example:"Bearer"`
	ExpiresIn    int    `json:"expires_in" example:"300"`
}

// RefreshHandler godoc
// @Summary Refresh access token
// @Description Uses a valid refresh token to issue a new access token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RefreshRequest true "Refresh token payload"
// @Success 200 {object} RefreshResponse
// @Failure 400 {object} map[string]string "Malformed request or missing refresh token"
// @Failure 401 {object} map[string]string "Invalid or expired refresh token"
// @Router /api/v1/auth/refresh [post]
func RefreshHandler(w http.ResponseWriter, r *http.Request) {
	var body RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		corestub.RespondWithTraceError(w, "refresh_parse_fail", http.StatusBadRequest)
		return
	}

	tokenResp, err := keycloak.RefreshWithToken(GlobalConfig, body.RefreshToken)
	if err != nil {
		corestub.RespondWithTraceError(w, "refresh_invalid", http.StatusUnauthorized)
		return
	}

	metrics.IncRefreshed()

	ctxUser := r.Context().Value(middleware.UserContextKey)
	if claims, ok := ctxUser.(map[string]interface{}); ok {
		if uname, ok := claims["preferred_username"].(string); ok {
			_ = sqlite.LogAuditEvent(uname, "token_refreshed", r.RemoteAddr, r.UserAgent())
		}
	}

	w.Header().Set("Content-Type", "application/json")
	accessToken, _ := tokenResp["access_token"].(string)
	refreshToken, _ := tokenResp["refresh_token"].(string)
	tokenType, _ := tokenResp["token_type"].(string)
	expiresIn, _ := tokenResp["expires_in"].(float64)

	json.NewEncoder(w).Encode(RefreshResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    tokenType,
		ExpiresIn:    int(expiresIn),
	})
}
