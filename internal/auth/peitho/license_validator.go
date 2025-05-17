package peitho

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	corestub "github.com/peithosecure/peitho-backend/internal/corestub"
)

type LicensePayload struct {
	Email            string `json:"email"`
	DeviceID         string `json:"device_id"`
	IssuedAt         string `json:"issued_at"`
	EngineHash       string `json:"engine_hash"`
	BrandingRequired bool   `json:"branding_required"`
}

func init() {
	if corestub.PeithoTrap() != "__peitho_signature__" {
		fmt.Fprintln(os.Stderr, "❌ PeithoCore stealth check failed — integrity compromised")
		os.Exit(77)
	}
}

// ValidateLicenseToken performs a local unlock.lic validation check,
// binding the engine hash and enforcing license integrity.
func ValidateLicenseToken(_ string) error {
	// ✅ Skip if already marked validated
	if os.Getenv("PEITHO_LICENSE_HASH_OK") == "true" {
		fmt.Println("🔁 License already validated — skipping revalidation")
		return nil
	}

	const unlockPath = "./peitho-core/unlock.lic"

	raw, err := os.ReadFile(unlockPath)
	if err != nil {
		return fmt.Errorf("🔒 unlock.lic missing: %w", err)
	}

	parts := strings.SplitN(string(raw), "||", 2)
	if len(parts) != 2 {
		return errors.New("🧨 invalid unlock.lic structure (missing || delimiter)")
	}

	sigB64 := strings.TrimSpace(parts[0])
	payload := []byte(strings.TrimSpace(parts[1]))

	// ✅ Verify signature
	if err := corestub.ValidateLicenseSignature(sigB64, payload); err != nil {
		return fmt.Errorf("🛑 license verification failed: %w", err)
	}

	// ✅ Decode payload
	var lic LicensePayload
	if err := json.Unmarshal(payload, &lic); err != nil {
		return errors.New("❌ failed to parse license payload JSON")
	}

	// 🔐 Assign hash for later runtime comparison
	corestub.AssignExpectedEngineHash(lic.EngineHash)

	// ✅ Mark validated early to avoid re-exec
	_ = os.Setenv("PEITHO_LICENSE_HASH_OK", "true")

	// 🧬 Hash check
	if !corestub.ValidateEngineHash() {
		corestub.ShutdownNow("trace engine tampering detected")
	}

	// ✅ Valid license — print diagnostics
	fmt.Println("🔓 License validated successfully.")
	fmt.Println("📧 Email         :", lic.Email)
	fmt.Println("🖥️  Licensed For :", lic.DeviceID)
	fmt.Println("🔐 Branding Lock :", lic.BrandingRequired)

	if os.Getenv("PEITHO_ALLOW_MULTI_DEVICE") == "true" {
		fmt.Println("⚠️  Device binding check is DISABLED (multi-device mode).")
		return nil
	}

	actualDeviceID := os.Getenv("PEITHO_DEVICE_ID")
	if actualDeviceID == "" {
		actualDeviceID = "web-default"
	}

	if lic.DeviceID != actualDeviceID {
		return fmt.Errorf("🚫 device mismatch: license bound to '%s', current is '%s'", lic.DeviceID, actualDeviceID)
	}

	return nil
}
