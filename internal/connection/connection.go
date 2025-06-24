package connection

import (
	"bytes"
	"sync"
	"time"
)

// Connection is ...
type Connection struct {
	Timestamp       time.Time
	Payload         *bytes.Buffer `json:"-"`
	Analyzers       map[string]any
	SourceIP        string
	DestinationIP   string
	TransportType   string
	State           string `json:",omitempty"`
	UID             uint64
	SourcePort      int
	DestinationPort int
	Duration        float64
	// Performance optimizations
	reused bool
}

// ConnectionPool provides a pool of Connection objects for reuse
type ConnectionPool struct {
	pool sync.Pool
}

// NewConnectionPool creates a new connection pool
func NewConnectionPool() *ConnectionPool {
	return &ConnectionPool{
		pool: sync.Pool{
			New: func() any {
				return &Connection{
					Analyzers: make(map[string]any),
					Payload:   new(bytes.Buffer),
				}
			},
		},
	}
}

// Get retrieves a Connection from the pool or creates a new one
func (cp *ConnectionPool) Get() *Connection {
	connInterface := cp.pool.Get()
	conn := connInterface.(*Connection)
	conn.reused = true
	// Reset state
	conn.Payload.Reset()
	conn.Analyzers = make(map[string]any)
	conn.Timestamp = time.Time{}
	conn.SourceIP = ""
	conn.DestinationIP = ""
	conn.TransportType = ""
	conn.State = ""
	conn.UID = 0
	conn.SourcePort = 0
	conn.DestinationPort = 0
	conn.Duration = 0
	return conn
}

// Put returns a Connection to the pool
func (cp *ConnectionPool) Put(conn *Connection) {
	if conn.reused {
		cp.pool.Put(conn)
	}
}

// GlobalConnectionPool is a global connection pool
var GlobalConnectionPool = NewConnectionPool()
