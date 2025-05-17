package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	TokensIssued = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "peitho_tokens_issued_total",
		Help: "Total tokens issued",
	})

	TokensRefreshed = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "peitho_tokens_refreshed_total",
		Help: "Total tokens refreshed",
	})

	TokensRevoked = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "peitho_tokens_revoked_total",
		Help: "Total tokens revoked (logout)",
	})
)

func RegisterTokenMetrics() {
	prometheus.MustRegister(TokensIssued)
	prometheus.MustRegister(TokensRefreshed)
	prometheus.MustRegister(TokensRevoked)
}
