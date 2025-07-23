# Netmoth Architecture Documentation

## Table of Contents

1. [System Overview](#system-overview)
2. [Architecture Components](#architecture-components)
3. [Data Flow](#data-flow)
4. [Deployment Models](#deployment-models)
5. [Performance Architecture](#performance-architecture)
6. [Security Architecture](#security-architecture)
7. [Scalability Design](#scalability-design)
8. [Monitoring and Observability](#monitoring-and-observability)

## System Overview

Netmoth is designed as a distributed network monitoring and intrusion detection system with a clear separation between data collection (agents) and data processing (manager). The architecture follows a microservices pattern with high-performance packet processing capabilities.

### Core Design Principles

- **High Performance**: Zero-copy packet processing with multiple capture strategies
- **Scalability**: Distributed architecture with lightweight agents
- **Reliability**: Fault-tolerant design with graceful degradation
- **Security**: Multi-layered security with signature-based detection
- **Observability**: Comprehensive logging and monitoring

## Architecture Components

### 1. Agent Component

The agent is a lightweight, high-performance packet capture and analysis component designed to run on network segments or endpoints.

#### Agent Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                        Agent                                │
├─────────────────────────────────────────────────────────────┤
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐         │
│  │   PCAP      │  │  AF_PACKET  │  │    eBPF     │         │
│  │  Strategy   │  │  Strategy   │  │  Strategy   │         │
│  └─────────────┘  └─────────────┘  └─────────────┘         │
│           │               │               │                │
│           └───────────────┼───────────────┘                │
│                           │                                │
│  ┌─────────────────────────────────────────────────────┐   │
│  │              Packet Processing Engine               │   │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐ │   │
│  │  │ Zero-Copy   │  │ TCP Stream  │  │ Protocol    │ │   │
│  │  │ Processing  │  │ Reassembly  │  │ Analyzers   │ │   │
│  │  └─────────────┘  └─────────────┘  └─────────────┘ │   │
│  └─────────────────────────────────────────────────────┘   │
│                           │                                │
│  ┌─────────────────────────────────────────────────────┐   │
│  │              Signature Detection                    │   │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐ │   │
│  │  │ IP Check    │  │ Botnet      │  │ Tracker     │ │   │
│  │  │             │  │ Detection   │  │ Detection   │ │   │
│  │  └─────────────┘  └─────────────┘  └─────────────┘ │   │
│  └─────────────────────────────────────────────────────┘   │
│                           │                                │
│  ┌─────────────────────────────────────────────────────┐   │
│  │              Agent Communication                    │   │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐ │   │
│  │  │ Registration│  │ Data        │  │ Health      │ │   │
│  │  │             │  │ Transmission│  │ Monitoring  │ │   │
│  │  └─────────────┘  └─────────────┘  └─────────────┘ │   │
│  └─────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
```

#### Key Features

- **Multiple Capture Strategies**: PCAP, AF_PACKET, PF_RING, eBPF
- **Zero-Copy Processing**: Minimizes memory allocations
- **TCP Stream Reassembly**: Complete connection reconstruction
- **Protocol Analysis**: HTTP, HTTPS, DNS, TLS analysis
- **Signature Detection**: Real-time threat detection
- **Agent Communication**: Registration and data transmission

### 2. Manager Component

The manager serves as the central coordination point, providing data aggregation, storage, and web interface.

#### Manager Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                        Manager                              │
├─────────────────────────────────────────────────────────────┤
│  ┌─────────────────────────────────────────────────────┐   │
│  │              Web Interface                          │   │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐ │   │
│  │  │ REST API    │  │ WebSocket   │  │ Static      │ │   │
│  │  │             │  │ Support     │  │ Files       │ │   │
│  │  └─────────────┘  └─────────────┘  └─────────────┘ │   │
│  └─────────────────────────────────────────────────────┘   │
│                           │                                │
│  ┌─────────────────────────────────────────────────────┐   │
│  │              Agent Management                       │   │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐ │   │
│  │  │ Agent       │  │ Data        │  │ Health      │ │   │
│  │  │ Registration│  │ Processing  │  │ Monitoring  │ │   │
│  │  └─────────────┘  └─────────────┘  └─────────────┘ │   │
│  └─────────────────────────────────────────────────────┘   │
│                           │                                │
│  ┌─────────────────────────────────────────────────────┐   │
│  │              Data Storage                           │   │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐ │   │
│  │  │ PostgreSQL  │  │ Redis       │  │ File        │ │   │
│  │  │ Database    │  │ Cache       │  │ Logging     │ │   │
│  │  └─────────────┘  └─────────────┘  └─────────────┘ │   │
│  └─────────────────────────────────────────────────────┘   │
│                           │                                │
│  ┌─────────────────────────────────────────────────────┐   │
│  │              Signature Management                   │   │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐ │   │
│  │  │ Signature   │  │ Threat      │  │ Update      │ │   │
│  │  │ Database    │  │ Intelligence│  │ Management  │ │   │
│  │  └─────────────┘  └─────────────┘  └─────────────┘ │   │
│  └─────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
```

#### Key Features

- **RESTful API**: Agent communication and data access
- **WebSocket Support**: Real-time data streaming
- **Agent Management**: Registration and health monitoring
- **Data Storage**: PostgreSQL for persistence, Redis for caching
- **Signature Management**: Threat database management

## Data Flow

### 1. Packet Capture Flow

```
Network Interface
       │
       ▼
┌─────────────┐
│   Capture   │ ← PCAP/AF_PACKET/eBPF Strategy
│   Strategy  │
└─────────────┘
       │
       ▼
┌─────────────┐
│ Zero-Copy   │ ← Packet Processing
│ Processing  │
└─────────────┘
       │
       ▼
┌─────────────┐
│ TCP Stream  │ ← Connection Reassembly
│ Reassembly  │
└─────────────┘
       │
       ▼
┌─────────────┐
│ Protocol    │ ← HTTP/DNS/TLS Analysis
│ Analyzers   │
└─────────────┘
       │
       ▼
┌─────────────┐
│ Signature   │ ← Threat Detection
│ Detection   │
└─────────────┘
       │
       ▼
┌─────────────┐
│ Data        │ ← Storage/Transmission
│ Output      │
└─────────────┘
```

### 2. Agent-Manager Communication Flow

```
Agent                    Manager
  │                        │
  │─── Register ──────────►│
  │                        │
  │◄─── ACK ───────────────│
  │                        │
  │─── Data ──────────────►│
  │                        │
  │◄─── ACK ───────────────│
  │                        │
  │─── Health ────────────►│
  │                        │
  │◄─── Status ────────────│
```

### 3. Data Processing Flow

```
Raw Packet Data
       │
       ▼
┌─────────────┐
│ Packet      │ ← Initial Processing
│ Decoding    │
└─────────────┘
       │
       ▼
┌─────────────┐
│ Connection  │ ← Stream Assembly
│ Tracking    │
└─────────────┘
       │
       ▼
┌─────────────┐
│ Protocol    │ ← Application Layer
│ Analysis    │   Analysis
└─────────────┘
       │
       ▼
┌─────────────┐
│ Signature   │ ← Threat Detection
│ Matching    │
└─────────────┘
       │
       ▼
┌─────────────┐
│ Data        │ ← Storage/Aggregation
│ Storage     │
└─────────────┘
```

## Deployment Models

### 1. Single Node Deployment

```
┌─────────────────────────────────────────────────────────────┐
│                    Single Node                             │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐         │
│  │   Agent     │  │  Manager    │  │ PostgreSQL  │         │
│  │             │  │             │  │             │         │
│  └─────────────┘  └─────────────┘  └─────────────┘         │
└─────────────────────────────────────────────────────────────┘
```

**Use Cases**: Development, testing, small networks

### 2. Distributed Deployment

```
┌─────────────────────────────────────────────────────────────┐
│                    Central Manager                         │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐         │
│  │   Manager   │  │ PostgreSQL  │  │    Redis    │         │
│  │             │  │             │  │             │         │
│  └─────────────┘  └─────────────┘  └─────────────┘         │
└─────────────────────────────────────────────────────────────┘
                              │
                              │ Network
                              ▼
┌─────────────┐  ┌─────────────┐  ┌─────────────┐
│   Agent 1   │  │   Agent 2   │  │   Agent N   │
│             │  │             │  │             │
└─────────────┘  └─────────────┘  └─────────────┘
```

**Use Cases**: Production environments, large networks

### 3. High Availability Deployment

```
┌─────────────────────────────────────────────────────────────┐
│                    Load Balancer                           │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────┐  ┌─────────────┐  ┌─────────────┐
│ Manager 1   │  │ Manager 2   │  │ Manager N   │
│             │  │             │  │             │
└─────────────┘  └─────────────┘  └─────────────┘
       │               │               │
       └───────────────┼───────────────┘
                       ▼
┌─────────────┐  ┌─────────────┐  ┌─────────────┐
│ PostgreSQL  │  │ PostgreSQL  │  │ PostgreSQL  │
│ Primary     │  │ Secondary   │  │ Secondary   │
└─────────────┘  └─────────────┘  └─────────────┘
```

**Use Cases**: Enterprise environments, critical infrastructure

## Performance Architecture

### 1. Packet Processing Pipeline

```
┌─────────────────────────────────────────────────────────────┐
│                    Performance Pipeline                    │
├─────────────────────────────────────────────────────────────┤
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐         │
│  │ Kernel      │  │ User Space  │  │ Application │         │
│  │ eBPF        │  │ Zero-Copy   │  │ Processing  │         │
│  │ Processing  │  │ Buffers     │  │             │         │
│  └─────────────┘  └─────────────┘  └─────────────┘         │
│           │               │               │                │
│           ▼               ▼               ▼                │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐         │
│  │ Packet      │  │ Memory      │  │ Worker      │         │
│  │ Filtering   │  │ Pool        │  │ Pool        │         │
│  └─────────────┘  └─────────────┘  └─────────────┘         │
└─────────────────────────────────────────────────────────────┘
```

### 2. Concurrency Model

```
┌─────────────────────────────────────────────────────────────┐
│                    Concurrency Model                       │
├─────────────────────────────────────────────────────────────┤
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐         │
│  │ Capture     │  │ Processing  │  │ Analysis    │         │
│  │ Goroutines  │  │ Goroutines  │  │ Goroutines  │         │
│  │             │  │             │  │             │         │
│  └─────────────┘  └─────────────┘  └─────────────┘         │
│           │               │               │                │
│           ▼               ▼               ▼                │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐         │
│  │ Packet      │  │ Connection  │  │ Signature   │         │
│  │ Channels    │  │ Channels    │  │ Channels    │         │
│  └─────────────┘  └─────────────┘  └─────────────┘         │
└─────────────────────────────────────────────────────────────┘
```

### 3. Memory Management

- **Sync.Pool**: Efficient buffer reuse
- **Zero-Copy**: Minimize memory allocations
- **Connection Pooling**: Database connection management
- **Garbage Collection**: Optimized for high-throughput

## Security Architecture

### 1. Multi-Layer Security

```
┌─────────────────────────────────────────────────────────────┐
│                    Security Layers                         │
├─────────────────────────────────────────────────────────────┤
│  ┌─────────────────────────────────────────────────────┐   │
│  │              Application Layer                      │   │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐ │   │
│  │  │ Input       │  │ SQL         │  │ Certificate │ │   │
│  │  │ Validation  │  │ Injection   │  │ Validation  │ │   │
│  │  │             │  │ Prevention  │  │             │ │   │
│  │  └─────────────┘  └─────────────┘  └─────────────┘ │   │
│  └─────────────────────────────────────────────────────┘   │
│                           │                                │
│  ┌─────────────────────────────────────────────────────┐   │
│  │              Network Layer                          │   │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐ │   │
│  │  │ TLS/SSL     │  │ Firewall    │  │ Network     │ │   │
│  │  │ Encryption  │  │ Rules       │  │ Isolation   │ │   │
│  │  └─────────────┘  └─────────────┘  └─────────────┘ │   │
│  └─────────────────────────────────────────────────────┘   │
│                           │                                │
│  ┌─────────────────────────────────────────────────────┐   │
│  │              Signature Layer                        │   │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐ │   │
│  │  │ IP          │  │ Botnet      │  │ Tracker     │ │   │
│  │  │ Blacklists  │  │ Detection   │  │ Detection   │ │   │
│  │  └─────────────┘  └─────────────┘  └─────────────┘ │   │
│  └─────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
```

### 2. Threat Detection

- **Signature-Based Detection**: Known threat patterns
- **Behavioral Analysis**: Anomaly detection
- **Certificate Validation**: TLS certificate analysis
- **IP Reputation**: Blacklist checking

### 3. Data Protection

- **Encryption**: TLS for data in transit
- **Access Control**: Authentication and authorization
- **Audit Logging**: Comprehensive security logging
- **Data Sanitization**: Input validation and sanitization

## Scalability Design

### 1. Horizontal Scaling

```
┌─────────────────────────────────────────────────────────────┐
│                    Scalability Model                       │
├─────────────────────────────────────────────────────────────┤
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐         │
│  │ Load        │  │ Manager     │  │ Database    │         │
│  │ Balancer    │  │ Cluster     │  │ Cluster     │         │
│  └─────────────┘  └─────────────┘  └─────────────┘         │
│           │               │               │                │
│           ▼               ▼               ▼                │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐         │
│  │ Agent       │  │ Manager     │  │ PostgreSQL  │         │
│  │ Pool        │  │ Instances   │  │ Replicas    │         │
│  └─────────────┘  └─────────────┘  └─────────────┘         │
└─────────────────────────────────────────────────────────────┘
```

### 2. Vertical Scaling

- **CPU Optimization**: Multi-core processing
- **Memory Management**: Efficient buffer usage
- **I/O Optimization**: Zero-copy operations
- **Database Optimization**: Connection pooling

### 3. Elastic Scaling

- **Auto-scaling**: Based on load metrics
- **Resource Management**: Dynamic resource allocation
- **Load Distribution**: Intelligent load balancing
- **Failover**: Automatic failover mechanisms

## Monitoring and Observability

### 1. Metrics Collection

```
┌─────────────────────────────────────────────────────────────┐
│                    Monitoring Stack                        │
├─────────────────────────────────────────────────────────────┤
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐         │
│  │ Application │  │ System      │  │ Network     │         │
│  │ Metrics     │  │ Metrics     │  │ Metrics     │         │
│  └─────────────┘  └─────────────┘  └─────────────┘         │
│           │               │               │                │
│           ▼               ▼               ▼                │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐         │
│  │ Prometheus  │  │ Grafana     │  │ Alerting    │         │
│  │ Collection  │  │ Dashboards  │  │ System      │         │
│  └─────────────┘  └─────────────┘  └─────────────┘         │
└─────────────────────────────────────────────────────────────┘
```

### 2. Logging Architecture

- **Structured Logging**: JSON-formatted logs
- **Log Levels**: Debug, Info, Warn, Error
- **Log Aggregation**: Centralized log collection
- **Log Retention**: Configurable retention policies

### 3. Health Monitoring

- **Agent Health**: Regular health checks
- **System Health**: Resource monitoring
- **Network Health**: Connectivity monitoring
- **Database Health**: Connection monitoring

### 4. Alerting

- **Performance Alerts**: High CPU/memory usage
- **Security Alerts**: Threat detection
- **Availability Alerts**: Service downtime
- **Capacity Alerts**: Resource exhaustion

## Configuration Management

### 1. Configuration Hierarchy

```
┌─────────────────────────────────────────────────────────────┐
│                    Configuration Model                     │
├─────────────────────────────────────────────────────────────┤
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐         │
│  │ Environment │  │ Command     │  │ Config      │         │
│  │ Variables   │  │ Line Args   │  │ Files       │         │
│  └─────────────┘  └─────────────┘  └─────────────┘         │
│           │               │               │                │
│           ▼               ▼               ▼                │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐         │
│  │ Defaults    │  │ Validation  │  │ Runtime     │         │
│  │             │  │             │  │ Updates     │         │
│  └─────────────┘  └─────────────┘  └─────────────┘         │
└─────────────────────────────────────────────────────────────┘
```

### 2. Configuration Types

- **Agent Configuration**: Capture settings, agent mode
- **Manager Configuration**: Server settings, database config
- **Network Configuration**: Interface settings, capture strategies
- **Security Configuration**: Signature settings, access control

This architecture documentation provides a comprehensive overview of the Netmoth system design, enabling developers and operators to understand the system's structure, data flow, and operational characteristics.