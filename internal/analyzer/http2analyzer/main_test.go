package http2analyzer

import (
	"bytes"
	"testing"
	"time"

	"github.com/netmoth/netmoth/internal/connection"
)

func TestAnalyzeHTTP2(t *testing.T) {
	// Create test data for HTTP/2
	http2Data := []byte{
		// HTTP/2 Preface
		0x50, 0x52, 0x49, 0x20, 0x2a, 0x20, 0x48, 0x54, 0x54, 0x50, 0x2f, 0x32, 0x2e, 0x30, 0x0d, 0x0a,
		0x0d, 0x0a, 0x53, 0x4d, 0x0d, 0x0a, 0x0d, 0x0a,

		// SETTINGS frame
		0x00, 0x00, 0x0c, // Length: 12
		0x04,                   // Type: SETTINGS
		0x00,                   // Flags: 0
		0x00, 0x00, 0x00, 0x00, // Stream ID: 0

		// Settings payload
		0x00, 0x01, 0x00, 0x00, 0x04, 0x00, // HEADER_TABLE_SIZE: 4096
		0x00, 0x02, 0x00, 0x00, 0x00, 0x01, // ENABLE_PUSH: 1

		// HEADERS frame
		0x00, 0x00, 0x25, // Length: 37
		0x01,                   // Type: HEADERS
		0x04,                   // Flags: END_HEADERS
		0x00, 0x00, 0x00, 0x01, // Stream ID: 1

		// Headers payload (simplified)
		0x00, 0x00, 0x00, 0x00, // Stream dependency: 0
		0x00, // Weight: 0
		// HPACK encoded headers would go here
		0x82, 0x84, 0x87, 0x41, 0x8a, 0xa0, 0xe4, 0x1d, 0x13, 0x9d, 0x09, 0xb8, 0xf0, 0x1e, 0x07, 0x42,
		0x9a, 0xa0, 0xe4, 0x1d, 0x13, 0x9d, 0x09, 0xb8, 0xf0, 0x1e, 0x07, 0x42, 0x9a, 0xa0, 0xe4, 0x1d,
		0x13, 0x9d, 0x09, 0xb8, 0xf0, 0x1e, 0x07, 0x42, 0x9a, 0xa0, 0xe4, 0x1d, 0x13, 0x9d, 0x09, 0xb8,
	}

	// Create connection
	conn := &connection.Connection{
		Timestamp:       time.Now(),
		SourceIP:        "192.168.1.100",
		DestinationIP:   "192.168.1.1",
		SourcePort:      12345,
		DestinationPort: 443,
		TransportType:   "tcp",
		Payload:         bytes.NewBuffer(http2Data),
		Analyzers:       make(map[string]interface{}),
	}

	// Analyze
	result, err := Analyze(conn)
	if err != nil {
		t.Fatalf("Analyze failed: %v", err)
	}

	// Check results
	if result == nil {
		t.Fatal("Result is nil")
	}

	if !result.HasHTTP2Preface {
		t.Error("Expected HasHTTP2Preface to be true")
	}

	if result.Protocol != "HTTP/2" {
		t.Errorf("Expected Protocol to be 'HTTP/2', got '%s'", result.Protocol)
	}

	if len(result.Frames) < 2 {
		t.Errorf("Expected at least 2 frames, got %d", len(result.Frames))
	}

	if len(result.SettingsFrames) < 1 {
		t.Error("Expected at least 1 SETTINGS frame")
	}

	if len(result.HeadersFrames) < 1 {
		t.Error("Expected at least 1 HEADERS frame")
	}

	if result.PayloadLength != len(http2Data) {
		t.Errorf("Expected PayloadLength to be %d, got %d", len(http2Data), result.PayloadLength)
	}
}

func TestAnalyzeHTTP2WithoutPreface(t *testing.T) {
	// Create test data for HTTP/2 without preface (just frames)
	http2Data := []byte{
		// SETTINGS frame
		0x00, 0x00, 0x0c, // Length: 12
		0x04,                   // Type: SETTINGS
		0x00,                   // Flags: 0
		0x00, 0x00, 0x00, 0x00, // Stream ID: 0

		// Settings payload
		0x00, 0x01, 0x00, 0x00, 0x04, 0x00, // HEADER_TABLE_SIZE: 4096
		0x00, 0x02, 0x00, 0x00, 0x00, 0x01, // ENABLE_PUSH: 1
	}

	// Create connection
	conn := &connection.Connection{
		Timestamp:       time.Now(),
		SourceIP:        "192.168.1.100",
		DestinationIP:   "192.168.1.1",
		SourcePort:      12345,
		DestinationPort: 443,
		TransportType:   "tcp",
		Payload:         bytes.NewBuffer(http2Data),
		Analyzers:       make(map[string]interface{}),
	}

	// Analyze
	result, err := Analyze(conn)
	if err != nil {
		t.Fatalf("Analyze failed: %v", err)
	}

	// Check results
	if result == nil {
		t.Fatal("Result is nil")
	}

	if result.HasHTTP2Preface {
		t.Error("Expected HasHTTP2Preface to be false")
	}

	if result.Protocol != "HTTP/2 (detected)" {
		t.Errorf("Expected Protocol to be 'HTTP/2 (detected)', got '%s'", result.Protocol)
	}

	if len(result.Frames) < 1 {
		t.Error("Expected at least 1 frame")
	}
}

func TestAnalyzeEmptyPayload(t *testing.T) {
	conn := &connection.Connection{
		Timestamp:       time.Now(),
		SourceIP:        "192.168.1.100",
		DestinationIP:   "192.168.1.1",
		SourcePort:      12345,
		DestinationPort: 443,
		TransportType:   "tcp",
		Payload:         bytes.NewBuffer([]byte{}),
		Analyzers:       make(map[string]interface{}),
	}

	_, err := Analyze(conn)
	if err == nil {
		t.Error("Expected error for empty payload")
	}
}

func TestAnalyzeNilPayload(t *testing.T) {
	conn := &connection.Connection{
		Timestamp:       time.Now(),
		SourceIP:        "192.168.1.100",
		DestinationIP:   "192.168.1.1",
		SourcePort:      12345,
		DestinationPort: 443,
		TransportType:   "tcp",
		Payload:         nil,
		Analyzers:       make(map[string]interface{}),
	}

	_, err := Analyze(conn)
	if err == nil {
		t.Error("Expected error for nil payload")
	}
}
