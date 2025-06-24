package http2analyzer

import (
	"fmt"

	"github.com/netmoth/netmoth/internal/connection"
)

// Analyze analyzes HTTP/2 traffic and returns structured information
func Analyze(c *connection.Connection) (*HTTP2, error) {
	if c.Payload == nil || c.Payload.Len() == 0 {
		return nil, fmt.Errorf("empty payload")
	}

	data := c.Payload.Bytes()
	result := &HTTP2{
		PayloadLength: len(data),
		Frames:        make([]Frame, 0),
	}

	// Check for HTTP/2 preface
	if len(data) >= 24 {
		preface := string(data[:24])
		if preface == "PRI * HTTP/2.0\r\n\r\nSM\r\n\r\n" {
			result.HasHTTP2Preface = true
			result.Protocol = "HTTP/2"

			// Parse frames after preface
			if err := parseFrames(data[24:], result); err != nil {
				return nil, fmt.Errorf("error parsing HTTP/2 frames: %w", err)
			}
		}
	}

	// Try to detect HTTP/2 by analyzing frame structure
	if !result.HasHTTP2Preface {
		if err := detectHTTP2Frames(data, result); err != nil {
			return nil, fmt.Errorf("error detecting HTTP/2 frames: %w", err)
		}
	}

	return result, nil
}

// parseFrames parses HTTP/2 frames from the data
func parseFrames(data []byte, result *HTTP2) error {
	offset := 0

	for offset < len(data) {
		if offset+9 > len(data) {
			break // Not enough data for frame header
		}

		// Parse frame header (9 bytes)
		length := int(data[offset])<<16 | int(data[offset+1])<<8 | int(data[offset+2])
		frameType := data[offset+3]
		flags := data[offset+4]
		streamID := uint32(data[offset+5])<<24 | uint32(data[offset+6])<<16 | uint32(data[offset+7])<<8 | uint32(data[offset+8])

		if offset+9+length > len(data) {
			break // Not enough data for frame payload
		}

		frameData := data[offset+9 : offset+9+length]

		frame := Frame{
			Length:   length,
			Type:     getFrameTypeName(frameType),
			Flags:    flags,
			StreamID: streamID,
			Data:     frameData,
		}

		// Parse specific frame types
		switch frameType {
		case 0x1: // HEADERS
			if err := parseHeadersFrame(frameData, &frame); err == nil {
				result.HeadersFrames = append(result.HeadersFrames, frame)
			}
		case 0x0: // DATA
			if err := parseDataFrame(frameData, &frame); err == nil {
				result.DataFrames = append(result.DataFrames, frame)
			}
		case 0x4: // SETTINGS
			if err := parseSettingsFrame(frameData, &frame); err == nil {
				result.SettingsFrames = append(result.SettingsFrames, frame)
			}
		case 0x5: // PUSH_PROMISE
			if err := parsePushPromiseFrame(frameData, &frame); err == nil {
				result.PushPromiseFrames = append(result.PushPromiseFrames, frame)
			}
		}

		result.Frames = append(result.Frames, frame)
		offset += 9 + length
	}

	return nil
}

// detectHTTP2Frames tries to detect HTTP/2 frames without preface
func detectHTTP2Frames(data []byte, result *HTTP2) error {
	offset := 0

	for offset < len(data) {
		if offset+9 > len(data) {
			break
		}

		// Try to parse as HTTP/2 frame
		length := int(data[offset])<<16 | int(data[offset+1])<<8 | int(data[offset+2])
		frameType := data[offset+3]

		// Validate frame type
		if frameType > 0x9 {
			break // Invalid frame type
		}

		// Check if length is reasonable
		if length > 16384 {
			break // Frame too large
		}

		if offset+9+length <= len(data) {
			// This looks like a valid HTTP/2 frame
			result.Protocol = "HTTP/2 (detected)"
			return parseFrames(data, result)
		}

		offset++
	}

	return nil
}

// parseHeadersFrame parses HTTP/2 HEADERS frame
func parseHeadersFrame(data []byte, frame *Frame) error {
	if len(data) < 4 {
		return fmt.Errorf("insufficient data for headers frame")
	}

	// Parse stream dependency and weight
	streamDependency := uint32(data[0])<<24 | uint32(data[1])<<16 | uint32(data[2])<<8 | uint32(data[3])
	frame.StreamDependency = streamDependency

	if len(data) > 4 {
		frame.Weight = data[4]
	}

	// Try to extract headers (HPACK encoded)
	frame.HeaderBlockFragment = data[5:]

	return nil
}

// parseDataFrame parses HTTP/2 DATA frame
func parseDataFrame(data []byte, frame *Frame) error {
	frame.DataPayload = data
	frame.DataLength = len(data)
	return nil
}

// parseSettingsFrame parses HTTP/2 SETTINGS frame
func parseSettingsFrame(data []byte, frame *Frame) error {
	if len(data)%6 != 0 {
		return fmt.Errorf("invalid settings frame length")
	}

	for i := 0; i < len(data); i += 6 {
		if i+6 <= len(data) {
			identifier := uint16(data[i])<<8 | uint16(data[i+1])
			value := uint32(data[i+2])<<24 | uint32(data[i+3])<<16 | uint32(data[i+4])<<8 | uint32(data[i+5])

			setting := Setting{
				Identifier: identifier,
				Value:      value,
			}
			frame.Settings = append(frame.Settings, setting)
		}
	}

	return nil
}

// parsePushPromiseFrame parses HTTP/2 PUSH_PROMISE frame
func parsePushPromiseFrame(data []byte, frame *Frame) error {
	if len(data) < 4 {
		return fmt.Errorf("insufficient data for push promise frame")
	}

	promisedStreamID := uint32(data[0])<<24 | uint32(data[1])<<16 | uint32(data[2])<<8 | uint32(data[3])
	frame.PromisedStreamID = promisedStreamID

	if len(data) > 4 {
		frame.HeaderBlockFragment = data[4:]
	}

	return nil
}

// getFrameTypeName returns the name of the frame type
func getFrameTypeName(frameType byte) string {
	frameTypes := map[byte]string{
		0x0: "DATA",
		0x1: "HEADERS",
		0x2: "PRIORITY",
		0x3: "RST_STREAM",
		0x4: "SETTINGS",
		0x5: "PUSH_PROMISE",
		0x6: "PING",
		0x7: "GOAWAY",
		0x8: "WINDOW_UPDATE",
		0x9: "CONTINUATION",
	}

	if name, exists := frameTypes[frameType]; exists {
		return name
	}
	return fmt.Sprintf("Unknown (%d)", frameType)
}

// Key returns the key for storing in analyzers
func (h *HTTP2) Key() string {
	return "http2"
}
