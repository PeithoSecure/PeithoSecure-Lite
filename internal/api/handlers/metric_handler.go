package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/peithosecure/peitho-backend/internal/config"
	"github.com/peithosecure/peitho-backend/internal/metrics"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var cfg *config.Config

// InitMetricsHandler initializes the config for metrics
func InitMetricsHandler(configLoaded *config.Config) {
	cfg = configLoaded
}

// MetricsHandler godoc
// @Summary Prometheus metrics endpoint
// @Description Returns raw Prometheus metrics for external monitoring
// @Tags metrics
// @Produce plain
// @Success 200 {string} string "Prometheus-formatted metrics"
// @Router /api/v1/metrics [get]
func MetricsHandler(w http.ResponseWriter, r *http.Request) {
	promhttp.Handler().ServeHTTP(w, r)
}

// AdminMetricsHandler godoc
// @Summary Admin-only Prometheus metrics
// @Description Returns full Prometheus metrics (Basic Auth protected)
// @Tags metrics
// @Produce plain
// @Success 200 {string} string "Prometheus-formatted metrics"
// @Failure 401 {string} string "Unauthorized"
// @Router /api/v1/admin-metrics [get]
func AdminMetricsHandler(w http.ResponseWriter, r *http.Request) {
	username, password, ok := r.BasicAuth()
	if !ok || username != cfg.AdminMetricsUsername || password != cfg.AdminMetricsPassword {
		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	promhttp.Handler().ServeHTTP(w, r)
}

// TokenMetricsHandler godoc
// @Summary Token metrics (usage stats)
// @Description Returns current issued, refreshed, revoked, and active token counts
// @Tags metrics
// @Produce json
// @Success 200 {object} TokenMetricsResponse
// @Router /api/v1/metrics/tokens [get]
func TokenMetricsHandler(w http.ResponseWriter, r *http.Request) {
	stats := TokenMetricsResponse{
		Issued:    metrics.IssuedTokens.Load(),
		Refreshed: metrics.RefreshedTokens.Load(),
		Revoked:   metrics.RevokedTokens.Load(),
		Active:    metrics.ActiveTokens.Load(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

// TokenMetricsResponse defines the JSON structure for token statistics
type TokenMetricsResponse struct {
	Issued    int64 `json:"issued" example:"142"`
	Refreshed int64 `json:"refreshed" example:"87"`
	Revoked   int64 `json:"revoked" example:"12"`
	Active    int64 `json:"active" example:"103"`
}
