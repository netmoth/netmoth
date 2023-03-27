package connection

import (
	"bytes"
	"time"
)

// Connection is ...
type Connection struct {
	Timestamp       time.Time
	UID             uint64
	SourceIP        string
	SourcePort      int
	DestinationIP   string
	DestinationPort int
	TransportType   string
	Duration        float64
	State           string        `json:",omitempty"`
	Payload         *bytes.Buffer `json:"-"`
	Analyzers       map[string]interface{}
}
