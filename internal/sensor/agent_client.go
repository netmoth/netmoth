package sensor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/netmoth/netmoth/internal/connection"
	"github.com/netmoth/netmoth/internal/signature"
	"github.com/netmoth/netmoth/internal/version"
)

// AgentClient represents a client for sending data to the central server
type AgentClient struct {
	serverURL string
	agentID   string
	hostname  string
	client    *http.Client
}

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

// NewAgentClient creates a new agent client
func NewAgentClient(serverURL, agentID, interfaceName string) *AgentClient {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	return &AgentClient{
		serverURL: serverURL,
		agentID:   agentID,
		hostname:  hostname,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Register registers the agent on the central server
func (ac *AgentClient) Register(interfaceName string) error {
	registration := AgentRegistration{
		AgentID:   ac.agentID,
		Hostname:  ac.hostname,
		Interface: interfaceName,
		Version:   version.Version(),
	}

	jsonData, err := json.Marshal(registration)
	if err != nil {
		return fmt.Errorf("failed to marshal registration data: %v", err)
	}

	resp, err := ac.client.Post(
		fmt.Sprintf("%s/api/agent/register", ac.serverURL),
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return fmt.Errorf("failed to register agent: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("registration failed with status: %d", resp.StatusCode)
	}

	var response AgentResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return fmt.Errorf("failed to decode response: %v", err)
	}

	if !response.Success {
		return fmt.Errorf("registration failed: %s", response.Error)
	}

	log.Printf("Agent %s registered successfully", ac.agentID)
	return nil
}

// SendData sends audit data to the central server
func (ac *AgentClient) SendData(connections []*connection.Connection, signatures []signature.Detect, stats AgentStats, interfaceName string) error {
	data := AgentData{
		AgentID:     ac.agentID,
		Hostname:    ac.hostname,
		Interface:   interfaceName,
		Timestamp:   time.Now(),
		Connections: connections,
		Signatures:  signatures,
		Stats:       stats,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal agent data: %v", err)
	}

	resp, err := ac.client.Post(
		fmt.Sprintf("%s/api/agent/data", ac.serverURL),
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return fmt.Errorf("failed to send data: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("data sending failed with status: %d", resp.StatusCode)
	}

	var response AgentResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return fmt.Errorf("failed to decode response: %v", err)
	}

	if !response.Success {
		return fmt.Errorf("data sending failed: %s", response.Error)
	}

	log.Printf("Data sent successfully to server: %d connections, %d signatures", len(connections), len(signatures))
	return nil
}

// SendHealth sends agent health check
func (ac *AgentClient) SendHealth() error {
	resp, err := ac.client.Get(
		fmt.Sprintf("%s/api/agent/health?agent_id=%s", ac.serverURL, ac.agentID),
	)
	if err != nil {
		return fmt.Errorf("failed to send health check: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("health check failed with status: %d", resp.StatusCode)
	}

	var response AgentResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return fmt.Errorf("failed to decode health response: %v", err)
	}

	if !response.Success {
		return fmt.Errorf("health check failed: %s", response.Error)
	}

	log.Printf("Health check successful: %s", response.Message)
	return nil
}
