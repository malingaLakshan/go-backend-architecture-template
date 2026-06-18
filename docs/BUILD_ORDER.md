# Build Order Guide

This document explains how to recreate this project from scratch, step by step. Follow this order when building a new Go backend using this architecture pattern.

---

## Step 1 — Create Go Module

```bash
mkdir go-backend-architecture-template
cd go-backend-architecture-template
go mod init go-backend-architecture-template
```

## Step 2 — Create Folder Structure

```bash
mkdir -p cmd/api
mkdir -p internal/app
mkdir -p internal/config
mkdir -p internal/logger
mkdir -p internal/database
mkdir -p internal/transport/httpserver
mkdir -p internal/customer
mkdir -p data
mkdir -p logs
mkdir -p docs
```

## Step 3 — Create Config Package

**File:** `internal/config/config.go`

- Define `Config`, `ServerConfig`, `DatabaseConfig`, and `LoggerConfig` structs.
- Create a `Load()` function that reads from environment variables.
- Provide sensible defaults for every setting.
- No external dependencies needed — use `os.Getenv` and `strconv`.

## Step 4 — Create Logger Package

**File:** `internal/logger/logger.go`

- Use Go's standard `log/slog` package.
- Create a `Setup()` function that returns a configured `*slog.Logger`.
- Use `slog.NewJSONHandler` with `io.MultiWriter` to write to both stdout and a file.
- Parse the log level from the configuration string.

## Step 5 — Create Database Package

**File:** `internal/database/sqlite.go`

- Use `modernc.org/sqlite` (pure Go SQLite driver, no CGO).
- Create a `Connect()` function that opens the SQLite file.
- Enable WAL mode and foreign keys via PRAGMA statements.
- Ensure the data directory exists before creating the database file.

Install the dependency:

```bash
go get modernc.org/sqlite
```

## Step 6 — Create Migrations

**File:** `internal/database/migrations.go`

- Define a `migration` struct with `Name` and `SQL` fields.
- Create an `allMigrations()` function that returns an ordered slice of migrations.
- Create a `RunMigrations()` function that:
  1. Creates a `schema_migrations` tracking table
  2. Checks which migrations have already been applied
  3. Applies pending migrations in order
  4. Records each applied migration

## Step 7 — Create Customer Model and DTOs

**File:** `internal/customer/model.go`

- Define the `Customer` struct with all database fields.
- Define status constants (`ACTIVE`, `INACTIVE`).

**File:** `internal/customer/dto.go`

- Define `CreateRequest` and `UpdateRequest` structs for incoming data.
- Define a `Response` struct for outgoing data.
- Create `ToResponse()` and `ToResponseList()` conversion functions.

**File:** `internal/customer/errors.go`

- Define sentinel errors: `ErrNotFound`, `ErrDuplicateEmail`, `ErrInvalidStatus`.

**File:** `internal/customer/validator.go`

- Create `ValidateCreateRequest()` and `ValidateUpdateRequest()` functions.
- Validate required fields and email format using `net/mail.ParseAddress`.

## Step 8 — Create Customer Repository

**File:** `internal/customer/repository.go`

- Define a `Repository` struct with a `*sql.DB` field.
- Create a `NewRepository()` constructor.
- Implement CRUD methods: `Create`, `GetByID`, `List`, `Update`, `Delete`.
- Use `context.Context` in all methods.
- Map SQL errors to domain errors (e.g., UNIQUE violation → `ErrDuplicateEmail`).
- Handle `time.Time` ↔ `TEXT` conversion for SQLite.

## Step 9 — Create Customer Service

**File:** `internal/customer/service.go`

- Define a `Service` struct with `*Repository` and `*slog.Logger` fields.
- Create a `NewService()` constructor.
- Implement business logic methods: `Create`, `GetByID`, `List`, `Update`, `Delete`.
- Generate UUIDs for new customers.
- Set timestamps (created_at, updated_at).
- Log important operations.

Install the UUID dependency:

```bash
go get github.com/google/uuid
```

## Step 10 — Create Customer Handler

**File:** `internal/customer/handler.go`

- Define a `Handler` struct with a `*Service` field.
- Create a `NewHandler()` constructor.
- Implement HTTP handler methods: `Create`, `List`, `GetByID`, `Update`, `Delete`.
- Parse JSON request bodies.
- Extract path parameters using `r.PathValue()` (Go 1.22+).
- Call validators, then the service.
- Map domain errors to HTTP status codes.
- Use the response helpers from `httpserver` package.

## Step 11 — Create Router and Middleware

**File:** `internal/transport/httpserver/response.go`

- Define `SuccessResponse` and `ErrorResponse` structs.
- Create `JSON()`, `Success()`, and `Error()` helper functions.

**File:** `internal/transport/httpserver/middleware.go`

- Create `RequestLogger()` — logs method, path, status, and duration.
- Create `PanicRecovery()` — catches panics and returns a 500 response.
- Create `CORS()` — adds basic CORS headers for development.

**File:** `internal/transport/httpserver/router.go`

- Create `NewRouter()` that accepts the logger and all handler dependencies.
- Register the health check route.
- Register all customer routes.
- Apply middleware in order: PanicRecovery → RequestLogger → CORS.

## Step 12 — Wire Dependencies in app.go

**File:** `internal/app/app.go`

- Create a `Run()` function that:
  1. Loads configuration
  2. Sets up the logger
  3. Connects to the database
  4. Runs migrations
  5. Creates repository → service → handler for each module
  6. Creates the router
  7. Starts the HTTP server
  8. Handles graceful shutdown on SIGINT/SIGTERM

## Step 13 — Start Server from main.go

**File:** `cmd/api/main.go`

- Call `app.Run()`.
- If it returns an error, log it with `log.Fatalf` and exit.
- This file should be ~10 lines of code.

## Step 14 — Test APIs

```bash
# Start the server
go run cmd/api/main.go

# Test health check
curl http://localhost:8080/health

# Create a customer
curl -X POST http://localhost:8080/api/customers \
  -H "Content-Type: application/json" \
  -d '{"name": "Jane Doe", "email": "jane@example.com", "phone": "+1-555-0100"}'

# List all customers
curl http://localhost:8080/api/customers

# Get a customer (replace <id> with actual UUID)
curl http://localhost:8080/api/customers/<id>

# Update a customer
curl -X PUT http://localhost:8080/api/customers/<id> \
  -H "Content-Type: application/json" \
  -d '{"name": "Jane Smith", "email": "jane@example.com", "phone": "+1-555-0100", "status": "INACTIVE"}'

# Delete a customer
curl -X DELETE http://localhost:8080/api/customers/<id>
```

---

## Summary

The build order follows the dependency graph bottom-up:

```
config → logger → database → migrations → model/dto/errors → repository → service → handler → router → app → main
```

Each layer only depends on the layer below it. Start from the foundations and work your way up.
