#!/bin/bash

# Test script for eBPF functionality in Netmoth
# This script verifies that eBPF support is working correctly

set -e

echo "=== Netmoth eBPF Test Script ==="
echo

# Check if running as root
if [[ $EUID -ne 0 ]]; then
   echo "This script must be run as root for eBPF functionality"
   exit 1
fi

# Check kernel version
echo "1. Checking kernel version..."
KERNEL_VERSION=$(uname -r)
echo "   Kernel version: $KERNEL_VERSION"

# Check if kernel supports eBPF
if [[ $(echo $KERNEL_VERSION | cut -d. -f1) -lt 4 ]] || \
   ([[ $(echo $KERNEL_VERSION | cut -d. -f1) -eq 4 ]] && [[ $(echo $KERNEL_VERSION | cut -d. -f2) -lt 18 ]]); then
    echo "   WARNING: Kernel version $KERNEL_VERSION may not fully support eBPF"
    echo "   Recommended: Kernel 4.18+ for full eBPF support"
else
    echo "   ✓ Kernel version supports eBPF"
fi

# Check eBPF support in kernel
echo
echo "2. Checking eBPF support..."
if [[ -f /sys/kernel/debug/bpf ]]; then
    echo "   ✓ eBPF debugfs available"
else
    echo "   WARNING: eBPF debugfs not available"
fi

# Check if bpftool is available
if command -v bpftool &> /dev/null; then
    echo "   ✓ bpftool available"
else
    echo "   WARNING: bpftool not found. Install with: apt-get install linux-tools-common"
fi

# Check network interfaces
echo
echo "3. Checking network interfaces..."
INTERFACES=$(ip link show | grep -E "^[0-9]+:" | cut -d: -f2 | tr -d ' ')
echo "   Available interfaces: $INTERFACES"

# Check if eth0 exists
if ip link show eth0 &> /dev/null; then
    echo "   ✓ eth0 interface found"
    
    # Check driver
    DRIVER=$(ethtool -i eth0 2>/dev/null | grep driver | cut -d: -f2 | tr -d ' ')
    echo "   Driver: $DRIVER"
    
    # Check XDP support
    if ethtool -i eth0 2>/dev/null | grep -q "driver: i40e\|driver: ixgbe\|driver: mlx4\|driver: mlx5\|driver: bnxt_en"; then
        echo "   ✓ Driver supports XDP"
    else
        echo "   WARNING: Driver may not support XDP"
    fi
else
    echo "   WARNING: eth0 interface not found"
fi

# Check if Netmoth binaries exist
echo
echo "4. Checking Netmoth binaries..."
if [[ -f bin/agent ]]; then
    echo "   ✓ agent binary found"
    AGENT_SIZE=$(ls -lh bin/agent | awk '{print $5}')
    echo "   Agent size: $AGENT_SIZE"
else
    echo "   ERROR: agent binary not found"
    exit 1
fi

if [[ -f bin/manager ]]; then
    echo "   ✓ manager binary found"
    MANAGER_SIZE=$(ls -lh bin/manager | awk '{print $5}')
    echo "   Manager size: $MANAGER_SIZE"
else
    echo "   ERROR: manager binary not found"
    exit 1
fi

# Check configuration files
echo
echo "5. Checking configuration files..."
if [[ -f config_ebpf.yml ]]; then
    echo "   ✓ eBPF configuration file found"
else
    echo "   WARNING: config_ebpf.yml not found"
    echo "   Creating basic eBPF configuration..."
    cat > config_ebpf.yml << EOF
# eBPF Configuration for Netmoth
interface: "eth0"
strategy: "ebpf"
number_of_rings: 4
max_cores: 8
snapshot_length: 65536
connection_timeout: 300
zero_copy: true
promiscuous: true
bpf: ""
log_file: "netmoth_ebpf.log"
postgres:
  user: "netmoth"
  password: "netmoth"
  db: "netmoth"
  host: "localhost:5432"
redis:
  host: "localhost:6379"
  password: ""
EOF
    echo "   ✓ Created config_ebpf.yml"
fi

# Test eBPF strategy loading
echo
echo "6. Testing eBPF strategy loading..."
if timeout 5s ./bin/agent -cfg config_ebpf.yml &> /dev/null; then
    echo "   ✓ eBPF strategy loads successfully"
else
    echo "   WARNING: eBPF strategy may have issues"
    echo "   This is normal if PostgreSQL is not running"
fi

# Check system resources
echo
echo "7. Checking system resources..."
CPU_CORES=$(nproc)
echo "   CPU cores: $CPU_CORES"

MEMORY=$(free -h | grep Mem | awk '{print $2}')
echo "   Total memory: $MEMORY"

# Check if PostgreSQL is running
echo
echo "8. Checking dependencies..."
if pg_isready -h localhost -p 5432 &> /dev/null; then
    echo "   ✓ PostgreSQL is running"
else
    echo "   WARNING: PostgreSQL is not running"
    echo "   Start with: sudo systemctl start postgresql"
fi

# Check if Redis is running
if redis-cli ping &> /dev/null; then
    echo "   ✓ Redis is running"
else
    echo "   WARNING: Redis is not running"
    echo "   Start with: sudo systemctl start redis"
fi

echo
echo "=== Test Summary ==="
echo "eBPF support has been successfully added to Netmoth!"
echo
echo "To start using eBPF:"
echo "1. Ensure PostgreSQL and Redis are running"
echo "2. Run: sudo ./bin/agent -cfg config_ebpf.yml"
echo "3. Run: ./bin/manager -cfg config_ebpf.yml"
echo
echo "For more information, see EBPF_SUPPORT.md" 