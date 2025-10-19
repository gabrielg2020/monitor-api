# monitor-api

A REST API for collecting and serving system metrics from homelab clusters. Built with Go and Gin framework.

## Features

- **Fast & Lightweight**: Built with Go for minimal resource usage
- **SQLite Backend**: No complex database setup required
- **RESTful API**: Clean, intuitive endpoints
- **CORS Support**: Configurable cross-origin access
- **Health Checks**: Built-in health monitoring endpoint
- **Docker Ready**: Pre-built container images available

## 📋 Table of Contents

- [Architecture](#architecture)
- [Prerequisites](#prerequisites)
- [Quick Start](#quick-start)
- [Development](#development)
- [API Documentation](#api-documentation)
- [Configuration](#configuration)
- [Deployment](#deployment)

## Architecture

## Prerequisites

- Go 1.24+
- Docker
- SQLite3

## Quick Start

### Using Docker

```bash
# Pull the image
docker pull ghcr.io/gabrielg2020/monitor-api:latest

# Run the container
docker run -d \
  -p 8191:8191 \
  -v $(pwd)/data:/app/data \
  -e DB_PATH=/app/data/monitoring.db \
  -e PORT=8191 \
  -e ALLOWED_ORIGINS=http://localhost:5173 \
  --name monitor-api \
  ghcr.io/gabrielg2020/monitor-api:latest
```

### Using Docker Compose

Create `docker-compose.yml`:
```yaml
services:
  api:
    image: ghcr.io/gabrielg2020/monitor-api:latest
    container_name: monitor-api
    restart: unless-stopped
    ports:
      - "8191:8191"
    volumes:
      - ../data:/app/data
    environment:
      - DB_PATH=/app/data/monitoring.db
      - PORT=8191
      - GIN_MODE=release
      - ALLOWED_ORIGINS=http://localhost:5173,https://monitoring.yourdomain.com
```

Run:
```bash
docker-compose up -d
```

### Test the API
```bash
# Health check
curl http://localhost:8191/health

# Get hosts
curl http://localhost:8191/api/v1/hosts

# Get metrics
curl "http://localhost:8191/api/v1/metrics?host_id=1&limit=10"
```

## Development

### Clone and Setup
```bash
# Clone repository
git clone https://github.com/gabrielg2020/monitor-api.git
cd monitor-api

# Install dependencies
go mod download

# Create .env file
cp .env.example .env
```

### Environment Variables

Edit `.env` file:
```bash
PORT=8191
DB_PATH=./data/monitoring.db
GIN_MODE=debug
ALLOWED_ORIGINS=http://localhost:5173
```

### Run Locally
```bash
go run cmd/main.go
```

### Build API Docs
```bash
swag init -g cmd/main.go -o docs --parseDependency --parseInternal --useStructName
```

## API Documentation

The API documentation is generated using Swagger and can be accessed at `/swagger/index.html` when the server is running.

## Configuration

### Environment Variables

| Variable          | Description                  | Default           | Required |
|-------------------|------------------------------|-------------------|----------|
| `PORT`            | Server port                  | `8191`            | Yes      |
| `DB_PATH`         | SQLite database file path    | `./monitoring.db` | Yes      |
| `GIN_MODE`        | Gin mode (debug/release)     | `debug`           | No       |
| `ALLOWED_ORIGINS` | Comma-separated CORS origins | `*`               | No       |

### CORS Configuration

To allow specific origins:
```bash
ALLOWED_ORIGINS=http://localhost:5173,https://monitoring.yourdomain.com
```

### Database Configuration

Follow the instructions in the [Monitor db](https://github.com/gabrielg2020/monitor-db)

```bash
git clone https://github.com/gabrielg2020/monitor-db
```

## Deployment

### Building Docker Image
```bash
# Build the image
docker build -t ghcr.io/yourusername/monitor-api:latest .

# Push to registry
docker push ghcr.io/yourusername/monitor-api:latest
```
### Production Deployment

See the [Deployment Guide](https://github.com/gabrielg2020/monitor-frontend/docs/DEPLOYMENT.md) in [Monitor Frontend](https://github.com/gabrielg2020/monitor-frontend) for detailed instructions.

## Related Projects

- [Monitor Frontend](https://github.com/gabrielg2020/monitor-frontend) - React dashboard
- [Monitor Agent](https://github.com/gabrielg2020/monitor-agent) - Python metric collector


## License

This project is licensed under the MIT License - see the LICENSE file for details.

---

Built with 💻 by Gabriel Guimaraes
