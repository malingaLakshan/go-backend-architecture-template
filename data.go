We need to restore and fix the RRE CLI config-based flow properly after the folder rename from replay-engine-mvp-cli to resonate-replay-engine.

Context:
This project is the Resonance Replay Engine CLI. The executable is rre.exe. The folder is now resonate-replay-engine, not replay-engine-mvp-cli. The CLI should support both direct flags and JSON config files for summary, validate, play, and mock-server.

Please review and update all related files so everything works consistently.

Main requirements:

1. Fix CLI args.go

Add Config string to the Flags struct.

Example:

type Flags struct {
	Command   string
	File      string
	Out       string
	TargetURL string
	SiteID    string
	Port      int
	Config    string
}

Add -config support for these commands:
- summary
- mock-server
- validate
- play

Keep existing direct flags working:
- summary -file data/recording_001.sqlite
- validate -file data/recording_001.sqlite -target-url http://localhost:8080 -site-id b3489888-aacf-4451-893c-d7d994240f93
- play -file data/recording_001.sqlite -target-url http://localhost:8080 -site-id b3489888-aacf-4451-893c-d7d994240f93
- mock-server -port 8080

Config examples should work:
- ./rre.exe summary -config configs/pass_config.json
- ./rre.exe validate -config configs/pass_config.json
- ./rre.exe play -config configs/pass_config.json
- ./rre.exe mock-server -config configs/pass_config.json
- ./rre.exe mock-server -config configs/fail_config.json

The -config flag must appear in:
- ./rre.exe summary -help
- ./rre.exe validate -help
- ./rre.exe play -help
- ./rre.exe mock-server -help

2. Fix config model/loading

Use config JSON fields:

{
  "recording_file": "data/recording_001.sqlite",
  "target_url": "http://localhost:8080",
  "site_id": "b3489888-aacf-4451-893c-d7d994240f93",
  "mock_mode": "sqlite",
  "mock_site_file": ""
}

The config model should support:
- RecordingFile string `json:"recording_file"`
- TargetURL string `json:"target_url"`
- SiteID string `json:"site_id"`
- MockMode string `json:"mock_mode"`
- MockSiteFile string `json:"mock_site_file"`

3. Fix summary command config support

summary must work with direct file and config file.

Required commands:

Direct:
./rre.exe summary -file data/recording_001.sqlite

Config:
./rre.exe summary -config configs/pass_config.json
./rre.exe summary -config configs/fail_config.json

If -config is provided:
- Load config JSON.
- Use recording_file as the SQLite file path.
- Do not require target_url, site_id, mock_mode, or mock_site_file for summary.

If both -file and -config are missing, show clear error:
[ERROR] -file or -config is required for summary

Expected summary output:
- Recording session summary
- RawReads count
- unique reader count
- unique tag count
- first injection time
- last injection time
- total duration

Important:
- summary must read the real Recorder SQLite schema.
- Use snake_case columns from the real Recorder DB.

4. Fix validate command config support

validate must work with direct flags and config file.

Direct:
./rre.exe validate -file data/recording_001.sqlite -target-url http://localhost:8080 -site-id b3489888-aacf-4451-893c-d7d994240f93

Config:
./rre.exe validate -config configs/pass_config.json
./rre.exe validate -config configs/fail_config.json

If -config is provided:
- Load config JSON.
- Set flags.File from cfg.RecordingFile.
- Set flags.TargetURL from cfg.TargetURL.
- Set flags.SiteID from cfg.SiteID.

If required values are missing, show clear error:
[ERROR] -file, -target-url, and -site-id are required for validate

Validate should:
- Load recorded site config from SQLite SiteInformation.site_json.
- Fetch target site config from GET /sites/{siteId}.
- Validate recorded site config against target site config.
- Print recorded site config summary.
- Print target site config summary.
- Return non-zero exit code when validation fails.

Update old error messages from RawSiteJSON to SiteJSON.

