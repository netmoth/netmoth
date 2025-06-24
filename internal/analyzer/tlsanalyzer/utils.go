package tlsanalyzer

import (
	"crypto/x509"
	"encoding/binary"
	"fmt"
)

// TLS versions
const (
	TLSVersion10 = 0x0301
	TLSVersion11 = 0x0302
	TLSVersion12 = 0x0303
	TLSVersion13 = 0x0304
)

// TLS record types
const (
	TLSRecordTypeChangeCipherSpec = 20
	TLSRecordTypeAlert            = 21
	TLSRecordTypeHandshake        = 22
	TLSRecordTypeApplicationData  = 23
)

// TLS handshake message types
const (
	TLSHandshakeTypeClientHello        = 1
	TLSHandshakeTypeServerHello        = 2
	TLSHandshakeTypeCertificate        = 11
	TLSHandshakeTypeServerKeyExchange  = 12
	TLSHandshakeTypeCertificateRequest = 13
	TLSHandshakeTypeServerHelloDone    = 14
	TLSHandshakeTypeCertificateVerify  = 15
	TLSHandshakeTypeClientKeyExchange  = 16
	TLSHandshakeTypeFinished           = 20
)

// TLS extension types
const (
	TLSExtensionServerName           = 0
	TLSExtensionSupportedGroups      = 10
	TLSExtensionECPointFormats       = 11
	TLSExtensionSignatureAlgorithms  = 13
	TLSExtensionALPN                 = 16
	TLSExtensionExtendedMasterSecret = 23
	TLSExtensionSessionTicket        = 35
	TLSExtensionRenegotiationInfo    = 65281
)

// getTLSVersion returns a string representation of the TLS version
func getTLSVersion(version uint16) string {
	switch version {
	case TLSVersion10:
		return "TLS 1.0"
	case TLSVersion11:
		return "TLS 1.1"
	case TLSVersion12:
		return "TLS 1.2"
	case TLSVersion13:
		return "TLS 1.3"
	default:
		return fmt.Sprintf("Unknown (0x%04x)", version)
	}
}

// getHandshakeType returns a string representation of the handshake type
func getHandshakeType(handshakeType uint8) string {
	switch handshakeType {
	case TLSHandshakeTypeClientHello:
		return "Client Hello"
	case TLSHandshakeTypeServerHello:
		return "Server Hello"
	case TLSHandshakeTypeCertificate:
		return "Certificate"
	case TLSHandshakeTypeServerKeyExchange:
		return "Server Key Exchange"
	case TLSHandshakeTypeCertificateRequest:
		return "Certificate Request"
	case TLSHandshakeTypeServerHelloDone:
		return "Server Hello Done"
	case TLSHandshakeTypeCertificateVerify:
		return "Certificate Verify"
	case TLSHandshakeTypeClientKeyExchange:
		return "Client Key Exchange"
	case TLSHandshakeTypeFinished:
		return "Finished"
	default:
		return fmt.Sprintf("Unknown (%d)", handshakeType)
	}
}

// getCipherSuiteName returns the name of the cipher by its ID
func getCipherSuiteName(cipherSuite uint16) string {
	cipherSuites := map[uint16]string{
		0x0000: "TLS_NULL_WITH_NULL_NULL",
		0x0001: "TLS_RSA_WITH_NULL_MD5",
		0x0002: "TLS_RSA_WITH_NULL_SHA",
		0x0004: "TLS_RSA_WITH_RC4_128_MD5",
		0x0005: "TLS_RSA_WITH_RC4_128_SHA",
		0x000A: "TLS_RSA_WITH_3DES_EDE_CBC_SHA",
		0x002F: "TLS_RSA_WITH_AES_128_CBC_SHA",
		0x0035: "TLS_RSA_WITH_AES_256_CBC_SHA",
		0x003C: "TLS_RSA_WITH_AES_128_CBC_SHA256",
		0x003D: "TLS_RSA_WITH_AES_256_CBC_SHA256",
		0x009C: "TLS_RSA_WITH_AES_128_GCM_SHA256",
		0x009D: "TLS_RSA_WITH_AES_256_GCM_SHA384",
		0xC02F: "TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256",
		0xC030: "TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384",
		0xCCA8: "TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256",
		0xCCA9: "TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256",
		0x1301: "TLS_AES_128_GCM_SHA256",
		0x1302: "TLS_AES_256_GCM_SHA384",
		0x1303: "TLS_CHACHA20_POLY1305_SHA256",
	}

	if name, exists := cipherSuites[cipherSuite]; exists {
		return name
	}
	return fmt.Sprintf("Unknown (0x%04x)", cipherSuite)
}

// parseCertificate parses a certificate and returns a Certificate struct
func parseCertificate(certData []byte) (*Certificate, error) {
	cert, err := x509.ParseCertificate(certData)
	if err != nil {
		return nil, err
	}

	return &Certificate{
		Subject:      cert.Subject.String(),
		Issuer:       cert.Issuer.String(),
		ValidFrom:    cert.NotBefore,
		ValidUntil:   cert.NotAfter,
		DNSNames:     cert.DNSNames,
		CommonName:   cert.Subject.CommonName,
		SerialNumber: cert.SerialNumber.String(),
	}, nil
}

// getExtensionName returns the name of the extension by its ID
func getExtensionName(extensionType uint16) string {
	extensions := map[uint16]string{
		TLSExtensionServerName:           "server_name",
		TLSExtensionSupportedGroups:      "supported_groups",
		TLSExtensionECPointFormats:       "ec_point_formats",
		TLSExtensionSignatureAlgorithms:  "signature_algorithms",
		TLSExtensionALPN:                 "application_layer_protocol_negotiation",
		TLSExtensionExtendedMasterSecret: "extended_master_secret",
		TLSExtensionSessionTicket:        "session_ticket",
		TLSExtensionRenegotiationInfo:    "renegotiation_info",
	}

	if name, exists := extensions[extensionType]; exists {
		return name
	}
	return fmt.Sprintf("unknown_%d", extensionType)
}

// parseSNI extracts Server Name Indication from extensions
func parseSNI(data []byte) (string, error) {
	if len(data) < 2 {
		return "", fmt.Errorf("insufficient data for SNI")
	}

	// Skip extension type (2 bytes) and length (2 bytes)
	offset := 4

	if len(data) < offset+2 {
		return "", fmt.Errorf("insufficient data for SNI list length")
	}

	sniListLength := binary.BigEndian.Uint16(data[offset:])
	offset += 2

	if len(data) < offset+int(sniListLength) {
		return "", fmt.Errorf("insufficient data for SNI list")
	}

	sniData := data[offset : offset+int(sniListLength)]

	if len(sniData) < 3 {
		return "", fmt.Errorf("insufficient data for SNI entry")
	}

	// Skip name type (1 byte) and name length (2 bytes)
	nameType := sniData[0]
	if nameType != 0 { // 0 = hostname
		return "", fmt.Errorf("unsupported SNI name type: %d", nameType)
	}

	nameLength := binary.BigEndian.Uint16(sniData[1:3])
	if len(sniData) < 3+int(nameLength) {
		return "", fmt.Errorf("insufficient data for SNI hostname")
	}

	hostname := string(sniData[3 : 3+nameLength])
	return hostname, nil
}
