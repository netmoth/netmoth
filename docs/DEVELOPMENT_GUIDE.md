# Netmoth Development Guide

## Table of Contents

1. [Prerequisites](#prerequisites)
2. [Development Environment Setup](#development-environment-setup)
3. [Building the Project](#building-the-project)
4. [Running the Application](#running-the-application)
5. [Testing](#testing)
6. [Code Style and Standards](#code-style-and-standards)
7. [Debugging](#debugging)
8. [Performance Profiling](#performance-profiling)
9. [Contributing Guidelines](#contributing-guidelines)
10. [Release Process](#release-process)

## Prerequisites

### Required Software

- **Go 1.24+**: [Download from golang.org](https://golang.org/dl/)
- **Git**: [Download from git-scm.com](https://git-scm.com/)
- **Make**: Usually pre-installed on Linux/macOS
- **Docker** (optional): For containerized development

### System Requirements

- **Linux**: Primary development platform (Ubuntu 20.04+ recommended)
- **macOS**: Supported for development (some features may be limited)
- **Windows**: Limited support (WSL2 recommended)

### Hardware Requirements

- **CPU**: Multi-core processor (4+ cores recommended)
- **RAM**: 8GB+ for development
- **Storage**: 10GB+ free space
- **Network**: Network interface for packet capture testing

### Optional Dependencies

- **PostgreSQL**: For database development
- **Redis**: For caching development
- **PF_RING**: For high-performance packet capture
- **eBPF tools**: For eBPF development

## Development Environment Setup

### 1. Clone the Repository

```bash
git clone https://github.com/netmoth/netmoth.git
cd netmoth
```

### 2. Install Go Dependencies

```bash
go mod download
go mod tidy
```

### 3. Install Development Tools

```bash
# Install golangci-lint for code linting
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Install gofumpt for code formatting
go install mvdan.cc/gofumpt@latest

# Install staticcheck for static analysis
go install honnef.co/go/tools/cmd/staticcheck@latest
```

### 4. Setup Database (Optional)

```bash
# Using Docker for PostgreSQL
docker run -d \
  --name netmoth-postgres \
  -e POSTGRES_PASSWORD=netmoth \
  -e POSTGRES_USER=netmoth \
  -e POSTGRES_DB=netmoth \
  -p 5432:5432 \
  postgres:15

# Using Docker for Redis
docker run -d \
  --name netmoth-redis \
  -p 6379:6379 \
  redis:7-alpine
```

### 5. Environment Variables

Create a `.env` file in the project root:

```bash
# Database Configuration
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=netmoth
POSTGRES_PASSWORD=netmoth
POSTGRES_DB=netmoth

# Redis Configuration
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=

# Development Settings
DEBUG=true
LOG_LEVEL=debug
```

## Building the Project

### 1. Basic Build Commands

```bash
# Build all components
make build

# Build specific component
make build agent
make build manager

# Build with optimizations
make build-optimized

# Build with eBPF support
make build-ebpf

# Build with race detector
make build-race
```

### 2. Build Flags and Options

```bash
# Standard build flags
GO_FLAGS="-ldflags=-s -w" -gcflags="-l=4" -trimpath

# Optimized build flags
GO_OPTIMIZE_FLAGS="-ldflags=-s -w -extldflags=-Wl,-z,relro,-z,now" -gcflags="-l=4 -B -N" -trimpath

# eBPF build flags
GO_EBPF_FLAGS="-ldflags=-s -w" -gcflags="-l=4" -trimpath -tags=ebpf
```

### 3. Cross-Platform Building

```bash
# Build for Linux
GOOS=linux GOARCH=amd64 make build

# Build for macOS
GOOS=darwin GOARCH=amd64 make build

# Build for Windows
GOOS=windows GOARCH=amd64 make build
```

### 4. Docker Builds

```bash
# Build agent Docker image
docker build -f docker/Dockerfile.agent -t netmoth/agent:latest .

# Build manager Docker image
docker build -f docker/Dockerfile.manager -t netmoth/manager:latest .
```

## Running the Application

### 1. Development Mode

```bash
# Run agent with default config
make run-agent

# Run agent with optimized config
make run-agent-optimized

# Run agent with eBPF config
make run-agent-ebpf

# Run manager
make run-manager

# Run manager with optimized config
make run-manager-optimized
```

### 2. Configuration Files

Copy and modify configuration files:

```bash
# Agent configurations
cp cmd/agent/config.yml.example cmd/agent/config.yml
cp cmd/agent/config_optimized.yml cmd/agent/config.yml
cp cmd/agent/config_ebpf.yml cmd/agent/config.yml

# Manager configurations
cp cmd/manager/config.yml.example cmd/manager/config.yml
cp cmd/manager/config_optimized.yml cmd/manager/config.yml
cp cmd/manager/config_ebpf.yml cmd/manager/config.yml
```

### 3. Running with Custom Configuration

```bash
# Run agent with custom config
./bin/agent -cfg /path/to/custom/config.yml

# Run manager with custom config
./bin/manager -cfg /path/to/custom/config.yml
```

### 4. Docker Compose (Development)

```bash
# Start all services
docker-compose -f docker-compose.dev.yml up -d

# View logs
docker-compose -f docker-compose.dev.yml logs -f

# Stop services
docker-compose -f docker-compose.dev.yml down
```

## Testing

### 1. Unit Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with verbose output
go test -v ./...

# Run specific test
go test -v ./internal/sensor

# Run tests with race detector
go test -race ./...
```

### 2. Integration Tests

```bash
# Run integration tests
go test -tags=integration ./...

# Run with database
go test -tags=integration,postgres ./...

# Run with Redis
go test -tags=integration,redis ./...
```

### 3. Performance Tests

```bash
# Run benchmarks
go test -bench=. ./...

# Run benchmarks with memory profiling
go test -bench=. -benchmem ./...

# Run specific benchmark
go test -bench=BenchmarkPacketProcessing ./internal/sensor
```

### 4. Test Coverage

```bash
# Generate coverage report
go test -coverprofile=coverage.out ./...

# View coverage in browser
go tool cover -html=coverage.out

# Generate coverage report for specific package
go test -coverprofile=sensor.out ./internal/sensor
go tool cover -html=sensor.out
```

### 5. Test Data

```bash
# Generate test PCAP files
./scripts/generate_test_data.sh

# Download sample PCAP files
wget https://wiki.wireshark.org/SampleCaptures -O testdata/
```

## Code Style and Standards

### 1. Go Code Style

Follow the official Go code style:

```bash
# Format code
gofumpt -w .

# Check formatting
gofumpt -d .

# Run linter
golangci-lint run

# Run static analysis
staticcheck ./...
```

### 2. Code Organization

```
project/
├── cmd/                    # Application entry points
│   ├── agent/             # Agent application
│   └── manager/           # Manager application
├── internal/              # Private application code
│   ├── analyzer/          # Protocol analyzers
│   ├── config/            # Configuration management
│   ├── connection/        # Connection handling
│   ├── sensor/            # Packet capture and processing
│   ├── signature/         # Signature detection
│   ├── storage/           # Data storage
│   ├── utils/             # Utility functions
│   ├── version/           # Version information
│   └── web/               # Web interface
├── docs/                  # Documentation
├── scripts/               # Build and utility scripts
├── testdata/              # Test data and fixtures
└── docker/                # Docker configurations
```

### 3. Naming Conventions

- **Packages**: Lowercase, single word
- **Functions**: MixedCaps or mixedCaps
- **Variables**: MixedCaps or mixedCaps
- **Constants**: MixedCaps
- **Types**: MixedCaps
- **Interfaces**: MixedCaps ending with 'er' if appropriate

### 4. Error Handling

```go
// Good error handling
if err != nil {
    return fmt.Errorf("failed to process packet: %w", err)
}

// Use wrapped errors for context
if err := processData(); err != nil {
    return fmt.Errorf("data processing failed: %w", err)
}
```

### 5. Documentation

```go
// Package documentation
// Package sensor provides packet capture and analysis capabilities.
package sensor

// Function documentation
// New creates a new sensor instance with the given configuration.
// It initializes packet capture, database connections, and starts
// the processing pipeline.
func New(config *config.Config) {
    // Implementation
}

// Type documentation
// Sensor represents a packet capture and analysis sensor.
type Sensor struct {
    // Fields
}
```

## Debugging

### 1. Debug Builds

```bash
# Build with debug symbols
go build -gcflags="all=-N -l" -o bin/agent cmd/agent/main.go

# Build with debug information
go build -ldflags="-X main.debug=true" -o bin/agent cmd/agent/main.go
```

### 2. Using Delve Debugger

```bash
# Install delve
go install github.com/go-delve/delve/cmd/dlv@latest

# Debug agent
dlv debug cmd/agent/main.go

# Debug manager
dlv debug cmd/manager/main.go

# Attach to running process
dlv attach <pid>
```

### 3. Logging and Debugging

```bash
# Enable debug logging
export DEBUG=true
export LOG_LEVEL=debug

# Run with debug output
./bin/agent -debug -cfg cmd/agent/config.yml
```

### 4. Profiling

```bash
# CPU profiling
go test -cpuprofile=cpu.prof -bench=.

# Memory profiling
go test -memprofile=mem.prof -bench=.

# Block profiling
go test -blockprofile=block.prof -bench=.

# Analyze profiles
go tool pprof cpu.prof
go tool pprof mem.prof
```

## Performance Profiling

### 1. Runtime Profiling

```bash
# Start profiling server
go tool pprof -http=:8080 http://localhost:6060/debug/pprof/profile

# Profile heap
go tool pprof -http=:8080 http://localhost:6060/debug/pprof/heap

# Profile goroutines
go tool pprof -http=:8080 http://localhost:6060/debug/pprof/goroutine
```

### 2. Benchmarking

```bash
# Run all benchmarks
go test -bench=. ./...

# Run specific benchmark
go test -bench=BenchmarkPacketCapture ./internal/sensor

# Run benchmarks with different strategies
go test -bench=BenchmarkCapture -benchmem ./internal/sensor/strategies
```

### 3. Performance Monitoring

```bash
# Monitor CPU usage
top -p $(pgrep netmoth)

# Monitor memory usage
ps aux | grep netmoth

# Monitor network usage
iftop -i eth0

# Monitor disk I/O
iotop -p $(pgrep netmoth)
```

## Contributing Guidelines

### 1. Development Workflow

1. **Fork the repository**
2. **Create a feature branch**
   ```bash
   git checkout -b feature/your-feature-name
   ```
3. **Make your changes**
4. **Write tests for new functionality**
5. **Run tests and linting**
   ```bash
   make test
   make lint
   ```
6. **Commit your changes**
   ```bash
   git commit -m "feat: add new feature description"
   ```
7. **Push to your fork**
   ```bash
   git push origin feature/your-feature-name
   ```
8. **Create a pull request**

### 2. Commit Message Format

Follow the [Conventional Commits](https://www.conventionalcommits.org/) format:

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

Types:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes
- `refactor`: Code refactoring
- `test`: Test changes
- `chore`: Build process or auxiliary tool changes

### 3. Pull Request Guidelines

- **Title**: Clear and descriptive
- **Description**: Explain what and why, not how
- **Tests**: Include tests for new functionality
- **Documentation**: Update relevant documentation
- **Breaking Changes**: Clearly mark and explain

### 4. Code Review Process

1. **Automated Checks**: CI/CD pipeline must pass
2. **Code Review**: At least one maintainer approval
3. **Testing**: All tests must pass
4. **Documentation**: Documentation updated if needed

## Release Process

### 1. Version Management

```bash
# Update version
echo "v1.2.3" > internal/version/version.txt

# Tag release
git tag -a v1.2.3 -m "Release v1.2.3"
git push origin v1.2.3
```

### 2. Building Releases

```bash
# Build release binaries
make build-release

# Build Docker images
make docker-build-release

# Create release artifacts
make release-artifacts
```

### 3. Release Checklist

- [ ] Update version number
- [ ] Update CHANGELOG.md
- [ ] Run full test suite
- [ ] Build all platforms
- [ ] Create release notes
- [ ] Tag release
- [ ] Upload artifacts
- [ ] Announce release

### 4. Release Automation

```bash
# Using goreleaser
goreleaser release --rm-dist

# Using custom scripts
./scripts/release.sh v1.2.3
```

## Troubleshooting

### Common Issues

1. **Permission Denied for Packet Capture**
   ```bash
   sudo setcap cap_net_raw,cap_net_admin=eip bin/agent
   ```

2. **Database Connection Issues**
   ```bash
   # Check PostgreSQL status
   sudo systemctl status postgresql
   
   # Check connection
   psql -h localhost -U netmoth -d netmoth
   ```

3. **eBPF Loading Issues**
   ```bash
   # Check eBPF support
   cat /sys/kernel/debug/bpf/verifier_log
   
   # Check kernel version
   uname -r
   ```

4. **Performance Issues**
   ```bash
   # Check CPU usage
   top -p $(pgrep netmoth)
   
   # Check memory usage
   cat /proc/$(pgrep netmoth)/status
   ```

### Getting Help

- **Issues**: Create an issue on GitHub
- **Discussions**: Use GitHub Discussions
- **Documentation**: Check the docs/ directory
- **Examples**: Look at the examples/ directory

This development guide provides comprehensive information for developers working on the Netmoth project, covering all aspects from initial setup to contributing and releasing.