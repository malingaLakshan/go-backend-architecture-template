# Resonate Replay Engine

The **Resonate Replay Engine (RRE)** is a deterministic command-line testing tool that replays recorded RFID and location-reader data from a Recorder SQLite database into a target Resonate HTTP instance.

Before replay begins, RRE retrieves the target SiteGraph and validates it against the site configuration captured during recording. Playback starts only when the recorded and target sites are structurally compatible.

---

## Table of Contents

1. [Overview](#overview)
2. [Key Features](#key-features)
3. [High-Level Architecture](#high-level-architecture)
4. [Replay Sequence](#replay-sequence)
5. [Validation Sequence](#validation-sequence)
6. [SiteGraph Structure](#sitegraph-structure)
7. [Validation Rules](#validation-rules)
8. [Component Architecture](#component-architecture)
9. [Project Structure](#project-structure)
10. [Configuration](#configuration)
11. [Command Reference](#command-reference)
12. [Mock Target Server](#mock-target-server)
13. [Recording Summary](#recording-summary)
14. [Playback and Pacing](#playback-and-pacing)
15. [Logging and Error Handling](#logging-and-error-handling)
16. [QA Testing Workflow](#qa-testing-workflow)
17. [Build and Test](#build-and-test)
18. [Release Build](#release-build)
19. [Security and Repository Hygiene](#security-and-repository-hygiene)
20. [Troubleshooting](#troubleshooting)

---

## Overview

RRE supports four primary workflows:

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
Read the recorded site from SQLite
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
- Displays recording metadata and table counts.
- Loads configuration from relative or absolute JSON file paths.
- Supports direct command-line arguments.
- Retrieves target SiteGraphs through HTTP.
- Provides a local mock target server for QA testing.
- Automatically discovers multiple SiteGraph JSON files.
- Preserves complete SiteGraph JSON without re-marshalling.
- Validates recursive nested regions.
- Validates Readers directly under Floors.
- Validates Antenna Ports under their correct Readers.
- Allows the target site to contain additional structures.
- Prevents playback when required recorded structures are missing.
- Reconstructs standard reader-bundle payloads.
- Replays RawReads according to recorded injection timestamps.
- Continues processing later records when an individual HTTP request fails.
- Produces clear validation and replay summaries.
- Maintains application logs and replay error information.
- Supports version injection during the build process.

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
    LogFiles[(Application Logs and JSONL Files)]

    User -->|summary, validate, play, serve| CLI

    CLI --> Config
    CLI --> Recording
    CLI --> Validator
    CLI --> Replay
    CLI --> MockServer

    Recording --> SQLite

    Validator --> Recording
    Validator -->|GET /sites/siteId| Target

    Replay --> Recording
    Replay --> Pacing
    Replay -->|POST /reader-bundles| Target

    MockServer --> SiteStore
    SiteStore --> SiteFiles

    Config --> Recording
    Config --> Validator
    Config --> Replay
    Config --> MockServer

    CLI --> Logger
    Validator --> Logger
    Replay --> Logger
    MockServer --> Logger
    Logger --> LogFiles
```

---

## Replay Sequence

The `play` command validates the recorded site before loading and replaying RawReads.

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

    User->>CLI: rre play --config path

    CLI->>Config: Load configuration file
    Config-->>CLI: Recording path, target URL and Site ID

    CLI->>DB: Read SiteInformation
    DB-->>CLI: Recorded Site ID and SiteGraph JSON

    CLI->>API: GET /sites/{siteId}
    API-->>CLI: Target SiteGraph JSON

    CLI->>Parser: Parse recorded SiteGraph
    Parser-->>CLI: Recorded site model

    CLI->>Parser: Parse target SiteGraph
    Parser-->>CLI: Target site model

    CLI->>Validator: Compare recorded and target sites

    Validator->>Validator: Validate configured and recorded Site IDs
    Validator->>Validator: Validate target root Site ID
    Validator->>Validator: Validate required Floors
    Validator->>Validator: Validate recursive Regions
    Validator->>Validator: Validate floor-level Readers
    Validator->>Validator: Validate Reader Antenna Ports

    alt Validation fails
        Validator-->>CLI: Validation failure with mismatch details
        CLI->>Logger: Log validation failure
        CLI-->>User: Display failure and abort playback
    else Validation passes
        Validator-->>CLI: Validation successful
        CLI->>Logger: Log playback start

        CLI->>DB: Read RawReads ordered by Injection Time
        DB-->>CLI: Chronologically ordered records

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
                Note over Replay,Pacing: Continue scheduling remaining records
            end
        end

        Replay-->>CLI: Replay result
        CLI->>Logger: Log playback completion
        CLI-->>User: Display replay summary
    end
```

---

## Validation Sequence

```mermaid
sequenceDiagram
    autonumber

    actor User
    participant CLI as RRE CLI
    participant Config as Config Loader
    participant DB as Recorder SQLite
    participant API as Target Resonate API
    participant Parser as SiteGraph Parser
    participant Validator as Validator
    participant Logger as Logger

    User->>CLI: rre validate --config path

    CLI->>Config: Load configuration
    Config-->>CLI: Recording file, target URL and Site ID

    CLI->>DB: Read SiteInformation
    DB-->>CLI: Recorded site configuration

    CLI->>API: GET /sites/{siteId}
    API-->>CLI: Target SiteGraph

    CLI->>Parser: Parse recorded configuration
    Parser-->>CLI: Recorded SiteGraph model

    CLI->>Parser: Parse target configuration
    Parser-->>CLI: Target SiteGraph model

    CLI->>Validator: Validate compatibility

    Validator->>Validator: Compare Site IDs
    Validator->>Validator: Compare required Floors
    Validator->>Validator: Compare recursive Regions
    Validator->>Validator: Compare floor-level Readers
    Validator->>Validator: Compare Antenna Ports

    alt Structurally compatible
        Validator-->>CLI: Validation passed
        CLI->>Logger: Log validation success
        CLI-->>User: Display matching result
    else Structurally incompatible
        Validator-->>CLI: Validation failed with missing structures
        CLI->>Logger: Log validation failure
        CLI-->>User: Display detailed mismatch information
    end
```

---

## SiteGraph Structure

A SiteGraph represents the complete configuration of one Resonate site.

The supported hierarchy is:

```text
Site
└── Floors
    ├── Readers
    │   └── Antenna Ports
    └── Regions
        └── Child Regions
            └── Child Regions recursively
```

Important structural rules:

- Readers belong directly to Floors.
- Readers are not stored inside Regions.
- Antenna Ports belong to Readers.
- Regions belong to Floors.
- A Region may contain additional child Regions recursively.
- The target SiteGraph may contain additional Floors, Regions, Readers or Antenna Ports.
- Every structure required by the recorded SiteGraph must exist in the target SiteGraph.

### SiteGraph hierarchy diagram

```mermaid
flowchart TB
    Site[Site]
    Floors[Floors]
    Floor[Floor]

    Readers[Readers]
    Reader[Reader]
    Ports[Antenna Ports]

    Regions[Regions]
    Region[Region]
    ChildRegions[Child Regions]
    NestedRegion[Nested Region]
    FurtherRegions[Further Nested Regions]

    Site --> Floors
    Floors --> Floor

    Floor --> Readers
    Readers --> Reader
    Reader --> Ports

    Floor --> Regions
    Regions --> Region
    Region --> ChildRegions
    ChildRegions --> NestedRegion
    NestedRegion --> FurtherRegions

    ReaderNote[Readers are direct children of a Floor]
    RegionNote[Regions may contain Regions recursively]

    Readers -.-> ReaderNote
    ChildRegions -.-> RegionNote
```

---

## Recorded Data and Playback Data

The Recorder SQLite database contains the data required for validation and playback.

| Table | Important data | Purpose |
|---|---|---|
| `RecordingSession` | Recording metadata and build information | Describes the recording session |
| `SiteInformation` | Recorded Site ID and SiteGraph JSON | Used for structural validation |
| `RawReads` | Raw payload and injection timestamp | Used for deterministic playback |
| `MLT_SOW_Locations` | Recorded location information | Preserved recording data |
| `ResonateEvents` | Recorded Resonate events | Preserved recording data |
| `Snapshots` | Captured snapshots | Preserved recording data |

Validation uses the SiteGraph stored in `SiteInformation`.

Playback uses the RawRead payloads stored in `RawReads`, ordered by their recorded injection timestamps.

---

## Validation Rules

Before playback starts, RRE verifies that the recorded site is compatible with the target site.

The following rules are applied:

1. The configured Site ID must match the Site ID stored in the recording.
2. The recorded SiteGraph root ID must match the recorded Site ID.
3. The target SiteGraph root ID must match the requested Site ID.
4. Every recorded Floor ID must exist in the target site.
5. Every recorded Region ID must exist under the correct Floor.
6. Nested Regions must be validated recursively.
7. Every recorded Reader ID must exist directly under the correct Floor.
8. Every recorded Antenna Port must exist under the correct Reader.
9. Additional target structures are allowed.
10. Playback must stop immediately when validation fails.

### Validation flow

```mermaid
flowchart TD
    Start([Validation starts])
    LoadConfig[Load configuration]
    ReadRecorded[Read SiteInformation from SQLite]
    FetchTarget[Retrieve target SiteGraph]
    Parse[Parse both SiteGraphs]

    CheckSite{Site IDs match?}
    CheckFloors{All required Floors exist?}
    CheckRegions{All recursive Regions exist under the correct Floor?}
    CheckReaders{All Readers exist under the correct Floor?}
    CheckPorts{All Antenna Ports exist under the correct Reader?}

    Passed[Validation passed]
    Failed[Validation failed]
    Details[Return detailed mismatch information]
    End([End])

    Start --> LoadConfig
    LoadConfig --> ReadRecorded
    ReadRecorded --> FetchTarget
    FetchTarget --> Parse

    Parse --> CheckSite

    CheckSite -- No --> Failed
    CheckSite -- Yes --> CheckFloors

    CheckFloors -- No --> Failed
    CheckFloors -- Yes --> CheckRegions

    CheckRegions -- No --> Failed
    CheckRegions -- Yes --> CheckReaders

    CheckReaders -- No --> Failed
    CheckReaders -- Yes --> CheckPorts

    CheckPorts -- No --> Failed
    CheckPorts -- Yes --> Passed

    Failed --> Details
    Details --> End
    Passed --> End
```

### Example successful validation

```text
Site Configuration Validation

Recorded Site
  Site ID:        b3489888-aacf-4451-893c-d7d994240f93
  Floors:        1
  Regions:       18
  Readers:       20
  Antenna Ports: 24

Target Site
  Site ID:        b3489888-aacf-4451-893c-d7d994240f93
  Floors:        1
  Regions:       18
  Readers:       20
  Antenna Ports: 24

Validation Results

✓ Site ID matched
✓ 1 of 1 required Floors matched
✓ 18 of 18 required Regions matched
✓ 20 of 20 required Readers matched
✓ 24 of 24 required Antenna Ports matched

Validation passed.
The recorded site is structurally compatible with the target site.
```

### Example failed validation

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

## Component Architecture

```mermaid
flowchart TB
    Main[cmd/rre/main.go]

    subgraph CLIPackage[internal/cli]
        Args[args.go<br/>Argument parsing]
        Commands[commands.go<br/>Command orchestration]
    end

    subgraph ConfigPackage[internal/config]
        ConfigModel[model.go<br/>Configuration loading and saving]
    end

    subgraph RecordingPackage[internal/recording]
        RecordingModel[model.go<br/>Recording models]
        Repository[repository.go<br/>SQLite repository]
        RecordingSummary[summary.go<br/>Recording summary]
    end

    subgraph ReplayPackage[internal/replay]
        ReplayModel[model.go]
        ReplayService[service.go<br/>Replay orchestration]
        Payload[payload.go<br/>Payload construction]
        Pacing[pacing.go<br/>Playback scheduling]
        Injector[injector.go<br/>HTTP injection]
        Status[status.go<br/>Result tracking]
    end

    subgraph SitePackage[internal/site]
        SiteClient[client.go<br/>Target SiteGraph client]
        SiteModel[model.go]
        SiteParser[parser.go<br/>Recursive SiteGraph parser]
        SiteValidator[validator.go<br/>Structural validation]
        ValidationModel[validation_model.go]
        SiteSummary[summary.go]
    end

    subgraph MockTargetPackage[internal/mocktarget]
        MockServer[server.go<br/>HTTP server]
        Handler[handler.go<br/>Request handling]
        SiteStore[site_store.go<br/>SiteGraph discovery and storage]
    end

    subgraph LoggerPackage[internal/logger]
        Logger[logger.go<br/>Application logging]
    end

    SQLitePackage[internal/sqlite]
    Database[(Recorder SQLite)]
    SiteFiles[(data/sites/*.json)]
    TargetAPI[Target Resonate API]
    Logs[(logs)]

    Main --> Args
    Main --> Commands

    Commands --> ConfigModel
    Commands --> Repository
    Commands --> RecordingSummary
    Commands --> SiteClient
    Commands --> SiteValidator
    Commands --> ReplayService
    Commands --> MockServer
    Commands --> Logger

    Repository --> SQLitePackage
    SQLitePackage --> Database

    SiteClient --> TargetAPI
    SiteClient --> SiteParser
    SiteValidator --> SiteParser
    SiteValidator --> ValidationModel
    SiteSummary --> SiteModel

    ReplayService --> Repository
    ReplayService --> Payload
    ReplayService --> Pacing
    ReplayService --> Injector
    ReplayService --> Status
    Injector --> TargetAPI

    MockServer --> Handler
    Handler --> SiteStore
    SiteStore --> SiteFiles

    Logger --> Logs
```

---

## Project Structure

```text
resonate-replay-engine/
├── .vscode/
│   ├── launch.json
│   └── settings.json
│
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
│   ├── rre.log
│   └── received_payloads.jsonl
│
├── payloadsDir/
├── .gitignore
├── go.mod
├── go.sum
└── README.md
```

Generated databases, logs, JSONL files and executables should not be committed unless explicitly approved.

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

### Configuration fields

| Field | Required | Used by | Description |
|---|---:|---|---|
| `recording_file` | Yes | summary, validate, play | Path to the Recorder SQLite database |
| `target_url` | Yes | validate, play | Base URL of the target Resonate HTTP instance |
| `site_id` | Yes | validate, play | Site ID used to retrieve and validate the target SiteGraph |
| `mock_port` | No | serve | Mock server port; defaults to `8080` |
| `site_graph_directory` | No | serve | Directory containing SiteGraph JSON files; defaults to `data/sites` |

### Config file location

The configuration file may be stored in a user-defined location.

Relative path:

```powershell
.\rre.exe validate --config .\configs\config.json
```

Absolute path:

```powershell
.\rre.exe validate --config "C:\QA\RRE\config.json"
```

The config file:

- must exist;
- must be a regular file;
- must have a `.json` extension;
- must contain valid JSON;
- must be readable by the current user.

Use quotation marks when the path contains spaces.

---

## Command Reference

Display all available commands:

```powershell
.\rre.exe help
```

Display version information:

```powershell
.\rre.exe version
```

### `summary`

Displays recording metadata, table counts and RawRead statistics.

Config mode:

```powershell
.\rre.exe summary --config .\configs\config.json
```

Direct argument mode:

```powershell
.\rre.exe summary `
  --file .\data\recording_001.sqlite
```

### `serve`

Starts the local mock target server.

```powershell
.\rre.exe serve --config .\configs\config.json
```

The server loads SiteGraph JSON files from `site_graph_directory`.

### `validate`

Compares the recorded site with the target SiteGraph.

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

### `play`

Validates the site and replays the recorded RawReads.

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

## Mock Target Server

The mock target server supports local validation and playback testing without requiring a live Resonate environment.

Start the server:

```powershell
.\rre.exe serve --config .\configs\config.json
```

Example startup output:

```text
Resonate Replay Engine Mock Server

SiteGraph Directory:
  data/sites

Loaded SiteGraphs: 1

Available Sites

1. Bentonville
   Site ID:       b3489888-aacf-4451-893c-d7d994240f93
   Floors:       1
   Regions:      18
   Readers:      20
   Antenna Ports: 24

Server:
  http://localhost:8080
```

### SiteGraph auto-discovery

When the mock server starts:

1. It scans the configured SiteGraph directory.
2. Every directly contained `.json` file is inspected.
3. The root `id` field is used as the Site ID.
4. Duplicate Site IDs are rejected.
5. Invalid JSON files are rejected.
6. Files without a root Site ID are rejected.
7. Non-JSON files are ignored or rejected according to validation rules.
8. The original JSON bytes are preserved.
9. Restarting the server reloads the directory contents.

To add a mock site:

1. Copy its SiteGraph JSON file into `data/sites`.
2. Confirm it has a unique root `id`.
3. Restart the mock server.

### Endpoints

#### `GET /sites`

Returns a summary of all loaded mock sites.

```powershell
Invoke-RestMethod http://localhost:8080/sites |
  ConvertTo-Json -Depth 10
```

Example response:

```json
{
  "sites": [
    {
      "id": "b3489888-aacf-4451-893c-d7d994240f93",
      "name": "Bentonville",
      "floors": 1,
      "regions": 18,
      "readers": 20,
      "antenna_ports": 24
    }
  ]
}
```

#### `GET /sites/{siteId}`

Returns the complete, original SiteGraph JSON for the requested Site ID.

```powershell
Invoke-RestMethod `
  http://localhost:8080/sites/b3489888-aacf-4451-893c-d7d994240f93 |
  ConvertTo-Json -Depth 100
```

Unknown Site ID response:

```json
{
  "error": "site not found",
  "site_id": "unknown-id"
}
```

#### `POST /reader-bundles`

Receives replay payloads sent by the `play` command.

For local QA verification, received payloads may be appended to:

```text
logs/received_payloads.jsonl
```

Generated payload files should not be committed.

---

## Recording Summary

The `summary` command helps confirm that the SQLite database contains the expected recording information before validation or playback.

```powershell
.\rre.exe summary --config .\configs\config.json
```

Example output:

```text
Resonate Site ID:   b3489888-aacf-4451-893c-d7d994240f93
Resonate Build:     mock-resonate-build
Firmware Build:     mock-firmware-build
Reader Apps Build:  mock-reader-apps-build

Table Counts:
  RecordingSession:   1
  SiteInformation:    1
  RawReads:            272
  MLT_SOW_Locations:   272
  ResonateEvents:      272
  Snapshots:           1

RawReads Details:
  Total Records:       272
  Unique Readers:      1
  Unique Tags:         100
  First InjectionTime: 2026-07-02 07:10:42.138
  Last InjectionTime:  2026-07-02 07:11:09.238
  Total Duration:      27.1s
```

The summary duration is calculated from:

```text
Last recorded injection timestamp - First recorded injection timestamp
```

It represents the original recorded event period, not the complete wall-clock execution time of the replay command.

---

## Playback and Pacing

RawReads are loaded in chronological order using their recorded injection timestamps.

The first record is sent at the beginning of playback. Every following record is scheduled according to its offset from the first recorded injection timestamp.

### Pacing flow

```mermaid
flowchart TD
    Start([Playback starts])
    Load[Load RawReads ordered by injection time]
    Establish[Store playback start time and first injection timestamp]
    More{More records available?}

    Offset[Calculate recorded offset from first record]
    TargetTime[Target send time equals playback start plus recorded offset]
    Wait[Wait until target send time]
    Send[POST payload to reader-bundles]
    Result{Request successful?}

    Success[Increment successful count]
    Failure[Increment failed count and log error]
    Summary[Generate replay summary]
    End([Playback completed])

    Start --> Load
    Load --> Establish
    Establish --> More

    More -- Yes --> Offset
    Offset --> TargetTime
    TargetTime --> Wait
    Wait --> Send
    Send --> Result

    Result -- Yes --> Success
    Result -- No --> Failure

    Success --> More
    Failure --> More

    More -- No --> Summary
    Summary --> End
```

Scheduling records relative to the playback start time reduces cumulative drift caused by:

- HTTP request processing;
- JSON formatting;
- operating-system scheduling;
- logging;
- local server processing;
- timer precision.

The wall-clock replay duration may still be slightly longer than the recorded duration because it includes application and HTTP-processing overhead.

Example:

```text
Recorded duration: 27.100s
Replay duration:   28.657s
Difference:         1.557s
```

For 272 HTTP requests, this is approximately 5.7 milliseconds of average overhead per record.

A small difference is expected. Large or continuously increasing drift should be investigated.

### Replay summary

```text
Playback finished, successful: 272, failed: 0

Replay Summary

Total Records: 272
Successful:    272
Failed:        0
Duration:      28.657s
Status:        Completed
```

---

## Logging and Error Handling

RRE records important execution events through the logger component.

Typical logged events include:

- application startup;
- configuration loading failures;
- validation start;
- validation success;
- validation failure;
- playback start;
- HTTP request failures;
- playback completion;
- successful and failed record totals.

Expected error-handling behavior:

- Invalid configuration stops command execution.
- Missing recording files return a clear error.
- Invalid SiteGraph JSON prevents the mock server from loading that file.
- SiteGraph mismatches stop playback before RawReads are loaded.
- Individual replay request failures are logged.
- A failed HTTP request does not terminate the full replay schedule.
- The replay summary shows successful and failed record totals.

Generated logs are stored under:

```text
logs/
```

---

## QA Testing Workflow

### 1. Build the executable

```powershell
go build -o rre.exe .\cmd\rre
```

### 2. Check the version

```powershell
.\rre.exe version
```

### 3. Inspect the recording

```powershell
.\rre.exe summary --config .\configs\config.json
```

### 4. Start the mock target server

Open a separate terminal:

```powershell
.\rre.exe serve --config .\configs\config.json
```

Keep this terminal running.

### 5. View all mock sites

```powershell
Invoke-RestMethod http://localhost:8080/sites |
  ConvertTo-Json -Depth 10
```

### 6. View the complete SiteGraph

```powershell
Invoke-RestMethod `
  http://localhost:8080/sites/b3489888-aacf-4451-893c-d7d994240f93 |
  ConvertTo-Json -Depth 100
```

### 7. Validate the recording

```powershell
.\rre.exe validate --config .\configs\config.json
```

### 8. Start playback

Run playback only after validation succeeds:

```powershell
.\rre.exe play --config .\configs\config.json
```

### 9. Test validation failure

1. Stop the mock server using `Ctrl+C`.
2. Create a temporary backup of the SiteGraph file.
3. Change one required Floor ID, Region ID, Reader ID or Antenna Port.
4. Restart the mock server.
5. Run `validate` again.
6. Confirm validation fails with the expected mismatch.
7. Confirm playback does not begin.
8. Restore the original SiteGraph after testing.

Do not commit intentionally modified QA SiteGraph files.

### 10. Test an external config path

```powershell
.\rre.exe validate --config "C:\QA\RRE\config.json"
```

Confirm that a valid JSON configuration file works from a user-defined location.

---

## Build and Test

### Requirements

- Go installed.
- A supported Recorder SQLite database.
- A target Resonate environment or local mock target.
- PowerShell for the documented Windows examples.

### Download dependencies

```powershell
go mod download
```

### Format the code

```powershell
gofmt -w .
```

### Run all tests

```powershell
go test ./...
```

Verbose execution:

```powershell
go test ./... -v
```

### Run static analysis

```powershell
go vet ./...
```

### Build

```powershell
go build -o rre.exe .\cmd\rre
```

---

## Release Build

Inject the application version using `ldflags`.

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

Expected output:

```text
v0.2.1
```

Build releases only from the latest merged `main` branch.

Recommended flow:

```text
Feature or bug-fix branch
        ↓
Tests and manual verification
        ↓
Commit and push branch
        ↓
Pull request and review
        ↓
Merge into main
        ↓
Pull latest main
        ↓
Final tests and release build
        ↓
Create version tag
        ↓
Publish release package
```

### Git tag example

From the monorepo root:

```powershell
git tag -a resonate-replay-engine/v0.2.1 `
  -m "Resonate Replay Engine v0.2.1"
```

Push the tag:

```powershell
git push origin resonate-replay-engine/v0.2.1
```

---

## Security and Repository Hygiene

Do not commit customer data, generated files or runtime output without explicit approval.

Avoid committing:

- real customer SiteGraph JSON files;
- real Recorder SQLite databases;
- SQLite shared-memory or write-ahead-log files;
- application logs;
- received replay payloads;
- generated executables;
- temporary files;
- backup files;
- credentials, tokens or private target URLs.

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

Important: `.gitignore` does not stop tracking a file that was committed earlier.

To stop tracking an already committed generated file:

```powershell
git rm --cached logs/received_payloads.jsonl
```

Commit that removal through the normal pull-request process.

---

## Troubleshooting

### Config file is not found

```text
config file not found
```

Confirm:

- the path is correct;
- the file exists;
- quotation marks are used when the path contains spaces;
- the file has a `.json` extension.

Example:

```powershell
.\rre.exe validate --config "C:\QA Files\RRE\config.json"
```

### Invalid config JSON

```text
failed to parse config file
```

Confirm:

- property names use double quotation marks;
- commas are correctly placed;
- there are no trailing commas;
- all required fields are present.

### Target site cannot be retrieved

Confirm:

- the target URL is correct;
- the mock server or target service is running;
- the configured Site ID exists;
- `GET /sites/{siteId}` returns a SiteGraph.

### Validation fails

Review the detailed missing structure output.

Common causes include:

- incorrect Site ID;
- missing Floor;
- missing nested Region;
- Reader placed under the wrong Floor;
- missing Reader;
- missing Antenna Port.

### Log file prevents Git operations

A running server may keep `logs/received_payloads.jsonl` open on Windows.

Stop the server:

```text
Ctrl+C
```

Then discard the generated change:

```powershell
git restore -- logs/received_payloads.jsonl
```

### Replay duration is slightly longer than recorded duration

The recorded duration only represents the difference between the first and last injection timestamps.

The replay duration also includes:

- JSON processing;
- HTTP request and response time;
- logging;
- timer precision;
- operating-system scheduling;
- target processing time.

A small difference is expected. Investigate only when drift is unexpectedly large or increases significantly for longer recordings.

---

## Summary

RRE provides a repeatable workflow for:

1. inspecting Recorder SQLite data;
2. hosting local mock SiteGraphs;
3. validating recorded and target site compatibility;
4. replaying RFID reader data with recorded timing;
5. collecting clear validation, logging and replay results.

The Replay Engine intentionally separates recorded-site validation from payload replay so incompatible target configurations are detected before data injection begins.