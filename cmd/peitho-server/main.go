package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	licenseguard "github.com/peithosecure/peitho-backend/internal/corestub"

	"github.com/joho/godotenv"
	_ "github.com/peithosecure/peitho-backend/docs"
	"github.com/peithosecure/peitho-backend/internal/api/handlers"
	passwordreset "github.com/peithosecure/peitho-backend/internal/api/handlers"
	"github.com/peithosecure/peitho-backend/internal/api/routes"
	"github.com/peithosecure/peitho-backend/internal/config"
	"github.com/peithosecure/peitho-backend/internal/db/sqlite"
	"github.com/peithosecure/peitho-backend/internal/metrics"
	"github.com/peithosecure/peitho-backend/internal/middleware"
)

const unlockPath = "/app/peitho-core/unlock.lic"

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  No .env file found. We're flying environmental freestyle.")
	}

	// 🥩 Roast-Only Mode
	if os.Getenv("PEITHO_ROAST_ONLY") == "true" {
		printAsciiBanner()
		log.Println("💥 Roast Engine™ stub initialized.")
		log.Println("🧠 Memory is empty. No unlock present.")
		log.Println("🔗 Branding not verified. We’re winging it.")
		log.Println("🔥 Proceeding to fry illusions... Goodbye.")
		return
	}

	// ⏳ Ensure license presence
	waitForLicense(unlockPath)

	// 🔐 Validate license if not already trusted
	unlocked, _ := licenseguard.UnlockStatus()
	if !unlocked {

		if payload, err := licenseguard.ValidateUnlock(); err != nil {
			printAsciiBanner()
			log.Println("🧨 License Validation Error:")
			log.Printf("🔥 %s", payload.Message)
			log.Printf("🧬 Runtime Engine Hash: %s", licenseguard.CurrentEngineHash())

			if raw, rerr := os.ReadFile(unlockPath); rerr == nil {
				log.Printf("📄 unlock.lic contents: %s", truncateMiddle(string(raw), 160))
				parts := splitParts(string(raw))
				if len(parts) == 2 {
					var obj map[string]interface{}
					if jerr := json.Unmarshal([]byte(parts[1]), &obj); jerr == nil {
						log.Printf("📦 Parsed License Payload: %+v", obj)
						if h, ok := obj["engine_hash"]; ok {
							log.Printf("🔗 License Engine Hash: %s", h)
						}
					}
				}
			}

			log.Println("💡 Try: `make force-license docker up` to regenerate a valid license.")
			log.Fatalf("💥 BOOM! The system self-terminated due to license fraud: %v", err)
		}

		os.Setenv("PEITHO_LICENSE_HASH_OK", "true")
		log.Println("✅ Post-Quantum License Check: Passed with flying colors 🪂")
		log.Println("🔐 PQC Mode: Dilithium2 (via CIRCL)")
		log.Println("🧠 Roast Engine Fingerprint: 🔒 Verified")
		log.Println("🚨 Lockdown Protocol: Active and Sassy")
		log.Println("📛 Branding Check: PeithoCore Authentic ✔️")
		log.Println("🔗 Secured by Peitho 🔐 — anything less is a joke")
	} else {
		log.Println("🧬 Skipping license validation — you already unlocked the boss room.")
	}

	sqlite.InitDB()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("🧱 Config load failed. Even IKEA gives better instructions: %v", err)
	}
	log.Printf("📦 Loaded Config: %+v", cfg)

	handlers.InitWithConfig(cfg)
	handlers.InitMetricsHandler(cfg)
	handlers.InitEmailService()
	handlers.InitIntegrationHandler(cfg)
	passwordreset.InjectConfig(cfg)

	metrics.RegisterTokenMetrics()

	log.Println("🔧 Initializing routes...")
	router := routes.SetupRoutes()
	handler := middleware.CorsMiddleware(middleware.LoggerMiddleware(router))
	log.Println("✅ Routes wired and ready for judgment day.")

	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS13,
		},
	}

	log.Printf("🚀 PeithoSecure Lite server armed and dangerous on port %s...", cfg.Port)

	if err := server.ListenAndServeTLS(
		os.Getenv("TLS_CERT_FILE"),
		os.Getenv("TLS_KEY_FILE"),
	); err != nil {
		log.Fatalf("💣 TLS Server Error — pulled the wrong wire: %v", err)
	}
}

func waitForLicense(path string) {
	maxAttempts := 10
	for i := 1; i <= maxAttempts; i++ {
		if _, err := os.Stat(path); err == nil {
			log.Printf("✅ Found unlock.lic on attempt %d/%d", i, maxAttempts)
			return
		}
		log.Printf("⏳ Waiting for license file at %s (attempt %d/%d)", path, i, maxAttempts)
		time.Sleep(2 * time.Second)
	}

	printAsciiBanner()
	log.Println("💀 Still no unlock.lic. This isn’t a security system—it’s a digital piñata.")
	log.Fatalln("🧨 Aborting launch. Your server is now officially ghosted. 👻")
}

func splitParts(data string) []string {
	return strings.SplitN(data, "||", 2)
}

func truncateMiddle(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max/2] + " … " + s[len(s)-max/2:]
}

func printAsciiBanner() {
	fmt.Println()
	fmt.Println("╔══════════════════════════════════════════╗")
	fmt.Println("║        🛡️  PeithoSecure Lite™ v1.0        ║")
	fmt.Println("║    Post-Quantum Identity Management      ║")
	fmt.Println("║     Fork it. Clone it. Cry anyway.       ║")
	fmt.Println("╚══════════════════════════════════════════╝")
	fmt.Println()
}
