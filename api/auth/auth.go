package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/auth0/go-jwt-middleware"
	"github.com/form3tech-oss/jwt-go"
)

// JSONWebKeys outlines the decode JWT token
type JSONWebKeys struct {
	Kty string   `json:"kty"`
	Kid string   `json:"kid"`
	Use string   `json:"use"`
	N   string   `json:"n"`
	E   string   `json:"e"`
	X5c []string `json:"x5c"`
}

// NewMiddleware creates a middleware checking JWT tokens using provided Auth0 endpoints
func NewMiddleware(aud, iss string) *jwtmiddleware.JWTMiddleware {
	return jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			if !token.Claims.(jwt.MapClaims).VerifyAudience(aud, false) {
				return token, fmt.Errorf("invalid token audience")
			}
			if !token.Claims.(jwt.MapClaims).VerifyIssuer(iss, false) {
				return token, fmt.Errorf("invalid token issuer")
			}

			cert, err := getCert(token, iss)
			if err != nil {
				return nil, fmt.Errorf("load cert: %s", err)
			}

			result, err := jwt.ParseRSAPublicKeyFromPEM(cert)
			if err != nil {
				return nil, fmt.Errorf("parse cert: %s", err)
			}
			return result, nil
		},
		SigningMethod: jwt.SigningMethodRS256,
	})
}

func getCert(token *jwt.Token, issuer string) ([]byte, error) {
	resp, err := http.Get(strings.TrimSuffix(issuer, "/") + "/.well-known/jwks.json")
	if err != nil {
		return nil, fmt.Errorf("load JWK: %s", err)
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}

	var payload struct {
		Keys []JSONWebKeys `json:"keys"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf("parse JWK JSON: %s", err)
	}

	for _, key := range payload.Keys {
		if token.Header["kid"] == key.Kid {
			return []byte("-----BEGIN CERTIFICATE-----\n" + key.X5c[0] + "\n-----END CERTIFICATE-----"), nil
		}
	}

	return nil, fmt.Errorf("unable to find matching key")
}
