package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	corestub "github.com/peithosecure/peitho-backend/internal/corestub"
)

// UniversalDeeplinkHandler godoc
// @Summary Universal deep link handler
// @Description Redirects users to the appropriate app or web route based on platform and deep link type (verify/reset).
// @Tags DeepLink
// @Produce html
// @Param type query string true "Link type (verify or reset)"
// @Param token query string true "Verification or reset token"
// @Success 302 {string} string "Redirect to app or browser path"
// @Failure 400 {object} map[string]string "Missing or invalid query parameters"
// @Router /api/v1/deeplink [get]
func UniversalDeeplinkHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	linkType := query.Get("type")
	token := query.Get("token")

	if linkType == "" || token == "" {
		corestub.RespondWithTraceError(w, "missing_deeplink_params", http.StatusBadRequest)
		return
	}

	var appPath, webPath string
	switch linkType {
	case "verify":
		appPath = fmt.Sprintf("peitho://verify?token=%s", token)
		webPath = fmt.Sprintf("%s/verify?token=%s", emailConfig.FrontendURL, token)
	case "reset":
		appPath = fmt.Sprintf("peitho://reset?token=%s", token)
		webPath = fmt.Sprintf("%s/reset-password?token=%s&reset=true", emailConfig.FrontendURL, token)
	default:
		corestub.RespondWithTraceError(w, "invalid_deeplink_type", http.StatusBadRequest)
		return
	}

	userAgent := strings.ToLower(r.UserAgent())
	isMobile := strings.Contains(userAgent, "iphone") ||
		strings.Contains(userAgent, "ipad") ||
		strings.Contains(userAgent, "android")
	isDesktop := strings.Contains(userAgent, "macintosh") ||
		strings.Contains(userAgent, "windows") ||
		strings.Contains(userAgent, "linux")

	switch {
	case isMobile:
		log.Printf("📡 [DEEPLINK] Mobile UA: %s → %s", userAgent, appPath)
		http.Redirect(w, r, appPath, http.StatusFound)
	case isDesktop:
		log.Printf("📡 [DEEPLINK] Desktop UA: %s → Launching %s", userAgent, appPath)
		w.Header().Set("Content-Type", "text/html")
		html := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
  <title>PeithoSecure Link</title>
  <meta http-equiv="refresh" content="2; url=%s" />
</head>
<body style="font-family:sans-serif;text-align:center;padding:2rem;">
  <p>Trying to open the PeithoSecure app...</p>
  <p><a href="%s">Click here if not redirected</a></p>
  <script>
    setTimeout(function() {
      window.location = "%s";
    }, 500);
  </script>
</body>
</html>
`, appPath, appPath, appPath)
		fmt.Fprint(w, html)
	default:
		log.Printf("📡 [DEEPLINK] Unknown UA: %s → %s", userAgent, webPath)
		http.Redirect(w, r, webPath, http.StatusFound)
	}
}
