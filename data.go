You are working inside this Go project:

C:\projects\ALTRFIDTools\resonate-replay-engine

Work only inside `resonate-replay-engine`.

Before editing, inspect the current implementation carefully. Do not remove or break existing working functionality.

Existing functionality that must remain working:

- help
- version
- generate-sample
- summary
- mock-server
- validate
- play
- real Recorder SQLite schema support
- config-file execution
- direct command-line argument execution
- SiteInformation.site_json validation
- RawReads.raw_payload replay
- InjectionTime-based pacing
- security fixes for Cycode findings
- README and CLI help
- version injection using Go ldflags

==================================================
CURRENT PROBLEM
==================================================

The real Recorder SQLite stores a complete and very large SiteGraph JSON in:

SiteInformation.site_json

The existing Go SiteConfig model is incomplete. Because Go ignores unknown JSON fields, the current mock server can print incorrect values such as:

Floors: 1
Regions: 0
Readers: 0
Antennas: 0

even though the recorded SiteGraph contains floors, regions, readers and antennas.

This can also cause a false validation pass because both recorded and target configurations may become incomplete objects with zero readers, regions or antennas.

Do not solve this by adding hundreds of unrelated SiteGraph fields to Go structs.

Implement a small structural validation model that contains only the fields required by the validation epic.

==================================================
CONFIRMED REAL SITEGRAPH HIERARCHY
==================================================

The confirmed real Recorder SiteGraph hierarchy is:

Site
  -> Floors
      -> Regions
          -> Readers
              -> Antennas

Confirmed details:

1. One site can have multiple floors.
2. Each floor has multiple regions.
3. Regions are directly inside a floor.
4. Regions are not recursively nested for this Recorder output.
5. Each region contains readers.
6. Each reader contains antennas.
7. Reader identity is the string field:

"id"

8. Antennas do not have an ID field.
9. Antennas are identified by the integer field:

"port"

10. One reader may contain antenna ports 1 through 8.
11. Unrelated fields are not needed for minimum validation, such as:
   - name
   - type
   - group
   - make
   - model
   - bounds
   - position
   - rotation
   - networking
   - timeouts
   - physicality
   - behaviors
   - inventoryType

==================================================
VALIDATION EPIC REQUIREMENT
==================================================

Before playback, RRE must:

1. Read the recorded complete SiteGraph from:

SiteInformation.site_json

2. Fetch the target SiteGraph from:

GET /sites/{siteId}

3. Validate structural compatibility.

At minimum validate:

- Site ID matches.
- Every recorded Floor ID exists in the target.
- Every recorded Region ID exists under the correct floor.
- Every recorded Reader ID exists under the correct region.
- Every recorded antenna port exists under the correct reader.

The target is allowed to contain additional floors, regions, readers or antenna ports.

The minimum rule is:

Everything required by the recorded SiteGraph must exist in the target SiteGraph.

If anything is missing, validation must fail immediately with a clear message, for example:

Configuration mismatch: target site is missing Floor ID <floor-id>

Configuration mismatch: target site is missing Region ID <region-id> under Floor ID <floor-id>

Configuration mismatch: target site is missing Reader ID <reader-id> under Region ID <region-id>

Configuration mismatch: target site is missing antenna port 4 under Reader ID <reader-id>

The `play` command must stop before reading or replaying RawReads when validation fails.

==================================================
SIMPLIFIED INTERNAL VALIDATION MODEL
==================================================

Create a focused internal model similar to:

type ValidationSite struct {
    SiteID string            `json:"site_id"`
    Floors []ValidationFloor `json:"floors"`
}

type ValidationFloor struct {
    ID      string             `json:"id"`
    Regions []ValidationRegion `json:"regions"`
}

type ValidationRegion struct {
    ID      string             `json:"id"`
    Readers []ValidationReader `json:"readers"`
}

type ValidationReader struct {
    ID           string `json:"id"`
    AntennaPorts []int  `json:"antenna_ports"`
}

Follow current project naming conventions.

Suggested new files:

internal/site/validation_model.go
internal/site/parser.go
internal/site/summary.go

Implement a parser that converts the complete SiteGraph JSON into this simplified model.

The parser must:

- read root `id` as SiteID
- read every floor in `floors`
- read each floor `id`
- read every region in floor `regions`
- read each region `id`
- read every reader in region `readers`
- read each reader `id`
- read every antenna `port` in reader `antennas`
- process all floors, regions, readers and antennas
- ignore unrelated unknown fields safely
- return clear errors for malformed required fields
- avoid false zero-value results

Preserve the complete relationships:

