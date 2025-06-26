<p align="center">
    <a href="https://netmoth.com" target="_blank" rel="noopener">
        <img src="https://github.com/netmoth/.github/raw/main/img/logo.svg" alt="A lightweight, fast, simple and complete solution for traffic analysis and intrusion detection" width="50%" />
    </a>
</p>

<p align="center">
    <a href="https://github.com/netmoth/netmoth/releases">
    <img src="https://img.shields.io/github/v/release/netmoth/netmoth?sort=semver&label=Release&color=651FFF" />
    </a>
    &nbsp;
    <a href="/LICENSE"><img src="https://img.shields.io/badge/license-Apache--2.0-green.svg"></a>
    &nbsp;
    <a href="https://goreportcard.com/report/github.com/netmoth/netmoth"><img src="https://goreportcard.com/badge/github.com/netmoth/netmoth"></a>
    &nbsp;
    <a href="https://www.codefactor.io/repository/github/netmoth/netmoth"><img src="https://www.codefactor.io/repository/github/netmoth/netmoth/badge" alt="CodeFactor" /></a>
    &nbsp;
    <a href="https://github.com/netmoth/netmoth"><img src="https://img.shields.io/badge/backend-go-orange.svg"></a>
    &nbsp;
    <a href="https://github.com/netmoth/netmoth/blob/main/go.mod"><img src="https://img.shields.io/github/go-mod/go-version/netmoth/netmoth?color=7fd5ea"></a>
</p>

---

## 🚀&nbsp;&nbsp;What is Netmoth?

**Netmoth** is a high-performance, distributed network monitoring and intrusion detection system designed for modern network security challenges. Built with Go and featuring advanced packet capture technologies, Netmoth provides real-time traffic analysis, threat detection, and comprehensive network visibility across your entire infrastructure.

### 🎯&nbsp;&nbsp;Who is Netmoth for?

**🔒 Security Professionals & SOC Teams**
- Monitor network traffic in real-time across multiple locations
- Detect malware, botnets, and suspicious activities automatically
- Get instant alerts on security threats and anomalies
- Maintain compliance with security standards and regulations

**🏢 Enterprise IT Teams**
- Gain complete visibility into network traffic patterns
- Monitor distributed infrastructure with lightweight agents
- Identify performance bottlenecks and network issues
- Ensure network security across multiple sites and cloud environments

**🛡️ Network Administrators**
- Deploy agents on critical network segments for continuous monitoring
- Analyze HTTP, HTTPS, DNS, and TLS traffic in detail
- Use advanced packet capture (eBPF, PF_RING, AF_PACKET) for maximum performance
- Get detailed insights into network behavior and traffic flows

**🔬 Security Researchers & Penetration Testers**
- Capture and analyze network traffic for security assessments
- Detect and track malicious activities in real-time
- Use signature-based detection for known threats
- Integrate with existing security tools and workflows

**☁️ DevOps & Cloud Engineers**
- Monitor containerized environments and microservices
- Deploy lightweight agents across cloud infrastructure
- Get real-time visibility into cloud network traffic
- Ensure security compliance in dynamic environments

### ✨&nbsp;&nbsp;Why Choose Netmoth?

- **🚀 High Performance**: Zero-copy packet processing with eBPF support for maximum throughput
- **🔧 Easy Deployment**: Minimal configuration, automatic agent registration, and simple setup
- **📊 Real-time Monitoring**: Live traffic analysis with instant threat detection and alerts
- **🌐 Distributed Architecture**: Scale across multiple locations with lightweight agents
- **🛡️ Advanced Detection**: Signature-based malware detection, botnet identification, and threat intelligence
- **📈 Modern Web Interface**: Beautiful Vue.js dashboard with real-time data visualization
- **🔍 Deep Analysis**: HTTP/2, TLS, DNS traffic analysis with automatic TCP stream reassembly
- **⚡ Low Resource Usage**: Optimized for production environments with minimal system impact

> ⚠️&nbsp;&nbsp;Current major version is zero (`v0.x.x`) to accommodate rapid development and fast iteration while getting early feedback from users. Please keep in mind that netmoth is still under active development and therefore full backward compatibility is not guaranteed before reaching v1.0.0.

## 🏆&nbsp;&nbsp;Features

