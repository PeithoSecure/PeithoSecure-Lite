package middleware

import (
	"net/http"

	corestub "github.com/peithosecure/peitho-backend/internal/corestub"
)

// UnlockGuardMiddleware blocks requests if unlock is invalid
func UnlockGuardMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		unlocked, _ := corestub.UnlockStatus()
		if !unlocked {
			// Stubbed event tracking and roast trigger
			corestub.RespondWithTraceError(w, "pqc_lock_enforced", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
