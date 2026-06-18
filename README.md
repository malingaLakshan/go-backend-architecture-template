# Go Backend Architecture Template

A clean, professional Go backend architecture template demonstrating layered design patterns with a sample Customer CRUD module. Built as a learning and reference project for teams adopting Go.

---

## Purpose

This project provides a **production-ready architecture pattern** for Go REST APIs. It follows a layered approach similar to Spring Boot, adapted to idiomatic Go conventions:

```
HTTP Request → Handler → Service → Repository → SQLite
```

The Customer CRUD module serves as a reference implementation. When adding new features, copy its structure and follow the same patterns.

---

## Quick Start

### Prerequisites

- [Go 1.23+](https://go.dev/dl/)

### Run the Server

```bash
# Clone or navigate to the project
cd go-backend-architecture-template

# Download dependencies
go mod tidy

# Start the server
go run cmd/api/main.go
```

The server starts on `http://localhost:8080` by default.

### Verify

```bash
curl http://localhost:8080/health
```

---

## Configuration

Configuration is loaded from environment variables. All settings have sensible defaults:

| Variable              | Default          | Description               |
|-----------------------|------------------|---------------------------|
| `SERVER_PORT`         | `8080`           | HTTP server port          |
| `SERVER_READ_TIMEOUT` | `15`             | Read timeout (seconds)    |
| `SERVER_WRITE_TIMEOUT`| `15`             | Write timeout (seconds)   |
| `DB_FILE_PATH`        | `./data/app.db`  | SQLite database file path |
| `LOG_LEVEL`           | `DEBUG`          | Log level (DEBUG, INFO, WARN, ERROR) |
| `LOG_FILE_PATH`       | `./logs/app.log` | Log file path             |

---

## API Endpoints

| Method   | Endpoint                | Description           |
|----------|-------------------------|-----------------------|
| `GET`    | `/health`               | Health check          |
| `POST`   | `/api/customers`        | Create a customer     |
| `GET`    | `/api/customers`        | List all customers    |
| `GET`    | `/api/customers/{id}`   | Get customer by ID    |
| `PUT`    | `/api/customers/{id}`   | Update a customer     |
| `DELETE` | `/api/customers/{id}`   | Delete a customer     |

For detailed request/response examples, see [docs/API.md](docs/API.md).

---

## Architecture Summary

| Layer        | Location                         | Responsibility                          |
|-------------|-----------------------------------|-----------------------------------------|
| Entry Point | `cmd/api/main.go`                | Start the application                   |
| App Wiring  | `internal/app/app.go`            | Wire dependencies, start server         |
| Config      | `internal/config/`               | Load environment configuration          |
| Logger      | `internal/logger/`               | Structured logging with `log/slog`      |
| Database    | `internal/database/`             | SQLite connection and migrations        |
| Transport   | `internal/transport/httpserver/`  | Router, middleware, response helpers    |
| Feature     | `internal/customer/`             | Handler → Service → Repository          |

For a deep dive, see [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md).

---

## Project Structure

```
go-backend-architecture-template/
├── cmd/
│   └── api/
│       └── main.go                  # Entry point
├── internal/
│   ├── app/
│   │   └── app.go                   # Dependency wiring
│   ├── config/
│   │   └── config.go                # Configuration
│   ├── logger/
│   │   └── logger.go                # Structured logging
│   ├── database/
│   │   ├── sqlite.go                # Database connection
│   │   └── migrations.go            # Schema migrations
│   ├── transport/
│   │   └── httpserver/
│   │       ├── router.go            # Route registration
│   │       ├── middleware.go         # HTTP middleware
│   │       └── response.go          # Response helpers
│   └── customer/
│       ├── model.go                 # Domain model
│       ├── dto.go                   # Request/Response DTOs
│       ├── handler.go               # HTTP handler
│       ├── service.go               # Business logic
│       ├── repository.go            # Database queries
│       ├── errors.go                # Domain errors
│       └── validator.go             # Input validation
├── data/                            # SQLite database (auto-created)
├── logs/                            # Log files (auto-created)
├── docs/
│   ├── ARCHITECTURE.md              # Architecture guide
│   ├── API.md                       # API documentation
│   └── BUILD_ORDER.md               # Step-by-step build guide
├── go.mod
└── README.md
```

---

## Technology Stack

| Technology                | Purpose                          |
|--------------------------|----------------------------------|
| Go 1.23+                 | Programming language             |
| `net/http` (Go 1.22+)   | HTTP server and routing          |
| `log/slog`               | Structured logging               |
| `modernc.org/sqlite`     | SQLite driver (pure Go, no CGO)  |
| `github.com/google/uuid` | UUID generation                  |

---

## Documentation

- [Architecture Guide](docs/ARCHITECTURE.md) — Layer responsibilities, patterns, how to add modules
- [API Documentation](docs/API.md) — Endpoints, request/response examples, cURL commands
- [Build Order](docs/BUILD_ORDER.md) — How to recreate this project from scratch

---

## License

This is a learning and reference project. Use freely.
