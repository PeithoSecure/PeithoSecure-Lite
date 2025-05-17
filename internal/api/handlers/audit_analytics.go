package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	corestub "github.com/peithosecure/peitho-backend/internal/corestub"
	"github.com/peithosecure/peitho-backend/internal/db/sqlite"
	"github.com/peithosecure/peitho-backend/internal/middleware"
)

// AuditAnalyticsHandler godoc
// @Summary View audit activity
// @Description Returns a timeline of recent audit events for this user
// @Tags Analytics
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.AuditEvent
// @Failure 401 {object} handlers.GenericErrorResponse
// @Failure 500 {object} handlers.GenericErrorResponse
// @Router /api/v1/analytics/audit [get]
func AuditAnalyticsHandler(w http.ResponseWriter, r *http.Request) {
	username, err := middleware.ExtractUsernameFromContext(r.Context())
	if err != nil || username == "" {
		corestub.RespondWithTraceError(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	log.Printf("🔍 AuditAnalyticsHandler called for: %s", username)

	events, err := sqlite.GetAuditEventsByUsername(username, 50)
	if err != nil {
		log.Printf("❌ Failed to fetch audit events for %s: %v", username, err)
		corestub.RespondWithTraceError(w, "audit_query_failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(events)
}