Floor -> Region -> Reader -> Antenna Port

Do not validate readers or antenna ports only as unrelated global counts.

Use structural keys where useful:

floorID

floorID + regionID

floorID + regionID + readerID

floorID + regionID + readerID + antennaPort

==================================================
ONE CONFIG FILE ONLY
==================================================

Use only one QA configuration file:

configs/config.json

Do not create:

- pass_config.json
- fail_config.json
- pass_site_config.json
- fail_site_config.json
- wrong_site_config.json

The same `config.json` must be used for both passing and failing tests.

It must contain:

- runtime inputs
- the visible editable target validation structure

Use a structure similar to:

{
  "recording_file": "data/recording_001.sqlite",
  "target_url": "http://localhost:8080",
  "site_id": "b3489888-aacf-4451-893c-d7d994240f93",
  "mock_port": 8080,
  "target_site": {
    "site_id": "",
    "floors": []
  }
}

Update `internal/config/model.go` so `RunConfig` includes the current runtime fields and:

TargetSite ValidationSite `json:"target_site"`

or a pointer if needed to detect whether it is missing.

The configuration class must remain central to the workflow.

==================================================
AUTOMATIC CORRECT CONFIG POPULATION
==================================================

When QA first runs:

.\rre.exe mock-server -config configs\config.json

the application must:

1. Load `configs/config.json`.
2. Read `recording_file`.
3. Open the Recorder SQLite database read-only.
4. Read `SiteInformation.site_json` for the configured site.
5. Parse the full SiteGraph into the simplified `ValidationSite`.
6. Check `target_site` inside config.json.

If `target_site` is missing or empty:

- populate `target_site` automatically using the correct values extracted from SQLite
- write those correct values back into the same `configs/config.json`
- preserve:
  - recording_file
  - target_url
  - site_id
  - mock_port
- format the JSON clearly for QA
- print that the config file was initialized
- start the mock server using the populated `target_site`

Expected message:

[INFO] Loaded config: configs/config.json
[INFO] Loaded recorded SiteGraph from SQLite
[INFO] Populated target_site with correct recorded values
[INFO] Updated configs/config.json
[INFO] QA can edit target_site and restart mock-server to test validation failures

Very important:

After `target_site` has been populated, normal mock-server startup must NOT overwrite QA’s edits.

Normal command:

.\rre.exe mock-server -config configs\config.json

must use the current visible `target_site` values exactly as QA edited them.

Add a refresh option:

.\rre.exe mock-server `
  -config configs\config.json `
  -refresh-site-config

When `-refresh-site-config` is used:

- re-read the real SiteInformation.site_json
- regenerate the correct simplified structure
- overwrite only `target_site`
- preserve the runtime fields
- rewrite config.json safely
- start mock-server using the restored correct values

This allows:

Correct values -> validation passes

QA-edited values -> validation fails

Refresh -> correct values are restored

No pass or fail behavior may be hardcoded.

==================================================
CONFIG.JSON EXAMPLE AFTER AUTO-POPULATION
==================================================

The generated config must look similar to this, but use actual values from SQLite:

{
  "recording_file": "data/recording_001.sqlite",
  "target_url": "http://localhost:8080",
  "site_id": "b3489888-aacf-4451-893c-d7d994240f93",
  "mock_port": 8080,
  "target_site": {
    "site_id": "b3489888-aacf-4451-893c-d7d994240f93",
    "floors": [
      {
        "id": "actual-floor-id",
        "regions": [
          {
            "id": "actual-region-id",
            "readers": [
              {
                "id": "actual-reader-id",
                "antenna_ports": [1, 2, 3, 4, 5, 6, 7, 8]
              }
            ]
          }
        ]
      }
    ]
  }
}

Only include structural values required for validation.

Do not include:

- names
- coordinates
- bounds
- models
- networking
- timeouts
- other unrelated fields

QA must be able to clearly see and edit:

- Site ID
- Floor IDs
- Region IDs
- Reader IDs
- Antenna ports

==================================================
QA PASS AND FAIL FLOW
==================================================

PASS FLOW:

1. QA runs:

.\rre.exe mock-server `
  -config configs\config.json `
  -refresh-site-config

2. Correct values from SQLite are written into config.json.
3. Mock-server starts using those correct values.
4. In another terminal:

.\rre.exe validate -config configs\config.json

Expected:

[OK] Validation passed: recorded site is structurally compatible with target site

FAIL FLOW:

1. QA stops mock-server.
2. QA opens configs/config.json.
3. QA changes one value inside `target_site`, for example:
   - change a floor ID
   - change a region ID
   - change a reader ID
   - remove antenna port 4
