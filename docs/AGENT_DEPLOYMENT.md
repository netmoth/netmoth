# Netmoth Agent Deployment Guide

This document describes how to deploy and configure Netmoth agents to collect audit data from different machines and send it to a central server.

## Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Agent #1      │    │   Agent #2      │    │   Agent #N      │
│   (Machine 1)   │    │   (Machine 2)   │    │   (Machine N)   │
└─────────┬───────┘    └─────────┬───────┘    └─────────┬───────┘
          │                      │                      │
          └──────────────────────┼──────────────────────┘
                                 │
                    ┌─────────────▼─────────────┐
                    │   Central Server          │
                    │   (Manager)               │
                    │   - API Server            │
                    │   - Database              │
                    │   - Web Interface         │
                    └───────────────────────────┘
```

## Components

### 1. Agent
- Collects network traffic from local interface
- Analyzes connections and detects signatures
- Sends data to central server
- Periodically sends health checks

### 2. Central Server (Manager)
- Receives data from agents via REST API
- Stores data in database
- Provides web interface for monitoring
- Manages agent registration and status

## Installation and Configuration

### 1. Building Components

```bash
# Build agent
make build-agent

# Build central server
make build-manager
```

### 2. Central Server Configuration

Choose the appropriate configuration file for your central server:

**Basic configuration:**
```bash
cp cmd/manager/config.yml.example cmd/manager/config.yml
```

**Optimized configuration (for high-load environments):**
```bash
cp cmd/manager/config_optimized.yml cmd/manager/config.yml
```

**eBPF configuration (for maximum performance):**
```bash
cp cmd/manager/config_ebpf.yml cmd/manager/config.yml
```

Example configuration:
```yaml
# Netmoth Manager Configuration
interface: eth0
strategy: pcap
number_of_rings: 1
zero_copy: true

snapshot_length: 512
promiscuous: true
connection_timeout: 0
bpf: ""
max_cores: 0
log_file: manager.log

# Manager configuration (no agent mode)
agent_mode: false

# Database configuration for storing agent data
postgres:
  user: netmoth
  password: postgresPassword
  db: netmoth
  host: localhost:5432
  max_conn: 50
  max_idlec_conn: 10
  max_lifetime_conn: 300

# Redis configuration for caching
redis:
  password: ""
  host: "localhost:6379"
```

### 3. Agent Configuration

Choose the appropriate configuration file for your agent:

**Basic configuration:**
```bash
cp cmd/agent/config.yml.example cmd/agent/config.yml
```

**Optimized configuration (for high-load environments):**
```bash
cp cmd/agent/config_optimized.yml cmd/agent/config.yml
```

**eBPF configuration (for maximum performance):**
```bash
cp cmd/agent/config_ebpf.yml cmd/agent/config.yml
```

Example configuration:
```yaml
# Netmoth Agent Configuration
interface: eth0
strategy: pcap
number_of_rings: 1
zero_copy: true

snapshot_length: 512
promiscuous: true
connection_timeout: 0
bpf: ""
max_cores: 0
log_file: agent.log

# Agent configuration
agent_mode: true
agent_id: "agent-001"  # unique agent identifier
server_url: "http://192.168.1.100:3000"  # central server URL
data_interval: 60  # send data every 60 seconds
health_interval: 300  # send health check every 5 minutes
```

## Running

### 1. Starting Central Server

```bash
# Start central server
make run-manager
```

The server will be available at `http://localhost:3000`

### 2. Starting Agent

```bash
# Start agent with default configuration
make run-agent

# Start agent with command line parameters
./scripts/run_agent.sh -i agent-001 -s http://192.168.1.100:3000 -n eth0

# Start agent with root privileges (required for packet capture)
sudo ./scripts/run_agent.sh -i agent-001 -s http://192.168.1.100:3000 -n eth0
```

## Deploying to Remote Machines

### 1. Automatic Deployment

```bash
# Deploy agent to remote machine
make deploy-agent HOST=user@192.168.1.100 CONFIG=config_agent.yml

# Start agent on remote machine
ssh user@192.168.1.100 'cd ~/netmoth && sudo ./run_agent.sh'
```

### 2. Manual Deployment

1. Copy files to remote machine:
```bash
scp bin/agent config_agent.yml scripts/run_agent.sh user@192.168.1.100:~/netmoth/
```

2. Set permissions:
```bash
ssh user@192.168.1.100 'cd ~/netmoth && chmod +x run_agent.sh'
```

3. Start agent:
```bash
ssh user@192.168.1.100 'cd ~/netmoth && sudo ./run_agent.sh'
```

## API Endpoints

### Agent Registration
```
POST /api/agent/register
Content-Type: application/json

{
  "agent_id": "agent-001",
  "hostname": "machine1.example.com",
  "interface": "eth0",
  "version": "1.0.0"
}
```

### Send Data
```
POST /api/agent/data
Content-Type: application/json

{
  "agent_id": "agent-001",
  "hostname": "machine1.example.com",
  "interface": "eth0",
  "timestamp": "2024-01-01T12:00:00Z",
  "connections": [...],
  "signatures": [...],
  "stats": {
    "packets_received": 1000,
    "packets_dropped": 10,
    "packets_processed": 990,
    "connections_found": 50
  }
}
```

### Health Check
```
GET /api/agent/health?agent_id=agent-001
```

## Monitoring

### 1. Agent Logs
```bash
tail -f agent.log
```

### 2. Central Server Logs
```bash
tail -f manager.log
```

### 3. Web Interface
Open `http://localhost:3000` in your browser to access the web interface.

### 4. API Status
```bash
curl http://localhost:3000/api/version
```

## Security

### 1. Network Security
- Use HTTPS for communication between agent and server
- Configure firewall to restrict API access
- Use VPN for communication between agents and server

### 2. Authentication
- Add authentication tokens for API
- Use certificates for mutual authentication

### 3. Data Encryption
- Encrypt data in transit
- Encrypt data in database

## Troubleshooting

### 1. Agent Cannot Connect to Server
- Check server URL in configuration
- Ensure server is running and accessible
- Check network settings and firewall

### 2. Agent Cannot Capture Packets
- Ensure agent is running with root privileges
- Check that specified interface exists
- Check permissions for network interface

### 3. High Resource Usage
- Decrease `data_interval` for more frequent data sending
- Configure BPF filters to limit traffic
- Increase `snapshot_length` for more efficient capture

## Scaling

### 1. Horizontal Scaling
- Run multiple agents on different machines
- Use load balancer for API server
- Configure database replication

### 2. Vertical Scaling
- Increase number of cores in `max_cores`
- Configure buffer sizes and connection pools
- Optimize database settings

## Configuration Examples

### Agent for High-Load Server
```yaml
interface: eth0
strategy: afpacket  # use AF_PACKET for better performance
zero_copy: true
max_cores: 8
data_interval: 30  # send every 30 seconds
bpf: "port 80 or port 443"  # only HTTP/HTTPS traffic
```

### Agent for DNS Monitoring
```yaml
interface: eth0
strategy: pcap
bpf: "port 53"  # only DNS traffic
data_interval: 120  # send every 2 minutes
```

### Agent for VPN Monitoring
```yaml
interface: tun0  # VPN interface
strategy: pcap
bpf: "not port 53"  # exclude DNS traffic
data_interval: 60
``` 