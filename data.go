# Resonate Replay Engine v0.1.0

This is the first QA release of the Resonate Replay Engine CLI.

The Replay Engine supports replaying recorded RFID/location raw read data from a Recorder SQLite database into a target Resonate-compatible HTTP endpoint. This release includes real Recorder SQLite schema support, site configuration validation, replay execution, mock target testing support, and CLI versioning.

## Added

### CLI Commands
- Added Replay Engine CLI command support for:
  - `help`
  - `generate-sample`
  - `summary`
  - `mock-server`
  - `validate`
  - `play`
  - `version`

### Real Recorder SQLite Support
- Added support for the real Recorder SQLite schema.
- Reads Recorder tables using snake_case columns, including:
  - `RecordingSession`
  - `SiteInformation`
  - `RawReads`
- Reads recorded site configuration from `SiteInformation.site_json`.
- Reads raw RFID/location data from `RawReads`.

### Site Configuration Validation
- Added validation before replay starts.
- Compares recorded site configuration against target site configuration.
- Validates important site config details such as:
  - Site ID
  - Readers
  - Antennas
  - Floors
  - Regions
- Aborts replay when site configuration does not match.

### Replay Flow
- Added RawReads replay support.
- Replays reads ordered by injection time.
- Supports timing-based replay pacing using recorded injection timestamps.
- Sends replay payloads to the target endpoint.

### Mock Target Server
- Added mock target server for local QA testing.
- Supports:
  - `GET /sites/{siteId}`
  - `POST /reader-bundles`
- Stores received payloads in JSONL format for verification.
- Shows mock server summary after shutdown.

### Config-Based QA Testing
- Added config-based execution support.
- Supports running commands using config files such as:
  - `configs/pass_config.json`
  - `configs/fail_config.json`
  - `configs/wrong_site_config.json`
- Added terminal output for recorded and target site configuration summaries.

### Versioning
- Added CLI version command:
  - `rre version`
  - `rre --version`
  - `rre -v`
- Default local build version is `dev`.
- Added build-time version injection using Go `ldflags`.

## Build

For local development build:

```powershell
go build -o rre.exe ./cmd/rre