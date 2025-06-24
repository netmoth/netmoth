package tlsanalyzer

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/netmoth/netmoth/internal/connection"
)

// Analyze analyzes TLS traffic and returns structured information
func Analyze(c *connection.Connection) (*TLS, error) {
	if c.Payload == nil || c.Payload.Len() == 0 {
		return nil, fmt.Errorf("empty payload")
	}

	data := c.Payload.Bytes()
	result := &TLS{
		PayloadLength: len(data),
		Extensions:    make(map[string]string),
	}

	// Parse TLS records
	offset := 0
	for offset < len(data) {
		if offset+5 > len(data) {
			break // Not enough data for TLS record header
		}

		// Read TLS record header
		recordType := data[offset]
		length := int(binary.BigEndian.Uint16(data[offset+3 : offset+5]))

		if offset+5+length > len(data) {
			break // Not enough data for record
		}

		recordData := data[offset+5 : offset+5+length]

		switch recordType {
		case TLSRecordTypeHandshake:
			if err := parseHandshake(recordData, result); err != nil {
				// Log error, but continue analysis
				fmt.Printf("Error parsing handshake: %v\n", err)
			}
		case TLSRecordTypeAlert:
			if err := parseAlert(recordData, result); err != nil {
				fmt.Printf("Error parsing alert: %v\n", err)
			}
		case TLSRecordTypeApplicationData:
			if err := parseApplicationData(recordData, result); err != nil {
				fmt.Printf("Error parsing application data: %v\n", err)
			}
		}

		offset += 5 + length
	}

	return result, nil
}

// parseHandshake parses TLS handshake messages
func parseHandshake(data []byte, result *TLS) error {
	offset := 0
	for offset < len(data) {
		if offset+4 > len(data) {
			break
		}

		handshakeType := data[offset]
		handshakeLength := readUint24(data[offset+1:])

		if offset+4+int(handshakeLength) > len(data) {
			break
		}

		handshakeData := data[offset+4 : offset+4+int(handshakeLength)]

		switch handshakeType {
		case TLSHandshakeTypeClientHello:
			result.IsClientHello = true
			result.HandshakeType = getHandshakeType(handshakeType)
			if err := parseClientHello(handshakeData, result); err != nil {
				return fmt.Errorf("error parsing client hello: %w", err)
			}
		case TLSHandshakeTypeServerHello:
			result.IsServerHello = true
			result.HandshakeType = getHandshakeType(handshakeType)
			if err := parseServerHello(handshakeData, result); err != nil {
				return fmt.Errorf("error parsing server hello: %w", err)
			}
		case TLSHandshakeTypeCertificate:
			result.IsCertificate = true
			result.HandshakeType = getHandshakeType(handshakeType)
			if err := parseCertificateMessage(handshakeData, result); err != nil {
				return fmt.Errorf("error parsing certificate: %w", err)
			}
		}

		offset += 4 + int(handshakeLength)
	}

	return nil
}

// parseClientHello parses Client Hello message
func parseClientHello(data []byte, result *TLS) error {
	if len(data) < 34 {
		return fmt.Errorf("insufficient data for client hello")
	}

	// Skip client version (2 bytes)
	clientVersion := binary.BigEndian.Uint16(data[0:])
	result.Version = getTLSVersion(clientVersion)

	// Skip random (32 bytes)
	offset := 34

	// Skip session ID length
	if offset+1 > len(data) {
		return fmt.Errorf("insufficient data for session ID length")
	}
	sessionIDLength := data[offset]
	offset += 1 + int(sessionIDLength)

	// Read supported cipher suites
	if offset+2 > len(data) {
		return fmt.Errorf("insufficient data for cipher suites length")
	}
	cipherSuitesLength := binary.BigEndian.Uint16(data[offset:])
	offset += 2

	if offset+int(cipherSuitesLength) > len(data) {
		return fmt.Errorf("insufficient data for cipher suites")
	}

	// Parse supported ciphers
	for i := 0; i < int(cipherSuitesLength); i += 2 {
		if offset+i+2 > len(data) {
			break
		}
		cipherSuite := binary.BigEndian.Uint16(data[offset+i:])
		result.SupportedCiphers = append(result.SupportedCiphers, getCipherSuiteName(cipherSuite))
	}
	offset += int(cipherSuitesLength)

	// Skip compression methods
	if offset+1 > len(data) {
		return fmt.Errorf("insufficient data for compression methods length")
	}
	compressionMethodsLength := data[offset]
	offset += 1 + int(compressionMethodsLength)

	// Parse extensions
	if offset+2 <= len(data) {
		extensionsLength := binary.BigEndian.Uint16(data[offset:])
		offset += 2

		if offset+int(extensionsLength) <= len(data) {
			parseExtensions(data[offset:offset+int(extensionsLength)], result)
		}
	}

	return nil
}

