package web

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/netmoth/netmoth/internal/connection"
	"github.com/netmoth/netmoth/internal/signature"
)

// AgentData represents data sent by the agent
type AgentData struct {
	AgentID     string                   `json:"agent_id"`
	Hostname    string                   `json:"hostname"`
	Interface   string                   `json:"interface"`
	Timestamp   time.Time                `json:"timestamp"`
	Connections []*connection.Connection `json:"connections"`
	Signatures  []signature.Detect       `json:"signatures"`
	Stats       AgentStats               `json:"stats"`
}

// AgentStats represents agent statistics
type AgentStats struct {
	PacketsReceived  uint64 `json:"packets_received"`
	PacketsDropped   uint64 `json:"packets_dropped"`
	PacketsProcessed uint64 `json:"packets_processed"`
	ConnectionsFound uint64 `json:"connections_found"`
}

// AgentRegistration represents agent registration
type AgentRegistration struct {
	AgentID   string `json:"agent_id"`
	Hostname  string `json:"hostname"`
	Interface string `json:"interface"`
	Version   string `json:"version"`
}

// AgentResponse represents server response to agent
type AgentResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

// agentDataHandler processes data from agents
func agentDataHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var agentData AgentData
	if err := json.NewDecoder(r.Body).Decode(&agentData); err != nil {
		log.Printf("Error decoding agent data: %v", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Data validation
	if agentData.AgentID == "" {
		http.Error(w, "Agent ID is required", http.StatusBadRequest)
		return
	}

	// Here will be logic for saving data to database
	log.Printf("Received data from agent %s: %d connections, %d signatures",
		agentData.AgentID, len(agentData.Connections), len(agentData.Signatures))

	// Send response to agent
	response := AgentResponse{
		Success: true,
		Message: "Data received successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// agentRegistrationHandler processes agent registration
func agentRegistrationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var registration AgentRegistration
	if err := json.NewDecoder(r.Body).Decode(&registration); err != nil {
		log.Printf("Error decoding agent registration: %v", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Data validation
	if registration.AgentID == "" || registration.Hostname == "" {
		http.Error(w, "Agent ID and hostname are required", http.StatusBadRequest)
		return
	}

	log.Printf("Agent registration: %s (%s) on interface %s, version %s",
		registration.AgentID, registration.Hostname, registration.Interface, registration.Version)

	// Here will be logic for saving agent information to database

	response := AgentResponse{
		Success: true,
		Message: "Agent registered successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// agentHealthHandler processes agent health checks
func agentHealthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	agentID := r.URL.Query().Get("agent_id")
	if agentID == "" {
		http.Error(w, "Agent ID is required", http.StatusBadRequest)
		return
	}

	// Here will be logic for checking agent status

	response := AgentResponse{
		Success: true,
		Message: fmt.Sprintf("Agent %s is healthy", agentID),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