4. QA saves config.json.
5. QA starts mock-server normally without refresh:

.\rre.exe mock-server -config configs\config.json

6. Mock-server must preserve and use the edited values.
7. QA runs:

.\rre.exe validate -config configs\config.json

Expected:

Validation fails with a clear structural mismatch.

To restore the correct pass values:

.\rre.exe mock-server `
  -config configs\config.json `
  -refresh-site-config

==================================================
MOCK SERVER BEHAVIOR
==================================================

The mock server must support the simplified `target_site` from config.json.

When handling:

GET /sites/{siteId}

return JSON that the target client and validator can normalize into the same internal ValidationSite model.

The validator must support both:

1. A complete real SiteGraph returned by a real Resonator endpoint.
2. The simplified QA `target_site` returned by mock-server.

Use one normalized internal `ValidationSite` structure for comparison.

Preserve existing full SQLite SiteGraph mock mode only if it is still required by existing functionality, but the main QA workflow must use `config.json.target_site`.

Do not lose or break:

POST /reader-bundles

received payload logging

existing mock-server summary

==================================================
MOCK SERVER TERMINAL DISPLAY
==================================================

When mock-server starts, print:

RRE Mock Target Server

Config file: configs/config.json
Recording file: data/recording_001.sqlite
Target source: config.json target_site

Site ID: <id>
Floors: <count>
Regions: <count>
Readers: <count>
Antenna Ports: <count>

Listening on: http://localhost:8080

Endpoints:
GET /sites/{siteId}
POST /reader-bundles

Also print a readable hierarchy:

Floor: <floor-id>
  Region: <region-id>
    Reader: <reader-id>
      Antenna Ports: 1, 2, 3, 4, 5, 6, 7, 8

Counts must include all structures.

Do not use incorrect top-level fields such as:

len(site.Regions)

when regions are inside floors.

==================================================
VALIDATE COMMAND DISPLAY
==================================================

When validate runs, display:

Recorded site:
  Site ID: ...
  Floors: ...
  Regions: ...
  Readers: ...
  Antenna Ports: ...

Target site:
  Site ID: ...
  Floors: ...
  Regions: ...
  Readers: ...
  Antenna Ports: ...

Then print either:

[OK] Validation passed: recorded site is structurally compatible with target site

or detailed errors.

Do not allow a false pass because both incomplete structures contain zero readers or regions.

==================================================
CONFIG MODE AND DIRECT ARGUMENT MODE
==================================================

Both execution methods must remain supported.

Config mode is the main QA workflow.

Validate:

.\rre.exe validate -config configs\config.json

Direct validate:

.\rre.exe validate `
  -file data\recording_001.sqlite `
  -target-url http://localhost:8080 `
  -site-id b3489888-aacf-4451-893c-d7d994240f93

Play:

.\rre.exe play -config configs\config.json

Direct play:

.\rre.exe play `
  -file data\recording_001.sqlite `
  -target-url http://localhost:8080 `
  -site-id b3489888-aacf-4451-893c-d7d994240f93

Summary:

.\rre.exe summary -config configs\config.json

Direct summary:

.\rre.exe summary -file data\recording_001.sqlite

Mock server:

.\rre.exe mock-server -config configs\config.json

Mock server refresh:

.\rre.exe mock-server `
  -config configs\config.json `
  -refresh-site-config

Keep direct mock-server flags working if currently supported.

Do not create separate implementations for config mode and direct mode.

Resolve both input styles into one shared internal configuration and call the same command logic.

Inspect and preserve:

internal/config/model.go
internal/cli/args.go
internal/cli/commands.go

==================================================
PLAY COMMAND ORDER
==================================================

The `play` command must continue in this exact order:

1. Resolve config or direct arguments.
2. Open the recording SQLite.
3. Read SiteInformation.site_json.
4. Fetch the target site data.
5. Normalize recorded and target structures.
6. Perform structural validation.
7. Abort immediately if validation fails.
8. Only when validation passes:
   - read RawReads
   - order by injection_time_utc
   - build replay payloads
   - use raw_payload or structured fallback
   - send to POST /reader-bundles

Do not break RawReads replay.

Remember:

SiteInformation.site_json is used for validation.

RawReads.raw_payload is used for playback.

==================================================
REAL SQLITE SCHEMA MUST REMAIN
==================================================

Do not revert the real Recorder SQLite schema support.

SiteInformation columns:

- site_information_id
- recording_session_id
- site_id
- site_name
- site_json

RawReads columns:

