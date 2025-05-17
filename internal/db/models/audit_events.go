package models

import "time"

// AuditEvent represents a user lifecycle or security-sensitive event
type AuditEvent struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	EventType string    `json:"event_type"` // e.g. login, logout, password_reset
	IPAddress string    `json:"ip_address"`
	UserAgent string    `json:"user_agent"`
	CreatedAt time.Time `json:"created_at"`
}
