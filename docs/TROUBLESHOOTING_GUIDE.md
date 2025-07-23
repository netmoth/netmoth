# Netmoth Troubleshooting Guide

## Table of Contents

1. [General Troubleshooting](#general-troubleshooting)
2. [Agent Issues](#agent-issues)
3. [Manager Issues](#manager-issues)
4. [Network Issues](#network-issues)
5. [Database Issues](#database-issues)
6. [Performance Issues](#performance-issues)
7. [Security Issues](#security-issues)
8. [Build and Installation Issues](#build-and-installation-issues)
9. [Log Analysis](#log-analysis)
10. [Getting Help](#getting-help)

## General Troubleshooting

### Basic Diagnostic Commands

```bash
# Check if Netmoth processes are running
ps aux | grep netmoth

# Check system resources
top -p $(pgrep netmoth)
free -h
df -h

# Check network interfaces
ip addr show
ifconfig

# Check system logs
journalctl -u netmoth -f
dmesg | tail -50

# Check file permissions
ls -la /var/log/netmoth/
ls -la /etc/netmoth/
```

### Common Error Messages

#### "Permission Denied"
```bash
# Solution: Set proper capabilities for packet capture
sudo setcap cap_net_raw,cap_net_admin=eip /usr/local/bin/netmoth-agent

# Or run with sudo (not recommended for production)
sudo ./netmoth-agent
```

#### "Interface Not Found"
```bash
# Check available interfaces
ip link show
ifconfig -a

# Verify interface name in configuration
cat config.yml | grep interface
```

#### "Configuration File Not Found"
```bash
# Check if config file exists
ls -la config.yml

# Verify config file path
./netmoth-agent -cfg /path/to/config.yml
```

## Agent Issues

### Agent Won't Start

#### Problem: Agent fails to start with "failed to initialize sensor"

**Symptoms:**
- Agent exits immediately after startup
- Error message: "failed to initialize sensor"
- No packet capture occurring

**Diagnosis:**
```bash
# Check configuration file
cat config.yml

# Check interface availability
ip link show

# Check permissions
ls -la /dev/bpf*
ls -la /proc/net/packet
```

**Solutions:**

1. **Interface Issues:**
   ```bash
   # Verify interface exists and is up
   ip link set eth0 up
   
   # Check interface status
   ip addr show eth0
   ```

2. **Permission Issues:**
   ```bash
   # Set capabilities for packet capture
   sudo setcap cap_net_raw,cap_net_admin=eip /usr/local/bin/netmoth-agent
   
   # Or add user to appropriate groups
   sudo usermod -a -G wireshark $USER
   ```

3. **Configuration Issues:**
   ```yaml
   # Ensure proper configuration
   interface: "eth0"  # Use correct interface name
   strategy: "pcap"   # Start with pcap for testing
   promiscuous: true
   ```

### Agent Not Capturing Packets

#### Problem: Agent starts but no packets are captured

**Symptoms:**
- Agent runs without errors
- No packet statistics in logs
- Zero packets processed

**Diagnosis:**
```bash
# Check if interface is receiving traffic
tcpdump -i eth0 -c 10

# Check BPF filter
tcpdump -i eth0 -c 10 "port 80 or port 443"

# Check system packet counters
cat /proc/net/dev | grep eth0
```

**Solutions:**

1. **Interface Configuration:**
   ```yaml
   # Ensure interface is properly configured
   interface: "eth0"
   promiscuous: true
   bpf: ""  # Remove BPF filter for testing
   ```

2. **Network Traffic:**
   ```bash
   # Generate test traffic
   curl http://example.com
   ping google.com
   
   # Check if traffic is visible
   tcpdump -i eth0 -c 5
   ```

3. **Strategy Issues:**
   ```yaml
   # Try different capture strategies
   strategy: "pcap"      # Standard libpcap
   strategy: "afpacket"  # Linux AF_PACKET
   strategy: "ebpf"      # eBPF (requires kernel support)
   ```

### Agent Performance Issues

#### Problem: High CPU usage or packet drops

**Symptoms:**
- High CPU utilization
- Packet drops reported in logs
- Poor performance under load

**Diagnosis:**
```bash
# Monitor CPU usage
top -p $(pgrep netmoth-agent)

# Check packet statistics
cat /proc/net/dev | grep eth0

# Monitor system interrupts
cat /proc/interrupts | grep eth0
```

**Solutions:**

1. **Optimize Configuration:**
   ```yaml
   # Enable zero-copy processing
   zero_copy: true
   
   # Increase snapshot length
   snapshot_length: 65536
   
   # Use multiple cores
   max_cores: 4
   ```

2. **Use High-Performance Strategy:**
   ```yaml
   # Switch to eBPF for maximum performance
   strategy: "ebpf"
   
   # Or use AF_PACKET
   strategy: "afpacket"
   ```

3. **System Tuning:**
   ```bash
   # Increase network buffer sizes
   echo 'net.core.rmem_max = 16777216' | sudo tee -a /etc/sysctl.conf
   echo 'net.core.wmem_max = 16777216' | sudo tee -a /etc/sysctl.conf
   sudo sysctl -p
   
   # Disable CPU frequency scaling
   echo performance | sudo tee /sys/devices/system/cpu/cpu*/cpufreq/scaling_governor
   ```

### Agent Communication Issues

#### Problem: Agent cannot connect to manager

**Symptoms:**
- Connection refused errors
- Agent not registered with manager
- Data not being transmitted

**Diagnosis:**
```bash
# Test network connectivity
ping manager.example.com
telnet manager.example.com 3000

# Check DNS resolution
nslookup manager.example.com

# Test HTTP connectivity
curl -v http://manager.example.com:3000/api/version
```

**Solutions:**

1. **Network Connectivity:**
   ```bash
   # Check firewall rules
   sudo iptables -L
   sudo ufw status
   
   # Test port connectivity
   nc -zv manager.example.com 3000
   ```

2. **Configuration Issues:**
   ```yaml
   # Verify manager URL
   server_url: "http://manager.example.com:3000"
   
   # Check agent ID
   agent_id: "unique-agent-id"
   ```

3. **TLS/SSL Issues:**
   ```yaml
   # For HTTPS connections
   server_url: "https://manager.example.com:3000"
   
   # Disable SSL verification for testing
   tls_insecure: true
   ```

## Manager Issues

### Manager Won't Start

#### Problem: Manager fails to start

**Symptoms:**
- Manager exits immediately
- Port already in use errors
- Database connection failures

**Diagnosis:**
```bash
# Check if port is in use
netstat -tlnp | grep :3000
lsof -i :3000

# Check database connectivity
psql -h localhost -U netmoth -d netmoth -c "SELECT 1;"

# Check configuration
cat config.yml
```

**Solutions:**

1. **Port Conflicts:**
   ```bash
   # Kill process using port 3000
   sudo fuser -k 3000/tcp
   
   # Or change port in configuration
   port: 3001
   ```

2. **Database Issues:**
   ```bash
   # Start PostgreSQL
   sudo systemctl start postgresql
   
   # Create database and user
   sudo -u postgres createdb netmoth
   sudo -u postgres createuser netmoth
   sudo -u postgres psql -c "GRANT ALL PRIVILEGES ON DATABASE netmoth TO netmoth;"
   ```

3. **Configuration Issues:**
   ```yaml
   # Verify database configuration
   postgres:
     host: "localhost"
     port: 5432
     user: "netmoth"
     password: "netmoth"
     db: "netmoth"
   ```

### Manager Performance Issues

#### Problem: Manager is slow or unresponsive

**Symptoms:**
- High response times
- Database connection timeouts
- Memory usage growing

**Diagnosis:**
```bash
# Check resource usage
top -p $(pgrep netmoth-manager)
free -h

# Check database performance
psql -h localhost -U netmoth -d netmoth -c "SELECT * FROM pg_stat_activity;"
```

**Solutions:**

1. **Database Optimization:**
   ```yaml
   # Increase connection pool
   postgres:
     max_conn: 100
     max_idlec_conn: 20
     max_lifetime_conn: 600
   ```

2. **Memory Management:**
   ```bash
   # Increase system memory
   # Add swap if needed
   sudo fallocate -l 2G /swapfile
   sudo chmod 600 /swapfile
   sudo mkswap /swapfile
   sudo swapon /swapfile
   ```

3. **Caching:**
   ```yaml
   # Enable Redis caching
   redis:
     host: "localhost"
     port: 6379
     pool_size: 20
   ```

## Network Issues

### Packet Capture Problems

#### Problem: Cannot capture packets on specific interface

**Symptoms:**
- No packets captured
- Interface errors
- Permission denied

**Diagnosis:**
```bash
# Check interface status
ip link show eth0

# Test packet capture with tcpdump
sudo tcpdump -i eth0 -c 10

# Check interface statistics
cat /proc/net/dev | grep eth0
```

**Solutions:**

1. **Interface Configuration:**
   ```bash
   # Bring interface up
   sudo ip link set eth0 up
   
   # Set promiscuous mode
   sudo ip link set eth0 promisc on
   ```

2. **Kernel Module Issues:**
   ```bash
   # Load required kernel modules
   sudo modprobe af_packet
   sudo modprobe pktgen
   ```

3. **Virtual Interface Issues:**
   ```bash
   # For virtual machines, ensure proper network configuration
   # Check VM network adapter settings
   # Ensure promiscuous mode is enabled in hypervisor
   ```

### eBPF Issues

#### Problem: eBPF strategy fails to load

**Symptoms:**
- "eBPF program load failed" errors
- Kernel version compatibility issues
- Permission denied for eBPF operations

**Diagnosis:**
```bash
# Check kernel version
uname -r

# Check eBPF support
cat /sys/kernel/debug/bpf/verifier_log

# Check eBPF filesystem
ls -la /sys/fs/bpf/
```

**Solutions:**

1. **Kernel Requirements:**
   ```bash
   # Ensure kernel version 4.18+ for full eBPF support
   # Update kernel if needed
   sudo apt update
   sudo apt install linux-generic
   ```

2. **eBPF Filesystem:**
   ```bash
   # Mount eBPF filesystem
   sudo mount -t bpf bpf /sys/fs/bpf/
   
   # Add to /etc/fstab for persistence
   echo "bpf /sys/fs/bpf bpf defaults 0 0" | sudo tee -a /etc/fstab
   ```

3. **Capabilities:**
   ```bash
   # Set required capabilities
   sudo setcap cap_bpf,cap_net_admin=eip /usr/local/bin/netmoth-agent
   ```

## Database Issues

### PostgreSQL Connection Problems

#### Problem: Cannot connect to PostgreSQL

**Symptoms:**
- "connection refused" errors
- Authentication failures
- Database not found errors

**Diagnosis:**
```bash
# Check PostgreSQL status
sudo systemctl status postgresql

# Check PostgreSQL logs
sudo tail -f /var/log/postgresql/postgresql-*.log

# Test connection
psql -h localhost -U netmoth -d netmoth
```

**Solutions:**

1. **PostgreSQL Service:**
   ```bash
   # Start PostgreSQL
   sudo systemctl start postgresql
   sudo systemctl enable postgresql
   
   # Check configuration
   sudo cat /etc/postgresql/*/main/postgresql.conf | grep listen_addresses
   ```

2. **Database Setup:**
   ```bash
   # Create database and user
   sudo -u postgres createdb netmoth
   sudo -u postgres createuser netmoth
   sudo -u postgres psql -c "ALTER USER netmoth WITH PASSWORD 'netmoth';"
   sudo -u postgres psql -c "GRANT ALL PRIVILEGES ON DATABASE netmoth TO netmoth;"
   ```

3. **Connection Configuration:**
   ```yaml
   # Verify connection settings
   postgres:
     host: "localhost"
     port: 5432
     user: "netmoth"
     password: "netmoth"
     db: "netmoth"
     ssl_mode: "disable"
   ```

### Redis Connection Problems

#### Problem: Cannot connect to Redis

**Symptoms:**
- Redis connection timeouts
- Authentication failures
- Redis service not running

**Diagnosis:**
```bash
# Check Redis status
sudo systemctl status redis

# Test Redis connection
redis-cli ping

# Check Redis configuration
sudo cat /etc/redis/redis.conf | grep bind
```

**Solutions:**

1. **Redis Service:**
   ```bash
   # Start Redis
   sudo systemctl start redis
   sudo systemctl enable redis
   
   # Check Redis logs
   sudo tail -f /var/log/redis/redis-server.log
   ```

2. **Redis Configuration:**
   ```bash
   # Edit Redis configuration
   sudo nano /etc/redis/redis.conf
   
   # Ensure bind address allows connections
   bind 127.0.0.1 ::1
   ```

3. **Connection Settings:**
   ```yaml
   # Verify Redis configuration
   redis:
     host: "localhost"
     port: 6379
     password: ""
     db: 0
   ```

## Performance Issues

### High CPU Usage

#### Problem: Netmoth is using too much CPU

**Symptoms:**
- CPU usage > 80%
- System becomes unresponsive
- High load average

**Diagnosis:**
```bash
# Monitor CPU usage
top -p $(pgrep netmoth)

# Check CPU frequency
cat /sys/devices/system/cpu/cpu*/cpufreq/scaling_governor

# Check system load
uptime
```

**Solutions:**

1. **Optimize Configuration:**
   ```yaml
   # Limit CPU usage
   max_cores: 4
   
   # Enable zero-copy
   zero_copy: true
   
   # Use efficient strategy
   strategy: "ebpf"
   ```

2. **System Tuning:**
   ```bash
   # Set CPU governor to performance
   echo performance | sudo tee /sys/devices/system/cpu/cpu*/cpufreq/scaling_governor
   
   # Increase process priority
   sudo renice -n -10 -p $(pgrep netmoth)
   ```

3. **Hardware Considerations:**
   - Use dedicated CPU cores for packet processing
   - Ensure sufficient RAM
   - Use high-performance network cards

### Memory Issues

#### Problem: High memory usage or out of memory errors

**Symptoms:**
- High memory usage
- Out of memory errors
- System swapping

**Diagnosis:**
```bash
# Check memory usage
free -h
cat /proc/meminfo

# Check process memory
ps aux | grep netmoth
cat /proc/$(pgrep netmoth)/status | grep VmRSS
```

**Solutions:**

1. **Memory Configuration:**
   ```yaml
   # Optimize memory pools
   memory_pool:
     packet_buffer_size: 32768
     packet_buffer_count: 500
     connection_buffer_size: 16384
     connection_buffer_count: 200
   ```

2. **System Memory:**
   ```bash
   # Add swap space
   sudo fallocate -l 4G /swapfile
   sudo chmod 600 /swapfile
   sudo mkswap /swapfile
   sudo swapon /swapfile
   ```

3. **Garbage Collection:**
   ```bash
   # Set Go garbage collection environment variables
   export GOGC=50
   export GOMEMLIMIT=1GiB
   ```

### Network Performance Issues

#### Problem: Packet drops or poor network performance

**Symptoms:**
- High packet drop rates
- Network timeouts
- Poor throughput

**Diagnosis:**
```bash
# Check network statistics
cat /proc/net/dev | grep eth0

# Check interface errors
ethtool -S eth0

# Monitor network usage
iftop -i eth0
```

**Solutions:**

1. **Network Tuning:**
   ```bash
   # Increase network buffer sizes
   echo 'net.core.rmem_max = 16777216' | sudo tee -a /etc/sysctl.conf
   echo 'net.core.wmem_max = 16777216' | sudo tee -a /etc/sysctl.conf
   echo 'net.core.rmem_default = 262144' | sudo tee -a /etc/sysctl.conf
   echo 'net.core.wmem_default = 262144' | sudo tee -a /etc/sysctl.conf
   sudo sysctl -p
   ```

2. **Interface Configuration:**
   ```bash
   # Optimize network interface
   sudo ethtool -G eth0 rx 4096 tx 4096
   sudo ethtool -C eth0 adaptive-rx on adaptive-tx on
   ```

3. **Hardware Optimization:**
   - Use dedicated network cards
   - Enable hardware offloading
   - Use high-speed interfaces (10Gbps+)

## Security Issues

### Permission Problems

#### Problem: Permission denied errors

**Symptoms:**
- Cannot access network interfaces
- Cannot write to log files
- Cannot bind to ports

**Diagnosis:**
```bash
# Check file permissions
ls -la /var/log/netmoth/
ls -la /etc/netmoth/

# Check user capabilities
getcap /usr/local/bin/netmoth-agent

# Check user groups
groups $USER
```

**Solutions:**

1. **File Permissions:**
   ```bash
   # Set proper permissions
   sudo chown -R netmoth:netmoth /var/log/netmoth/
   sudo chmod 755 /var/log/netmoth/
   sudo chmod 644 /var/log/netmoth/*.log
   ```

2. **Capabilities:**
   ```bash
   # Set required capabilities
   sudo setcap cap_net_raw,cap_net_admin=eip /usr/local/bin/netmoth-agent
   sudo setcap cap_bpf=eip /usr/local/bin/netmoth-agent
   ```

3. **User Configuration:**
   ```bash
   # Add user to required groups
   sudo usermod -a -G wireshark $USER
   sudo usermod -a -G netdev $USER
   ```

### TLS/SSL Issues

#### Problem: TLS certificate or SSL connection problems

**Symptoms:**
- Certificate validation errors
- SSL handshake failures
- Connection timeouts

**Diagnosis:**
```bash
# Test SSL connection
openssl s_client -connect manager.example.com:3000

# Check certificate
openssl x509 -in server.crt -text -noout

# Verify certificate chain
openssl verify server.crt
```

**Solutions:**

1. **Certificate Issues:**
   ```bash
   # Generate self-signed certificate for testing
   openssl req -x509 -newkey rsa:4096 -keyout server.key -out server.crt -days 365 -nodes
   
   # Update configuration
   tls:
     enabled: true
     cert_file: "server.crt"
     key_file: "server.key"
   ```

2. **Certificate Validation:**
   ```yaml
   # Disable certificate verification for testing
   tls_insecure: true
   
   # Or specify CA certificate
   tls:
     ca_file: "ca.crt"
   ```

## Build and Installation Issues

### Compilation Problems

#### Problem: Build fails with compilation errors

**Symptoms:**
- Go compilation errors
- Missing dependencies
- Version compatibility issues

**Diagnosis:**
```bash
# Check Go version
go version

# Check dependencies
go mod tidy
go mod download

# Check build environment
go env
```

**Solutions:**

1. **Go Version:**
   ```bash
   # Update Go to required version (1.24+)
   wget https://golang.org/dl/go1.24.linux-amd64.tar.gz
   sudo tar -C /usr/local -xzf go1.24.linux-amd64.tar.gz
   export PATH=$PATH:/usr/local/go/bin
   ```

2. **Dependencies:**
   ```bash
   # Clean and rebuild
   go clean -modcache
   go mod download
   go mod tidy
   go build ./...
   ```

3. **System Dependencies:**
   ```bash
   # Install required system packages
   sudo apt update
   sudo apt install build-essential libpcap-dev libpfring-dev
   ```

### Installation Problems

#### Problem: Installation fails or binary doesn't work

**Symptoms:**
- Installation script errors
- Binary not found
- Runtime library errors

**Diagnosis:**
```bash
# Check binary location
which netmoth-agent

# Check binary permissions
ls -la /usr/local/bin/netmoth-agent

# Check shared libraries
ldd /usr/local/bin/netmoth-agent
```

**Solutions:**

1. **Installation Path:**
   ```bash
   # Install to system path
   sudo cp bin/agent /usr/local/bin/netmoth-agent
   sudo cp bin/manager /usr/local/bin/netmoth-manager
   sudo chmod +x /usr/local/bin/netmoth-*
   ```

2. **System Dependencies:**
   ```bash
   # Install runtime dependencies
   sudo apt install libpcap0.8 libc6
   ```

3. **Manual Installation:**
   ```bash
   # Build from source
   make build
   sudo make install
   ```

## Log Analysis

### Understanding Log Messages

#### Common Log Patterns

```bash
# Agent startup logs
INFO Starting Netmoth Agent v1.0.0
INFO Initializing packet capture on interface eth0
INFO Using capture strategy: pcap
INFO Agent registered with manager

# Packet processing logs
DEBUG Processing packet: src=192.168.1.100, dst=8.8.8.8
INFO Connection established: 192.168.1.100:12345 -> 8.8.8.8:53
WARN Signature detected: IP=192.168.1.100, Type=malware

# Error logs
ERROR Failed to capture packet: permission denied
ERROR Database connection failed: connection refused
ERROR Agent registration failed: timeout
```

#### Log Level Configuration

```yaml
# Configure log levels
logging:
  level: "info"  # debug, info, warn, error
  
  # Component-specific levels
  components:
    sensor: "debug"
    analyzer: "info"
    signature: "warn"
```

### Debugging with Logs

#### Enable Debug Logging

```bash
# Set debug environment variable
export DEBUG=true
export LOG_LEVEL=debug

# Run with debug output
./netmoth-agent -debug -cfg config.yml
```

#### Log Analysis Tools

```bash
# Monitor logs in real-time
tail -f /var/log/netmoth/agent.log

# Search for specific errors
grep -i error /var/log/netmoth/agent.log

# Count log entries by level
grep -o "level=error" /var/log/netmoth/agent.log | wc -l

# Analyze log patterns
awk '/ERROR/ {print $0}' /var/log/netmoth/agent.log
```

## Getting Help

### Self-Help Resources

1. **Documentation:**
   - Check the `docs/` directory
   - Read the README.md file
   - Review configuration examples

2. **Logs and Debugging:**
   - Enable debug logging
   - Check system logs
   - Monitor resource usage

3. **Community Resources:**
   - GitHub Issues
   - GitHub Discussions
   - Stack Overflow

### Reporting Issues

When reporting issues, include:

1. **System Information:**
   ```bash
   # OS and version
   uname -a
   cat /etc/os-release
   
   # Go version
   go version
   
   # Netmoth version
   ./netmoth-agent -version
   ```

2. **Configuration:**
   ```bash
   # Configuration file (sanitized)
   cat config.yml
   ```

3. **Logs:**
   ```bash
   # Relevant log entries
   tail -100 /var/log/netmoth/agent.log
   ```

4. **Steps to Reproduce:**
   - Detailed steps to reproduce the issue
   - Expected vs actual behavior
   - Any error messages

### Contact Information

- **GitHub Issues**: https://github.com/netmoth/netmoth/issues
- **GitHub Discussions**: https://github.com/netmoth/netmoth/discussions
- **Documentation**: https://github.com/netmoth/netmoth/tree/main/docs

This troubleshooting guide provides comprehensive solutions for common Netmoth issues, enabling users to diagnose and resolve problems independently.