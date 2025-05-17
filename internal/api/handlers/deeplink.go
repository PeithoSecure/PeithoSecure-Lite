package handlers

import (
	"net/http"

	corestub "github.com/peithosecure/peitho-backend/internal/corestub"
)

// DeepLinkRedirectHandler godoc
// @Summary Legacy deep link redirect
// @Description Redirects to PeithoSecure app using deep link for verify/reset tokens (legacy path)
// @Tags DeepLink
// @Produce plain
// @Param type query string true "Link type (verify/reset)"
// @Param token query string true "Verification or reset token"
// @Success 302 {string} string "Redirect to app URI"
// @Failure 400 {object} map[string]string "Missing or invalid query parameters"
// @Router /api/v1/deeplink/legacy [get]
func DeepLinkRedirectHandler(w http.ResponseWriter, r *http.Request) {
	linkType := r.URL.Query().Get("type")
	token := r.URL.Query().Get("token")

	if linkType == "" || token == "" {
		corestub.RespondWithTraceError(w, "missing_deeplink_params", http.StatusBadRequest)
		return
	}

	var frontendURL string
	switch linkType {
	case "verify":
		frontendURL = "peitho://verify?token=" + token
	case "reset":
		frontendURL = "peitho://setup-password?token=" + token + "&reset=true"
	default:
		corestub.RespondWithTraceError(w, "invalid_deeplink_type", http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, frontendURL, http.StatusFound)
}
