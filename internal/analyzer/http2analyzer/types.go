package http2analyzer

// HTTP2 represents the result of HTTP/2 traffic analysis
type HTTP2 struct {
	Protocol           string  `json:"protocol,omitempty"`
	HasHTTP2Preface    bool    `json:"has_http2_preface"`
	PayloadLength      int     `json:"payload_length"`
	Frames             []Frame `json:"frames,omitempty"`
	HeadersFrames      []Frame `json:"headers_frames,omitempty"`
	DataFrames         []Frame `json:"data_frames,omitempty"`
	SettingsFrames     []Frame `json:"settings_frames,omitempty"`
	PushPromiseFrames  []Frame `json:"push_promise_frames,omitempty"`
	TotalFrames        int     `json:"total_frames"`
	TotalDataFrames    int     `json:"total_data_frames"`
	TotalHeadersFrames int     `json:"total_headers_frames"`
}

// Frame represents an HTTP/2 frame
type Frame struct {
	Length              int               `json:"length"`
	Type                string            `json:"type"`
	Flags               byte              `json:"flags"`
	StreamID            uint32            `json:"stream_id"`
	Data                []byte            `json:"data,omitempty"`
	StreamDependency    uint32            `json:"stream_dependency,omitempty"`
	Weight              byte              `json:"weight,omitempty"`
	HeaderBlockFragment []byte            `json:"header_block_fragment,omitempty"`
	DataPayload         []byte            `json:"data_payload,omitempty"`
	DataLength          int               `json:"data_length,omitempty"`
	Settings            []Setting         `json:"settings,omitempty"`
	PromisedStreamID    uint32            `json:"promised_stream_id,omitempty"`
	Headers             map[string]string `json:"headers,omitempty"`
}

// Setting represents an HTTP/2 setting
type Setting struct {
	Identifier uint16 `json:"identifier"`
	Value      uint32 `json:"value"`
}
