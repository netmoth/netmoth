package sensor

import (
	"compress/gzip"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/netmoth/netmoth/internal/api"
	"github.com/netmoth/netmoth/internal/connection"
	"github.com/netmoth/netmoth/internal/signature"
)

func readerFromRequest(r *http.Request) io.Reader {
	var rd io.Reader = r.Body
	if r.Header.Get("Content-Encoding") == "gzip" {
		gz, err := gzip.NewReader(r.Body)
		if err == nil {
			return gz
		}
	}
	return rd
}

func TestAgentClient_Register_SendData_Health(t *testing.T) {
	token := "tok"
	mux := http.NewServeMux()

	mux.HandleFunc("/api/agent/register", func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "Bearer "+token {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		dec := json.NewDecoder(readerFromRequest(r))
		var reg api.AgentRegistration
		if err := dec.Decode(&reg); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		_ = json.NewEncoder(w).Encode(api.AgentResponse{Success: true, Message: "ok"})
	})
	mux.HandleFunc("/api/agent/data", func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "Bearer "+token {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		dec := json.NewDecoder(readerFromRequest(r))
		var data api.AgentData
		if err := dec.Decode(&data); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		_ = json.NewEncoder(w).Encode(api.AgentResponse{Success: true, Message: "ok"})
	})
	mux.HandleFunc("/api/agent/health", func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "Bearer "+token {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		_ = json.NewEncoder(w).Encode(api.AgentResponse{Success: true, Message: "ok"})
	})

	ts := httptest.NewServer(mux)
	defer ts.Close()

	ac := NewAgentClient(ts.URL, "a1", "eth0", token)
	if err := ac.Register("eth0"); err != nil {
		t.Fatalf("register: %v", err)
	}

	conns := []*connection.Connection{}
	sigs := []signature.Detect{}
	stats := AgentStats{PacketsReceived: 1}
	if err := ac.SendData(conns, sigs, stats, "eth0"); err != nil {
		t.Fatalf("send data: %v", err)
	}

	if err := ac.SendHealth(); err != nil {
		t.Fatalf("health: %v", err)
	}
}

func TestAgentClient_SendHealth_Unauthorized(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/agent/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	})
	ts := httptest.NewServer(mux)
	defer ts.Close()
	ac := NewAgentClient(ts.URL, "a1", "eth0", "tok")
	start := time.Now()
	err := ac.SendHealth()
	if err == nil {
		t.Fatalf("expected error")
	}
	if time.Since(start) > time.Second*2 {
		t.Fatalf("health should not block too long on unauthorized")
	}
}