5. Fix play command config support

play must work with direct flags and config file.

Direct:
./rre.exe play -file data/recording_001.sqlite -target-url http://localhost:8080 -site-id b3489888-aacf-4451-893c-d7d994240f93

Config:
./rre.exe play -config configs/pass_config.json

If -config is provided:
- Load config JSON.
- Set flags.File from cfg.RecordingFile.
- Set flags.TargetURL from cfg.TargetURL.
- Set flags.SiteID from cfg.SiteID.

If required values are missing, show clear error:
[ERROR] -file, -target-url, and -site-id are required for play

play should:
- Validate recorded site config against target site config first.
- Abort replay if validation fails.
- Load RawReads from real Recorder SQLite schema.
- Replay RawReads to POST /reader-bundles.
- Use InjectionTime-based pacing.
- Write sender-side replay log to logs/replay_output.jsonl.

6. Fix mock-server config behavior

mock-server must work with direct port and config file.

Direct:
./rre.exe mock-server -port 8080

Config:
./rre.exe mock-server -config configs/pass_config.json
./rre.exe mock-server -config configs/fail_config.json

If no -config is provided:
- Start default mock server on selected port.
- Default port should remain 8080.

If -config is provided:
- Load config JSON.
- Use cfg.MockMode.
- Use cfg.RecordingFile when mock_mode is sqlite.
- Use cfg.MockSiteFile when mock_mode is file.
- Use cfg.SiteID as the site ID.

mock_mode = sqlite:
- Load site config from the recording SQLite SiteInformation.site_json by site_id.
- Serve that exact config from GET /sites/{siteId}.

mock_mode = file:
- Load mock_site_file JSON.
- Serve that file content from GET /sites/{siteId}.

If mock_mode is missing or invalid, show clear error:
[ERROR] config mock_mode must be "sqlite" or "file"

If mock_mode is sqlite and recording_file is missing, show:
[ERROR] config missing required field for sqlite mode: recording_file

If mock_mode is file and mock_site_file is missing, show:
[ERROR] config missing required field for file mode: mock_site_file

If mock_site_file JSON is invalid, show clear error.

mock-server terminal output should clearly show:
- loaded config path
- mock mode
- recording file or mock site file
- site config source
- site ID
- site name
- reader count
- antenna count
- floor count
- region count
- listening URL
- endpoints
- received payload log path

7. Fix mocktarget server/handler

Ensure server.go supports:
- StartServer(port int) for default server
- StartServerWithConfig(port int, cfg *site.SiteConfig) for config-based server

Ensure handler.go supports:
- NewHandler()
- NewHandlerWithConfig(cfg *site.SiteConfig)
- GET /sites/{siteId}
- POST /reader-bundles

GET /sites/{siteId}:
- If config-based server has siteConfig, return only when requested site ID matches.
- If requested site ID does not match, return 404.
- If no config-based siteConfig is given, use default fallback mock config.

POST /reader-bundles:
- Read request body.
- Save each received payload to logs/received_payloads.jsonl.
- Print first payload in pretty JSON.
- Print received count, site_id, reader_id, read count, and payload size.
- Return JSON response:
  {"status":"accepted"}

8. Fix config JSON files

Create or restore these files under configs:
- configs/pass_config.json
- configs/fail_config.json
- configs/wrong_site_config.json

configs/pass_config.json:

{
  "recording_file": "data/recording_001.sqlite",
  "target_url": "http://localhost:8080",
  "site_id": "b3489888-aacf-4451-893c-d7d994240f93",
  "mock_mode": "sqlite",
  "mock_site_file": ""
}

configs/fail_config.json:

{
  "recording_file": "data/recording_001.sqlite",
  "target_url": "http://localhost:8080",
  "site_id": "b3489888-aacf-4451-893c-d7d994240f93",
  "mock_mode": "file",
  "mock_site_file": "configs/wrong_site_config.json"
}

