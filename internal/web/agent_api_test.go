package web

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/netmoth/netmoth/internal/api"
	"github.com/netmoth/netmoth/internal/config"
)

func makeGzipBody(t *testing.T, v any) *bytes.Buffer {
	t.Helper()
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	enc := json.NewEncoder(gz)
	if err := enc.Encode(v); err != nil {
		t.Fatalf("encode: %v", err)
	}
	_ = gz.Close()
	return &buf
}

func TestAgentRegistration_Unauthorized(t *testing.T) {
	cfg := &config.Config{AgentToken: "secret"}
	req := httptest.NewRequest(http.MethodPost, "/api/agent/register", bytes.NewBufferString(`{}`))
	rec := httptest.NewRecorder()
	makeAgentRegistrationHandler(cfg)(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rec.Code)
	}
}

func TestAgentRegistration_Success(t *testing.T) {
	cfg := &config.Config{AgentToken: "secret"}
	body := api.AgentRegistration{AgentID: "a1", Hostname: "h", Interface: "eth0", Version: "v"}
	req := httptest.NewRequest(http.MethodPost, "/api/agent/register", makeGzipBody(t, body))
	req.Header.Set("Authorization", "Bearer secret")
	req.Header.Set("Content-Encoding", "gzip")
	rec := httptest.NewRecorder()
	makeAgentRegistrationHandler(cfg)(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestAgentData_InvalidJSON(t *testing.T) {
	cfg := &config.Config{AgentToken: "secret"}
	req := httptest.NewRequest(http.MethodPost, "/api/agent/data", bytes.NewBufferString("{"))
	req.Header.Set("Authorization", "Bearer secret")
	rec := httptest.NewRecorder()
	makeAgentDataHandler(cfg)(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestAgentData_Success(t *testing.T) {
	cfg := &config.Config{AgentToken: "secret"}
	body := api.AgentData{AgentID: "a1"}
	req := httptest.NewRequest(http.MethodPost, "/api/agent/data", makeGzipBody(t, body))
	req.Header.Set("Authorization", "Bearer secret")
	req.Header.Set("Content-Encoding", "gzip")
	rec := httptest.NewRecorder()
	makeAgentDataHandler(cfg)(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}
