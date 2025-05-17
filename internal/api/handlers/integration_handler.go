package handlers

import (
	"net/http"

	"github.com/peithosecure/peitho-backend/internal/auth/keycloak"
	"github.com/peithosecure/peitho-backend/internal/config"
	corestub "github.com/peithosecure/peitho-backend/internal/corestub"
	"github.com/peithosecure/peitho-backend/internal/utils"
)

var integrationCfg *config.Config

// InitIntegrationHandler sets global config for integration queries
func InitIntegrationHandler(cfg *config.Config) {
	integrationCfg = cfg
}

// IntegrationClient represents a client application registered in Keycloak
type IntegrationClient struct {
	ClientID     string `json:"client_id" example:"peitho-dashboard"`
	Name         string `json:"name" example:"Peitho Dashboard"`
	Protocol     string `json:"protocol" example:"openid-connect"`
	Enabled      bool   `json:"enabled" example:"true"`
	PublicClient bool   `json:"public_client" example:"true"`
}

// GetAppIntegrations godoc
// @Summary List active app integrations
// @Description Returns a list of enabled Keycloak clients integrated with the platform
// @Tags Integrations
// @Produce json
// @Success 200 {array} IntegrationClient
// @Failure 401 {object} map[string]string "Unauthorized or token expired"
// @Router /api/v1/integrations [get]
func GetAppIntegrations(w http.ResponseWriter, r *http.Request) {
	clients, err := keycloak.GetRealmClients(integrationCfg)
	if err != nil {
		corestub.RespondWithTraceError(w, "integration_fetch_fail", http.StatusUnauthorized)
		return
	}

	var filtered []IntegrationClient
	for _, c := range clients {
		if c.Enabled {
			filtered = append(filtered, IntegrationClient{
				ClientID:     c.ClientID,
				Name:         c.Name,
				Protocol:     c.Protocol,
				Enabled:      c.Enabled,
				PublicClient: c.PublicClient,
			})
		}
	}

	utils.JSON(w, http.StatusOK, filtered)
}
