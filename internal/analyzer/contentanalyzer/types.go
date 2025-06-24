package contentanalyzer

// Content represents the result of content analysis
type Content struct {
	PayloadLength       int      `json:"payload_length"`
	ContentType         string   `json:"content_type,omitempty"`
	DataType            string   `json:"data_type,omitempty"`
	IsText              bool     `json:"is_text"`
	IsBinary            bool     `json:"is_binary"`
	FileType            string   `json:"file_type,omitempty"`
	CompressionType     string   `json:"compression_type,omitempty"`
	TextContent         string   `json:"text_content,omitempty"`
	RawData             []byte   `json:"raw_data,omitempty"`
	DecompressedContent []byte   `json:"decompressed_content,omitempty"`
	StructuredData      any      `json:"structured_data,omitempty"`
	XMLTags             []string `json:"xml_tags,omitempty"`
	URLs                []string `json:"urls,omitempty"`
	Domains             []string `json:"domains,omitempty"`
	HasCompression      bool     `json:"has_compression"`
	EstimatedSize       int      `json:"estimated_size,omitempty"`
}
