package contentanalyzer

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"sync"

	"github.com/netmoth/netmoth/internal/connection"
)

// BufferPool пул буферов для переиспользования
var BufferPool = sync.Pool{
	New: func() any {
		return new(bytes.Buffer)
	},
}

// Analyze analyzes content from various protocols and extracts useful information
func Analyze(c *connection.Connection) (*Content, error) {
	if c.Payload == nil || c.Payload.Len() == 0 {
		return nil, fmt.Errorf("empty payload")
	}

	// Получаем данные без копирования
	data := c.Payload.Bytes()
	result := &Content{
		PayloadLength: len(data),
		RawData:       data,
	}

	// Try to detect content type
	result.ContentType = detectContentType(data)

	// Extract text content if possible
	if err := extractTextContent(data, result); err != nil {
		return nil, fmt.Errorf("error extracting text content: %w", err)
	}

	// Try to decompress if compressed
	if err := decompressContent(data, result); err != nil {
		return nil, fmt.Errorf("error decompressing content: %w", err)
	}

	// Extract structured data
	if err := extractStructuredData(data, result); err != nil {
		return nil, fmt.Errorf("error extracting structured data: %w", err)
	}

	// Extract URLs and domains
	extractURLsAndDomains(data, result)

	// Extract file signatures
	extractFileSignatures(data, result)

	return result, nil
}

// detectContentType tries to detect the content type
func detectContentType(data []byte) string {
	if len(data) == 0 {
		return "unknown"
	}

	// Check for JSON
	if detectJSONPattern(data) {
		return "application/json"
	}

	// Check for XML
	if detectXMLPattern(data) {
		return "application/xml"
	}

	// Check for HTML
	if detectHTMLPattern(data) {
		return "text/html"
	}

	// Check for plain text
	if detectTextPattern(data) {
		return "text/plain"
	}

	// Check for binary data
	if detectBinaryPattern(data) {
		return "application/octet-stream"
	}

	return "unknown"
}

// extractTextContent extracts readable text from the data
func extractTextContent(data []byte, result *Content) error {
	if len(data) == 0 {
		return nil
	}

	// Try to extract as UTF-8 text
	text := string(data)

	// Check if it's mostly printable
	printableCount := 0
	for _, b := range data {
		if b >= 32 && b <= 126 || b == 9 || b == 10 || b == 13 {
			printableCount++
		}
	}

	if float64(printableCount)/float64(len(data)) > 0.7 {
		result.TextContent = text
		result.IsText = true
	}

	return nil
}

// decompressContent tries to decompress compressed content
func decompressContent(data []byte, result *Content) error {
	// Try gzip decompression
	if len(data) >= 2 && data[0] == 0x1f && data[1] == 0x8b {
		reader, err := gzip.NewReader(bytes.NewReader(data))
		if err == nil {
			defer reader.Close()
			// Используем пул буферов
			buffer := BufferPool.Get().(*bytes.Buffer)
			buffer.Reset()
			defer BufferPool.Put(buffer)

			_, err := io.Copy(buffer, reader)
			if err == nil {
				result.DecompressedContent = make([]byte, buffer.Len())
				copy(result.DecompressedContent, buffer.Bytes())
				result.CompressionType = "gzip"
				return nil
			}
		}
	}

	// Try zlib decompression
	if len(data) >= 2 && (data[0] == 0x78 && (data[1] == 0x01 || data[1] == 0x9c || data[1] == 0xda)) {
		reader, err := zlib.NewReader(bytes.NewReader(data))
		if err == nil {
			defer reader.Close()
			// Используем пул буферов
			buffer := BufferPool.Get().(*bytes.Buffer)
			buffer.Reset()
			defer BufferPool.Put(buffer)

			_, err := io.Copy(buffer, reader)
			if err == nil {
				result.DecompressedContent = make([]byte, buffer.Len())
				copy(result.DecompressedContent, buffer.Bytes())
				result.CompressionType = "zlib"
				return nil
			}
		}
	}

	return nil
}

// extractStructuredData extracts structured data like JSON, XML
func extractStructuredData(data []byte, result *Content) error {
	// Try to parse as JSON
	if detectJSONPattern(data) {
		var jsonData any
		if err := json.Unmarshal(data, &jsonData); err == nil {
			result.StructuredData = jsonData
			result.DataType = "json"
		}
	}

	// Try to parse as XML (simplified)
	if detectXMLPattern(data) {
		result.DataType = "xml"
		// Extract XML tags
		extractXMLTags(data, result)
	}

	return nil
}

