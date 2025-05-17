package config

import (
	"os"
)

type Config struct {
	Port                  string
	LicenseToken          string
	SQLitePath            string
	KeycloakIssuerURL     string
	KeycloakURL           string
	KeycloakInternalURL   string
	KeycloakClientID      string
	KeycloakClientSecret  string
	KeycloakAdmin         string
	KeycloakAdminPassword string
	AdminMetricsUsername  string
	AdminMetricsPassword  string
}

func LoadConfig() (*Config, error) {
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080" // Default port
	}

	return &Config{
		Port:                  port,
		LicenseToken:          os.Getenv("PEITHO_LICENSE_TOKEN"),
		SQLitePath:            os.Getenv("PEITHO_SQLITE_PATH"),
		KeycloakIssuerURL:     os.Getenv("KEYCLOAK_ISSUER_URL"),
		KeycloakURL:           os.Getenv("KEYCLOAK_URL"),
		KeycloakInternalURL:   os.Getenv("KEYCLOAK_INTERNAL_URL"),
		KeycloakClientID:      os.Getenv("KEYCLOAK_CLIENT_ID"),
		KeycloakClientSecret:  os.Getenv("KEYCLOAK_CLIENT_SECRET"),
		KeycloakAdmin:         os.Getenv("KEYCLOAK_ADMIN"),
		KeycloakAdminPassword: os.Getenv("KEYCLOAK_ADMIN_PASSWORD"),
		AdminMetricsUsername:  os.Getenv("ADMIN_METRICS_USERNAME"),
		AdminMetricsPassword:  os.Getenv("ADMIN_METRICS_PASSWORD"),
	}, nil
}
