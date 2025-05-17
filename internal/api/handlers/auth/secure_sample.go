package handlers

import (
	"net/http"

	corestub "github.com/peithosecure/peitho-backend/internal/corestub"
	"github.com/peithosecure/peitho-backend/internal/middleware"
)

// SecureSampleHandler is a sample protected route
func SecureSampleHandler(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(middleware.UserContextKey)
	if user == nil {
		corestub.RespondWithTraceError(w, "unauthorized_access", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Secure data accessed. Well done, operator."}`))
}