- [x] Monitors traffic on all interfaces
- [x] Minimal configuration
- [x] PCAP
- [x] AF_PACKET 
- [x] PF_RING
- [x] eBPF
- [x] Zero copy packet processing
- [x] Automatic TCP stream reassembly
- [x] Berkeley Packet Filter
- [x] Check IP on blocklist
- [x] Checking botnet on blocklist
- [ ] Checking certificate on blocklist
- [x] Checking tracker on blocklist
- [ ] Web-interface
- [ ] Rules
- [x] Agents

## 📚&nbsp;&nbsp;Documentation

- [eBPF Support](docs/EBPF_SUPPORT.md) - Detailed information about eBPF support
- [eBPF README](docs/README_EBPF.md) - Guide to using eBPF
- [Performance Optimizations](docs/PERFORMANCE_OPTIMIZATIONS.md) - Performance optimization details
- [Optimizations README](docs/README_OPTIMIZATIONS.md) - Optimization usage guide

## 🏁&nbsp;&nbsp;Installation

Simple agent installation
```bash
mkdir netmoth
cd ./netmoth
curl -L https://raw.githubusercontent.com/netmoth/netmoth/main/config_example.yml > config.yml
curl -L https://github.com/netmoth/netmoth/releases/latest/download/netmoth_agent_Linux_x86_64 > netmoth_agent
sudo chmod u+x netmoth_agent
```

if necessary, make changes to the `config.yml` file, then run the agent
```bash
./netmoth_agent
```


## 👑&nbsp;&nbsp;Community

... coming soon ...

## 👍&nbsp;&nbsp;Contribute

We would for you to get involved with netmoth development! If you want to say **thank you** and/or support the active development of `netmoth`:

1. Add a [GitHub Star](https://github.com/netmoth/netmoth/stargazers) to the project.
2. Tweet about the project [on your Twitter](https://twitter.com/intent/tweet?text=%F0%9F%9A%80%20netmoth%20-%20a%20lightweight%2C%20fast%2C%20simple%20and%20complete%20solution%20for%20traffic%20analysis%20and%20intrusion%20detection%20on%20%23Go%20https%3A//github.com/netmoth/netmoth).
3. Write a review or tutorial on [Medium](https://medium.com/), [Dev.to](https://dev.to/) or personal blog.

You can learn more about how you can contribute to this project in the [contribution guide](CONTRIBUTING.md).

## 🚨&nbsp;&nbsp;Security

... coming soon ...

## Quick Start

### 1. Build the System
```bash
# Build agent and manager
make build-agent
make build-manager
```

### 2. Configure
```bash
# For agent (choose one):
cp cmd/agent/config.yml.example cmd/agent/config.yml          # Basic
cp cmd/agent/config_optimized.yml cmd/agent/config.yml        # Optimized
cp cmd/agent/config_ebpf.yml cmd/agent/config.yml             # eBPF

# For manager (choose one):
cp cmd/manager/config.yml.example cmd/manager/config.yml      # Basic
cp cmd/manager/config_optimized.yml cmd/manager/config.yml    # Optimized
cp cmd/manager/config_ebpf.yml cmd/manager/config.yml         # eBPF
```

### 3. Run
```bash
# Start central server
make run-manager

# Start agent (in another terminal)
make run-agent
```

## Configuration Options

### Agent Configurations
- **Basic** (`config.yml.example`): Standard configuration for most environments
- **Optimized** (`config_optimized.yml`): High-performance settings for busy networks
- **eBPF** (`config_ebpf.yml`): Maximum performance using eBPF packet capture

### Manager Configurations
- **Basic** (`config.yml.example`): Standard configuration with database
- **Optimized** (`config_optimized.yml`): High-performance settings with optimized database
- **eBPF** (`config_ebpf.yml`): Maximum performance using eBPF with database

## Quick Commands

```bash
# Build
make build-agent build-manager

# Run with specific configurations
make run-agent-optimized    # Agent with optimized config
make run-agent-ebpf         # Agent with eBPF config
make run-manager-optimized  # Manager with optimized config
make run-manager-ebpf       # Manager with eBPF config

# Deploy agent to remote machine
make deploy-agent HOST=user@192.168.1.100 CONFIG=cmd/agent/config.yml
```

