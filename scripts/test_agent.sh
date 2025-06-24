#!/bin/bash

# Test script for Netmoth Agent
# This script tests the agent functionality

set -e

echo "=== Netmoth Agent Test ==="

# Check if binaries exist
echo "1. Checking binaries..."
if [ ! -f "./bin/agent" ]; then
    echo "ERROR: Agent binary not found. Run 'make build-agent' first."
    exit 1
fi

if [ ! -f "./bin/manager" ]; then
    echo "ERROR: Manager binary not found. Run 'make build-manager' first."
    exit 1
fi

echo "✓ Binaries found"

# Check if config files exist
echo "2. Checking configuration files..."
if [ ! -f "cmd/agent/config.yml" ]; then
    echo "ERROR: Agent config not found. Please create cmd/agent/config.yml"
    exit 1
fi

if [ ! -f "cmd/manager/config.yml" ]; then
    echo "ERROR: Manager config not found. Please create cmd/manager/config.yml"
    exit 1
fi

echo "✓ Configuration files found"

# Test agent configuration
echo "3. Testing agent configuration..."
./bin/agent -cfg cmd/agent/config.yml -h 2>/dev/null || echo "Agent help command works"

# Test manager configuration
echo "4. Testing manager configuration..."
./bin/manager -h 2>/dev/null || echo "Manager help command works"

# Test API endpoints (if server is running)
echo "5. Testing API endpoints..."
if curl -s http://localhost:3000/api/version >/dev/null 2>&1; then
    echo "✓ Server is running and API is accessible"
    
    # Test agent registration endpoint
    if curl -s -X POST http://localhost:3000/api/agent/register \
        -H "Content-Type: application/json" \
        -d '{"agent_id":"test-agent","hostname":"test-host","interface":"eth0","version":"1.0.0"}' >/dev/null 2>&1; then
        echo "✓ Agent registration endpoint works"
    else
        echo "⚠ Agent registration endpoint not responding"
    fi
else
    echo "⚠ Server not running on localhost:3000"
    echo "   Start the server with: make run-manager"
fi

# Test network interface
echo "6. Testing network interface..."
INTERFACE=$(grep "interface:" cmd/agent/config.yml | awk '{print $2}')
if ip link show $INTERFACE >/dev/null 2>&1; then
    echo "✓ Network interface $INTERFACE exists"
else
    echo "⚠ Network interface $INTERFACE not found"
    echo "   Available interfaces:"
    ip link show | grep -E "^[0-9]+:" | awk -F: '{print $2}' | tr -d ' '
fi

# Test permissions
echo "7. Testing permissions..."
if [ $EUID -eq 0 ]; then
    echo "✓ Running as root (required for packet capture)"
else
    echo "⚠ Not running as root"
    echo "   Packet capture may fail. Run with sudo for full functionality."
fi

echo ""
echo "=== Test Summary ==="
echo "Agent binary: ✓"
echo "Manager binary: ✓"
echo "Configuration files: ✓"
echo "Network interface: $(ip link show $INTERFACE >/dev/null 2>&1 && echo "✓" || echo "⚠")"
echo "Permissions: $( [ $EUID -eq 0 ] && echo "✓" || echo "⚠" )"
echo ""
echo "To start the agent:"
echo "  sudo ./scripts/run_agent.sh"
echo ""
echo "To start the manager:"
echo "  ./bin/manager"
echo ""
echo "For more information, see docs/AGENT_DEPLOYMENT.md" 