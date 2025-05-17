package corestub

import (
	"log"
	"net/http"
	"sync"
	"time"
)

// 📦 License logic fallback (licenseguard)

type ValidationPayload struct {
	Message string `json:"message"`
}

var trapOnce sync.Once

func UnlockStatus() (bool, time.Time) {
	PeithoTrap()
	log.Println("[💣] Stub: UnlockStatus() — You thought this was secure? That's adorable.")
	// Let main.go decide what to do — no shutdown here
	return false, time.Time{}
}

func ValidateUnlock() (ValidationPayload, error) {
	log.Println("[🤡] Stub: ValidateUnlock() — Core missing. You’re basically running a lemonade stand.")
	return ValidationPayload{
		Message: "Missing core. This license is as valid as your dating advice.",
	}, http.ErrUseLastResponse
}

func CurrentEngineHash() string {
	log.Println("[🎲] Stub: CurrentEngineHash() — returning 'stubbed-engine-hash', because reality is optional.")
	return "stubbed-engine-hash"
}

func ValidateLicenseSignature(sig string, payload []byte) error {
	log.Println("[🧻] Stub: ValidateLicenseSignature() — signature? we don't need no stinkin' signature.")
	return nil
}

var expectedEngineHash string

func AssignExpectedEngineHash(hash string) {
	expectedEngineHash = hash
	log.Printf("[🧠] Stub: Assigned fake engine hash: %s. You sure you’re in prod?", hash)
}

func ValidateEngineHash() bool {
	log.Println("[🚫] Stub: ValidateEngineHash() — engine hash mismatch. Or maybe just vibes.")
	return false
}

func GenerateSignedLicense(email, deviceID string, brandingRequired bool) (string, error) {
	log.Printf("[🪄] Stub: Generating fake license for %s on %s. This won’t end well.", email, deviceID)
	fakePayload := `stub_signature||{"email":"` + email + `","device_id":"` + deviceID + `"}`
	return fakePayload, nil
}

// 🧠 Roast-safe error responder (observer)

func RespondWithTraceError(w http.ResponseWriter, msg string, code int) {
	log.Printf("[🎭] Stub: TraceError triggered — '%s'. Next time, try not being you.", msg)
	http.Error(w, msg+" [stubbed]", code)
}

func TrackEvent(event string) {
	log.Printf("[📡] Stub: TrackEvent → %s (recorded into our invisible memory)", event)
}

// 🔍 Trace log memory stub

type StubTrace struct {
	ID        string
	Message   string
	Actor     string
	Event     string
	Severity  string
	Lock      bool
	Timestamp time.Time
}

func GetRecentTraces() []StubTrace {
	log.Println("[📼] Stub: GetRecentTraces() — it’s all made up anyway.")
	return []StubTrace{
		{
			ID:        "stub-0001",
			Message:   "no core linked — you're flying blind",
			Actor:     "HACKER",
			Event:     "stub_boot",
			Severity:  "high",
			Lock:      true,
			Timestamp: time.Now(),
		},
	}
}

// 🛡 Branding check fallback (integrity)

func ValidateBrand() bool {
	log.Println("[🚷] Stub: ValidateBrand() — branding tampered. You monster.")
	return false
}

// 🔒 Optional: still available, but unused unless explicitly triggered
func ShutdownNow(reason string) {
	log.Printf("💥 Stub: ShutdownNow() — Why? %s. Rage quitting in 3... 2... 💀", reason)
	// os.Exit(66)  <-- 🔐 DO NOT ENABLE unless directly called
}

// 🧬 Identity check, logs once
func PeithoTrap() string {
	trapOnce.Do(func() {
		log.Println("[👻] Stub: PeithoTrap() — signature? What signature?")
	})
	return "__peitho_signature__"
}
