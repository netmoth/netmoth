# Netmoth Documentation

Welcome to the comprehensive documentation for Netmoth, a high-performance network monitoring and intrusion detection system.

## 📚 Documentation Overview

This documentation is organized to help you quickly find the information you need, whether you're a new user, developer, or system administrator.

## 🚀 Quick Start

### For New Users
1. **Installation**: Follow the [Installation Guide](README.md#installation) in the main README
2. **Basic Configuration**: See [Configuration Examples](CONFIGURATION_REFERENCE.md#configuration-examples)
3. **First Run**: Use the [Quick Start Guide](README.md#quick-start) in the main README

### For Developers
1. **Setup**: Follow the [Development Guide](DEVELOPMENT_GUIDE.md#development-environment-setup)
2. **Building**: See [Building the Project](DEVELOPMENT_GUIDE.md#building-the-project)
3. **Testing**: Refer to [Testing](DEVELOPMENT_GUIDE.md#testing)

### For System Administrators
1. **Deployment**: Review [Architecture Documentation](ARCHITECTURE.md#deployment-models)
2. **Configuration**: See [Configuration Reference](CONFIGURATION_REFERENCE.md)
3. **Monitoring**: Check [Troubleshooting Guide](TROUBLESHOOTING_GUIDE.md#monitoring-and-observability)

## 📖 Documentation Index

### Core Documentation

| Document | Description | Audience |
|----------|-------------|----------|
| **[API Documentation](API_DOCUMENTATION.md)** | Complete API reference for all components | Developers |
| **[Architecture Documentation](ARCHITECTURE.md)** | System design, components, and data flow | Architects, Developers |
| **[Development Guide](DEVELOPMENT_GUIDE.md)** | Setup, building, testing, and contributing | Developers |
| **[Configuration Reference](CONFIGURATION_REFERENCE.md)** | All configuration options and examples | Users, Administrators |
| **[Troubleshooting Guide](TROUBLESHOOTING_GUIDE.md)** | Common issues and solutions | Users, Administrators |

### Specialized Documentation

| Document | Description | Audience |
|----------|-------------|----------|
| **[eBPF Support](EBPF_SUPPORT.md)** | eBPF packet capture implementation | Developers, Advanced Users |
| **[Performance Optimizations](PERFORMANCE_OPTIMIZATIONS.md)** | Performance tuning and optimization | Administrators, Developers |
| **[Agent Deployment](AGENT_DEPLOYMENT.md)** | Agent deployment and management | Administrators |
| **[Agent API Reference](AGENT_API_REFERENCE.md)** | Agent API endpoints and usage | Developers |
| **[Agent Troubleshooting](AGENT_TROUBLESHOOTING.md)** | Agent-specific issues and solutions | Users, Administrators |

### Quick Reference Guides

| Document | Description | Audience |
|----------|-------------|----------|
| **[README eBPF](README_EBPF.md)** | Quick eBPF usage guide | Users |
| **[README Optimizations](README_OPTIMIZATIONS.md)** | Quick optimization guide | Users |
| **[README Agent](README_AGENT.md)** | Quick agent usage guide | Users |

## 🎯 Documentation by Use Case

### Getting Started
- **New to Netmoth?** Start with the [main README](../README.md)
- **Need to install?** Follow the [Installation Guide](../README.md#installation)
- **Want to configure?** See [Configuration Reference](CONFIGURATION_REFERENCE.md)

### Development
- **Setting up development environment?** See [Development Guide](DEVELOPMENT_GUIDE.md)
- **Need API reference?** Check [API Documentation](API_DOCUMENTATION.md)
- **Want to contribute?** Follow [Contributing Guidelines](DEVELOPMENT_GUIDE.md#contributing-guidelines)

### Deployment and Operations
- **Planning deployment?** Review [Architecture Documentation](ARCHITECTURE.md)
- **Configuring for production?** See [Configuration Reference](CONFIGURATION_REFERENCE.md)
- **Troubleshooting issues?** Check [Troubleshooting Guide](TROUBLESHOOTING_GUIDE.md)

### Performance and Optimization
- **Need maximum performance?** See [Performance Optimizations](PERFORMANCE_OPTIMIZATIONS.md)
- **Using eBPF?** Check [eBPF Support](EBPF_SUPPORT.md)
- **Optimizing configuration?** Review [README Optimizations](README_OPTIMIZATIONS.md)

## 🔧 System Components

### Agent Component
The agent is responsible for packet capture and analysis:

- **Packet Capture**: Multiple strategies (PCAP, AF_PACKET, PF_RING, eBPF)
- **Protocol Analysis**: HTTP, HTTPS, DNS, TLS analysis
- **Signature Detection**: Real-time threat detection
- **Data Transmission**: Sends data to central manager

**Key Documentation:**
- [Agent Deployment](AGENT_DEPLOYMENT.md)
- [Agent API Reference](AGENT_API_REFERENCE.md)
- [Agent Troubleshooting](AGENT_TROUBLESHOOTING.md)

### Manager Component
The manager provides central coordination and data aggregation:

- **Web Interface**: RESTful API and web dashboard
- **Agent Management**: Registration and health monitoring
- **Data Storage**: PostgreSQL and Redis integration
- **Signature Management**: Threat database management

**Key Documentation:**
- [API Documentation](API_DOCUMENTATION.md#web-package)
- [Architecture Documentation](ARCHITECTURE.md#manager-component)

## 📋 Configuration Quick Reference

### Essential Configuration Options

```yaml
# Agent Configuration
interface: "eth0"              # Network interface
strategy: "pcap"              # Capture strategy
agent_mode: true              # Enable agent mode
agent_id: "agent-001"         # Unique agent ID
server_url: "http://localhost:3000"  # Manager URL

# Performance Settings
zero_copy: true               # Enable zero-copy processing
max_cores: 0                  # Auto-detect CPU cores
snapshot_length: 512          # Packet capture length

# Database Configuration
postgres:
  host: "localhost"
  port: 5432
  user: "netmoth"
  password: "netmoth"
  db: "netmoth"
```

### Common Configuration Scenarios

| Scenario | Configuration File | Description |
|----------|-------------------|-------------|
| **Development** | `config.yml.example` | Basic configuration for development |
| **High Performance** | `config_optimized.yml` | Optimized for high-throughput networks |
| **eBPF** | `config_ebpf.yml` | Maximum performance using eBPF |

## 🚨 Troubleshooting Quick Reference

### Common Issues and Solutions

| Issue | Quick Fix | Full Documentation |
|-------|-----------|-------------------|
| **Permission Denied** | `sudo setcap cap_net_raw,cap_net_admin=eip /usr/local/bin/netmoth-agent` | [Troubleshooting Guide](TROUBLESHOOTING_GUIDE.md#permission-problems) |
| **Interface Not Found** | Check interface name with `ip link show` | [Network Issues](TROUBLESHOOTING_GUIDE.md#packet-capture-problems) |
| **Database Connection** | Start PostgreSQL: `sudo systemctl start postgresql` | [Database Issues](TROUBLESHOOTING_GUIDE.md#postgresql-connection-problems) |
| **High CPU Usage** | Enable zero-copy: `zero_copy: true` | [Performance Issues](TROUBLESHOOTING_GUIDE.md#high-cpu-usage) |

### Diagnostic Commands

```bash
# Check if Netmoth is running
ps aux | grep netmoth

# Check system resources
top -p $(pgrep netmoth)

# Check network interfaces
ip addr show

# Check logs
tail -f /var/log/netmoth/agent.log
```

## 🔗 Related Resources

### External Documentation
- **[Go Documentation](https://golang.org/doc/)**: Go programming language
- **[gopacket Documentation](https://pkg.go.dev/github.com/google/gopacket)**: Packet processing library
- **[eBPF Documentation](https://ebpf.io/)**: Extended Berkeley Packet Filter
- **[PostgreSQL Documentation](https://www.postgresql.org/docs/)**: Database system

### Community Resources
- **[GitHub Repository](https://github.com/netmoth/netmoth)**: Source code and issues
- **[GitHub Discussions](https://github.com/netmoth/netmoth/discussions)**: Community discussions
- **[GitHub Issues](https://github.com/netmoth/netmoth/issues)**: Bug reports and feature requests

## 📝 Contributing to Documentation

We welcome contributions to improve the documentation:

1. **Report Issues**: Create an issue for documentation problems
2. **Submit Improvements**: Fork the repository and submit a pull request
3. **Suggest Topics**: Use GitHub Discussions to suggest new documentation topics

### Documentation Standards

- Use clear, concise language
- Include practical examples
- Provide step-by-step instructions
- Include troubleshooting information
- Keep information up-to-date

## 📞 Getting Help

### Self-Help Resources
1. **Check this documentation** for your specific use case
2. **Search existing issues** on GitHub
3. **Review troubleshooting guides** for common problems

### Community Support
1. **GitHub Discussions**: For questions and general help
2. **GitHub Issues**: For bug reports and feature requests
3. **Documentation**: This comprehensive documentation suite

### Reporting Issues
When reporting issues, please include:
- System information (OS, version, etc.)
- Configuration files (sanitized)
- Relevant logs
- Steps to reproduce the issue

---

**Note**: This documentation is continuously updated. For the latest version, always check the [GitHub repository](https://github.com/netmoth/netmoth/tree/main/docs).

**Version**: This documentation corresponds to Netmoth v0.x.x 