package utils

import (
	"math/rand"
	"time"
)

// GenerateErrorHaiku returns a random error-themed haiku
func GenerateErrorHaiku() string {
	haikus := []string{
		"Server lost in mist,\nChaos brews in silent night,\nRestart, brave admin.",
		"Whispers in cables,\nPackets wander, dreams decay,\nRetry in the dawn.",
		"Keys forgotten, cold,\nLocks unguarded by old code,\nHope flickers anew.",
	}

	rand.Seed(time.Now().UnixNano())
	return haikus[rand.Intn(len(haikus))]
}
