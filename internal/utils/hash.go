package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

// HashString returns SHA-256 hash of input string
func HashString(input string) string {
	hasher := sha256.New()
	hasher.Write([]byte(input))
	return hex.EncodeToString(hasher.Sum(nil))
}