configs/wrong_site_config.json must be valid SiteConfig JSON, not invalid strings.

It should contain full structure with wrong IDs so validation fails because of ID mismatch, not JSON parse error.

Use this JSON:

{
  "id": "WRONG-SITE-ID",
  "name": "Wrong Site",
  "readers": [
    {
      "id": "WRONG-READER-01",
      "name": "Wrong Reader",
      "type": "RFID",
      "ipAddress": "192.168.1.200",
      "floorId": "WRONG-FLOOR-01",
      "x": 1,
      "y": 1,
      "antennas": [
        {
          "antenna_id": 99,
          "antenna_type": 2,
          "reader_id": "WRONG-READER-01",
          "x": 1,
          "y": 1
        }
      ]
    }
  ],
  "floors": [
    {
      "id": "WRONG-FLOOR-01",
      "name": "Wrong Floor",
      "number": 1,
      "width": 100,
      "height": 100,
      "regions": [
        {
          "id": "WRONG-REGION-01",
          "name": "Wrong Region",
          "type": "WRONG_TYPE",
          "physicality": "VIRTUAL",
          "inventoryType": "OTHER"
        }
      ]
    }
  ],
  "regions": [
    {
      "id": "WRONG-REGION-01",
      "name": "Wrong Region",
      "type": "WRONG_TYPE",
      "physicality": "VIRTUAL",
      "inventoryType": "OTHER"
    }
  ],
  "antennas": [
    {
      "antenna_id": 99,
      "antenna_type": 2,
      "reader_id": "WRONG-READER-01",
      "x": 1,
      "y": 1
    }
  ]
}

9. Fix validator.go

Validation must not only compare counts.

It must validate IDs and return clear errors.

Validate:
- site ID equality
- reader ID equality
- antenna ID equality
- floor ID equality
- region ID equality

If counts match but IDs are different, validation must fail.

Rules:
- If recorded site ID != target site ID, fail.
- Every recorded reader ID must exist in target reader IDs.
- Every target reader ID must exist in recorded reader IDs.
- Every recorded antenna ID must exist in target antenna IDs.
- Every target antenna ID must exist in recorded antenna IDs.
- Every recorded floor ID must exist in target floor IDs.
- Every target floor ID must exist in recorded floor IDs.
- Every recorded region ID must exist in target region IDs.
- Every target region ID must exist in recorded region IDs.

Error examples:
- Site ID mismatch: recorded=..., target=...
- Reader ID missing in target: ...
- Reader ID missing in recorded: ...
- Antenna ID missing in target: ...
- Antenna ID missing in recorded: ...
- Floor ID missing in target: ...
- Floor ID missing in recorded: ...
- Region ID missing in target: ...
- Region ID missing in recorded: ...

The validator result should include:
- Passed bool
- Errors []string

10. Fix SiteConfig model

Ensure site model supports real Recorder site_json structure enough for validation.

SiteConfig should include:
- id
- name
- readers
- floors
- regions
- antennas

Reader should include:
- id
- name
- type
- ipAddress
- floorId
- x
- y
- antennas

Antenna should include:
- antenna_id
- antenna_type
- reader_id
- x
- y

Floor should include:
- id
- name
- number
- width
- height
- regions

Region should include:
- id
- name
- type
- physicality
- inventoryType

Go must ignore extra JSON fields safely.

11. Fix real Recorder SQLite schema support

Repository code must use real Recorder snake_case columns.

SiteInformation table:
- site_information_id
- recording_session_id
- site_id
- site_name
- site_json

RawReads table:
- read_id
- recording_session_id
- tag_id
- z
- reader_id
- antenna_id
- antenna_type
- source_timestamp_utc
- injection_time_utc
- confidence
- rssi
- tag_x
- tag_y
- floor_id
- raw_payload

RecordingSession table should use real snake_case fields where applicable.

