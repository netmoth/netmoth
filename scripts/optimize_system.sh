#!/bin/bash

# System optimization script for Netmoth
# Recommended to run as root

set -e

echo "=== System optimization for Netmoth ==="

# Check for root privileges
if [[ $EUID -ne 0 ]]; then
   echo "This script must be run as root"
   exit 1
fi

# Function to get network interface
get_interface() {
    if [ -z "$1" ]; then
        # Automatically detect main interface
        INTERFACE=$(ip route | grep default | awk '{print $5}' | head -1)
        echo "Automatically detected interface: $INTERFACE"
    else
        INTERFACE=$1
    fi
}

# Function to optimize network interface
optimize_network_interface() {
    local interface=$1
    echo "Optimizing network interface $interface..."
    
    # Increase buffer size
    ethtool -G $interface rx 4096 tx 4096 2>/dev/null || echo "Failed to change buffer size"
    
    # Disable offloading for better control
    ethtool -K $interface tso off gso off gro off 2>/dev/null || echo "Failed to disable offloading"
    
    # Set promiscuous mode
    ip link set $interface promisc on
    
    echo "Network interface $interface optimized"
}

# Function to set up IRQ affinity
setup_irq_affinity() {
    local interface=$1
    echo "Setting up IRQ affinity for $interface..."
    
    # Get IRQs for the interface
    IRQS=$(grep $interface /proc/interrupts | awk '{print $1}' | sed 's/://')
    
    if [ -n "$IRQS" ]; then
        CPU=0
        for irq in $IRQS; do
            echo $((1 << CPU)) > /proc/irq/$irq/smp_affinity
            CPU=$((CPU + 1))
            if [ $CPU -ge $(nproc) ]; then
                CPU=0
            fi
        done
        echo "IRQ affinity set"
    else
        echo "No IRQ found for interface $interface"
    fi
}

# Function to optimize system parameters
optimize_system() {
    echo "Optimizing system parameters..."
    
    # Increase file descriptor limits
    echo "* soft nofile 65536" >> /etc/security/limits.conf
    echo "* hard nofile 65536" >> /etc/security/limits.conf
    
    # Increase kernel buffer sizes
    echo "net.core.rmem_max = 134217728" >> /etc/sysctl.conf
    echo "net.core.wmem_max = 134217728" >> /etc/sysctl.conf
    echo "net.core.rmem_default = 262144" >> /etc/sysctl.conf
    echo "net.core.wmem_default = 262144" >> /etc/sysctl.conf
    echo "net.core.netdev_max_backlog = 5000" >> /etc/sysctl.conf
    
    # TCP optimization
    echo "net.ipv4.tcp_rmem = 4096 87380 134217728" >> /etc/sysctl.conf
    echo "net.ipv4.tcp_wmem = 4096 65536 134217728" >> /etc/sysctl.conf
    echo "net.ipv4.tcp_congestion_control = bbr" >> /etc/sysctl.conf
    
    # Disable IPv6 if not used
    echo "net.ipv6.conf.all.disable_ipv6 = 1" >> /etc/sysctl.conf
    echo "net.ipv6.conf.default.disable_ipv6 = 1" >> /etc/sysctl.conf
    
    # Apply changes
    sysctl -p
    
    echo "System parameters optimized"
}

# Function to set up NUMA (if available)
setup_numa() {
    if command -v numactl &> /dev/null; then
        echo "Setting up NUMA..."
        # Detect number of NUMA nodes
        NUMA_NODES=$(numactl --hardware | grep "available:" | awk '{print $2}')
        if [ "$NUMA_NODES" -gt 1 ]; then
            echo "$NUMA_NODES NUMA nodes detected"
            echo "It is recommended to run Netmoth bound to a specific node:"
            echo "numactl --cpunodebind=0 --membind=0 ./bin/agent"
        fi
    fi
}

