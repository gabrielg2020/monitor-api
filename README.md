# monitor-api

A go api used to interact with the [monitor database](https://github.com/gabrielg2020/monitor-db).

## Installation

### 1. Install Go 1.25.x from [here](https://golang.org/dl/).

    # ensure go is installed
    go version

### 2. Clone the repository

    git clone https://github.com/gabrielg2020/pi-monitor-api
    cd pi-monitor-api

### 3. Setup monitoring database

Follow the instructions [here](https://github.com/gabrielg2020/monitor-db)

### 4. Install dependencies

    go mod tidy

### 5. Configure environment variables

    cp .env.local .env
    vim .env

## Usage

### With Docker

    docker build -t monitor-api .
    docker run -d -p 8191:8191 --env-file .env --volume <PATH TO monitor-db>:/app/data --name monitor-api monitor-api

### With Docker Compose
    
    # change volume path
    vim docker-compose.yml
    docker compose up -d --build

### Without Docker

    go run main.go

## API Endpoints

The API will be available at `http://localhost:8191` (or the host machine's IP address if running in Docker).

### Available Endpoints

#### - `GET /health` - Check the health status of the API.
  - Response:
    ```json
    {
      "checks": {
        "database": "healthy"
      },
      "status": "healthy",
      "timestamp": "2025-10-16T23:52:51+01:00"
     }
    ```

#### - `GET /api/v1/metrics` - grab all monitoring data from the database.
  - Response:
    ```json
    {
      "meta": {
        "count": 1,
        "limit": 100
      },
      "records": [
        {
          "id": 1,
          "host_id": 1,
          "timestamp": 1760663031,
          "cpu_usage": 45.5,
          "memory_usage_percent": 68.2,
          "memory_total_bytes": 16777216000,
          "memory_used_bytes": 11442954240,
          "memory_available_bytes": 5334261760,
          "disk_usage_percent": 72.3,
          "disk_total_bytes": 512110190592,
          "disk_used_bytes": 370191697920,
          "disk_available_bytes": 141918492672
        },
        {...}
      ]
    }
    ```
  - Query Parameters:
    - `host_id` (optional): Filter records by host ID.
    - `start_time` (optional): Filter records with a timestamp greater than or equal to
    - `end_time` (optional): Filter records with a timestamp less than or equal to
    - `limit` (optional): Limit the number of records returned (default: 100).
    - `order` (optional): Order of records by timestamp, either `asc` or `desc` (default: `desc`).

#### - `GET /api/v1/metrics` - grab all monitoring data from the database.
- Response:
  ```json
  {
    "metric": {
      "id": 1,
      "host_id": 1,
      "timestamp": 1760663031,
      "cpu_usage": 45.5,
      "memory_usage_percent": 68.2,
      "memory_total_bytes": 16777216000,
      "memory_used_bytes": 11442954240,
      "memory_available_bytes": 5334261760,
      "disk_usage_percent": 72.3,
      "disk_total_bytes": 512110190592,
      "disk_used_bytes": 370191697920,
      "disk_available_bytes": 141918492672
     }
  }
  ```
- Query Parameters:
    - `host_id` (optional): Filter records by host ID.

#### - `POST /api/v1/metrics` - Push new monitoring data to the database.
  - Request Body (JSON):
    ```json
    {
      "record": {
        "host_id": <integer>,
        "timestamp": <integer>,
        "cpu_usage": <float>,
        "memory_usage_percent": <float>,
        "memory_total_bytes": <integer>,
        "memory_used_bytes": <integer>,
        "memory_available_bytes": <integer>,
        "disk_usage_percent": <float>,
        "disk_total_bytes": <integer>,
        "disk_used_bytes": <integer>,
        "disk_available_bytes": <integer>
      }
    }
    ```

## License

This project is licensed under the MIT License - see the LICENSE file for details.

---

Built with 💻 by Gabriel Guimaraes
