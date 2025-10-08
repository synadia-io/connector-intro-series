# Temperature Analysis NATS Microservice

A NATS microservice that analyzes temperature data and classifies it as cold, warm, or hot.

## Features

- Accepts temperature in Celsius or Fahrenheit
- Returns classification: cold (<10°C), warm (10-25°C), or hot (>25°C)
- Built using NATS micro framework for service discovery and monitoring

## Running the Service

### 1. Start NATS Server
```bash
# Using Docker
docker run -p 4222:4222 nats:latest

# Or install and run locally
nats-server
```

### 2. Run the Temperature Service
```bash
./temperature-service

# Or with custom NATS URL
NATS_URL=nats://localhost:4222 ./temperature-service
```

### 3. Test with the Client
```bash
./temperature-client
```

## API

**Subject:** `temperature.analyze`

### Request Format
```json
{
  "temperature": 25.5,
  "unit": "celsius"  // Optional: "celsius" or "fahrenheit", defaults to celsius
}
```

### Response Format
```json
{
  "status": "warm",
  "temperature": 25.5,
  "unit": "celsius",
  "description": "25.5°C is warm"
}
```

## Temperature Classifications

- **Cold**: < 10°C (50°F)
- **Warm**: 10-25°C (50-77°F)
- **Hot**: > 25°C (77°F)

## Using with NATS CLI

```bash
# Send a request
nats req temperature.analyze '{"temperature": 20, "unit": "celsius"}'

# Monitor the service
nats micro info temperature-analyzer
nats micro stats temperature-analyzer
```

## Building from Source

```bash
go mod tidy
go build -o temperature-service temperature-service.go
go build -o temperature-client temperature-client.go
```