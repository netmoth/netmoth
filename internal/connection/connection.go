package connection

import (
	"bytes"
	"time"
)

// Connection is ...
type Connection struct {
	Timestamp       time.Time
	Payload         *bytes.Buffer `json:"-"`
	Analyzers       map[string]interface{}
	SourceIP        string
	DestinationIP   string
	TransportType   string
	State           string `json:",omitempty"`
	UID             uint64
	SourcePort      int
	DestinationPort int
	Duration        float64
}
