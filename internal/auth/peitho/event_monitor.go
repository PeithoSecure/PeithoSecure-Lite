package peitho

import (
	"log"
)

// TriggerEventHook is a public-stubbed event monitor (real roast hooks private)
func TriggerEventHook(userID string, reason string) {
	log.Printf("[EventMonitor] Activity detected for user %s: %s", userID, reason)
}