- read_id
- recording_session_id
- tag_id
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

RecordingSession uses real snake_case fields, including:

- recording_session_id
- test_name
- environment
- tester_name
- description
- start_time_utc
- end_time_utc
- resonate_build_number
- firmware_build_number
- reader_apps_build_number
- resonate_site_id
- state

Do not use:

- test_description
- RawSiteJSON
- old PascalCase SQL column names

Preserve nullable field handling such as nullable RSSI.

==================================================
SECURITY REQUIREMENTS
==================================================

The repository uses Cycode scanning.

Do not reintroduce previous violations.

Do not use:

- dynamic SQL
- fmt.Sprintf for SQL identifiers
- unrestricted arbitrary file reads
- unsafe path traversal
- unrestricted arbitrary HTTP URLs
- removal of SSRF protection
- ignored security findings

Use:

- parameterized hardcoded SQL
- safe read-only SQLite access
- config path validation
- filepath.Clean
- allowed config directory checks
- safe URL validation already used by the project

When rewriting config.json:

- only allow the expected file inside the approved configs directory
- clean and validate the path
- serialize to formatted JSON
- write to a temporary file in the same safe directory
- replace the original only after successful write
- preserve runtime fields
- return clear errors
- do not permit `..` traversal

==================================================
FILES TO INSPECT
==================================================

Inspect and update the relevant files, likely including:

cmd/rre/main.go
internal/cli/args.go
internal/cli/commands.go
internal/config/model.go
internal/recording/model.go
internal/recording/repository.go
internal/site/model.go
internal/site/client.go
internal/site/validator.go
internal/site/validator_test.go
internal/mocktarget/handler.go
internal/mocktarget/server.go
README.md
configs/config.json
.gitignore

Create focused new files if useful:

internal/site/validation_model.go
internal/site/parser.go
internal/site/parser_test.go
internal/site/summary.go

Use the exact Go module path from go.mod.

Do not modify:

../resonate-recorder
../resonate-analyzer

==================================================
TESTS
==================================================

Add or update tests for:

1. One site with multiple floors.
2. Multiple regions inside a floor.
3. Multiple readers inside a region.
4. One reader with antenna ports 1 through 8.
5. Matching recorded and target structures pass.
6. Missing floor fails.
7. Missing region under the correct floor fails.
8. Reader under the wrong region fails.
9. Missing reader fails.
10. Missing antenna port fails.
11. Extra target structures do not fail.
12. Invalid full SiteGraph JSON returns a clear error.
13. Empty target_site is auto-populated.
14. Populated config.json uses real values from SQLite.
15. Normal mock-server start does not overwrite edited target_site.
16. -refresh-site-config restores correct target_site values.
17. Config rewrite preserves runtime fields.
18. Mock-server prints correct counts.
19. Validate prints correct recorded and target counts.
20. Validate works with config mode.
21. Validate works with direct arguments.
22. Play works with config mode.
23. Play works with direct arguments.
24. Summary works with config mode.
25. Summary works with direct arguments.
26. Full real SiteGraph and simplified target_site normalize into the same internal model.
27. Unknown SiteGraph fields do not break parsing.
28. No false pass occurs when readers, regions or antenna ports are missing.

==================================================
README
==================================================

Update README.md with:

- Difference between SiteInformation.site_json and RawReads.raw_payload.
- Real SiteGraph hierarchy:

Site -> Floors -> Regions -> Readers -> Antenna Ports

- Explanation of the single configs/config.json file.
- How config.json is automatically initialized.
- How to refresh correct pass values.
- How QA edits target_site to test failures.
- Pass workflow.
- Fail workflow.
- Commands for mock-server, validate, play and summary.
- Direct argument examples.
- Expected pass and fail output.
- Warning not to commit real SQLite or customer data.

==================================================
IMPLEMENTATION PROCESS
==================================================

First inspect the code and provide a brief plan containing:

1. Current implementation.
2. Incorrect current assumptions.
3. Files to change.
4. New files to add.
5. How config.json auto-population works.
6. How QA edits are preserved.
7. How refresh restores correct values.
8. How config and direct argument modes remain supported.

Then implement the changes.

Do not stop after only explaining.

Make the actual code changes.

After implementation run:

gofmt on all changed Go files

go test ./...

go vet ./...

go build -o rre.exe ./cmd/rre

Do not commit or stage rre.exe.

Finally report:

1. Files changed.
2. Design decisions.
3. Exact PowerShell commands to test.
4. Expected pass output.
5. Expected fail output.
6. Any remaining assumptions requiring architect confirmation.