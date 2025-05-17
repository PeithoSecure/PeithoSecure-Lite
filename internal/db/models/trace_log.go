package models

import "time"

// TraceLog represents an in-memory or persistent trace event
type TraceLog struct {
	ID        int64     `db:"id"`
	UserID    int64     `db:"user_id"`
	Event     string    `db:"event"`
	CreatedAt time.Time `db:"created_at"`
}
