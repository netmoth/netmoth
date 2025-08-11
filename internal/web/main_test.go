package web

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/netmoth/netmoth/internal/config"
)

func TestVersionHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/version", nil)
	rec := httptest.NewRecorder()
	versionHandler(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var body map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("json: %v", err)
	}
	if body["version"] == "" {
		t.Fatalf("empty version")
	}
}

func TestCORS_AllOriginsOrExact(t *testing.T) {
	// allow specific origin
	h := corsMiddleware([]string{"http://a"}, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodOptions, "/", nil)
	req.Header.Set("Origin", "http://a")
	h.ServeHTTP(rec, req)
	if rec.Header().Get("Access-Control-Allow-Origin") != "http://a" {
		t.Fatalf("expected allow exact origin")
	}

	// allow all when list is empty
	h2 := corsMiddleware(nil, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))
	rec2 := httptest.NewRecorder()
	req2 := httptest.NewRequest(http.MethodOptions, "/", nil)
	req2.Header.Set("Origin", "http://any")
	h2.ServeHTTP(rec2, req2)
	if rec2.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Fatalf("expected allow all for empty list")
	}
}

func TestMakeHandlers_Unauthorized(t *testing.T) {
	cfg := &config.Config{AgentToken: "tok"}
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/agent/data", nil)
	makeAgentDataHandler(cfg)(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rec.Code)
	}
}

