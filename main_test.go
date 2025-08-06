package jwtfieldsheader_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hiasr/jwtfieldsheader" //nolint:depguard
)

func TestJwtFieldsHeader(t *testing.T) {
	cfg := jwtfieldsheader.CreateConfig()
	cfg.HeaderName = "X-User-Info"
	cfg.JwtClaims = []string{"sub", "role"}

	ctx := context.Background()
	next := http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {})

	handler, err := jwtfieldsheader.New(ctx, next, cfg, "jwtfieldsheader-plugin")
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	// Create a test JWT token with claims (without signature)
	testToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwicm9sZSI6ImFkbWluIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
	req.Header.Set("Authorization", "Bearer "+testToken)

	if err != nil {
		t.Fatal(err)
	}

	handler.ServeHTTP(recorder, req)

	assertHeader(t, req, "X-User-Info", "1234567890-admin")
}

func TestJwtFieldsHeaderNoToken(t *testing.T) {
	cfg := jwtfieldsheader.CreateConfig()
	cfg.HeaderName = "X-User-Info"
	cfg.JwtClaims = []string{"sub", "role"}

	ctx := context.Background()
	next := http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {})

	handler, err := jwtfieldsheader.New(ctx, next, cfg, "jwtfieldsheader-plugin")
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
	// No Authorization header

	if err != nil {
		t.Fatal(err)
	}

	handler.ServeHTTP(recorder, req)

	// Should not set any header
	if req.Header.Get("X-User-Info") != "" {
		t.Errorf("Expected no header, but got: %s", req.Header.Get("X-User-Info"))
	}
}

func TestJwtFieldsHeaderNoClaims(t *testing.T) {
	cfg := jwtfieldsheader.CreateConfig()
	cfg.HeaderName = "X-User-Info"
	cfg.JwtClaims = []string{} // No claims configured

	ctx := context.Background()
	next := http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {})

	handler, err := jwtfieldsheader.New(ctx, next, cfg, "jwtfieldsheader-plugin")
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	testToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwicm9sZSI6ImFkbWluIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
	req.Header.Set("Authorization", "Bearer "+testToken)

	if err != nil {
		t.Fatal(err)
	}

	handler.ServeHTTP(recorder, req)

	// Should not set any header when no claims configured
	if req.Header.Get("X-User-Info") != "" {
		t.Errorf("Expected no header, but got: %s", req.Header.Get("X-User-Info"))
	}
}

func assertHeader(t *testing.T, req *http.Request, key, expected string) {
	t.Helper()

	if req.Header.Get(key) != expected {
		t.Errorf("invalid header value: %s", req.Header.Get(key))
	}
}