// parseServerHello parses Server Hello message
func parseServerHello(data []byte, result *TLS) error {
	if len(data) < 34 {
		return fmt.Errorf("insufficient data for server hello")
	}

	// Skip server version (2 bytes)
	serverVersion := binary.BigEndian.Uint16(data[0:])
	result.Version = getTLSVersion(serverVersion)

	// Skip random (32 bytes)
	offset := 34

	// Skip session ID
	if offset+1 > len(data) {
		return fmt.Errorf("insufficient data for session ID length")
	}
	sessionIDLength := data[offset]
	offset += 1 + int(sessionIDLength)

	// Read selected cipher
	if offset+2 > len(data) {
		return fmt.Errorf("insufficient data for selected cipher suite")
	}
	selectedCipher := binary.BigEndian.Uint16(data[offset:])
	result.SelectedCipher = getCipherSuiteName(selectedCipher)
	offset += 2

	// Skip compression method
	if offset+1 > len(data) {
		return fmt.Errorf("insufficient data for compression method")
	}
	offset += 1

	// Parse extensions
	if offset+2 <= len(data) {
		extensionsLength := binary.BigEndian.Uint16(data[offset:])
		offset += 2

		if offset+int(extensionsLength) <= len(data) {
			parseExtensions(data[offset:offset+int(extensionsLength)], result)
		}
	}

	return nil
}

// parseCertificateMessage parses Certificate message
func parseCertificateMessage(data []byte, result *TLS) error {
	if len(data) < 3 {
		return fmt.Errorf("insufficient data for certificate message")
	}

	certificatesLength := readUint24(data[0:])
	offset := 3

	if offset+int(certificatesLength) > len(data) {
		return fmt.Errorf("insufficient data for certificates")
	}

	// Parse the first certificate (usually the server certificate)
	for offset < 3+int(certificatesLength) {
		if offset+3 > len(data) {
			break
		}

		certLength := readUint24(data[offset:])
		offset += 3

		if offset+int(certLength) > len(data) {
			break
		}

		certData := data[offset : offset+int(certLength)]
		cert, err := parseCertificate(certData)
		if err == nil {
			result.Certificate = cert
			break // Only take the first certificate
		}

		offset += int(certLength)
	}

	return nil
}

// parseAlert parses TLS Alert messages
func parseAlert(data []byte, result *TLS) error {
	if len(data) < 2 {
		return fmt.Errorf("insufficient data for alert")
	}

	result.IsAlert = true

	alertLevel := data[0]
	alertDescription := data[1]

	switch alertLevel {
	case 1:
		result.AlertLevel = "Warning"
	case 2:
		result.AlertLevel = "Fatal"
	default:
		result.AlertLevel = fmt.Sprintf("Unknown (%d)", alertLevel)
	}

	// Basic alert descriptions
	alertDescriptions := map[uint8]string{
		0:   "Close Notify",
		10:  "Unexpected Message",
		20:  "Bad Record MAC",
		21:  "Decryption Failed",
		22:  "Record Overflow",
		30:  "Decompression Failure",
		40:  "Handshake Failure",
		41:  "No Certificate",
		42:  "Bad Certificate",
		43:  "Unsupported Certificate",
		44:  "Certificate Revoked",
		45:  "Certificate Expired",
		46:  "Certificate Unknown",
		47:  "Illegal Parameter",
		48:  "Unknown CA",
		49:  "Access Denied",
		50:  "Decode Error",
		51:  "Decrypt Error",
		60:  "Export Restriction",
		70:  "Protocol Version",
		71:  "Insufficient Security",
		80:  "Internal Error",
		86:  "Inappropriate Fallback",
		90:  "User Canceled",
		100: "No Renegotiation",
		110: "Unsupported Extension",
		111: "Certificate Unobtainable",
		112: "Unrecognized Name",
		113: "Bad Certificate Status Response",
		114: "Bad Certificate Hash Value",
		115: "Unknown PSK Identity",
		116: "Certificate Required",
		120: "No Application Protocol",
	}

	if desc, exists := alertDescriptions[alertDescription]; exists {
		result.AlertDescription = desc
	} else {
		result.AlertDescription = fmt.Sprintf("Unknown (%d)", alertDescription)
	}

	return nil
}

// parseExtensions parses TLS extensions
func parseExtensions(data []byte, result *TLS) error {
	offset := 0
	for offset < len(data) {
		if offset+4 > len(data) {
			break
		}

		extensionType := binary.BigEndian.Uint16(data[offset:])
		extensionLength := binary.BigEndian.Uint16(data[offset+2:])
		offset += 4

		if offset+int(extensionLength) > len(data) {
			break
		}

		extensionData := data[offset : offset+int(extensionLength)]

		switch extensionType {
		case TLSExtensionServerName:
			if sni, err := parseSNI(extensionData); err == nil {
				result.SNI = sni
			}
		case TLSExtensionALPN:
			if alpn, err := parseALPN(extensionData); err == nil {
				result.Extensions["alpn"] = alpn
			}
		}

		offset += int(extensionLength)
	}

	return nil
}

