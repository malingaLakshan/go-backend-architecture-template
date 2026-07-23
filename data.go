# Resonate Replay Engine

The **Resonate Replay Engine (RRE)** is a command-line testing tool that replays recorded RFID and location-reader data from a Recorder SQLite database into a target Resonate HTTP instance.

Before replay begins, RRE retrieves the target SiteGraph and validates it against the site configuration captured during recording. Playback starts only when the recorded and target sites are structurally compatible.

---

## Table of Contents

1. [Overview](#overview)
2. [Key Features](#key-features)
3. [High-Level Architecture](#high-level-architecture)
4. [SiteGraph Structure](#sitegraph-structure)
5. [Project Structure](#project-structure)
6. [Configuration](#configuration)
7. [Command Reference](#command-reference)
8. [Validation Rules](#validation-rules)
9. [Replay Sequence](#replay-sequence)
10. [Mock Target Server](#mock-target-server)
11. [Playback and Pacing](#playback-and-pacing)
12. [Logging and Error Handling](#logging-and-error-handling)
13. [Build and Test](#build-and-test)
14. [Security and Repository Hygiene](#security-and-repository-hygiene)

---

## Overview

RRE supports four main workflows:

| Workflow | Purpose |
|---|---|
| Summary | Inspect the recording database and display session statistics |
| Mock target | Host SiteGraph data and receive replay payloads locally |
| Validation | Compare the recorded SiteGraph with the target SiteGraph |
| Playback | Replay recorded RawReads using their original timing |

The normal execution flow is:

```text
Load configuration
        ↓
Read recorded site information from SQLite
        ↓
Retrieve the target SiteGraph
        ↓
Validate structural compatibility
        ↓
Load RawReads in chronological order
        ↓
Replay payloads using recorded timing
        ↓
Display and log the replay result
```

---

## Key Features

- Reads Recorder SQLite databases.
- Displays recording metadata and RawRead statistics.
- Supports configuration-file and direct argument modes.
- Accepts relative and absolute config file paths.
- Retrieves target SiteGraphs through HTTP.
- Provides a local mock target server for QA testing.
- Automatically discovers SiteGraph JSON files.
- Preserves original SiteGraph JSON content.
- Validates recursive nested Regions.
- Validates Readers directly under Floors.
- Validates Antenna Ports under the correct Reader.
- Allows the target site to contain additional structures.
- Stops playback when required structures are missing.
- Reconstructs reader-bundle payloads.
- Replays RawReads using their recorded injection timing.
- Records validation, playback and error information through the logger.
- Supports build-time version injection.

---

## High-Level Architecture

```mermaid
flowchart LR
    User[QA Engineer or Developer]

    subgraph RRE[Resonate Replay Engine]
        CLI[CLI Layer]
        Config[Configuration Loader]
        Recording[Recording Repository]
        Validator[SiteGraph Validator]
        Replay[Replay Service]
        Pacing[Pacing Engine]
        Logger[Logger]
        MockServer[Mock Target Server]
        SiteStore[SiteGraph Store]
    end

    SQLite[(Recorder SQLite Database)]
    SiteFiles[(SiteGraph JSON Files)]
    Target[Target Resonate HTTP Instance]
    Logs[(Logs and JSONL Output)]

    User -->|summary, validate, play, serve| CLI

    CLI --> Config
    CLI --> Recording
    CLI --> Validator
    CLI --> Replay
    CLI --> MockServer

    Recording --> SQLite

    Validator --> Recording
    Validator -->|GET /sites/{siteId}| Target

    Replay --> Recording
    Replay --> Pacing
    Replay -->|POST /reader-bundles| Target

    MockServer --> SiteStore
    SiteStore --> SiteFiles

    CLI --> Logger
    Validator --> Logger
    Replay --> Logger
    MockServer --> Logger
    Logger --> Logs
```

---

## SiteGraph Structure

A SiteGraph represents the complete configuration of one Resonate site.

The supported structure is:

```text
Site
└── Floors
    ├── Readers
    │   └── Antenna Ports
    └── Regions
        └── Child Regions recursively
```

Important rules:

- Readers belong directly to Floors.
- Readers are not nested inside Regions.
- Antenna Ports belong to Readers.
- Regions belong to Floors.
- Regions may contain nested child Regions recursively.
- The target SiteGraph may contain additional Floors, Regions, Readers or Antenna Ports.
- Every structure required by the recorded SiteGraph must exist in the target SiteGraph.

---

## Project Structure

```text
resonate-replay-engine/
├── cmd/
│   └── rre/
│       └── main.go
│
├── configs/
│   └── config.json
│
├── data/
│   ├── recording_001.sqlite
│   └── sites/
│       └── Bentonville_SiteGraph.json
│
├── docs/
│   └── ARCHITECTURE.md
│
├── internal/
│   ├── cli/
│   │   ├── args.go
│   │   └── commands.go
│   │
│   ├── config/
│   │   └── model.go
│   │
│   ├── logger/
│   │   └── logger.go
│   │
│   ├── mocktarget/
│   │   ├── handler.go
│   │   ├── handler_test.go
│   │   ├── server.go
│   │   ├── site_store.go
│   │   └── site_store_test.go
│   │
│   ├── recording/
│   │   ├── model.go
│   │   ├── repository.go
│   │   ├── repository_test.go
│   │   └── summary.go
│   │
│   ├── replay/
│   │   ├── injector.go
│   │   ├── injector_test.go
│   │   ├── model.go
│   │   ├── pacing.go
│   │   ├── pacing_test.go
│   │   ├── payload.go
│   │   ├── payload_test.go
│   │   ├── service.go
│   │   └── status.go
│   │
│   ├── site/
│   │   ├── client.go
│   │   ├── model.go
│   │   ├── parser.go
│   │   ├── parser_test.go
│   │   ├── summary.go
│   │   ├── validation_model.go
│   │   ├── validator.go
│   │   └── validator_test.go
│   │
│   └── sqlite/
│
├── logs/
├── .gitignore
├── go.mod
├── go.sum
└── README.md
```

---

## Configuration

RRE supports:

1. Configuration-file mode.
2. Direct argument mode.

Configuration-file mode is recommended for repeatable QA testing.

### Example `config.json`

```json
{
  "recording_file": "data/recording_001.sqlite",
  "target_url": "http://localhost:8080",
  "site_id": "b3489888-aacf-4451-893c-d7d994240f93",
  "mock_port": 8080,
  "site_graph_directory": "data/sites"
}
```

### Configuration Fields

| Field | Required | Used by | Description |
|---|---:|---|---|
| `recording_file` | Yes | summary, validate, play | Path to the Recorder SQLite database |
| `target_url` | Yes | validate, play | Base URL of the target Resonate HTTP instance |
| `site_id` | Yes | validate, play | Site ID used to retrieve and validate the target SiteGraph |
| `mock_port` | No | serve | Mock server port; defaults to `8080` |
| `site_graph_directory` | No | serve | Directory containing SiteGraph JSON files; defaults to `data/sites` |

The config file may be stored in any valid readable location.

Relative path:

```powershell
.\rre.exe validate --config .\configs\config.json
```

Absolute path:

```powershell
.\rre.exe validate --config "C:\QA\RRE\config.json"
```

Use quotation marks when the path contains spaces.

---

## Command Reference

Display available commands:

```powershell
.\rre.exe help
```

Display version:

```powershell
.\rre.exe version
```

### Summary

Config mode:

```powershell
.\rre.exe summary --config .\configs\config.json
```

Direct argument mode:

```powershell
.\rre.exe summary `
  --file .\data\recording_001.sqlite
```

### Serve

Starts the local mock target server.

```powershell
.\rre.exe serve --config .\configs\config.json
```

### Validate

Config mode:

```powershell
.\rre.exe validate --config .\configs\config.json
```

Direct argument mode:

```powershell
.\rre.exe validate `
  --file .\data\recording_001.sqlite `
  --target-url http://localhost:8080 `
  --site-id b3489888-aacf-4451-893c-d7d994240f93
```

### Play

Config mode:

```powershell
.\rre.exe play --config .\configs\config.json
```

Direct argument mode:

```powershell
.\rre.exe play `
  --file .\data\recording_001.sqlite `
  --target-url http://localhost:8080 `
  --site-id b3489888-aacf-4451-893c-d7d994240f93
```

---

## Validation Rules

Before playback starts, RRE validates the recorded site against the target SiteGraph.

The following checks are performed:

1. The configured Site ID must match the Site ID stored in the recording.
2. The recorded SiteGraph root ID must match the recorded Site ID.
3. The target SiteGraph root ID must match the requested Site ID.
4. Every recorded Floor ID must exist in the target.
5. Every recorded Region ID must exist under the correct Floor.
6. Nested Regions must be validated recursively.
7. Every recorded Reader ID must exist directly under the correct Floor.
8. Every recorded Antenna Port must exist under the correct Reader.
9. Additional structures in the target are allowed.
10. Playback is aborted immediately when validation fails.

Example successful result:

```text
Validation Results

✓ Site ID matched
✓ 1 of 1 required Floors matched
✓ 18 of 18 required Regions matched
✓ 20 of 20 required Readers matched
✓ 24 of 24 required Antenna Ports matched

Validation passed.
```

Example failure:

```text
Validation Results

✓ Site ID matched
✓ 1 of 1 required Floors matched
✗ 17 of 18 required Regions matched

Validation failed.

Missing Region
  Region ID: 3caf648f-fce3-424e-bc08-e711451ddaab
  Floor ID:  b2c98296-6c4e-4a52-8380-17a11ddb2b2c
```

---

## Replay Sequence

```mermaid
sequenceDiagram
    autonumber

    actor User as QA Engineer or Developer
    participant CLI as RRE CLI
    participant Config as Config Loader
    participant DB as Recorder SQLite
    participant API as Target Resonate API
    participant Parser as SiteGraph Parser
    participant Validator as SiteGraph Validator
    participant Replay as Replay Service
    participant Pacing as Pacing Engine
    participant Logger as Logger

    User->>CLI: rre play --config <path>

    CLI->>Config: Load configuration
    Config-->>CLI: Recording file, target URL and Site ID

    CLI->>DB: Read recorded SiteInformation
    DB-->>CLI: Recorded Site ID and SiteGraph JSON

    CLI->>API: GET /sites/{siteId}
    API-->>CLI: Target SiteGraph JSON

    CLI->>Parser: Parse recorded SiteGraph
    Parser-->>CLI: Recorded site model

    CLI->>Parser: Parse target SiteGraph
    Parser-->>CLI: Target site model

    CLI->>Validator: Validate structural compatibility

    Validator->>Validator: Validate Site IDs
    Validator->>Validator: Validate Floors
    Validator->>Validator: Validate recursive Regions
    Validator->>Validator: Validate floor-level Readers
    Validator->>Validator: Validate Antenna Ports

    alt Validation fails
        Validator-->>CLI: Validation failure details
        CLI->>Logger: Log validation failure
        CLI-->>User: Abort playback and display mismatches
    else Validation passes
        Validator-->>CLI: Validation successful
        CLI->>Logger: Log playback start

        CLI->>DB: Read RawReads ordered by injection time
        DB-->>CLI: Ordered RawReads

        CLI->>Replay: Start replay

        loop For every RawRead
            Replay->>Pacing: Calculate scheduled send time
            Pacing-->>Replay: Wait until scheduled time
            Replay->>API: POST /reader-bundles
            API-->>Replay: HTTP response

            alt Request succeeds
                Replay->>Logger: Record successful injection
            else Request fails
                Replay->>Logger: Record error and timestamp
                Note over Replay,Pacing: Continue with remaining records
            end
        end

        Replay-->>CLI: Replay result
        CLI->>Logger: Log playback completion
        CLI-->>User: Display replay summary
    end
```

---

## Mock Target Server

The mock target server supports local testing without requiring a live Resonate environment.

Start the server:

```powershell
.\rre.exe serve --config .\configs\config.json
```

### SiteGraph Auto-Discovery

When the server starts:

1. It scans the configured SiteGraph directory.
2. Every directly contained `.json` file is inspected.
3. The root `id` field is used as the Site ID.
4. Duplicate Site IDs are rejected.
5. Invalid JSON files are rejected.
6. Files without a root Site ID are rejected.
7. The original JSON content is preserved.
8. Restarting the server reloads the directory.

### Endpoints

#### `GET /sites`

Returns a summary of all loaded sites.

```powershell
Invoke-RestMethod http://localhost:8080/sites |
  ConvertTo-Json -Depth 10
```

#### `GET /sites/{siteId}`

Returns the complete SiteGraph for the requested Site ID.

```powershell
Invoke-RestMethod `
  http://localhost:8080/sites/b3489888-aacf-4451-893c-d7d994240f93 |
  ConvertTo-Json -Depth 100
```

#### `POST /reader-bundles`

Receives replay payloads from the `play` command.

---

## Playback and Pacing

RawReads are loaded in chronological order using their recorded injection timestamps.

The first record is sent when playback starts. Each following record is scheduled according to its original offset from the first recorded injection timestamp.

This preserves the original timing pattern while reducing cumulative drift caused by request processing and operating-system scheduling.

The replay duration may be slightly longer than the recorded duration because it includes:

- JSON processing;
- HTTP request and response time;
- logging;
- timer precision;
- operating-system scheduling;
- target processing time.

A small difference is expected.

---

## Logging and Error Handling

RRE logs important execution events, including:

- configuration loading failures;
- validation start and result;
- playback start;
- HTTP request failures;
- playback completion;
- successful and failed record totals.

Expected behavior:

- Invalid configuration stops command execution.
- Missing recording files return a clear error.
- Invalid SiteGraph files are rejected.
- Validation failure stops playback before RawReads are replayed.
- Individual HTTP failures are logged.
- Later records continue to be scheduled.
- The final summary shows successful and failed totals.

Generated logs are written under:

```text
logs/
```

---

## Build and Test

### Download Dependencies

```powershell
go mod download
```

### Format Code

```powershell
gofmt -w .
```

### Run Tests

```powershell
go test ./...
```

Verbose mode:

```powershell
go test ./... -v
```

### Run Static Analysis

```powershell
go vet ./...
```

### Build

```powershell
go build -o rre.exe .\cmd\rre
```

### Build with Version Injection

```powershell
go build `
  -ldflags="-X main.version=v0.2.1" `
  -o rre.exe `
  .\cmd\rre
```

Verify:

```powershell
.\rre.exe version
```

---

## Security and Repository Hygiene

Do not commit customer data, generated output or runtime files without explicit approval.

Avoid committing:

- real customer SiteGraph files;
- real Recorder SQLite databases;
- SQLite shared-memory and write-ahead-log files;
- application logs;
- received replay payloads;
- generated executables;
- temporary files;
- backup files;
- credentials or private target URLs.

Recommended `.gitignore` entries:

```gitignore
# Logs
logs/*.log
logs/*.jsonl

# Recorder data
data/*.sqlite
data/*.db
data/*.sqlite-shm
data/*.sqlite-wal

# Generated binaries
*.exe
*.bin

# Temporary files
*.tmp
*.backup
```

Note: `.gitignore` does not stop tracking a file that was committed earlier.

To stop tracking an existing generated file:

```powershell
git rm --cached logs/received_payloads.jsonl
```

---

## Summary

RRE provides a repeatable workflow for:

1. inspecting Recorder SQLite data;
2. hosting mock SiteGraphs;
3. validating recorded and target site compatibility;
4. replaying reader data using recorded timing;
5. collecting clear logs and replay results.

Validation is completed before playback begins, ensuring that incompatible target configurations are detected before payload injection.