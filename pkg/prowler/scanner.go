package prowler

import (
	"log"
	"net/http"
	"os"
)

// ScanResult holds the scan summary
type ScanResult struct {
	DBExists            bool `json:"db_exists"`
	SMTPConfigured      bool `json:"smtp_configured"`
	LicenseTokenPresent bool `json:"license_token_present"`
	DockerDetected      bool `json:"docker_detected"`
}

// RunSecurityScan performs a limited lightweight security audit
func RunSecurityScan() ScanResult {
	log.Println("[Prowler] Starting limited security scan...")

	result := ScanResult{
		DBExists:            checkSQLiteExists(),
		SMTPConfigured:      checkSMTPEnv(),
		LicenseTokenPresent: checkLicenseToken(),
		DockerDetected:      detectDocker(),
	}

	// 🚨 Optional basic alerting
	if !result.DBExists {
		TriggerAlert("SQLite database file missing")
	}
	if !result.SMTPConfigured {
		TriggerAlert("SMTP settings not properly configured")
	}
	if !result.LicenseTokenPresent {
		TriggerAlert("License token missing in environment")
	}

	log.Println("[Prowler] Security scan completed.")

	return result
}

// --- Internal Check Helpers ---

func checkSQLiteExists() bool {
	_, err := os.Stat("/data/peitho_secure.db")
	return err == nil
}

func checkSMTPEnv() bool {
	return os.Getenv("SMTP_HOST") != "" &&
		os.Getenv("SMTP_PORT") != "" &&
		os.Getenv("SMTP_USERNAME") != "" &&
		os.Getenv("SMTP_PASSWORD") != ""
}

func checkLicenseToken() bool {
	return os.Getenv("PEITHO_LICENSE_TOKEN") != ""
}

func detectDocker() bool {
	// Check if running inside Docker by checking cgroup
	if _, err := os.Stat("/.dockerenv"); err == nil {
		return true
	}
	resp, err := http.Get("http://docker.for.mac.localhost")
	if err == nil && resp.StatusCode == 200 {
		return true
	}
	return false
}
