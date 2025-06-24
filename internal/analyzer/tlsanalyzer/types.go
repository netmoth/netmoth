package tlsanalyzer

import (
	"time"
)

// TLS represents the result of TLS traffic analysis
type TLS struct {
	Version          string            `json:"version,omitempty"`
	SNI              string            `json:"sni,omitempty"`
	SupportedCiphers []string          `json:"supported_ciphers,omitempty"`
	SelectedCipher   string            `json:"selected_cipher,omitempty"`
	Certificate      *Certificate      `json:"certificate,omitempty"`
	HandshakeType    string            `json:"handshake_type,omitempty"`
	IsClientHello    bool              `json:"is_client_hello"`
	IsServerHello    bool              `json:"is_server_hello"`
	IsCertificate    bool              `json:"is_certificate"`
	IsAlert          bool              `json:"is_alert"`
	AlertLevel       string            `json:"alert_level,omitempty"`
	AlertDescription string            `json:"alert_description,omitempty"`
	Extensions       map[string]string `json:"extensions,omitempty"`
	PayloadLength    int               `json:"payload_length"`
	ApplicationData  *ApplicationData  `json:"application_data,omitempty"`
}

// Certificate represents information about a certificate
type Certificate struct {
	Subject      string    `json:"subject,omitempty"`
	Issuer       string    `json:"issuer,omitempty"`
	ValidFrom    time.Time `json:"valid_from,omitempty"`
	ValidUntil   time.Time `json:"valid_until,omitempty"`
	DNSNames     []string  `json:"dns_names,omitempty"`
	CommonName   string    `json:"common_name,omitempty"`
	SerialNumber string    `json:"serial_number,omitempty"`
}

// ApplicationData represents information about encrypted application data
type ApplicationData struct {
	Length              int    `json:"length"`
	IsEncrypted         bool   `json:"is_encrypted"`
	DetectedContentType string `json:"detected_content_type,omitempty"`
	DetectedProtocol    string `json:"detected_protocol,omitempty"`
	HasCompression      bool   `json:"has_compression"`
	EstimatedSize       int    `json:"estimated_size,omitempty"`
}

// Key returns the key for storing in analyzers
func (t *TLS) Key() string {
	return "tls"
}
