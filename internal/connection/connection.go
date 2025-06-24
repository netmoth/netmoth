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
	// Оптимизации производительности
	reused bool
}

// ConnectionPool предоставляет пул объектов Connection для переиспользования
type ConnectionPool struct {
	pool sync.Pool
}

// NewConnectionPool создает новый пул соединений
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

// Get получает Connection из пула или создает новый
func (cp *ConnectionPool) Get() *Connection {
	connInterface := cp.pool.Get()
	conn := connInterface.(*Connection)
	conn.reused = true
	// Сбрасываем состояние
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

// Put возвращает Connection в пул
func (cp *ConnectionPool) Put(conn *Connection) {
	if conn.reused {
		cp.pool.Put(conn)
	}
}

// GlobalConnectionPool глобальный пул соединений
var GlobalConnectionPool = NewConnectionPool()