# Function to create systemd service
create_systemd_service() {
    echo "Creating systemd service..."
    
    cat > /etc/systemd/system/netmoth.service << EOF
[Unit]
Description=Netmoth Network Traffic Analyzer
After=network.target postgresql.service redis.service

[Service]
Type=simple
User=netmoth
Group=netmoth
WorkingDirectory=/opt/netmoth
ExecStart=/opt/netmoth/bin/agent
ExecReload=/bin/kill -HUP \$MAINPID
Restart=always
RestartSec=5
LimitNOFILE=65536
LimitNPROC=65536

# Performance optimizations
Nice=-10
IOSchedulingClass=1
IOSchedulingPriority=4

# CPU binding (uncomment and configure if needed)
# CPUAffinity=0 1 2 3 4 5 6 7

[Install]
WantedBy=multi-user.target
EOF

    # Create separate service for web interface
    cat > /etc/systemd/system/netmoth-web.service << EOF
[Unit]
Description=Netmoth Web Interface
After=network.target postgresql.service redis.service
Wants=netmoth.service

[Service]
Type=simple
User=netmoth
Group=netmoth
WorkingDirectory=/opt/netmoth
ExecStart=/opt/netmoth/bin/manager
ExecReload=/bin/kill -HUP \$MAINPID
Restart=always
RestartSec=5
LimitNOFILE=65536
LimitNPROC=65536

# Performance optimizations
Nice=-5
IOSchedulingClass=1
IOSchedulingPriority=4

[Install]
WantedBy=multi-user.target
EOF

    echo "Systemd services created:"
    echo "- /etc/systemd/system/netmoth.service (agent)"
    echo "- /etc/systemd/system/netmoth-web.service (web interface)"
    echo "To enable, run: systemctl enable netmoth netmoth-web"
}

# Function to create user
create_user() {
    echo "Creating user netmoth..."
    
    if ! id "netmoth" &>/dev/null; then
        useradd -r -s /bin/false -d /opt/netmoth netmoth
        mkdir -p /opt/netmoth
        chown netmoth:netmoth /opt/netmoth
        echo "User netmoth created"
    else
        echo "User netmoth already exists"
    fi
}

# Main function
main() {
    local interface=$1
    
    get_interface $interface
    
    echo "Starting optimization for interface: $INTERFACE"
    
    # Create user
    create_user
    
    # Optimize system
    optimize_system
    
    # Optimize network interface
    optimize_network_interface $INTERFACE
    
    # Set up IRQ affinity
    setup_irq_affinity $INTERFACE
    
    # Set up NUMA
    setup_numa
    
    # Create systemd service
    create_systemd_service
    
    echo ""
    echo "=== Optimization complete ==="
    echo ""
    echo "Recommendations:"
    echo "1. Reboot the system to apply all changes"
    echo "2. Use the config_optimized.yml configuration"
    echo "3. Run Netmoth as root or via systemd"
    echo "4. Monitor performance via pprof on port 6060"
    echo ""
    echo "Monitoring commands:"
    echo "- CPU: htop"
    echo "- Network: watch -n 1 'cat /proc/net/dev | grep $INTERFACE'"
    echo "- Profiling: go tool pprof http://localhost:6060/debug/pprof/profile"
    echo ""
    echo "Service management:"
    echo "- Start services: systemctl start netmoth netmoth-web"
    echo "- Enable services: systemctl enable netmoth netmoth-web"
    echo "- Check status: systemctl status netmoth netmoth-web"
}

# Argument handling
case "${1:-}" in
    -h|--help)
        echo "Usage: $0 [interface]"
        echo "  interface - network interface name (default: auto-detect)"
        echo ""
        echo "Examples:"
        echo "  $0           # Auto-detect interface"
        echo "  $0 eth0      # Use eth0 interface"
        echo "  $0 enp0s3    # Use enp0s3 interface"
        exit 0
        ;;
    *)
        main "$@"
        ;;
esac 