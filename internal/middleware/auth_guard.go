package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	corestub "github.com/peithosecure/peitho-backend/internal/corestub"
)

// Context key for passing claims
type contextKey string

const UserContextKey contextKey = "user"

var (
	jwksURL = "http://keycloak:8080/realms/peitho/protocol/openid-connect/certs"
	issuer  = "http://keycloak:8080/realms/peitho"
)

// AuthGuard validates Bearer JWT properly
func AuthGuard(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			corestub.RespondWithTraceError(w, "missing_auth_header", http.StatusUnauthorized)
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			corestub.RespondWithTraceError(w, "invalid_auth_format", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]
		claims, err := validateJWT(tokenString)
		if err != nil {
			code := http.StatusUnauthorized
			if strings.Contains(err.Error(), "no key found") || strings.Contains(err.Error(), "issuer") {
				code = http.StatusForbidden
			}
			corestub.RespondWithTraceError(w, "invalid_token", code)
			return
		}

		// Inject claims into request context
		ctx := context.WithValue(r.Context(), UserContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func validateJWT(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Enforce RSA signing method
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// Extract key ID from header
		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, errors.New("no kid found")
		}

		// Fetch and parse JWKS
		jwks, err := fetchJWKS()
		if err != nil {
			return nil, err
		}

		key := findKey(jwks, kid)
		if key == nil {
			return nil, fmt.Errorf("no key found for kid %s", kid)
		}

		return keyToPublicKey(*key)
	})
	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	// Enforce issuer match
	if iss, ok := claims["iss"].(string); !ok || iss != issuer {
		return nil, errors.New("invalid issuer")
	}

	return claims, nil
}

// ExtractUsernameFromContext gets the username from JWT claims
func ExtractUsernameFromContext(ctx context.Context) (string, error) {
	claims, ok := ctx.Value(UserContextKey).(jwt.MapClaims)
	if !ok {
		return "", errors.New("no user claims in context")
	}
	username, ok := claims["preferred_username"].(string)
	if !ok || username == "" {
		return "", errors.New("username not found in claims")
	}
	return username, nil
}
