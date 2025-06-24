# Netmoth Agent API Reference

This document provides detailed information about the Netmoth Agent API endpoints.

## Base URL

All API endpoints are relative to the central server URL: `http://localhost:3000`

## API Endpoints

### 1. Agent Registration

**Endpoint:** `POST /api/agent/register`

**Request:**
```json
{
  "agent_id": "agent-001",
  "hostname": "machine1.example.com",
  "interface": "eth0",
  "version": "1.0.0"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Agent registered successfully"
}
```

### 2. Send Data

**Endpoint:** `POST /api/agent/data`

**Request:**
```json
{
  "agent_id": "agent-001",
  "hostname": "machine1.example.com",
  "interface": "eth0",
  "timestamp": "2024-01-01T12:00:00Z",
  "connections": [...],
  "signatures": [...],
  "stats": {
    "packets_received": 1000,
    "packets_dropped": 10,
    "packets_processed": 990,
    "connections_found": 50
  }
}
```

### 3. Health Check

**Endpoint:** `GET /api/agent/health?agent_id=agent-001`

**Response:**
```json
{
  "success": true,
  "message": "Agent agent-001 is healthy"
}
```

### 4. Version

**Endpoint:** `GET /api/version`

**Response:**
```json
{
  "version": "1.0.0"
}
```

## Error Handling

All endpoints return consistent error responses:

```json
{
  "success": false,
  "message": "Error description",
  "error": "Detailed error information"
}
```

## Examples

### Register Agent
```bash
curl -X POST http://localhost:3000/api/agent/register \
  -H "Content-Type: application/json" \
  -d '{"agent_id":"agent-001","hostname":"test","interface":"eth0","version":"1.0.0"}'
```

### Send Data
```bash
curl -X POST http://localhost:3000/api/agent/data \
  -H "Content-Type: application/json" \
  -d '{"agent_id":"agent-001","connections":[],"signatures":[],"stats":{"packets_received":100}}'
```

### Health Check
```bash
curl "http://localhost:3000/api/agent/health?agent_id=agent-001"
``` 