# Netmoth Agent - Distributed Network Monitoring System

## Overview

A distributed network monitoring system has been implemented where agents can run on different machines and send audit data to a central server.

## What's Implemented

### 1. Agent
- **Data Collection**: Network traffic capture and analysis
- **Buffering**: Local data accumulation before sending
- **Data Transmission**: Periodic sending to central server
- **Registration**: Automatic registration with server
- **Health Checks**: Periodic health monitoring

### 2. Central Server (Manager)
- **REST API**: Receiving data from agents
- **Agent Registration**: Managing connected agents
- **Web Interface**: Monitoring and management
- **Database**: Storing audit data

### 3. API Endpoints
- `POST /api/agent/register` - Agent registration
- `POST /api/agent/data` - Sending audit data
- `GET /api/agent/health` - Agent health check

## Quick Start

### 1. Building
```bash
# Build agent and server
make build-agent
make build-manager
```

### 2. Configuration
```bash
# Create agent configuration
cp cmd/agent/config.yml.example cmd/agent/config.yml
# Edit cmd/agent/config.yml

# Create server configuration
cp cmd/manager/config.yml.example cmd/manager/config.yml
# Edit cmd/manager/config.yml
```

### 3. Running
```bash
# Start central server
make run-manager

# In another terminal - start agent
make run-agent
```

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

## Key Files

### Agent
- `cmd/agent/main.go` - Agent entry point
- `cmd/agent/config.yml` - Agent configuration
- `internal/sensor/agent_client.go` - Client for sending data
- `internal/sensor/sensor.go` - Main sensor logic with agent support

### Server
- `cmd/manager/main.go` - Server entry point
- `cmd/manager/config.yml` - Server configuration
- `internal/web/agent_api.go` - API for agents
- `internal/web/main.go` - Web server with API endpoints

### Utilities
- `scripts/run_agent.sh` - Agent startup script
- `scripts/test_agent.sh` - Agent testing
- `Makefile` - Build and deployment commands

## Agent Configuration

```yaml
# Agent configuration
agent_mode: true
agent_id: "agent-001"  # unique identifier
server_url: "http://192.168.1.100:3000"  # server URL
data_interval: 60  # send data every 60 seconds
health_interval: 300  # health check every 5 minutes
```

## Deploying to Remote Machines

```bash
# Automatic deployment
make deploy-agent HOST=user@192.168.1.100 CONFIG=config_agent.yml

# Manual deployment
scp bin/agent config_agent.yml scripts/run_agent.sh user@192.168.1.100:~/netmoth/
ssh user@192.168.1.100 'cd ~/netmoth && sudo ./run_agent.sh'
```

## Monitoring

- **Agent Logs**: `tail -f agent.log`
- **Server Logs**: `tail -f manager.log`
- **Web Interface**: `http://localhost:3000`
- **API Status**: `curl http://localhost:3000/api/version`

## Security

- Use HTTPS for communication between agent and server
- Configure firewall to restrict API access
- Use VPN for communication between agents and server
- Add authentication tokens for API

## Scaling

- **Horizontal**: Run multiple agents on different machines
- **Vertical**: Increase number of cores and optimize settings

## Documentation

Detailed documentation is available in `docs/AGENT_DEPLOYMENT.md`

## Testing

```bash
# Run tests
./scripts/test_agent.sh
```

## Support

For help and additional information, refer to the documentation in the `docs/` folder. 