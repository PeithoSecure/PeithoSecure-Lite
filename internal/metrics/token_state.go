package metrics

import (
	"sync/atomic"
)

var (
	IssuedTokens    atomic.Int64
	RefreshedTokens atomic.Int64
	RevokedTokens   atomic.Int64
	ActiveTokens    atomic.Int64 // optional, for session-aware tracking
)

func IncIssued() {
	IssuedTokens.Add(1)
}

func IncRefreshed() {
	RefreshedTokens.Add(1)
}

func IncRevoked() {
	RevokedTokens.Add(1)
}
