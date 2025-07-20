# Netmoth API Documentation

## Table of Contents

1. [Overview](#overview)
2. [Core Components](#core-components)
3. [Sensor Package](#sensor-package)
4. [Configuration Package](#configuration-package)
5. [Analyzer Package](#analyzer-package)
6. [Signature Package](#signature-package)
7. [Storage Package](#storage-package)
8. [Web Package](#web-package)
9. [Connection Package](#connection-package)
10. [Utils Package](#utils-package)
11. [Version Package](#version-package)

## Overview

Netmoth is a high-performance network monitoring and intrusion detection system built in Go. The system consists of two main components:

- **Agent**: Lightweight packet capture and analysis component
- **Manager**: Central server for data aggregation and web interface

## Core Components

### Main Entry Points

#### Agent (`cmd/agent/main.go`)
```go
func main()
```
Entry point for the Netmoth agent. Initializes configuration and starts the sensor.

**Parameters:**
- `-cfg string`: Configuration file path (default: "cmd/agent/config.yml")

#### Manager (`cmd/manager/main.go`)
```go
func main()
```
Entry point for the Netmoth manager (central server). Starts the web server with local configuration.

## Sensor Package

The sensor package is the core component responsible for packet capture, analysis, and data processing.

### Sensor (`internal/sensor/sensor.go`)

#### Types

```go
type sensor struct {
    strategy      strategies.PacketsCaptureStrategy
    db            *postgres.Connect
    detector      signature.Detector
    sensorMeta    *Metadata
    streamFactory *connection.TCPStreamFactory
    connections   chan *connection.Connection
    packets       []strategies.PacketDataSource
    packetPool    sync.Pool
    workerPool    chan struct{}
    statsMutex    sync.RWMutex
    packetStats   struct {
        received  uint64
        dropped   uint64
        processed uint64
    }
    // Agent mode fields
    agentClient       *AgentClient
    agentMode         bool
    dataInterval      time.Duration
    healthInterval    time.Duration
    connectionsBuffer []*connection.Connection
    signaturesBuffer  []signature.Detect
    bufferMutex       sync.RWMutex
}
```

#### Functions

```go
func New(config *config.Config)
```
Creates and initializes a new sensor instance.

**Parameters:**
- `config *config.Config`: Configuration object

**Features:**
- Initializes packet capture strategy
- Sets up database connection (if not in agent mode)
- Configures TCP stream factory
- Starts packet capture goroutines
- Handles graceful shutdown

```go
func (s *sensor) capturePacketsZeroCopy(source gopacket.ZeroCopyPacketDataSource, exit <-chan bool)
```
Captures packets using zero-copy mode for maximum performance.

**Parameters:**
- `source gopacket.ZeroCopyPacketDataSource`: Packet data source
- `exit <-chan bool`: Exit channel for graceful shutdown

```go
func (s *sensor) capturePackets(source gopacket.PacketDataSource, exit <-chan bool)
```
Captures packets using standard copy mode.

**Parameters:**
- `source gopacket.PacketDataSource`: Packet data source
- `exit <-chan bool`: Exit channel for graceful shutdown

```go
func (s *sensor) processPacket(packet gopacket.Packet)
```
Processes individual packets and extracts connection information.

**Parameters:**
- `packet gopacket.Packet`: Packet to process

```go
func (s *sensor) processConnections(logger *logSave)
```
Processes completed TCP connections and performs analysis.

**Parameters:**
- `logger *logSave`: Logger instance for saving results

### Metadata (`internal/sensor/sensor.go`)

```go
type Metadata struct {
    NetworkInterface string
    NetworkAddress   []string
}
```

Contains metadata about the sensor's network configuration.

### Agent Client (`internal/sensor/agent_client.go`)

```go
type AgentClient struct {
    serverURL string
    agentID   string
    interface string
    client    *http.Client
}
```

Handles communication between agent and central server.

#### Functions

```go
func NewAgentClient(serverURL, agentID, interface string) *AgentClient
```
Creates a new agent client instance.

```go
func (ac *AgentClient) Register(interface string) error
```
Registers the agent with the central server.

```go
func (ac *AgentClient) SendData(data []*connection.Connection) error
```
Sends connection data to the central server.

```go
func (ac *AgentClient) SendHealth(health *HealthStatus) error
```
Sends health status to the central server.

## Configuration Package

### Config (`internal/config/config.go`)

```go
type Config struct {
    Redis         Redis
    Interface     string
    Strategy      string
    Bpf           string
    LogFile       string `yaml:"log_file"`
    Postgres      Postgres
    NumberOfRings int  `yaml:"number_of_rings"`
    MaxCores      int  `yaml:"max_cores"`
    SnapLen       int  `yaml:"snapshot_length"`
    ConnTimeout   int  `yaml:"connection_timeout"`
    ZeroCopy      bool `yaml:"zero_copy"`
    Promiscuous   bool
    // Agent configuration
    AgentMode      bool   `yaml:"agent_mode"`
    AgentID        string `yaml:"agent_id"`
    ServerURL      string `yaml:"server_url"`
    DataInterval   int    `yaml:"data_interval"`   // seconds
    HealthInterval int    `yaml:"health_interval"` // seconds
}
```

#### Functions

```go
func New(cf string) (cfg *Config, err error)
```
Loads and validates configuration from a YAML file.

**Parameters:**
- `cf string`: Configuration file path

**Returns:**
- `cfg *Config`: Configuration object
- `err error`: Error if any

### Postgres (`internal/config/config.go`)

```go
type Postgres struct {
    User            string
    Password        string
    DB              string
    Host            string
    MaxConn         int `yaml:"max_conn"`
    MaxIDLecConn    int `yaml:"max_idlec_conn"`
    MaxLifeTimeConn int `yaml:"max_lifetime_conn"`
}
```

### Redis (`internal/config/config.go`)

```go
type Redis struct {
    Password string
    Host     string
}
```

## Analyzer Package

The analyzer package contains protocol-specific analyzers for different network protocols.

### HTTP Analyzer (`internal/analyzer/httpanalyzer/main.go`)

```go
type HTTP struct {
    Request  Request
    Response []signature.Detect
}

type Request struct {
    Method           string
    URL              string
    Headers          http.Header
    ContentLength    int64
    TransferEncoding []string
    Host             string
}
```

#### Functions

```go
func Analyze(conn *connection.Connection, detector signature.Detector) (*HTTP, error)
```
Analyzes HTTP traffic from a TCP connection.

**Parameters:**
- `conn *connection.Connection`: TCP connection
- `detector signature.Detector`: Signature detector

**Returns:**
- `*HTTP`: HTTP analysis results
- `error`: Error if any

### DNS Analyzer (`internal/analyzer/dnsanalyzer/`)

Analyzes DNS traffic and extracts domain information.

### TLS Analyzer (`internal/analyzer/tlsanalyzer/`)

Analyzes TLS/SSL traffic and extracts certificate information.

### HTTP/2 Analyzer (`internal/analyzer/http2analyzer/`)

Analyzes HTTP/2 traffic using the HTTP/2 protocol.

### Content Analyzer (`internal/analyzer/contentanalyzer/`)

Analyzes packet content for patterns and signatures.

## Signature Package

### Detector (`internal/signature/signature.go`)

```go
type Detector struct {
    postgres.Connect
}
```

Performs signature-based detection against various threat databases.

#### Types

```go
type Request struct {
    IP         string
    TrackerURL string
    CertSHA1   string `json:",omitempty"`
    Port       int
}

type Detect struct {
    Type        string
    Provider    string
    SignatureID int
}
```

#### Functions

```go
func New(conn postgres.Connect) Detector
```
Creates a new signature detector.

```go
func (d *Detector) Scan(req *Request) ([]Detect, error)
```
Scans a request against signature databases.

**Parameters:**
- `req *Request`: Request to scan

**Returns:**
- `[]Detect`: Detection results
- `error`: Error if any

**Detection Types:**
- IP address blacklists
- Botnet signatures
- Tracker URLs
- Certificate blacklists

### Update (`internal/signature/update.go`)

Handles signature database updates and synchronization.

## Storage Package

### PostgreSQL (`internal/storage/postgres/postgres.go`)

```go
type Connect struct {
    Conn *sql.DB
}

type PgSQLConfig struct {
    DSN             string
    MaxConn         int
    MaxIdleConn     int
    MaxLifetimeConn int
}
```

#### Functions

```go
func New(ctx context.Context, conf *PgSQLConfig) (*Connect, error)
```
Creates a new PostgreSQL connection.

**Parameters:**
- `ctx context.Context`: Context for connection
- `conf *PgSQLConfig`: PostgreSQL configuration

**Returns:**
- `*Connect`: Database connection
- `error`: Error if any

### Redis (`internal/storage/redis/redis.go`)

Handles Redis connections for caching and real-time data.

### Sanitize (`internal/storage/postgres/sanitize/sanitize.go`)

Provides SQL sanitization utilities for safe database operations.

## Web Package

### Web Server (`internal/web/main.go`)

```go
func New(configPath string)
```
Starts the web server with the specified configuration.

**Parameters:**
- `configPath string`: Configuration file path

**Features:**
- RESTful API endpoints
- WebSocket support
- Static file serving
- CORS middleware
- Profiling endpoints

#### API Endpoints

- `GET /api/version`: Get server version
- `POST /api/agent/register`: Register agent
- `POST /api/agent/data`: Receive agent data
- `GET /api/agent/health`: Agent health check
- `GET /ws`: WebSocket endpoint

### Agent API (`internal/web/agent_api.go`)

Handles agent-specific API endpoints for registration, data transmission, and health monitoring.

## Connection Package

### TCP Stream Factory (`internal/connection/`)

```go
type TCPStreamFactory struct {
    Connections chan *connection.Connection
    ConnTimeout int
    Assembler   *tcpassembly.Assembler
    Ticker      *time.Ticker
}
```

Manages TCP stream reassembly and connection tracking.

#### Functions

```go
func (f *TCPStreamFactory) CreateAssembler()
```
Creates a TCP stream assembler for packet reassembly.

```go
func (f *TCPStreamFactory) ProcessStream(net, transport gopacket.Flow, r io.ReadCloser)
```
Processes reassembled TCP streams.

## Utils Package

### Interface Utilities (`internal/utils/`)

Provides utilities for network interface management and address resolution.

## Version Package

### Version (`internal/version/version.go`)

```go
func Version() string
```
Returns the current version of the application.

## Packet Capture Strategies

### Strategy Interface (`internal/sensor/strategies/strategies.go`)

```go
type PacketsCaptureStrategy interface {
    New(config *config.Config) ([]gopacket.PacketDataSource, error)
    Destroy()
}
```

### Available Strategies

1. **PCAP** (`internal/sensor/strategies/pcap.go`): Standard libpcap-based capture
2. **AF_PACKET** (`internal/sensor/strategies/afpacket.go`): Linux AF_PACKET socket capture
3. **PF_RING** (`internal/sensor/strategies/pfring.go`): PF_RING high-performance capture
4. **eBPF** (`internal/sensor/strategies/ebpf.go`): Extended Berkeley Packet Filter capture

Each strategy implements the `PacketsCaptureStrategy` interface and provides different performance characteristics and features.

## Error Handling

The codebase uses standard Go error handling patterns:

- Functions return `error` types for error conditions
- Critical errors are logged and may cause program termination
- Graceful degradation is implemented where possible
- Context cancellation is supported for long-running operations

## Concurrency

The system uses Go's concurrency primitives extensively:

- **Goroutines**: For packet capture, processing, and communication
- **Channels**: For inter-goroutine communication
- **Mutexes**: For protecting shared data structures
- **Context**: For cancellation and timeout management
- **Sync.Pool**: For efficient memory management

## Performance Considerations

- **Zero-copy packet processing**: Minimizes memory allocations
- **Connection pooling**: Efficient database connection management
- **Worker pools**: Controlled concurrency for packet processing
- **Buffer management**: Efficient memory usage with sync.Pool
- **eBPF support**: Kernel-level packet processing for maximum performance

## Security Features

- **SQL injection prevention**: Parameterized queries and sanitization
- **Input validation**: Configuration and network interface validation
- **Secure defaults**: Safe configuration defaults
- **Certificate validation**: TLS certificate analysis
- **Signature-based detection**: Multiple threat detection mechanisms