// extractXMLTags extracts XML tags from the data
func extractXMLTags(data []byte, result *Content) {
	text := string(data)
	tags := make([]string, 0)

	// Simple XML tag extraction
	words := strings.Fields(text)
	for _, word := range words {
		if strings.HasPrefix(word, "<") && strings.HasSuffix(word, ">") {
			tags = append(tags, word)
		}
	}

	result.XMLTags = tags
}

// extractURLsAndDomains extracts URLs and domains from the content
func extractURLsAndDomains(data []byte, result *Content) {
	text := string(data)

	// Extract URLs
	urls := extractURLs(text)
	result.URLs = urls

	// Extract domains
	domains := extractDomains(text)
	result.Domains = domains
}

// extractURLs extracts URLs from text
func extractURLs(text string) []string {
	urls := make([]string, 0)

	// Simple URL extraction
	words := strings.Fields(text)
	for _, word := range words {
		if strings.HasPrefix(word, "http://") || strings.HasPrefix(word, "https://") {
			urls = append(urls, word)
		}
	}

	return urls
}

// extractDomains extracts domains from text
func extractDomains(text string) []string {
	domains := make([]string, 0)

	// Simple domain extraction
	words := strings.Fields(text)
	for _, word := range words {
		if strings.Contains(word, ".") && !strings.HasPrefix(word, "http") {
			// Remove common suffixes
			word = strings.TrimSuffix(word, ",")
			word = strings.TrimSuffix(word, ";")
			word = strings.TrimSuffix(word, ".")

			if len(word) > 3 && strings.Contains(word, ".") {
				domains = append(domains, word)
			}
		}
	}

	return domains
}

// extractFileSignatures extracts file signatures
func extractFileSignatures(data []byte, result *Content) {
	if len(data) < 4 {
		return
	}

	// Common file signatures
	signatures := map[string][]byte{
		"JPEG": {0xFF, 0xD8, 0xFF},
		"PNG":  {0x89, 0x50, 0x4E, 0x47},
		"GIF":  {0x47, 0x49, 0x46},
		"BMP":  {0x42, 0x4D},
		"PDF":  {0x25, 0x50, 0x44, 0x46},
		"ZIP":  {0x50, 0x4B, 0x03, 0x04},
		"RAR":  {0x52, 0x61, 0x72, 0x21},
		"EXE":  {0x4D, 0x5A},
		"ELF":  {0x7F, 0x45, 0x4C, 0x46},
	}

	for fileType, signature := range signatures {
		if len(data) >= len(signature) && bytes.Equal(data[:len(signature)], signature) {
			result.FileType = fileType
			result.IsBinary = true
			break
		}
	}
}

// Helper functions (same as in TLS analyzer)
func detectJSONPattern(data []byte) bool {
	if len(data) < 2 {
		return false
	}
	start := string(data[:2])
	return start == "{\"" || start == "[{" || start == "[\"" || start == "{\n" || start == "{\r"
}

func detectXMLPattern(data []byte) bool {
	if len(data) < 5 {
		return false
	}
	start := string(data[:5])
	return start == "<?xml" || start == "<root" || start == "<html" || start == "<soap"
}

func detectHTMLPattern(data []byte) bool {
	if len(data) < 6 {
		return false
	}
	start := strings.ToLower(string(data[:6]))
	return start == "<html" || start == "<!doct" || start == "<head>"
}

func detectTextPattern(data []byte) bool {
	if len(data) == 0 {
		return false
	}
	printableCount := 0
	for _, b := range data {
		if b >= 32 && b <= 126 || b == 9 || b == 10 || b == 13 {
			printableCount++
		}
	}
	return float64(printableCount)/float64(len(data)) > 0.8
}

func detectBinaryPattern(data []byte) bool {
	if len(data) == 0 {
		return false
	}
	printableCount := 0
	for _, b := range data {
		if b >= 32 && b <= 126 || b == 9 || b == 10 || b == 13 {
			printableCount++
		}
	}
	return float64(printableCount)/float64(len(data)) < 0.7
}

// Key returns the key for storing in analyzers
func (c *Content) Key() string {
	return "content"
}
