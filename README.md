# Synadia Connect Introduction Series

A comprehensive development environment demonstrating NATS messaging, JetStream, and Synadia Connect data pipelines with practical examples.

## Overview

This repository provides a complete dev container setup with sample applications showing:
- **NATS JetStream** - Event streaming and persistence
- **Synadia Connect** - Data pipeline orchestration between various systems
- **NATS Microservices** - Request-response patterns with service discovery
- **Multi-language support** - Examples in Go and Python

## Quick Start

### 1. Open in Dev Container

This project is designed to run in a VS Code Dev Container with all tools pre-installed.

**Prerequisites:**
- [VS Code](https://code.visualstudio.com/)
- [Docker Desktop](https://www.docker.com/products/docker-desktop)
- VS Code "Dev Containers" extension

**Steps:**
1. Clone this repository
2. Open the folder in VS Code
3. When prompted, click "Reopen in Container"
4. Or use Command Palette (F1) → "Dev Containers: Reopen in Container"

The container automatically installs:
- Go 1.x
- Python 3.11
- Node.js LTS
- NATS Server & CLI
- Synadia Connect
- Task Runner
- ngrok

### 2. Add NATS Credentials

You need credentials to connect to Synadia Cloud (NGS).

1. Get your credentials from [Synadia Cloud](https://cloud.synadia.com)
2. Place the `.creds` file in the `credentials/` folder:
   ```
   credentials/NGS-Premium-CLI.creds
   ```
3. The credentials file format:
   ```
   -----BEGIN NATS USER JWT-----
   eyJ0eXAiOiJKV1QiLCJhbGc...
   -----END NATS USER JWT-----

   ************************* IMPORTANT *************************
   NKEY Seed printed below can be used to sign and prove identity.
   NKEYs are sensitive and should be treated as secrets.

   -----BEGIN USER NKEY SEED-----
   SUABC123...
   -----END USER NKEY SEED-----
   ```

> **Security Note:** The `credentials/` folder is gitignored. Never commit credentials to version control.

### 3. Configure NATS Context

Set up a NATS context for connection management:

```bash
task nats-context
```

Or manually:
```bash
nats context add \
   "NGS-Default-CLI" \
   --server "tls://connect.ngs.global" \
   --creds ./credentials/NGS-Premium-CLI.creds \
   --select
```

Verify the context:
```bash
nats context ls
```

## Components

### Temperature Publisher (Go)

Publishes random temperature readings to NATS JetStream every 3 seconds.

- **Stream:** `Temperatures`
- **Subject:** `telemetry.sensors.temperature`
- **Data:** JSON temperature readings from 5 sensors
- **Storage:** Creates JetStream stream with 7-day retention

**Run:**
```bash
task publisher
```

Example output:
```
Published: Sensor TEMP-003 at Warehouse - 18.4°C (Stream: Temperatures, Seq: 42)
```

### HTTP Server (Go)

Simple HTTP server providing temperature data via REST API.

- **Port:** 8080
- **Endpoints:**
  - `GET /` - Home endpoint
  - `GET /temperature` - Get random temperature reading
  - `POST /temperature` - Submit temperature data

**Run:**
```bash
task server
```

**With ngrok:**
```bash
task server-ngrok
```

### Temperature Analysis Microservice (Python)

NATS microservice that analyzes temperature data from MongoDB change streams.

- **Service:** `temperature-analyzer` v1.0.0
- **Subject:** `temperature.analyze`
- **Input:** MongoDB change stream events
- **Output:** Temperature classification (cold, warm, hot)

**Classification Rules:**
- Cold: < 10°C
- Warm: 10-25°C
- Hot: > 25°C

**Run:**
```bash
task python-temperature
```

## Synadia Connect Resources

The `resources/` folder contains configuration examples for various connectors:

- **MongoDBInlet** - Read from MongoDB change streams
- **MongoDBOutlet** - Write to MongoDB collections
- **HTTPInlet** - Expose HTTP endpoints
- **HTTPOutlet** - Send HTTP requests
- **AmazonS3Inlet** - Read from S3 buckets
- **AmazonS3Outlet** - Write to S3 buckets
- **SQSPolicy** - AWS SQS queue integration
- **MappingTransformer.js** - Data transformation logic

## Project Structure

```
connector-intro-series/
├── .devcontainer/              # Dev container configuration
│   ├── devcontainer.json
│   └── setup.sh                # Auto-setup script
├── credentials/                # NATS credentials (gitignored)
│   └── NGS-Premium-CLI.creds
├── publisher/                  # Go JetStream publisher
│   ├── publisher.go
│   └── go.mod
├── server/                     # Go HTTP server
│   ├── server.go
│   └── go.mod
├── temperature-microservice/   # Python NATS microservice
│   ├── temperature_service.py
│   └── requirements.txt
├── resources/                  # Synadia Connect configs
│   ├── MongoDBInlet
│   ├── MongoDBOutlet
│   ├── HTTPInlet
│   ├── HTTPOutlet
│   ├── AmazonS3Inlet
│   ├── AmazonS3Outlet
│   ├── SQSPolicy
│   └── MappingTransformer.js
├── Taskfile.yaml              # Task automation
└── README.md
```

## Available Tasks

| Command | Description |
|---------|-------------|
| `task deps` | Install all Go dependencies |
| `task nats-context` | Set up NATS context for Synadia Cloud |
| `task publisher` | Run the JetStream temperature publisher |
| `task server` | Run the HTTP server |
| `task server-ngrok` | Run HTTP server with ngrok tunnel |
| `task python-temperature` | Run Python temperature microservice |
| `task all` | Run publisher, server, and ngrok simultaneously |

## Common Workflows

### 1. Stream Temperature Data to NATS

```bash
# Start the publisher
task publisher

# In another terminal, subscribe to the stream
nats stream view Temperatures
```

### 2. Test HTTP Endpoints

```bash
# Start the server
task server

# Get temperature reading
curl http://localhost:8080/temperature

# Post temperature data
curl -X POST http://localhost:8080/temperature \
  -H "Content-Type: application/json" \
  -d '{"sensor_id":"TEMP-001","temperature":22.5,"unit":"celsius","location":"Lab"}'
```

### 3. Run Temperature Analysis Service

```bash
# Start the microservice
task python-temperature

# Test with a MongoDB change stream event
nats req temperature.analyze '{
  "operationType": "insert",
  "fullDocument": {
    "temperature": 28.5,
    "sensor_id": "TEMP-001",
    "location": "Server Room"
  }
}'
```

### 4. Deploy Synadia Connect Pipeline

```bash
# Example: MongoDB to NATS pipeline
connect deploy \
  --inlet resources/MongoDBInlet \
  --outlet resources/MongoDBOutlet \
  --transformer resources/MappingTransformer.js
```

## Development Tools Included

- **Go** - NATS client development
- **Python 3.11** - Microservice development
- **NATS Server** - Local message broker
- **NATS CLI** - Stream and service management
- **Synadia Connect** - Data pipeline orchestration
- **Task Runner** - Build automation
- **ngrok** - Public tunnel for webhooks
- **Git & GitHub CLI** - Version control

## Troubleshooting

### Credentials Not Found

If you see `Credentials file not found`, ensure:
1. Your `.creds` file is in `credentials/NGS-Premium-CLI.creds`
2. The file has proper permissions (readable)
3. The path matches what's in the Taskfile

### Connection Failed

If connection to NGS fails:
1. Check credentials are valid and not expired
2. Verify internet connectivity
3. Ensure you're using `tls://connect.ngs.global`
4. Check context with `nats context ls`

### Python Dependencies

If the Python microservice fails to start:
```bash
cd temperature-microservice
pip install -r requirements.txt
```

### Go Module Issues

If Go dependencies fail:
```bash
task deps-publisher
task deps-server
```

## Resources

- [Synadia Cloud](https://cloud.synadia.com) - Get NATS credentials
- [NATS Documentation](https://docs.nats.io/) - Complete NATS guide
- [Synadia Connect Docs](https://docs.synadia.com/connect) - Connect configuration
- [NATS Microservices](https://docs.nats.io/using-nats/developer/services) - Service patterns
- [VS Code Dev Containers](https://code.visualstudio.com/docs/devcontainers/containers) - Container development

## License

This is an educational example repository for learning Synadia Connect and NATS patterns.
