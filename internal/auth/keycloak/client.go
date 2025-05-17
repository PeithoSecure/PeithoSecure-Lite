package keycloak

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/peithosecure/peitho-backend/internal/config"
)

type OIDCClient struct {
	Provider *oidc.Provider
	Verifier *oidc.IDTokenVerifier
}

func NewOIDCClient(cfg *config.Config) (*OIDCClient, error) {
	ctx := context.Background()
	provider, err := oidc.NewProvider(ctx, cfg.KeycloakIssuerURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create OIDC provider: %w", err)
	}
	return &OIDCClient{
		Provider: provider,
		Verifier: provider.Verifier(&oidc.Config{ClientID: cfg.KeycloakClientID}),
	}, nil
}

// ResetPassword sets or resets a user's password in Keycloak
func ResetPassword(cfg *config.Config, username, newPassword string) error {
	userID, err := GetUserIDByUsername(cfg, username)
	if err != nil {
		return fmt.Errorf("cannot find user ID: %v", err)
	}

	payload := map[string]interface{}{
		"type":      "password",
		"value":     newPassword,
		"temporary": false,
	}

	return sendAdminRequest(cfg, "PUT", "/admin/realms/peitho/users/"+userID+"/reset-password", payload)
}

// Optional: alias for clarity in Option B
func SetPassword(cfg *config.Config, username, newPassword string) error {
	return ResetPassword(cfg, username, newPassword)
}

// GetUserIDByUsername fetches the Keycloak user ID from the username
func GetUserIDByUsername(cfg *config.Config, username string) (string, error) {
	var users []struct {
		ID string `json:"id"`
	}
	err := sendAdminRequest(cfg, "GET", "/admin/realms/peitho/users?username="+username, nil, &users)
	if err != nil {
		return "", err
	}
	if len(users) == 0 {
		return "", fmt.Errorf("no user found for username %s", username)
	}
	return users[0].ID, nil
}

// KeycloakClient represents a Keycloak application client
type KeycloakClient struct {
	ClientID     string `json:"clientId"`
	Name         string `json:"name"`
	Protocol     string `json:"protocol"`
	Enabled      bool   `json:"enabled"`
	PublicClient bool   `json:"publicClient"`
}

// GetRealmClients fetches all registered clients in the realm
func GetRealmClients(cfg *config.Config) ([]KeycloakClient, error) {
	var clients []KeycloakClient
	err := sendAdminRequest(cfg, "GET", "/admin/realms/peitho/clients", nil, &clients)
	if err != nil {
		return nil, fmt.Errorf("failed to get clients: %w", err)
	}
	return clients, nil
}

// sendAdminRequest sends a Keycloak admin API request and optionally parses response
func sendAdminRequest(cfg *config.Config, method, path string, body interface{}, result ...interface{}) error {
	token, err := fetchAdminToken(cfg)
	if err != nil {
		return err
	}

	var bodyBytes []byte
	if body != nil {
		bodyBytes, err = json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request: %w", err)
		}
	}

	req, err := http.NewRequest(method, cfg.KeycloakInternalURL+path, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		msg, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("keycloak admin API error: %s\nResponse: %s", resp.Status, string(msg))
	}

	if len(result) > 0 {
		return json.NewDecoder(resp.Body).Decode(result[0])
	}
	return nil
}

// fetchAdminToken gets a short-lived access token using admin credentials
func fetchAdminToken(cfg *config.Config) (string, error) {
	data := "grant_type=password" +
		"&client_id=admin-cli" +
		"&username=" + cfg.KeycloakAdmin +
		"&password=" + cfg.KeycloakAdminPassword

	req, err := http.NewRequest("POST", cfg.KeycloakInternalURL+"/realms/master/protocol/openid-connect/token", bytes.NewBufferString(data))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		msg, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("token fetch failed: %s\n%s", resp.Status, string(msg))
	}

	var res struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", err
	}
	return res.AccessToken, nil
}
