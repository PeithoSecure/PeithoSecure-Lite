package middleware

import (
	"net/http"
	"strings"

	corestub "github.com/peithosecure/peitho-backend/internal/corestub"
)

// TamperDetectorMiddleware blocks requests with shady User-Agent or branding violations
func TamperDetectorMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ua := r.Header.Get("User-Agent")
		uaLower := strings.ToLower(ua)

		if strings.Contains(uaLower, "postman") ||
			strings.Contains(uaLower, "curl") ||
			strings.Contains(uaLower, "fiddler") ||
			strings.Contains(uaLower, "httpclient") {

			corestub.TrackEvent("tamper_detected")
			corestub.RespondWithTraceError(w, "tamper_detected", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
