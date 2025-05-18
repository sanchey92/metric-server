# Metric Server

A production-ready metrics collection server written in Go. It receives runtime metrics from clients via HTTP, stores them in memory, and flushes them asynchronously to PostgreSQL.

## 🧩 Features

- Receives batched metrics via `POST /update`
- Accepts compressed (gzip) JSON payloads
- In-memory storage for fast ingestion
- Periodic asynchronous flushing to PostgreSQL
- Configurable via YAML and environment variables
- Graceful shutdown on SIGINT/SIGTERM
- Clean architecture with modular components


## 🛠 Requirements

- Go 1.20+
- PostgreSQL 13+
- `golangci-lint` for linting (optional)

## 📦 Installation

```bash
    git clone https://github.com/sanchey92/metric-server.git
    cd metric-server
    go build -o metric-server ./cmd/server
```

## ⚙️ Configuration

**Create a .env file at src level:**
```bash
    CONFIG_PATH=./config/config.yaml
    
    MIGRATION_DIR=./migrations
    PG_DSN="host=localhost port=5432 dbname=metrics_db user=metric_user password=metric_password sslmode=disable"
    
    POSTGRES_DB=metrics_db
    POSTGRES_USER=metric_user
    POSTGRES_PASSWORD=metric_password
    
    HTTP_HOST=localhost
    HTTP_PORT=8080
```
## 🚀 Run
```bash
    make init-deps
    make docker-up
    make local-migration-up
    make run
```

## 📄 License
This project is licensed under the MIT License (or specify another if applicable).