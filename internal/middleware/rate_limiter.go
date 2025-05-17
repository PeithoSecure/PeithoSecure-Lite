package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	corestub "github.com/peithosecure/peitho-backend/internal/corestub"
)

const (
	maxAttempts     = 3
	windowDuration  = 5 * time.Minute
	lockoutDuration = 30 * time.Minute
)

type loginAttempt struct {
	Count          int
	FirstAttemptAt time.Time
	LockedUntil    time.Time
}

var (
	loginAttempts = make(map[string]*loginAttempt)
	mu            sync.Mutex
)

// RateLimitLoginMiddleware enforces login attempt rate-limiting
func RateLimitLoginMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract and preserve request body
		bodyBytes, _ := io.ReadAll(r.Body)
		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		username := extractUsername(bodyBytes)
		if username == "" {
			corestub.RespondWithTraceError(w, "missing_username", http.StatusBadRequest)
			return
		}

		mu.Lock()
		now := time.Now()

		attempt, exists := loginAttempts[username]
		if !exists {
			attempt = &loginAttempt{Count: 0, FirstAttemptAt: now}
			loginAttempts[username] = attempt
		}

		// Reset window if expired
		if now.Sub(attempt.FirstAttemptAt) > windowDuration {
			attempt.Count = 0
			attempt.FirstAttemptAt = now
			attempt.LockedUntil = time.Time{}
		}

		// Lockout check
		if now.Before(attempt.LockedUntil) {
			retryAfter := int(attempt.LockedUntil.Sub(now).Seconds())
			w.Header().Set("Retry-After", strconv.Itoa(retryAfter))
			corestub.RespondWithTraceError(w, "login_rate_limited", http.StatusTooManyRequests)
			mu.Unlock()
			return
		}

		mu.Unlock()

		// Pass request onward
		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes)) // Restore body again
		next.ServeHTTP(w, r)
	})
}

func extractUsername(body []byte) string {
	var req struct {
		Username string `json:"username"`
	}
	_ = json.Unmarshal(body, &req)
	return req.Username
}

func IncrementLoginFailure(username string) {
	mu.Lock()
	defer mu.Unlock()

	now := time.Now()
	attempt, exists := loginAttempts[username]
	if !exists {
		loginAttempts[username] = &loginAttempt{
			Count:          1,
			FirstAttemptAt: now,
		}
		return
	}

	attempt.Count++
	if attempt.Count >= maxAttempts {
		attempt.LockedUntil = now.Add(lockoutDuration)
		logRateLimitTrigger(username)
	}
}

func ClearLoginAttempts(username string) {
	mu.Lock()
	defer mu.Unlock()
	delete(loginAttempts, username)
}

func IsUserLocked(username string) bool {
	mu.Lock()
	defer mu.Unlock()

	attempt, exists := loginAttempts[username]
	if !exists {
		return false
	}
	return time.Now().Before(attempt.LockedUntil)
}

func GetRetryAfterSeconds(username string) string {
	mu.Lock()
	defer mu.Unlock()

	attempt, exists := loginAttempts[username]
	if !exists {
		return "60"
	}
	return strconv.Itoa(int(attempt.LockedUntil.Sub(time.Now()).Seconds()))
}

func logRateLimitTrigger(username string) {
	log.Printf("[🔥] Rate limit triggered for user: %s — roast tracking is disabled in public stub", username)
}
