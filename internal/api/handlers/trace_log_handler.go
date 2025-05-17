package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	corestub "github.com/peithosecure/peitho-backend/internal/corestub"
)

// TraceEventView is a stripped-down representation of a trace event
type TraceEventView struct {
	ID        string `json:"id" example:"7c8bfc1d"`
	Message   string `json:"message" example:"unauthorized access detected"`
	Actor     string `json:"actor" example:"USER"`
	Event     string `json:"event" example:"auth_failed"`
	Severity  string `json:"severity" example:"low"`
	Lock      bool   `json:"lock" example:"false"`
	Timestamp string `json:"timestamp" example:"2025-05-16T09:12:00Z"`
}

// TraceLogHandler godoc
// @Summary View in-memory trace logs
// @Description Returns recent trace events stored in memory (for debugging, audits, or escalations)
// @Tags logs
// @Produce json
// @Success 200 {array} TraceEventView
// @Security ApiKeyAuth
// @Router /api/v1/log/trace [get]
func TraceLogHandler(w http.ResponseWriter, r *http.Request) {
	traces := corestub.GetRecentTraces()

	payload := make([]TraceEventView, 0, len(traces))
	for _, t := range traces {
		payload = append(payload, TraceEventView{
			ID:        t.ID,
			Message:   t.Message,
			Actor:     t.Actor,
			Event:     t.Event,
			Severity:  t.Severity,
			Lock:      t.Lock,
			Timestamp: t.Timestamp.Format(time.RFC3339),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(payload)
}
