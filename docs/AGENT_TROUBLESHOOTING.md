# Netmoth Agent Troubleshooting Guide

This guide helps you diagnose and resolve common issues with Netmoth agents.

## Common Issues

### 1. Agent Cannot Connect to Server

**Symptoms:**
- Agent fails to register
- Data sending fails
- Connection timeout errors

**Diagnosis:**
```bash
# Check if server is running
curl http://localhost:3000/api/version

# Check network connectivity
ping <server-ip>

# Check firewall
sudo iptables -L
```

**Solutions:**
- Verify server URL in `config_agent.yml`
- Ensure server is running (`make run-manager`)
- Check firewall settings
- Verify network connectivity

### 2. Agent Cannot Capture Packets

**Symptoms:**
- No packets being processed
- Permission denied errors
- Interface not found errors

**Diagnosis:**
```bash
# Check if running as root
whoami

# Check interface exists
ip link show eth0

# Check interface permissions
ls -la /dev/eth0
```

**Solutions:**
- Run agent with root privileges: `sudo ./run_agent.sh`
- Verify interface name in configuration
- Check interface permissions
- Ensure interface is up: `sudo ip link set eth0 up`

### 3. High Resource Usage

**Symptoms:**
- High CPU usage
- High memory usage
- System becomes unresponsive

**Diagnosis:**
```bash
# Check resource usage
top -p $(pgrep agent)

# Check memory usage
ps aux | grep agent

# Check disk I/O
iotop
```

**Solutions:**
- Reduce `data_interval` in configuration
- Configure BPF filters to limit traffic
- Increase `snapshot_length`
- Optimize `max_cores` setting

### 4. Data Not Being Sent

**Symptoms:**
- No data in server logs
- Empty connections buffer
- No network activity

**Diagnosis:**
```bash
# Check agent logs
tail -f agent.log

# Check if data is being captured
tcpdump -i eth0 -c 10

# Check agent configuration
cat config_agent.yml
```

**Solutions:**
- Verify `agent_mode: true` in configuration
- Check `server_url` is correct
- Ensure network traffic exists
- Check BPF filters

## Log Analysis

### Agent Logs

**Location:** `agent.log`

**Common Log Messages:**
```
Agent agent-001 registered successfully
Data sent successfully to server: 50 connections, 5 signatures
Failed to send data to server: connection refused
Health check successful: Agent agent-001 is healthy
```

### Server Logs

**Location:** `manager.log`

**Common Log Messages:**
```
Received data from agent agent-001: 50 connections, 5 signatures
Agent registration: agent-001 (machine1) on interface eth0, version 1.0.0
```

## Performance Tuning

### 1. Optimize Packet Capture

```yaml
# For high-performance servers
strategy: afpacket
zero_copy: true
max_cores: 8
snapshot_length: 2048
```

### 2. Optimize Data Sending

```yaml
# Send data more frequently
data_interval: 30

# Send health checks less frequently
health_interval: 600
```

### 3. Filter Traffic

```yaml
# Only capture specific traffic
bpf: "port 80 or port 443 or port 53"
```

## Network Issues

### 1. Firewall Configuration

**Allow agent traffic:**
```bash
# Allow outbound HTTP/HTTPS
sudo iptables -A OUTPUT -p tcp --dport 3000 -j ACCEPT

# Allow agent registration
sudo iptables -A INPUT -p tcp --dport 3000 -j ACCEPT
```

### 2. DNS Resolution

**Check DNS:**
```bash
nslookup <server-hostname>
dig <server-hostname>
```

### 3. Proxy Configuration

If using a proxy, configure the agent client:

```go
// In agent_client.go
proxyURL, _ := url.Parse("http://proxy:8080")
client.Transport = &http.Transport{
    Proxy: http.ProxyURL(proxyURL),
}
```

## Security Issues

### 1. Certificate Errors

**Symptoms:**
- SSL/TLS handshake failures
- Certificate validation errors

**Solutions:**
- Use HTTP for testing (not recommended for production)
- Configure proper SSL certificates
- Add certificate to trusted store

### 2. Authentication Issues

**Symptoms:**
- 401 Unauthorized errors
- Authentication token failures

**Solutions:**
- Implement proper authentication
- Use API tokens
- Configure mutual TLS

## Monitoring and Alerts

### 1. Health Monitoring

Create a monitoring script:

```bash
#!/bin/bash
# monitor_agent.sh

AGENT_ID="agent-001"
SERVER_URL="http://localhost:3000"

# Check agent health
response=$(curl -s "$SERVER_URL/api/agent/health?agent_id=$AGENT_ID")
if [[ $response == *"healthy"* ]]; then
    echo "Agent is healthy"
else
    echo "Agent health check failed"
    # Send alert
fi
```

### 2. Resource Monitoring

Monitor system resources:

```bash
#!/bin/bash
# monitor_resources.sh

# Check CPU usage
cpu_usage=$(top -bn1 | grep "Cpu(s)" | awk '{print $2}' | cut -d'%' -f1)
if (( $(echo "$cpu_usage > 80" | bc -l) )); then
    echo "High CPU usage: $cpu_usage%"
fi

# Check memory usage
mem_usage=$(free | grep Mem | awk '{printf("%.2f", $3/$2 * 100.0)}')
if (( $(echo "$mem_usage > 80" | bc -l) )); then
    echo "High memory usage: $mem_usage%"
fi
```

## Recovery Procedures

### 1. Agent Restart

```bash
# Stop agent
sudo pkill -f agent

# Start agent
sudo ./scripts/run_agent.sh
```

### 2. Server Restart

```bash
# Stop server
pkill -f manager

# Start server
make run-manager
```

### 3. Configuration Reset

```bash
# Backup current config
cp config_agent.yml config_agent.yml.backup

# Reset to defaults
cp config_agent.yml.example config_agent.yml

# Edit configuration
nano config_agent.yml
```

## Getting Help

### 1. Enable Debug Logging

```bash
# Set debug level
export NETMOTH_LOG_LEVEL=debug

# Run agent with debug
sudo NETMOTH_LOG_LEVEL=debug ./bin/agent -cfg config_agent.yml
```

### 2. Collect Diagnostic Information

```bash
#!/bin/bash
# collect_diagnostics.sh

echo "=== Netmoth Agent Diagnostics ==="
echo "Date: $(date)"
echo "Agent Version: $(./bin/agent -version 2>/dev/null || echo 'Unknown')"
echo "System: $(uname -a)"
echo "Network Interfaces:"
ip link show
echo "Configuration:"
cat config_agent.yml
echo "Recent Logs:"
tail -20 agent.log
```

### 3. Common Commands

```bash
# Check agent status
ps aux | grep agent

# Check server status
ps aux | grep manager

# Check network interfaces
ip link show

# Check listening ports
netstat -tlnp

# Check firewall rules
sudo iptables -L
```

## Prevention

### 1. Regular Maintenance

- Monitor logs regularly
- Update configurations as needed
- Test connectivity periodically
- Backup configurations

### 2. Monitoring Setup

- Set up health checks
- Configure resource monitoring
- Implement alerting
- Regular performance reviews

### 3. Documentation

- Document configuration changes
- Keep troubleshooting notes
- Update runbooks
- Share knowledge with team 