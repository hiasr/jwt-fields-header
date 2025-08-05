// Package jwtfieldsheader provides a Traefik middleware that extracts JWT claims and creates headers
package jwtfieldsheader

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// Config the plugin configuration.
type Config struct {
	HeaderName string   `json:"headerName,omitempty"`
	JwtClaims  []string `json:"jwtClaims,omitempty"`
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{
		HeaderName: "X-JWT-Claims",
		JwtClaims:  []string{},
	}
}

// JwtFieldsHeader a JWT fields header plugin.
type JwtFieldsHeader struct {
	next       http.Handler
	name       string
	headerName string
	jwtClaims  []string
}

// New created a new JwtFieldsHeader plugin.
func New(_ context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	return &JwtFieldsHeader{
		next:       next,
		name:       name,
		headerName: config.HeaderName,
		jwtClaims:  config.JwtClaims,
	}, nil
}

func (j *JwtFieldsHeader) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	// Skip if no claims configured
	if len(j.jwtClaims) == 0 {
		j.next.ServeHTTP(rw, req)
		return
	}

	// Extract Bearer token from Authorization header
	authHeader := req.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		j.next.ServeHTTP(rw, req)
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == "" {
		j.next.ServeHTTP(rw, req)
		return
	}

	// Parse JWT token (without signature verification)
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		j.next.ServeHTTP(rw, req)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		j.next.ServeHTTP(rw, req)
		return
	}

	// Extract claim values
	var claimValues []string
	for _, claimName := range j.jwtClaims {
		if value, exists := claims[claimName]; exists {
			if strValue, ok := value.(string); ok {
				claimValues = append(claimValues, strValue)
			} else {
				// Convert non-string values to JSON string
				if jsonBytes, err := json.Marshal(value); err == nil {
					claimValues = append(claimValues, string(jsonBytes))
				}
			}
		}
	}

	// Set header with concatenated claim values
	if len(claimValues) > 0 {
		headerValue := strings.Join(claimValues, "-")
		req.Header.Set(j.headerName, headerValue)
	}

	j.next.ServeHTTP(rw, req)
}
