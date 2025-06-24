#!/bin/bash

# Netmoth Agent Runner Script
# This script runs the Netmoth agent on different machines

set -e

# Default values
CONFIG_FILE="cmd/agent/config.yml"
AGENT_ID=""
SERVER_URL=""
INTERFACE="eth0"

# Function to display usage
usage() {
    echo "Usage: $0 [OPTIONS]"
    echo "Options:"
    echo "  -c, --config FILE     Configuration file (default: cmd/agent/config.yml)"
    echo "  -i, --id ID           Agent ID (overrides config file)"
    echo "  -s, --server URL      Server URL (overrides config file)"
    echo "  -n, --interface IFACE Network interface (default: eth0)"
    echo "  -h, --help            Display this help message"
    echo ""
    echo "Examples:"
    echo "  $0 -c cmd/agent/config.yml"
    echo "  $0 -i agent-001 -s http://192.168.1.100:3000 -n eth0"
    echo ""
}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -c|--config)
            CONFIG_FILE="$2"
            shift 2
            ;;
        -i|--id)
            AGENT_ID="$2"
            shift 2
            ;;
        -s|--server)
            SERVER_URL="$2"
            shift 2
            ;;
        -n|--interface)
            INTERFACE="$2"
            shift 2
            ;;
        -h|--help)
            usage
            exit 0
            ;;
        *)
            echo "Unknown option: $1"
            usage
            exit 1
            ;;
    esac
done

# Check if config file exists
if [[ ! -f "$CONFIG_FILE" ]]; then
    echo "Error: Configuration file '$CONFIG_FILE' not found"
    exit 1
fi

# Check if agent binary exists
if [[ ! -f "./bin/agent" ]]; then
    echo "Error: Agent binary './bin/agent' not found"
    echo "Please build the agent first: make build-agent"
    exit 1
fi

# Check if running as root (required for packet capture)
if [[ $EUID -ne 0 ]]; then
    echo "Warning: This script should be run as root for packet capture"
    echo "You may need to run: sudo $0 $@"
fi

# Generate unique agent ID if not provided
if [[ -z "$AGENT_ID" ]]; then
    AGENT_ID="agent-$(hostname)-$(date +%s)"
fi

echo "Starting Netmoth Agent..."
echo "  Config file: $CONFIG_FILE"
echo "  Agent ID: $AGENT_ID"
echo "  Interface: $INTERFACE"
if [[ -n "$SERVER_URL" ]]; then
    echo "  Server URL: $SERVER_URL"
fi
echo ""

# Create temporary config file with overrides
TEMP_CONFIG=$(mktemp)
cp "$CONFIG_FILE" "$TEMP_CONFIG"

# Update config with command line overrides
if [[ -n "$AGENT_ID" ]]; then
    sed -i "s/agent_id:.*/agent_id: \"$AGENT_ID\"/" "$TEMP_CONFIG"
fi

if [[ -n "$SERVER_URL" ]]; then
    sed -i "s|server_url:.*|server_url: \"$SERVER_URL\"|" "$TEMP_CONFIG"
fi

if [[ -n "$INTERFACE" ]]; then
    sed -i "s/interface:.*/interface: $INTERFACE/" "$TEMP_CONFIG"
fi

# Run the agent
echo "Starting agent with config: $TEMP_CONFIG"
./bin/agent -cfg "$TEMP_CONFIG"

# Cleanup
rm -f "$TEMP_CONFIG" 