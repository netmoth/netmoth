# Netmoth Configuration Reference

## Table of Contents

1. [Overview](#overview)
2. [Configuration File Format](#configuration-file-format)
3. [Agent Configuration](#agent-configuration)
4. [Manager Configuration](#manager-configuration)
5. [Network Configuration](#network-configuration)
6. [Database Configuration](#database-configuration)
7. [Performance Configuration](#performance-configuration)
8. [Security Configuration](#security-configuration)
9. [Logging Configuration](#logging-configuration)
10. [Configuration Examples](#configuration-examples)

## Overview

Netmoth uses YAML configuration files to control all aspects of the system behavior. Configuration files are loaded at startup and can be specified via command-line arguments.

### Configuration File Locations

- **Agent**: `cmd/agent/config.yml` (default)
- **Manager**: `cmd/manager/config.yml` (default)
- **Custom**: Specify with `-cfg` flag

### Configuration File Structure

```yaml
# Network Interface Configuration
interface: "eth0"
strategy: "pcap"
promiscuous: true

# Performance Settings
zero_copy: true
max_cores: 0
snapshot_length: 512

# Database Configuration
postgres:
  host: "localhost"
  port: 5432
  user: "netmoth"
  password: "netmoth"
  db: "netmoth"

# Agent Configuration
agent_mode: false
agent_id: "agent-001"
server_url: "http://localhost:3000"
```

## Configuration File Format

### YAML Syntax

```yaml
# Comments start with #
# Use spaces for indentation (not tabs)
# String values can be quoted or unquoted
# Boolean values: true, false, yes, no, on, off

# Basic configuration
setting: value
boolean_setting: true
numeric_setting: 42

# Nested configuration
section:
  subsection:
    key: value
```

### Environment Variable Override

Configuration values can be overridden using environment variables:

```bash
# Override database host
export NETMOTH_POSTGRES_HOST=production-db.example.com

# Override agent ID
export NETMOTH_AGENT_ID=prod-agent-01
```

## Agent Configuration

### Basic Agent Settings

```yaml
# Agent mode enables distributed operation
agent_mode: true

# Unique identifier for this agent
agent_id: "agent-001"

# Central server URL for data transmission
server_url: "http://localhost:3000"

# Data transmission interval (seconds)
data_interval: 60

# Health check interval (seconds)
health_interval: 300
```

### Network Interface Configuration

```yaml
# Network interface to capture packets from
interface: "eth0"

# Packet capture strategy
strategy: "pcap"  # Options: pcap, afpacket, pfring, ebpf

# Enable promiscuous mode
promiscuous: true

# Berkeley Packet Filter (BPF) filter
bpf: "port 80 or port 443"

# Number of rings for cluster mode (PF_RING)
number_of_rings: 1
```

### Performance Settings

```yaml
# Enable zero-copy packet processing
zero_copy: true

# Maximum number of CPU cores to use (0 = auto-detect)
max_cores: 0

# Snapshot length for packet capture
snapshot_length: 512

# Connection timeout (seconds, 0 = no timeout)
connection_timeout: 0
```

### Logging Configuration

```yaml
# Log file path
log_file: "agent.log"

# Log level (debug, info, warn, error)
log_level: "info"

# Enable debug mode
debug: false
```

## Manager Configuration

### Server Settings

```yaml
# Server host and port
host: "0.0.0.0"
port: 3000

# TLS configuration
tls:
  enabled: false
  cert_file: "server.crt"
  key_file: "server.key"

# CORS settings
cors:
  allowed_origins: ["*"]
  allowed_methods: ["GET", "POST", "PUT", "DELETE"]
  allowed_headers: ["Content-Type", "Authorization"]
```

### Agent Management

```yaml
# Agent registration settings
agent_registration:
  enabled: true
  auto_approve: true
  max_agents: 1000

# Agent health monitoring
agent_health:
  check_interval: 60
  timeout: 30
  max_failures: 3
```

### Data Processing

```yaml
# Data processing settings
data_processing:
  batch_size: 1000
  batch_timeout: 5
  max_workers: 10

# Data retention
data_retention:
  connection_data: "30d"
  signature_data: "90d"
  log_data: "7d"
```

## Network Configuration

### Packet Capture Strategies

#### PCAP Strategy

```yaml
strategy: "pcap"
pcap:
  # PCAP-specific settings
  buffer_size: 262144
  timeout: 1000
  immediate_mode: true
```

#### AF_PACKET Strategy

```yaml
strategy: "afpacket"
afpacket:
  # AF_PACKET-specific settings
  frame_size: 65536
  block_size: 262144
  num_blocks: 128
  timeout: 1000
```

#### PF_RING Strategy

```yaml
strategy: "pfring"
pfring:
  # PF_RING-specific settings
  cluster_id: 99
  cluster_type: "cluster_round_robin"
  num_rings: 4
```

#### eBPF Strategy

```yaml
strategy: "ebpf"
ebpf:
  # eBPF-specific settings
  program_file: "packet_capture.o"
  map_size: 65536
  max_entries: 10000
```

### Network Filtering

```yaml
# Berkeley Packet Filter (BPF) filter
bpf: "port 80 or port 443 or port 53"

# IP address filtering
ip_filter:
  include:
    - "192.168.1.0/24"
    - "10.0.0.0/8"
  exclude:
    - "192.168.1.100"
    - "10.0.0.1"

# Port filtering
port_filter:
  include:
    - 80
    - 443
    - 53
  exclude:
    - 22
    - 3389
```

## Database Configuration

### PostgreSQL Configuration

```yaml
postgres:
  # Connection settings
  host: "localhost"
  port: 5432
  user: "netmoth"
  password: "netmoth"
  db: "netmoth"
  ssl_mode: "disable"

  # Connection pool settings
  max_conn: 50
  max_idlec_conn: 10
  max_lifetime_conn: 300

  # Query timeout
  query_timeout: 30

  # Migration settings
  migrations:
    enabled: true
    path: "./migration"
    table: "schema_migrations"
```

### Redis Configuration

```yaml
redis:
  # Connection settings
  host: "localhost"
  port: 6379
  password: ""
  db: 0

  # Connection pool settings
  pool_size: 10
  min_idle_conns: 5
  max_retries: 3

  # Timeout settings
  dial_timeout: 5
  read_timeout: 3
  write_timeout: 3

  # Cache settings
  cache:
    default_ttl: 3600
    max_memory: "1gb"
    eviction_policy: "allkeys-lru"
```

## Performance Configuration

### Memory Management

```yaml
# Memory pool settings
memory_pool:
  # Packet buffer pool
  packet_buffer_size: 65536
  packet_buffer_count: 1000

  # Connection buffer pool
  connection_buffer_size: 32768
  connection_buffer_count: 500

  # TCP stream buffer pool
  stream_buffer_size: 16384
  stream_buffer_count: 200
```

### Concurrency Settings

```yaml
# Worker pool configuration
worker_pool:
  # Packet processing workers
  packet_workers: 4
  # Connection processing workers
  connection_workers: 2
  # Analysis workers
  analysis_workers: 2

# Channel buffer sizes
channels:
  packet_channel: 10000
  connection_channel: 5000
  analysis_channel: 1000
```

### Processing Pipeline

```yaml
# Processing pipeline settings
pipeline:
  # Enable parallel processing
  parallel: true
  
  # Batch processing settings
  batch:
    size: 100
    timeout: 100
    
  # Stream processing settings
  stream:
    buffer_size: 8192
    flush_interval: 1000
```

## Security Configuration

### Authentication and Authorization

```yaml
# Authentication settings
auth:
  enabled: true
  type: "jwt"  # Options: jwt, basic, oauth2
  
  # JWT settings
  jwt:
    secret: "your-secret-key"
    expires_in: 3600
    issuer: "netmoth"
    
  # Basic auth settings
  basic:
    users:
      - username: "admin"
        password: "hashed-password"
        role: "admin"
      - username: "user"
        password: "hashed-password"
        role: "user"
```

### Network Security

```yaml
# Network security settings
security:
  # TLS/SSL settings
  tls:
    enabled: true
    cert_file: "server.crt"
    key_file: "server.key"
    ca_file: "ca.crt"
    min_version: "1.2"
    
  # Firewall rules
  firewall:
    allowed_ips:
      - "192.168.1.0/24"
      - "10.0.0.0/8"
    blocked_ips:
      - "192.168.1.100"
      
  # Rate limiting
  rate_limit:
    enabled: true
    requests_per_minute: 1000
    burst_size: 100
```

### Data Protection

```yaml
# Data protection settings
data_protection:
  # Data encryption
  encryption:
    enabled: true
    algorithm: "AES-256-GCM"
    key_file: "encryption.key"
    
  # Data anonymization
  anonymization:
    enabled: false
    ip_masking: "partial"  # Options: none, partial, full
    port_masking: false
    
  # Data retention
  retention:
    connection_data: "30d"
    signature_data: "90d"
    log_data: "7d"
```

## Logging Configuration

### Log Levels and Output

```yaml
# Logging configuration
logging:
  # Log level
  level: "info"  # Options: debug, info, warn, error
  
  # Log format
  format: "json"  # Options: json, text
  
  # Log output
  output: "file"  # Options: file, stdout, stderr
  
  # Log file settings
  file:
    path: "netmoth.log"
    max_size: 100  # MB
    max_age: 30    # days
    max_backups: 10
    compress: true
```

### Structured Logging

```yaml
# Structured logging fields
structured_logging:
  # Add custom fields to all log entries
  fields:
    service: "netmoth"
    version: "1.0.0"
    environment: "production"
    
  # Log specific events
  events:
    packet_capture:
      enabled: true
      level: "debug"
    connection_analysis:
      enabled: true
      level: "info"
    signature_detection:
      enabled: true
      level: "warn"
```

## Configuration Examples

### Basic Agent Configuration

```yaml
# Basic agent configuration for development
interface: "eth0"
strategy: "pcap"
promiscuous: true
zero_copy: true
max_cores: 0
snapshot_length: 512
connection_timeout: 0
log_file: "agent.log"

# Agent mode settings
agent_mode: true
agent_id: "dev-agent-001"
server_url: "http://localhost:3000"
data_interval: 60
health_interval: 300

# Database settings (for local development)
postgres:
  host: "localhost"
  port: 5432
  user: "netmoth"
  password: "netmoth"
  db: "netmoth"
  max_conn: 10
  max_idlec_conn: 5
  max_lifetime_conn: 300
```

### High-Performance Agent Configuration

```yaml
# High-performance agent configuration
interface: "eth0"
strategy: "ebpf"
promiscuous: true
zero_copy: true
max_cores: 8
snapshot_length: 65536
connection_timeout: 0
log_file: "agent.log"

# Agent mode settings
agent_mode: true
agent_id: "prod-agent-001"
server_url: "https://netmoth.example.com"
data_interval: 30
health_interval: 60

# Performance optimizations
ebpf:
  program_file: "packet_capture.o"
  map_size: 131072
  max_entries: 50000

# Memory pool settings
memory_pool:
  packet_buffer_size: 131072
  packet_buffer_count: 2000
  connection_buffer_size: 65536
  connection_buffer_count: 1000

# Worker pool settings
worker_pool:
  packet_workers: 8
  connection_workers: 4
  analysis_workers: 4
```

### Production Manager Configuration

```yaml
# Production manager configuration
host: "0.0.0.0"
port: 3000

# TLS settings
tls:
  enabled: true
  cert_file: "/etc/netmoth/server.crt"
  key_file: "/etc/netmoth/server.key"
  ca_file: "/etc/netmoth/ca.crt"

# Database settings
postgres:
  host: "db.example.com"
  port: 5432
  user: "netmoth"
  password: "secure-password"
  db: "netmoth"
  ssl_mode: "require"
  max_conn: 100
  max_idlec_conn: 20
  max_lifetime_conn: 600

# Redis settings
redis:
  host: "redis.example.com"
  port: 6379
  password: "redis-password"
  db: 0
  pool_size: 20
  min_idle_conns: 10

# Security settings
security:
  tls:
    enabled: true
    cert_file: "/etc/netmoth/server.crt"
    key_file: "/etc/netmoth/server.key"
  firewall:
    allowed_ips:
      - "10.0.0.0/8"
      - "192.168.0.0/16"
  rate_limit:
    enabled: true
    requests_per_minute: 10000
    burst_size: 1000

# Logging settings
logging:
  level: "info"
  format: "json"
  output: "file"
  file:
    path: "/var/log/netmoth/manager.log"
    max_size: 100
    max_age: 30
    max_backups: 10
    compress: true
```

### Development Configuration

```yaml
# Development configuration
interface: "lo"  # Loopback interface for testing
strategy: "pcap"
promiscuous: false
zero_copy: false
max_cores: 2
snapshot_length: 256
connection_timeout: 0
log_file: "dev.log"

# Debug settings
debug: true
log_level: "debug"

# Local database settings
postgres:
  host: "localhost"
  port: 5432
  user: "netmoth"
  password: "netmoth"
  db: "netmoth_dev"
  max_conn: 5
  max_idlec_conn: 2
  max_lifetime_conn: 300

# Local Redis settings
redis:
  host: "localhost"
  port: 6379
  password: ""
  db: 0
  pool_size: 5
  min_idle_conns: 2

# Development server settings
host: "localhost"
port: 3000
debug: true
```

## Configuration Validation

### Validation Rules

```yaml
# Configuration validation rules
validation:
  # Required fields
  required:
    - interface
    - strategy
    - log_file
    
  # Field validation
  rules:
    interface:
      type: "string"
      pattern: "^[a-zA-Z0-9._-]+$"
    strategy:
      type: "string"
      enum: ["pcap", "afpacket", "pfring", "ebpf"]
    max_cores:
      type: "integer"
      min: 0
      max: 64
    data_interval:
      type: "integer"
      min: 1
      max: 3600
```

### Configuration Testing

```bash
# Validate configuration file
./bin/netmoth validate -cfg config.yml

# Test configuration
./bin/netmoth test-config -cfg config.yml

# Generate configuration template
./bin/netmoth generate-config > config.template.yml
```

This configuration reference provides comprehensive documentation for all Netmoth configuration options, enabling users to customize the system behavior according to their specific requirements.