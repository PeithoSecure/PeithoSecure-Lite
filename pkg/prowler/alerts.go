package prowler

import (
	"log"
)

// TriggerAlert sends a simulated security alert
func TriggerAlert(reason string) {
	log.Printf("[Prowler] Security alert triggered: %s", reason)
}
