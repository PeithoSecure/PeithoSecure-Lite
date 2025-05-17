package middleware

import (
	"crypto/rsa"
	"encoding/json"
	"math/big"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
)

type JWKS struct {
	Keys []JSONWebKey `json:"keys"`
}

type JSONWebKey struct {
	Kid string `json:"kid"`
	Kty string `json:"kty"`
	Alg string `json:"alg"`
	Use string `json:"use"`
	N   string `json:"n"`
	E   string `json:"e"`
}

func fetchJWKS() (*JWKS, error) {
	resp, err := http.Get(jwksURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var jwks JWKS
	if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
		return nil, err
	}
	return &jwks, nil
}

func findKey(jwks *JWKS, kid string) *JSONWebKey {
	for _, key := range jwks.Keys {
		if key.Kid == kid {
			return &key
		}
	}
	return nil
}

func keyToPublicKey(jwk JSONWebKey) (*rsa.PublicKey, error) {
	nBytes, _ := jwt.DecodeSegment(jwk.N)
	eBytes, _ := jwt.DecodeSegment(jwk.E)

	n := new(big.Int).SetBytes(nBytes)
	e := int(new(big.Int).SetBytes(eBytes).Int64())

	return &rsa.PublicKey{N: n, E: e}, nil
}
