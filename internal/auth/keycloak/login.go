package keycloak

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/peithosecure/peitho-backend/internal/config"
)

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

var client = &http.Client{}

// RegisterNewUserWithEmail creates a new Keycloak user and optionally sets their password
func RegisterNewUserWithEmail(cfg *config.Config, username, password, email string) error {
	token, err := getAdminToken(cfg)
	if err != nil {
		return fmt.Errorf("failed to get admin token: %w", err)
	}

	user := map[string]interface{}{
		"username":        username,
		"enabled":         true,
		"email":           email,
		"emailVerified":   true,
		"requiredActions": []string{},
	}

	bodyBytes, _ := json.Marshal(user)
	req, err := http.NewRequest("POST", cfg.KeycloakInternalURL+"/admin/realms/peitho/users", bytes.NewBuffer(bodyBytes))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to create user: %s", string(body))
	}

	userID, err := findUserIDByUsername(cfg, token, username)
	if err != nil {
		return fmt.Errorf("failed to find created user: %w", err)
	}

	// ✅ Skip setting password if empty (for deferred flow)
	if password != "" {
		if err := resetUserPassword(cfg, token, userID, password); err != nil {
			return fmt.Errorf("user created but failed to set password: %w", err)
		}
	} else {
		fmt.Printf("🔒 Skipping password setup for user %s — will be set later.\n", username)
	}

	return nil
}

// UserExists checks if a user already exists in Keycloak
func UserExists(cfg *config.Config, username string) (bool, error) {
	token, err := getAdminToken(cfg)
	if err != nil {
		return false, err
	}

	_, err = findUserIDByUsername(cfg, token, username)
	if err != nil {
		if err.Error() == "user not found" {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// LoginWithPassword authenticates user
func LoginWithPassword(cfg *config.Config, username, password string) (*TokenResponse, error) {
	data := url.Values{}
	data.Set("grant_type", "password")
	data.Set("client_id", cfg.KeycloakClientID)
	data.Set("client_secret", cfg.KeycloakClientSecret)
	data.Set("username", username)
	data.Set("password", password)

	req, err := http.NewRequest("POST", cfg.KeycloakURL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create login request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send login request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("login failed: %s, body: %s", resp.Status, string(bodyBytes))
	}

	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, fmt.Errorf("failed to decode token response: %w", err)
	}

	return &tokenResp, nil
}

// RefreshWithToken refreshes access token
func RefreshWithToken(cfg *config.Config, refreshToken string) (map[string]interface{}, error) {
	form := url.Values{}
	form.Set("grant_type", "refresh_token")
	form.Set("refresh_token", refreshToken)
	form.Set("client_id", cfg.KeycloakClientID)
	form.Set("client_secret", cfg.KeycloakClientSecret)

	return sendTokenRequest(cfg, form)
}

// RevokeRefreshToken invalidates refresh token
func RevokeRefreshToken(cfg *config.Config, refreshToken string) error {
	form := url.Values{}
	form.Set("client_id", cfg.KeycloakClientID)
	form.Set("client_secret", cfg.KeycloakClientSecret)
	form.Set("refresh_token", refreshToken)

	logoutURL := cfg.KeycloakIssuerURL + "/protocol/openid-connect/logout"
	resp, err := client.PostForm(logoutURL, form)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to revoke token: %s", string(body))
	}

	return nil
}

// Helpers

func sendTokenRequest(cfg *config.Config, form url.Values) (map[string]interface{}, error) {
	req, err := http.NewRequest("POST", cfg.KeycloakURL, bytes.NewBufferString(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create token request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send token request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("token request failed: %s, body: %s", resp.Status, string(bodyBytes))
	}

	var respData map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return nil, fmt.Errorf("failed to decode token response: %w", err)
	}

	return respData, nil
}

func getAdminToken(cfg *config.Config) (string, error) {
	data := url.Values{}
	data.Set("grant_type", "password")
	data.Set("client_id", "admin-cli")
	data.Set("username", cfg.KeycloakAdmin)
	data.Set("password", cfg.KeycloakAdminPassword)

	req, err := http.NewRequest("POST", cfg.KeycloakInternalURL+"/realms/master/protocol/openid-connect/token", bytes.NewBufferString(data.Encode()))
	if err != nil {
		return "", fmt.Errorf("failed to create admin login request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send admin login request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("admin login failed: %s", string(body))
	}

	var tokenResp struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", fmt.Errorf("failed to decode admin token: %w", err)
	}

	return tokenResp.AccessToken, nil
}

func findUserIDByUsername(cfg *config.Config, token, username string) (string, error) {
	req, err := http.NewRequest("GET", cfg.KeycloakInternalURL+"/admin/realms/peitho/users?username="+url.QueryEscape(username), nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to lookup user: %s", string(body))
	}

	var users []struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&users); err != nil {
		return "", err
	}

	if len(users) == 0 {
		return "", errors.New("user not found")
	}

	return users[0].ID, nil
}

func resetUserPassword(cfg *config.Config, token, userID, newPassword string) error {
	payload := map[string]interface{}{
		"type":      "password",
		"value":     newPassword,
		"temporary": false,
	}
	bodyBytes, _ := json.Marshal(payload)

	req, err := http.NewRequest("PUT", cfg.KeycloakInternalURL+"/admin/realms/peitho/users/"+userID+"/reset-password", bytes.NewBuffer(bodyBytes))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to reset password: %s", string(body))
	}

	return nil
}

// DeleteUser removes a Keycloak user by username
func DeleteUser(cfg *config.Config, username string) error {
	token, err := getAdminToken(cfg)
	if err != nil {
		return fmt.Errorf("failed to get admin token: %w", err)
	}

	userID, err := findUserIDByUsername(cfg, token, username)
	if err != nil {
		return fmt.Errorf("failed to find user ID: %w", err)
	}

	req, err := http.NewRequest("DELETE", cfg.KeycloakInternalURL+"/admin/realms/peitho/users/"+userID, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to delete user: %s", string(body))
	}

	return nil
}
