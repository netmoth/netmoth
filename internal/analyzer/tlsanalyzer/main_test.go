package tlsanalyzer

import (
	"bytes"
	"testing"
	"time"

	"github.com/netmoth/netmoth/internal/connection"
)

func TestAnalyzeTLS(t *testing.T) {
	// Create test data for TLS Client Hello
	// This is a minimal TLS Client Hello packet
	tlsData := []byte{
		// TLS Record Header
		0x16,       // Handshake
		0x03, 0x01, // TLS 1.0
		0x00, 0x00, // Length (will be set later)

		// Handshake Header
		0x01,             // Client Hello
		0x00, 0x00, 0x00, // Length (will be set later)

		// Client Version
		0x03, 0x03, // TLS 1.2

		// Random (32 bytes)
		0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
		0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10,
		0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18,
		0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f, 0x20,

		// Session ID
		0x00, // Length

		// Cipher Suites
		0x00, 0x04, // Length (2 cipher suites)
		0x00, 0x2f, // TLS_RSA_WITH_AES_128_CBC_SHA
		0x00, 0x35, // TLS_RSA_WITH_AES_256_CBC_SHA

		// Compression Methods
		0x01, // Length
		0x00, // No compression

		// Extensions
		0x00, 0x00, // Length (no extensions for simplicity)
	}

	// Set correct lengths
	handshakeLength := len(tlsData) - 9    // 9 = 5 (TLS record header) + 4 (handshake header)
	tlsRecordLength := handshakeLength + 4 // handshake header + handshake data

	tlsData[3] = byte(tlsRecordLength >> 8)
	tlsData[4] = byte(tlsRecordLength & 0xff)
	tlsData[6] = byte(handshakeLength >> 16)
	tlsData[7] = byte((handshakeLength >> 8) & 0xff)
	tlsData[8] = byte(handshakeLength & 0xff)

	// Create connection
	conn := &connection.Connection{
		Timestamp:       time.Now(),
		SourceIP:        "192.168.1.100",
		DestinationIP:   "192.168.1.1",
		SourcePort:      12345,
		DestinationPort: 443,
		TransportType:   "tcp",
		Payload:         bytes.NewBuffer(tlsData),
		Analyzers:       make(map[string]any),
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

	if !result.IsClientHello {
		t.Error("Expected IsClientHello to be true")
	}

	if result.HandshakeType != "Client Hello" {
		t.Errorf("Expected HandshakeType to be 'Client Hello', got '%s'", result.HandshakeType)
	}

	if result.Version != "TLS 1.2" {
		t.Errorf("Expected Version to be 'TLS 1.2', got '%s'", result.Version)
	}

	if len(result.SupportedCiphers) != 2 {
		t.Errorf("Expected 2 supported ciphers, got %d", len(result.SupportedCiphers))
	}

	expectedCiphers := []string{
		"TLS_RSA_WITH_AES_128_CBC_SHA",
		"TLS_RSA_WITH_AES_256_CBC_SHA",
	}

	for i, expected := range expectedCiphers {
		if i >= len(result.SupportedCiphers) {
			t.Errorf("Missing cipher at index %d", i)
			continue
		}
		if result.SupportedCiphers[i] != expected {
			t.Errorf("Expected cipher '%s' at index %d, got '%s'", expected, i, result.SupportedCiphers[i])
		}
	}

	if result.PayloadLength != len(tlsData) {
		t.Errorf("Expected PayloadLength to be %d, got %d", len(tlsData), result.PayloadLength)
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
		Analyzers:       make(map[string]any),
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
		Analyzers:       make(map[string]any),
	}

	_, err := Analyze(conn)
	if err == nil {
		t.Error("Expected error for nil payload")
	}
}
