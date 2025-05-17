package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	corestub "github.com/peithosecure/peitho-backend/internal/corestub"
)

// EngineEventPayload defines the structure for custom trace logs
type EngineEventPayload struct {
	Event   string `json:"event" example:"custom_event"`           // Custom event identifier
	Actor   string `json:"actor" example:"USER"`                   // One of: DEV, USER, HACKER
	Message string `json:"message" example:"Manual log injection"` // Human-readable description
	Lock    bool   `json:"lock" example:"false"`                   // Whether to trigger lockdown
}

// EngineEventHandler godoc
// @Summary Submit system trace event
// @Description Log a custom trace event (devtools, external agent, honeypot, etc.)
// @Tags System
// @Accept json
// @Produce json
// @Param payload body EngineEventPayload true "Custom trace event to log"
// @Success 200 {object} GenericMessageResponse "Event logged successfully"
// @Failure 400 {object} map[string]string "Invalid request body"
// @Failure 422 {object} map[string]string "Unknown actor type"
// @Router /api/v1/events/log [post]
func EngineEventHandler(w http.ResponseWriter, r *http.Request) {
	var payload EngineEventPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		corestub.RespondWithTraceError(w, "event_payload_parse_fail", http.StatusBadRequest)
		return
	}

	actor := strings.ToUpper(payload.Actor)
	if actor != "DEV" && actor != "USER" && actor != "HACKER" {
		corestub.RespondWithTraceError(w, "invalid_actor", http.StatusUnprocessableEntity)
		return
	}

	logEntry := "[TRACE] " + time.Now().Format(time.RFC3339) + " | [" + actor + "] " + payload.Event + " → " + payload.Message
	if payload.Lock {
		logEntry += " 🔒 LOCK TRIGGERED"
	}
	corestub.TrackEvent(logEntry)

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(GenericMessageResponse{
		Message: "Event logged successfully",
	})
}