Fix any mixed old camelCase WHERE clauses, for example:
- RecordingSessionID should become recording_session_id
- ReaderID should become reader_id
- TagID should become tag_id

12. Fix payload.go

BuildPayload should support real Recorder DB path.

If raw_payload is present:
- Parse it if possible.

If raw_payload is empty or binary/unusable:
- Build payload from structured RawRead columns.
- Use RawRead Timestamp / InjectionTime safely.
- Do not create zero timestamps silently.

Add guard:
- If raw_payload is empty AND timestamp fields are empty/zero, return clear error:
  empty RawPayload and no structured timestamp fields for ReadID ...

13. Fix tests

Update tests to match real Recorder schema and config behavior.

Add or update tests for:
- summary with direct file
- summary with config
- validate with pass config
- validate with fail config
- validator fails when counts match but IDs differ
- mock-server config loading sqlite mode
- mock-server config loading file mode
- BuildPayload empty RawPayload guard

14. Fix README/help text

Update README and PrintHelp so config examples are shown.

Include direct examples:

./rre.exe summary -file data/recording_001.sqlite
./rre.exe mock-server -port 8080
./rre.exe validate -file data/recording_001.sqlite -target-url http://localhost:8080 -site-id b3489888-aacf-4451-893c-d7d994240f93
./rre.exe play -file data/recording_001.sqlite -target-url http://localhost:8080 -site-id b3489888-aacf-4451-893c-d7d994240f93

Include config examples:

./rre.exe summary -config configs/pass_config.json
./rre.exe mock-server -config configs/pass_config.json
./rre.exe validate -config configs/pass_config.json
./rre.exe play -config configs/pass_config.json

Fail validation demo:

Terminal 1:
./rre.exe mock-server -config configs/fail_config.json

Terminal 2:
./rre.exe validate -config configs/fail_config.json

Expected:
Validation should fail with mismatch errors.

15. Fix .gitignore

Keep generated runtime files out of git.

Do not commit:
- logs/received_payloads.jsonl
- logs/replay_output.jsonl
- *.sqlite-shm
- *.sqlite-wal
- temporary duplicate files like *.884371..., *.893632...
- rre.exe

Update .gitignore if needed.

16. Folder rename/import cleanup

The old folder name replay-engine-mvp-cli should not be referenced in code imports or README paths.

Use current module/import path from go.mod.

If go.mod module path is:
resonate-replay-engine

Then imports should use:
resonate-replay-engine/internal/...

If go.mod module path is different, use the exact module path from go.mod consistently.

Do not create temporary duplicate files like:
- handler.go.893632...
- repository.go.884371...
- model.go.421414...

Remove accidental temp duplicate files if present.

17. Build/test commands

After changes:
- Run gofmt on changed Go files.
- Run go test ./...
- Run go build -o rre.exe ./cmd/rre

18. Manual expected test flow

Test 1: Help

./rre.exe summary -help
./rre.exe validate -help
./rre.exe play -help
./rre.exe mock-server -help

Expected:
-config appears in all four commands.

Test 2: Summary with config

./rre.exe summary -config configs/pass_config.json

Expected:
Summary prints recording details from data/recording_001.sqlite.

Test 3: Pass validation

Terminal 1:
./rre.exe mock-server -config configs/pass_config.json

Terminal 2:
./rre.exe validate -config configs/pass_config.json

Expected:
Validation passed.

Test 4: Fail validation

Stop mock server.

Terminal 1:
./rre.exe mock-server -config configs/fail_config.json

Terminal 2:
./rre.exe validate -config configs/fail_config.json

Expected:
Validation failed with clear mismatch errors.

Test 5: Play

Stop mock server.

Terminal 1:
./rre.exe mock-server -config configs/pass_config.json

Terminal 2:
./rre.exe play -config configs/pass_config.json

Expected:
Replay sends payloads and mock server receives them.

Please implement all related changes carefully and consistently.