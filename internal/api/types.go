package api

import (
	"time"

	"github.com/netmoth/netmoth/internal/connection"
	"github.com/netmoth/netmoth/internal/signature"
)

// AgentData represents data sent by the agent
// Shared between agent client and server handlers to avoid contract drift.
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

