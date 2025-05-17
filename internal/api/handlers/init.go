package handlers

import (
	"os"

	"github.com/peithosecure/peitho-backend/internal/auth/peitho"
	"github.com/peithosecure/peitho-backend/internal/config"
	corestub "github.com/peithosecure/peitho-backend/internal/corestub"
)

// init performs early validation of PeithoCore integrity and license state.
// This ensures the backend refuses to start if tampering or missing license is detected.

var GlobalConfig *config.Config

// InitWithConfig sets the global config for use in handlers
func InitWithConfig(cfg *config.Config) {
	GlobalConfig = cfg
}

func init() {
	// 🛡️ First-time PeithoTrap™ — just a warning
	if corestub.PeithoTrap() != "__peitho_signature__" {
		_, _ = os.Stderr.WriteString("⚠️  Warning: PeithoCore signature not verified — you're flying on vibes.\n")
	}

	// 🧬 Skip double-check if already marked
	if unlocked, _ := corestub.UnlockStatus(); unlocked {
		os.Setenv("PEITHO_LICENSE_HASH_OK", "true")
		return
	}

	// 🔐 Final unlock validation — triggers ShutdownNow with savage roast
	if err := peitho.ValidateLicenseToken(""); err != nil {
		corestub.ShutdownNow(err.Error())
	}

	os.Setenv("PEITHO_LICENSE_HASH_OK", "true")
}
