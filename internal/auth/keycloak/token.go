package keycloak

import (
	"context"
	"fmt"

	"github.com/peithosecure/peitho-backend/internal/config"
)

// TokenClaims represents parsed JWT claims
type TokenClaims struct {
	Sub   string `json:"sub"`
	Email string `json:"email"`
	Exp   int64  `json:"exp"`
}

// VerifyIDToken verifies an ID token string and returns parsed claims
func VerifyIDToken(cfg *config.Config, tokenString string) (*TokenClaims, error) {
	client, err := NewOIDCClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize OIDC client: %w", err)
	}

	ctx := context.Background()
	idToken, err := client.Verifier.Verify(ctx, tokenString)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	var claims TokenClaims
	if err := idToken.Claims(&claims); err != nil {
		return nil, fmt.Errorf("failed to parse claims: %w", err)
	}

	return &claims, nil
}
