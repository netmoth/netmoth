package contentanalyzer

import (
	"bytes"
	"testing"
	"time"

	"github.com/netmoth/netmoth/internal/connection"
)

func TestAnalyzeJSONContent(t *testing.T) {
	// Create test JSON data
	jsonData := []byte(`{"name": "test", "value": 123, "active": true}`)

	// Create connection
	conn := &connection.Connection{
		Timestamp:       time.Now(),
		SourceIP:        "192.168.1.100",
		DestinationIP:   "192.168.1.1",
		SourcePort:      12345,
		DestinationPort: 80,
		TransportType:   "tcp",
		Payload:         bytes.NewBuffer(jsonData),
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

	if result.ContentType != "application/json" {
		t.Errorf("Expected ContentType to be 'application/json', got '%s'", result.ContentType)
	}

	if result.DataType != "json" {
		t.Errorf("Expected DataType to be 'json', got '%s'", result.DataType)
	}

	if !result.IsText {
		t.Error("Expected IsText to be true")
	}

	if result.StructuredData == nil {
		t.Error("Expected StructuredData to be not nil")
	}

	if result.PayloadLength != len(jsonData) {
		t.Errorf("Expected PayloadLength to be %d, got %d", len(jsonData), result.PayloadLength)
	}
}

func TestAnalyzeXMLContent(t *testing.T) {
	// Create test XML data
	xmlData := []byte(`<?xml version="1.0"?><root><item>test</item></root>`)

	// Create connection
	conn := &connection.Connection{
		Timestamp:       time.Now(),
		SourceIP:        "192.168.1.100",
		DestinationIP:   "192.168.1.1",
		SourcePort:      12345,
		DestinationPort: 80,
		TransportType:   "tcp",
		Payload:         bytes.NewBuffer(xmlData),
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

	if result.ContentType != "application/xml" {
		t.Errorf("Expected ContentType to be 'application/xml', got '%s'", result.ContentType)
	}

	if result.DataType != "xml" {
		t.Errorf("Expected DataType to be 'xml', got '%s'", result.DataType)
	}

	if !result.IsText {
		t.Error("Expected IsText to be true")
	}

	if len(result.XMLTags) == 0 {
		t.Error("Expected XMLTags to be not empty")
	}
}

func TestAnalyzeTextContent(t *testing.T) {
	// Create test text data
	textData := []byte("Hello, this is a test message with https://example.com and domain.com")

	// Create connection
	conn := &connection.Connection{
		Timestamp:       time.Now(),
		SourceIP:        "192.168.1.100",
		DestinationIP:   "192.168.1.1",
		SourcePort:      12345,
		DestinationPort: 80,
		TransportType:   "tcp",
		Payload:         bytes.NewBuffer(textData),
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

	if result.ContentType != "text/plain" {
		t.Errorf("Expected ContentType to be 'text/plain', got '%s'", result.ContentType)
	}

	if !result.IsText {
		t.Error("Expected IsText to be true")
	}

	if result.TextContent == "" {
		t.Error("Expected TextContent to be not empty")
	}

	if len(result.URLs) == 0 {
		t.Error("Expected URLs to be not empty")
	}

	if len(result.Domains) == 0 {
		t.Error("Expected Domains to be not empty")
	}
}

func TestAnalyzeBinaryContent(t *testing.T) {
	// Create test binary data (PNG signature)
	binaryData := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}

	// Create connection
	conn := &connection.Connection{
		Timestamp:       time.Now(),
		SourceIP:        "192.168.1.100",
		DestinationIP:   "192.168.1.1",
		SourcePort:      12345,
		DestinationPort: 80,
		TransportType:   "tcp",
		Payload:         bytes.NewBuffer(binaryData),
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

	if result.ContentType != "application/octet-stream" {
		t.Errorf("Expected ContentType to be 'application/octet-stream', got '%s'", result.ContentType)
	}

	if !result.IsBinary {
		t.Error("Expected IsBinary to be true")
	}

	if result.FileType != "PNG" {
		t.Errorf("Expected FileType to be 'PNG', got '%s'", result.FileType)
	}
}

func TestAnalyzeEmptyPayload(t *testing.T) {
	conn := &connection.Connection{
		Timestamp:       time.Now(),
		SourceIP:        "192.168.1.100",
		DestinationIP:   "192.168.1.1",
		SourcePort:      12345,
		DestinationPort: 80,
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
		DestinationPort: 80,
		TransportType:   "tcp",
		Payload:         nil,
		Analyzers:       make(map[string]any),
	}

	_, err := Analyze(conn)
	if err == nil {
		t.Error("Expected error for nil payload")
	}
}