// parseALPN parses Application Layer Protocol Negotiation extension
func parseALPN(data []byte) (string, error) {
	if len(data) < 2 {
		return "", fmt.Errorf("insufficient data for ALPN")
	}

	protocolsLength := binary.BigEndian.Uint16(data[0:])
	offset := 2

	if offset+int(protocolsLength) > len(data) {
		return "", fmt.Errorf("insufficient data for ALPN protocols")
	}

	// Parse the first protocol
	for offset < 2+int(protocolsLength) {
		if offset+1 > len(data) {
			break
		}

		protocolLength := data[offset]
		offset += 1

		if offset+int(protocolLength) > len(data) {
			break
		}

		protocol := string(data[offset : offset+int(protocolLength)])
		return protocol, nil // Return the first protocol
	}

	return "", fmt.Errorf("no ALPN protocols found")
}

// parseApplicationData parses TLS Application Data (encrypted content)
func parseApplicationData(data []byte, result *TLS) error {
	// Application Data is encrypted, but we can analyze its structure
	// and attempt to detect protocols like HTTP/2, HTTP/1.1, etc.

	if len(data) == 0 {
		return nil
	}

	// Try to detect HTTP/2 by looking for HTTP/2 preface
	if len(data) >= 24 && string(data[:24]) == "PRI * HTTP/2.0\r\n\r\nSM\r\n\r\n" {
		result.Extensions["detected_protocol"] = "HTTP/2"
		return nil
	}

	// Try to detect HTTP/1.1 by looking for common HTTP methods
	httpMethods := []string{"GET ", "POST ", "PUT ", "DELETE ", "HEAD ", "OPTIONS ", "PATCH "}
	for _, method := range httpMethods {
		if len(data) >= len(method) && string(data[:len(method)]) == method {
			result.Extensions["detected_protocol"] = "HTTP/1.1"
			break
		}
	}

	// Analyze data patterns
	result.ApplicationData = &ApplicationData{
		Length:      len(data),
		IsEncrypted: true,
	}

	// Try to detect content type based on patterns
	if detectJSONPattern(data) {
		result.ApplicationData.DetectedContentType = "application/json"
	} else if detectXMLPattern(data) {
		result.ApplicationData.DetectedContentType = "application/xml"
	} else if detectImagePattern(data) {
		result.ApplicationData.DetectedContentType = "image"
	} else if detectTextPattern(data) {
		result.ApplicationData.DetectedContentType = "text"
	}

	return nil
}

// detectJSONPattern tries to detect JSON content
func detectJSONPattern(data []byte) bool {
	if len(data) < 2 {
		return false
	}

	// Look for JSON start/end patterns
	start := string(data[:2])
	return start == "{\"" || start == "[{" || start == "[\"" || start == "{\n" || start == "{\r"
}

// detectXMLPattern tries to detect XML content
func detectXMLPattern(data []byte) bool {
	if len(data) < 5 {
		return false
	}

	// Look for XML declaration or root element
	start := string(data[:5])
	return start == "<?xml" || start == "<root" || start == "<html" || start == "<soap"
}

// detectImagePattern tries to detect image content
func detectImagePattern(data []byte) bool {
	if len(data) < 4 {
		return false
	}

	// Common image file signatures
	signatures := [][]byte{
		{0xFF, 0xD8, 0xFF},       // JPEG
		{0x89, 0x50, 0x4E, 0x47}, // PNG
		{0x47, 0x49, 0x46},       // GIF
		{0x42, 0x4D},             // BMP
		{0x52, 0x49, 0x46, 0x46}, // WebP
	}

	for _, sig := range signatures {
		if len(data) >= len(sig) && bytes.Equal(data[:len(sig)], sig) {
			return true
		}
	}

	return false
}

// detectTextPattern tries to detect text content
func detectTextPattern(data []byte) bool {
	if len(data) == 0 {
		return false
	}

	// Check if data contains mostly printable ASCII characters
	printableCount := 0
	for _, b := range data {
		if b >= 32 && b <= 126 || b == 9 || b == 10 || b == 13 {
			printableCount++
		}
	}

	// If more than 80% are printable, consider it text
	return float64(printableCount)/float64(len(data)) > 0.8
}

// readUint24 is a helper function to read 3-byte values
func readUint24(data []byte) uint32 {
	return uint32(data[0])<<16 | uint32(data[1])<<8 | uint32(data[2])
}
