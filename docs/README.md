# Netmoth Documentation

Welcome to the Netmoth documentation. This guide will help you understand and use the Netmoth network monitoring system.

## Table of Contents

### Getting Started
- [Agent Overview](README_AGENT.md) - Quick overview of the distributed agent system
- [Agent Deployment Guide](AGENT_DEPLOYMENT.md) - Complete deployment and configuration guide
- [Agent API Reference](AGENT_API_REFERENCE.md) - API documentation for agent communication
- [Agent Troubleshooting](AGENT_TROUBLESHOOTING.md) - Common issues and solutions

### Core Documentation
- [eBPF Support](EBPF_SUPPORT.md) - eBPF packet capture implementation
- [Performance Optimizations](PERFORMANCE_OPTIMIZATIONS.md) - Performance tuning guide
- [eBPF README](README_EBPF.md) - eBPF-specific documentation
- [Optimizations README](README_OPTIMIZATIONS.md) - Optimization techniques

## Quick Start

### 1. Build the System
```bash
# Build agent and manager
make build-agent
make build-manager
```

### 2. Configure
```bash
# Copy and edit configuration files
cp config_agent.yml.example config_agent.yml
cp config_example.yml config.yml
```

### 3. Run
```bash
# Start central server
make run-manager

# Start agent (in another terminal)
make run-agent
```

## Architecture Overview

Netmoth is a distributed network monitoring system consisting of:

- **Agents**: Collect network traffic on remote machines
- **Central Server**: Receives and processes data from agents
- **Web Interface**: Provides monitoring and management interface
- **Database**: Stores collected network data and signatures

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

## Key Features

### Agent Features
- **Packet Capture**: Multiple capture strategies (pcap, afpacket, eBPF)
- **Traffic Analysis**: HTTP, HTTPS, DNS, TLS analysis
- **Signature Detection**: Malware and threat detection
- **Data Buffering**: Local buffering before transmission
- **Health Monitoring**: Periodic health checks
- **Automatic Registration**: Self-registration with central server

### Server Features
- **REST API**: HTTP API for agent communication
- **Agent Management**: Registration and status tracking
- **Data Storage**: PostgreSQL database integration
- **Web Interface**: Vue.js-based monitoring dashboard
- **Real-time Monitoring**: Live data visualization

## Configuration

### Agent Configuration
```yaml
agent_mode: true
agent_id: "agent-001"
server_url: "http://192.168.1.100:3000"
data_interval: 60
health_interval: 300
```

### Server Configuration
```yaml
interface: eth0
strategy: pcap
postgres:
  host: localhost:5432
  user: netmoth
  password: password
  db: netmoth
```

## API Endpoints

- `POST /api/agent/register` - Agent registration
- `POST /api/agent/data` - Send audit data
- `GET /api/agent/health` - Health check
- `GET /api/version` - Server version

## Deployment

### Local Deployment
```bash
# Build and run locally
make build-agent build-manager
make run-manager
make run-agent
```

### Remote Deployment
```bash
# Deploy to remote machine
make deploy-agent HOST=user@192.168.1.100 CONFIG=config_agent.yml

# Run on remote machine
ssh user@192.168.1.100 'cd ~/netmoth && sudo ./run_agent.sh'
```

## Monitoring

### Logs
- Agent logs: `tail -f agent.log`
- Server logs: `tail -f manager.log`

### Web Interface
- Access: `http://localhost:3000`
- Features: Real-time monitoring, agent status, data visualization

### Health Checks
```bash
# Check agent health
curl "http://localhost:3000/api/agent/health?agent_id=agent-001"

# Check server status
curl http://localhost:3000/api/version
```

## Security

### Network Security
- Use HTTPS for agent-server communication
- Configure firewall rules
- Use VPN for secure communication

### Authentication
- Implement API tokens
- Use mutual TLS authentication
- Configure proper access controls

## Performance

### Optimization Tips
- Use eBPF for high-performance capture
- Configure appropriate buffer sizes
- Optimize database queries
- Use BPF filters to limit traffic

### Scaling
- **Horizontal**: Deploy multiple agents
- **Vertical**: Optimize resource usage
- **Load Balancing**: Use load balancers for API

## Troubleshooting

### Common Issues
1. **Connection Problems**: Check network connectivity and firewall
2. **Permission Errors**: Run agent with root privileges
3. **High Resource Usage**: Optimize configuration settings
4. **Data Not Sending**: Verify agent mode and server URL

### Getting Help
- Check the [Troubleshooting Guide](AGENT_TROUBLESHOOTING.md)
- Review logs for error messages
- Test connectivity with curl commands
- Enable debug logging for detailed information

## Development

### Building from Source
```bash
# Clone repository
git clone <repository-url>
cd netmoth

# Build with optimizations
make build-optimized

# Build with eBPF support
make build-ebpf
```

### Testing
```bash
# Run tests
./scripts/test_agent.sh

# Run specific tests
go test ./internal/sensor/...
```

### Contributing
1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## Support

For additional help and support:

1. Check the documentation in this folder
2. Review the troubleshooting guide
3. Check GitHub issues
4. Contact the development team

## License

This project is licensed under the terms specified in the LICENSE file. 