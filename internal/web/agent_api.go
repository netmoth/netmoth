package web

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/netmoth/netmoth/internal/api"
	"github.com/netmoth/netmoth/internal/config"
)

// readJSON reads JSON with size limit and gzip support
func readJSON(w http.ResponseWriter, r *http.Request, v any, maxBytes int64) error {
	var rc io.ReadCloser = r.Body
	if r.Header.Get("Content-Encoding") == "gzip" {
		gz, err := gzip.NewReader(r.Body)
		if err != nil {
			return err
		}
		rc = gz
	}
	defer rc.Close()
	limited := http.MaxBytesReader(w, rc, maxBytes)
	dec := json.NewDecoder(limited)
	dec.DisallowUnknownFields()
	return dec.Decode(v)
}

// simple auth: expects header Authorization: Bearer <token>
func authorize(r *http.Request, expectedToken string) bool {
	auth := r.Header.Get("Authorization")
	if expectedToken == "" {
		return false
	}
	if !strings.HasPrefix(auth, "Bearer ") {
		return false
	}
	return strings.TrimPrefix(auth, "Bearer ") == expectedToken
}

func makeAgentDataHandler(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		if !authorize(r, cfg.AgentToken) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		var agentData api.AgentData
		if err := readJSON(w, r, &agentData, 10<<20); err != nil { // 10MB cap
			log.Printf("Error decoding agent data: %v", err)
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if agentData.AgentID == "" {
			http.Error(w, "Agent ID is required", http.StatusBadRequest)
			return
		}

		log.Printf("Received data from agent %s: %d connections, %d signatures",
			agentData.AgentID, len(agentData.Connections), len(agentData.Signatures))

		response := api.AgentResponse{
			Success: true,
			Message: "Data received successfully",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func makeAgentRegistrationHandler(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		if !authorize(r, cfg.AgentToken) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		var registration api.AgentRegistration
		if err := readJSON(w, r, &registration, 1<<20); err != nil { // 1MB cap
			log.Printf("Error decoding agent registration: %v", err)
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if registration.AgentID == "" || registration.Hostname == "" {
			http.Error(w, "Agent ID and hostname are required", http.StatusBadRequest)
			return
		}

		log.Printf("Agent registration: %s (%s) on interface %s, version %s",
			registration.AgentID, registration.Hostname, registration.Interface, registration.Version)

		response := api.AgentResponse{
			Success: true,
			Message: "Agent registered successfully",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func makeAgentHealthHandler(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		agentID := r.URL.Query().Get("agent_id")
		if agentID == "" {
			http.Error(w, "Agent ID is required", http.StatusBadRequest)
			return
		}

		response := api.AgentResponse{
			Success: true,
			Message: fmt.Sprintf("Agent %s is healthy", agentID),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
