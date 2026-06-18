# Architecture Guide

This document describes the architectural decisions, folder structure, and patterns used in this project. Follow this guide when adding new features or modules.

---

## Overview

This project follows a **layered architecture** inspired by Spring Boot, adapted to idiomatic Go conventions. The core principle is **separation of concerns** — each layer has a single responsibility and communicates only with its adjacent layer.

```
HTTP Request → Handler → Service → Repository → SQLite
```

---

## Folder Structure

```
go-backend-architecture-template/
├── cmd/
│   └── api/
│       └── main.go                  # Application entry point
├── internal/
│   ├── app/
│   │   └── app.go                   # Dependency wiring (composition root)
│   ├── config/
│   │   └── config.go                # Configuration loading
│   ├── logger/
│   │   └── logger.go                # Structured logging setup
│   ├── database/
│   │   ├── sqlite.go                # Database connection
│   │   └── migrations.go            # Schema migrations
│   ├── transport/
│   │   └── httpserver/
│   │       ├── router.go            # Route definitions
│   │       ├── middleware.go         # HTTP middleware
│   │       └── response.go          # JSON response helpers
│   └── customer/                    # Feature module (example)
│       ├── model.go                 # Domain model
│       ├── dto.go                   # Request/Response DTOs
│       ├── handler.go               # HTTP handler (controller)
│       ├── service.go               # Business logic
│       ├── repository.go            # Database queries
│       ├── errors.go                # Domain-specific errors
│       └── validator.go             # Input validation
├── data/
│   └── app.db                       # SQLite database file
├── logs/
│   └── app.log                      # Application log file
├── docs/                            # Documentation
├── go.mod
└── README.md
```

---

## Layer Responsibilities

### `cmd/api/main.go` — Entry Point

- The **only** job of `main.go` is to call `app.Run()`.
- Contains zero business logic, zero configuration, zero wiring.
- If the application fails to start, it logs a fatal error and exits.

### `internal/app/app.go` — Composition Root

- Wires all dependencies together using **constructor injection**.
- Follows this startup sequence:
  1. Load configuration
  2. Set up logger
  3. Connect to database
  4. Run migrations
  5. Create repositories, services, and handlers
  6. Create router and start HTTP server
- This is the **only place** that knows about all packages.

### `internal/config/` — Configuration

- Reads configuration from environment variables.
- Provides sensible defaults for all settings.
- Returns a typed `Config` struct — no magic strings elsewhere.

### `internal/logger/` — Logging

- Configures Go's standard `log/slog` with JSON output.
- Writes to both stdout and a log file simultaneously.
- The configured logger is passed as a dependency — not used as a global.

### `internal/database/` — Database Layer

- `sqlite.go` — Opens the SQLite connection with WAL mode and foreign keys.
- `migrations.go` — Runs schema migrations using a simple tracking table. New migrations are appended to a slice — never modify existing migrations.

### `internal/transport/httpserver/` — HTTP Transport

- `router.go` — Registers all API routes using Go 1.22+ enhanced `ServeMux`.
- `middleware.go` — Provides request logging, panic recovery, and CORS.
- `response.go` — Defines standard JSON response formats (`Success()` and `Error()`) used by all handlers.

---

## Feature Module Pattern

Each feature (e.g., `customer`) is a self-contained package inside `internal/`. A module contains these files:

| File             | Responsibility                                      |
|------------------|------------------------------------------------------|
| `model.go`       | Domain struct that maps to the database table        |
| `dto.go`         | Request and response structs for the API contract    |
| `handler.go`     | HTTP handler — parses requests, calls service, writes responses |
| `service.go`     | Business logic — validation rules, data transformation |
| `repository.go`  | Database queries — the only layer that touches SQL   |
| `errors.go`      | Domain-specific sentinel errors                      |
| `validator.go`   | Input validation functions                           |

### Handler (HTTP Layer)

The handler is responsible for:
- Parsing the HTTP request (JSON body, path params, query params)
- Calling validators
- Calling the service
- Mapping domain errors to HTTP status codes
- Writing the JSON response

The handler **never** contains business logic or SQL.

### Service (Business Logic Layer)

The service is responsible for:
- Orchestrating business operations
- Generating IDs and timestamps
- Applying business rules
- Logging business events
- Calling the repository

The service **never** knows about HTTP or SQL.

### Repository (Data Access Layer)

The repository is responsible for:
- Executing SQL queries
- Mapping between Go types and database columns
- Returning domain errors (e.g., `ErrNotFound`) instead of raw SQL errors

The repository **never** knows about HTTP or business rules.

---

## Adding a New Module

To add a new module (e.g., `product`), follow these steps:

### 1. Create the package

```
internal/product/
├── model.go
├── dto.go
├── handler.go
├── service.go
├── repository.go
├── errors.go
└── validator.go
```

### 2. Define the model (`model.go`)

```go
package product

type Product struct {
    ID        string
    Name      string
    Price     float64
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

### 3. Add migration (`database/migrations.go`)

Append a new entry to the `allMigrations()` slice:

```go
{
    Name: "002_create_products_table",
    SQL:  `CREATE TABLE IF NOT EXISTS products (...)`,
},
```

### 4. Wire dependencies (`app/app.go`)

Add three lines in the dependency wiring section:

```go
productRepo := product.NewRepository(db)
productService := product.NewService(productRepo, log)
productHandler := product.NewHandler(productService)
```

### 5. Register routes (`transport/httpserver/router.go`)

Add the handler as a parameter and register routes:

```go
mux.HandleFunc("GET /api/products", productHandler.List)
mux.HandleFunc("POST /api/products", productHandler.Create)
```

---

## Dependency Injection

This project uses **constructor-based dependency injection** without any DI framework.

Each struct declares its dependencies as fields, and the constructor accepts them as parameters:

```go
type Service struct {
    repo   *Repository
    logger *slog.Logger
}

func NewService(repo *Repository, logger *slog.Logger) *Service {
    return &Service{repo: repo, logger: logger}
}
```

All wiring happens in `app.go`, making the dependency graph explicit and easy to understand.

---

## Error Handling

- **Domain errors** are defined as sentinel errors in each module's `errors.go`.
- **Repository** returns domain errors (e.g., `ErrNotFound`) — not raw SQL errors.
- **Service** propagates errors from the repository or creates new domain errors.
- **Handler** maps domain errors to HTTP status codes using `handleServiceError()`.

```
Repository → ErrNotFound → Service → ErrNotFound → Handler → 404 Not Found
```

---

## JSON Response Format

All API responses follow a consistent envelope format:

**Success:**
```json
{
    "success": true,
    "message": "Customer created successfully",
    "data": { ... }
}
```

**Error:**
```json
{
    "success": false,
    "message": "Validation failed",
    "error": "email is required"
}
```
