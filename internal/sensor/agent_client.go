package sensor

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/netmoth/netmoth/internal/api"
	"github.com/netmoth/netmoth/internal/connection"
	"github.com/netmoth/netmoth/internal/signature"
	"github.com/netmoth/netmoth/internal/version"
)

// AgentClient represents a client for sending data to the central server
type AgentClient struct {
	serverURL string
	agentID   string
	hostname  string
	token     string
	client    *http.Client
}

// AgentStats alias to shared type for backward compatibility
type AgentStats = api.AgentStats

// NewAgentClient creates a new agent client
func NewAgentClient(serverURL, agentID, interfaceName, token string) *AgentClient {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	return &AgentClient{
		serverURL: serverURL,
		agentID:   agentID,
		hostname:  hostname,
		token:     token,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// doJSON is a helper to POST JSON with gzip and basic retry
func (ac *AgentClient) doJSON(path string, payload any) (*http.Response, error) {
	var body bytes.Buffer
	gz := gzip.NewWriter(&body)
	enc := json.NewEncoder(gz)
	if err := enc.Encode(payload); err != nil {
		_ = gz.Close()
		return nil, fmt.Errorf("failed to encode payload: %v", err)
	}
	_ = gz.Close()

	url := fmt.Sprintf("%s%s", ac.serverURL, path)

	var resp *http.Response
	var err error
	backoff := 500 * time.Millisecond
	for attempt := 0; attempt < 5; attempt++ {
		req, rerr := http.NewRequest(http.MethodPost, url, &body)
		if rerr != nil {
			return nil, rerr
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Content-Encoding", "gzip")
		if ac.token != "" {
			req.Header.Set("Authorization", "Bearer "+ac.token)
		}

		resp, err = ac.client.Do(req)
		if err == nil && resp.StatusCode >= 200 && resp.StatusCode < 300 {
			return resp, nil
		}
		if resp != nil {
			io.Copy(io.Discard, resp.Body)
			resp.Body.Close()
		}
		time.Sleep(backoff)
		backoff *= 2
	}
	if err == nil && resp != nil {
		return resp, fmt.Errorf("request failed with status: %d", resp.StatusCode)
	}
	return nil, fmt.Errorf("request failed: %v", err)
}

// Register registers the agent on the central server
func (ac *AgentClient) Register(interfaceName string) error {
	registration := api.AgentRegistration{
		AgentID:   ac.agentID,
		Hostname:  ac.hostname,
		Interface: interfaceName,
		Version:   version.Version(),
	}

	resp, err := ac.doJSON("/api/agent/register", registration)
	if err != nil {
		return fmt.Errorf("failed to register agent: %v", err)
	}
	defer resp.Body.Close()

	var response api.AgentResponse
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
	data := api.AgentData{
		AgentID:     ac.agentID,
		Hostname:    ac.hostname,
		Interface:   interfaceName,
		Timestamp:   time.Now(),
		Connections: connections,
		Signatures:  signatures,
		Stats:       stats,
	}

	resp, err := ac.doJSON("/api/agent/data", data)
	if err != nil {
		return fmt.Errorf("failed to send data: %v", err)
	}
	defer resp.Body.Close()

	var response api.AgentResponse
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
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/agent/health?agent_id=%s", ac.serverURL, ac.agentID), nil)
	if ac.token != "" {
		req.Header.Set("Authorization", "Bearer "+ac.token)
	}
	resp, err := ac.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send health check: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("health check failed with status: %d", resp.StatusCode)
	}

	var response api.AgentResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return fmt.Errorf("failed to decode health response: %v", err)
	}
	if !response.Success {
		return fmt.Errorf("health check failed: %s", response.Error)
	}

	log.Printf("Health check successful: %s", response.Message)
	return nil
}
