package moodreactor

import (
	"math/rand"
	"net/http"
	"time"
)

// GetMood returns a random fake mood state
func GetMood() string {
	moods := []string{"happy", "angry", "sassy", "paranoid", "unstable"}
	rand.Seed(time.Now().UnixNano())
	return moods[rand.Intn(len(moods))]
}

// MoodAPIHandler provides mood via API (debug/testing only)
func MoodAPIHandler(w http.ResponseWriter, r *http.Request) {
	mood := GetMood()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"mood":"` + mood + `"}`))
}
