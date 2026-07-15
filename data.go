You are working inside this Go project:

C:\projects\ALTRFIDTools\resonate-replay-engine

Work only inside `resonate-replay-engine`.

Do not modify sibling projects such as:

../resonate-recorder
../resonate-analyzer

Before changing anything, inspect the current implementation carefully.

Do not remove or break existing working functionality.

The Replay Engine currently supports:

- help
- version
- generate-sample
- summary
- mock-server
- validate
- play
- config-file execution
- direct command-line argument execution
- real Recorder SQLite schema
- SiteInformation.site_json validation
- RawReads.raw_payload playback
- injection_time_utc replay pacing
- build-time version injection
- JSONL replay/mock output
- Cycode security fixes

Do not commit, push, create tags, or create a PR.

==================================================
DOMAIN AND EPIC REQUIREMENT
==================================================

A SiteGraph is the complete configuration of one real site in Resonator.

We currently have approximately four real SiteGraph JSON files.

Each SiteGraph represents a different target site and contains a unique root field:

"id"

The real target API behaves like:

GET /sites/{siteId}

and returns the complete SiteGraph for that requested Site ID.

The Replay Engine validation epic requires these inputs:

- source Recorder SQLite file
- target URL
- Site ID

Therefore, `site_id` must remain inside config.json and must remain available as a direct CLI argument.

The Site ID is required to select the correct target SiteGraph.

The Recorder SQLite contains:

SiteInformation.site_json
→ the complete site configuration stored when the recording was created

RawReads.raw_payload
→ the actual recorded data used during playback

Validation must compare:

Recorded SiteInformation.site_json

against:

Target SiteGraph returned by GET /sites/{siteId}

Playback must begin only after validation passes.

==================================================
CURRENT PROBLEM
==================================================

The existing implementation appears to unmarshal the large real SiteGraph into incomplete Go models.

Because Go ignores unknown JSON fields, the current mock-server terminal output can incorrectly show:

Floors: 1
Regions: 0
Readers: 0
Antennas: 0

even though the real SiteGraph contains floors, regions, readers and antenna ports.

This can cause false validation success because both recorded and target configurations may become incomplete structures containing zero regions, readers or antennas.

Do not fix this by adding hundreds of unrelated SiteGraph fields to Go structs.

Instead:

1. Preserve the complete raw SiteGraph JSON.
2. Extract only the structural values required for validation.
3. Compare the normalized structural models.

==================================================
CONFIRMED REAL SITEGRAPH HIERARCHY
==================================================

The confirmed hierarchy is:

Site
  -> Floors
      -> Regions
          -> Readers
              -> Antennas

Confirmed rules:

1. One site can contain multiple floors.
2. Each floor contains multiple regions.
3. Regions are directly inside a floor for the current Recorder output.
4. Regions are not recursively nested for this Recorder output.
5. Each region contains readers.
6. Each reader contains antennas.
7. Reader identity is the string field:

"id"

8. Antennas do not have a separate ID.
9. Antennas are identified using the integer field:

"port"

10. One reader may contain ports 1 through 8.
11. Other SiteGraph fields are not required for minimum structural validation, including:
   - name
   - type
   - group
   - make
   - model
   - category
   - coordinates
   - position
   - rotation
   - bounds
   - networking
   - timeouts
   - physicality
   - behaviors
   - inventoryType

Do not validate those unrelated fields unless existing project requirements explicitly require them.

==================================================
PROJECT CONFIGURATION ARCHITECTURE
==================================================

Use one central runtime configuration file:

configs/config.json

Store each complete real SiteGraph as a separate JSON file inside:

configs/sites/

Recommended structure:

resonate-replay-engine/
├── configs/
│   ├── config.json
│   └── sites/
│       ├── bentonville.json
│       ├── site-b.json
│       ├── site-c.json
│       └── site-d.json
├── data/
├── logs/
├── cmd/
├── internal/
└── README.md

Do not combine all large SiteGraphs into one giant JSON file.

Do not manually list every SiteGraph file in config.json.

Instead, config.json should contain an approved SiteGraph directory.

Recommended config.json:

{
  "recording_file": "data/recording_001.sqlite",
  "target_url": "http://localhost:8080",
  "site_id": "b3489888-aacf-4451-893c-d7d994240f93",
  "mock_port": 8080,
  "site_graph_directory": "configs/sites"
}

The meaning of each field:

recording_file
→ Recorder SQLite file used by summary, validate and play

target_url
→ target/mock server base URL

site_id
→ Site ID required by the epic; validate/play request this site from the target

mock_port
→ port used by mock-server

site_graph_directory
→ approved directory containing all mocked real SiteGraph JSON files

The root `id` inside each SiteGraph file is the source of truth for that mocked site.

The file name is not the Site ID.

==================================================
SITEGRAPH AUTO-DISCOVERY
==================================================

When QA runs:

.\rre.exe mock-server -config configs\config.json

the mock server must:

1. Load config.json.
2. Read site_graph_directory.
3. Safely find every `.json` file directly inside that approved directory.
4. Read each complete SiteGraph JSON.
5. Parse its root `id`.
6. Parse its site name if available for display only.
7. Reject:
   - invalid JSON
   - empty root Site ID
   - duplicate root Site IDs
   - unsafe paths
   - files outside the approved directory
   - non-JSON files
8. Store each complete raw SiteGraph by Site ID.

Conceptually:

map[string][]byte

where:

key
→ root SiteGraph ID

value
→ complete original SiteGraph JSON bytes

Do not unmarshal the full SiteGraph into an incomplete model and marshal it again before returning it.

Preserve the original JSON bytes.

Adding another SiteGraph later should only require:

1. Copy a new JSON file into configs/sites/
2. Restart mock-server

No config.json update should be needed.

==================================================
MOCK TARGET ENDPOINTS
==================================================

The mock server must support:

GET /sites

GET /sites/{siteId}

POST /reader-bundles

GET /sites

must return summaries of every loaded mocked site.

Recommended JSON response:

{
  "sites": [
    {
      "id": "site-id",
      "name": "Bentonville",
      "floors": 2,
      "regions": 8,
      "readers": 6,
      "antenna_ports": 48
    }
  ],
  "total": 4
}

GET /sites/{siteId}

must:

- search the loaded SiteGraph map using the requested Site ID
- return the exact complete raw SiteGraph JSON for that site
- return HTTP 404 when no site matches

Example error:

{
  "error": "site not found",
  "site_id": "unknown-id"
}

POST /reader-bundles

must preserve its current working replay/mock behavior.

Do not break payload receiving or output files.

Do not add a separate `rre sites` command in this implementation.

`GET /sites` is the core site-listing API.

QA can view available sites using:

Invoke-RestMethod http://localhost:8080/sites |
  ConvertTo-Json -Depth 10

==================================================
RECORDED SITE SELECTION
==================================================

The epic requires Site ID as an input.

Validation must use the provided `site_id` from config.json or the direct `-site-id` argument.

The validation flow must be:

1. Open the configured Recorder SQLite.
2. Read SiteInformation using the provided Site ID.
3. Read:
   - SiteInformation.site_id
   - SiteInformation.site_json
4. Parse the root SiteGraph ID from site_json.
5. Verify consistency between:
   - configured Site ID
   - SiteInformation.site_id
   - recorded site_json root ID
6. Request:

GET {target_url}/sites/{siteId}

7. Mock server finds the matching loaded SiteGraph.
8. RRE validates the recorded SiteGraph against that target SiteGraph.

QA should not manually select a SiteGraph filename.

The Site ID automatically selects the correct mocked SiteGraph.

If no matching target SiteGraph exists, validation must fail clearly:

Target site configuration was not found for Site ID <id>

==================================================
NORMALIZED VALIDATION MODEL
==================================================

Create a small focused internal validation model similar to:

type ValidationSite struct {
    SiteID string
    Floors []ValidationFloor
}

type ValidationFloor struct {
    ID      string
    Regions []ValidationRegion
}

type ValidationRegion struct {
    ID      string
    Readers []ValidationReader
}

type ValidationReader struct {
    ID           string
    AntennaPorts []int
}

Follow existing project naming conventions.

Suggested new files:

internal/site/validation_model.go
internal/site/parser.go
internal/site/summary.go

The parser must convert a complete real SiteGraph JSON into this normalized structure.

The parser must:

- read root `id`
- process every floor in `floors`
- read every floor `id`
- process every region inside each floor
- read every region `id`
- process every reader inside each region
- read every reader `id`
- process every antenna inside each reader
- read every antenna integer `port`
- process all floors, regions, readers and ports
- safely ignore unrelated unknown fields
- return descriptive errors when required structural data is malformed
- never silently return false empty structures when the JSON contains data

Preserve the relationships:

Floor
→ Region
→ Reader
→ Antenna Port

Do not compare only global counts.

Use structural keys where helpful:

floorID

floorID + regionID

floorID + regionID + readerID

floorID + regionID + readerID + antennaPort

==================================================
VALIDATION RULES
==================================================

Before playback, validate:

1. Configured Site ID matches the recorded SiteInformation Site ID.
2. Configured Site ID matches the root ID in recorded site_json.
3. Target SiteGraph root ID matches the requested Site ID.
4. Every recorded Floor ID exists in the target.
5. Every recorded Region ID exists under the correct Floor ID.
6. Every recorded Reader ID exists under the correct Region ID.
7. Every recorded antenna port exists under the correct Reader ID.

The target is allowed to contain additional:

- floors
- regions
- readers
- antenna ports

The required rule is:

Every structure required by the recording must exist in the target.

Clear failure examples:

Configuration mismatch: target site is missing Floor ID <floor-id>

Configuration mismatch: target site is missing Region ID <region-id> under Floor ID <floor-id>

Configuration mismatch: target site is missing Reader ID <reader-id> under Region ID <region-id>

Configuration mismatch: target site is missing antenna port 4 under Reader ID <reader-id>

Do not start playback when validation fails.

==================================================
STRUCTURED VALIDATION RESULTS
==================================================

Do not store validation results only as plain error strings.

Create structured results that can report category totals and exact mismatches.

Conceptual design:

type ValidationCategoryResult struct {
    Required int
    Matched  int
    Passed   bool
}

type ValidationMismatch struct {
    Type        string
    SiteID      string
    FloorID     string
    RegionID    string
    ReaderID    string
    AntennaPort int
    Message     string
}

type ValidationResult struct {
    SiteID       ValidationCategoryResult
    Floors       ValidationCategoryResult
    Regions      ValidationCategoryResult
    Readers      ValidationCategoryResult
    AntennaPorts ValidationCategoryResult
    Mismatches   []ValidationMismatch
    Passed       bool
}

Adjust to project conventions.

==================================================
VALIDATION TERMINAL DISPLAY
==================================================

When validation passes, show useful details.

Recommended output:

Site Configuration Validation

Recorded Site
  Site ID: ...
  Floors: 2
  Regions: 8
  Readers: 6
  Antenna Ports: 48

Target Site
  Site ID: ...
  Floors: 3
  Regions: 10
  Readers: 8
  Antenna Ports: 64

Validation Results

✓ Site ID matched
✓ 2 of 2 required floors matched
✓ 8 of 8 required regions matched
✓ 6 of 6 required readers matched
✓ 48 of 48 required antenna ports matched

Validation passed.
The recorded site is structurally compatible with the target site.

When validation fails:

Site Configuration Validation

Validation Results

✓ Site ID matched
✓ 2 of 2 required floors matched
✗ 7 of 8 required regions matched
✗ 5 of 6 required readers matched
✗ 47 of 48 required antenna ports matched

Validation failed.

Missing Region
  Region ID: ...
  Floor ID: ...

Missing Reader
  Reader ID: ...
  Region ID: ...
  Floor ID: ...

Missing Antenna Port
  Port: 4
  Reader ID: ...
  Region ID: ...
  Floor ID: ...

The terminal output must be readable for QA.

==================================================
TERMINAL OUTPUT VERSUS LOG FILES
==================================================

Do not show log-style prefixes in normal terminal output.

Avoid terminal lines such as:

[INFO]
[DEBUG]
[OK]
[ERROR]

Use clean messages:

Loaded 4 SiteGraphs.

Mock server is running at:
http://localhost:8080

Validation passed.

Validation failed.

Technical log levels should remain inside log files only.

Log files may contain:

INFO
DEBUG
WARN
ERROR

Separate:

user-facing terminal output

from:

technical file logging

Do not remove useful technical logs from log files.

==================================================
MOCK-SERVER TERMINAL DISPLAY
==================================================

Recommended startup output:

Resonate Replay Engine Mock Server

Configuration:
  configs\config.json

SiteGraph Directory:
  configs\sites

Loaded SiteGraphs: 4

Available Sites

1. Bentonville
   Site ID: ...
   Floors: 2
   Regions: 8
   Readers: 6
   Antenna Ports: 48

2. Site B
   Site ID: ...
   Floors: ...
   Regions: ...
   Readers: ...
   Antenna Ports: ...

Server:
  http://localhost:8080

Endpoints:
  GET  /sites
  GET  /sites/{siteId}
  POST /reader-bundles

Press Ctrl+C to stop.

Do not print complete SiteGraph JSON to the terminal.

==================================================
CONFIG MODE AND DIRECT ARGUMENT MODE
==================================================

The config class is central to the workflow.

Do not remove config mode.

Do not remove direct argument mode.

Both modes must call the same shared internal implementation.

Config validate:

.\rre.exe validate -config configs\config.json

Direct validate:

.\rre.exe validate `
  -file data\recording_001.sqlite `
  -target-url http://localhost:8080 `
  -site-id b3489888-aacf-4451-893c-d7d994240f93

Config play:

.\rre.exe play -config configs\config.json

Direct play:

.\rre.exe play `
  -file data\recording_001.sqlite `
  -target-url http://localhost:8080 `
  -site-id b3489888-aacf-4451-893c-d7d994240f93

Config summary:

.\rre.exe summary -config configs\config.json

Direct summary:

.\rre.exe summary `
  -file data\recording_001.sqlite

Config mock-server:

.\rre.exe mock-server -config configs\config.json

Keep current direct mock-server arguments working where already supported.

Update RunConfig to support:

- recording_file
- target_url
- site_id
- mock_port
- site_graph_directory

Input resolution should be:

1. Parse command arguments.
2. Load config when `-config` is supplied.
3. Apply supported direct values.
4. Validate the resolved values.
5. Call the same internal command implementation.

Do not maintain separate business logic for config and direct modes.

==================================================
PLAY COMMAND ORDER
==================================================

Preserve this exact order:

1. Resolve config or direct arguments.
2. Open Recorder SQLite.
3. Read SiteInformation.site_json.
4. Verify configured and recorded Site IDs.
5. Request the matching target SiteGraph.
6. Normalize recorded and target SiteGraphs.
7. Perform structural validation.
8. Abort immediately when validation fails.
9. Only after validation passes:
   - read RawReads
   - order by injection_time_utc
   - use raw_payload or the existing structured fallback
   - apply existing pacing
   - send to POST /reader-bundles

Do not break existing replay behavior.

Remember:

SiteInformation.site_json
→ validation

RawReads.raw_payload
→ playback

==================================================
RRE HELP COMMAND FOR QA
==================================================

Update:

.\rre.exe help

QA will use this as the main quick testing guide.

Keep it short and practical.

Add a section similar to:

QA SiteGraph Testing

1. Start the mock target:

   rre mock-server -config configs\config.json

2. View all mocked sites:

   Invoke-RestMethod http://localhost:8080/sites |
     ConvertTo-Json -Depth 10

3. Validate the recording against the configured Site ID:

   rre validate -config configs\config.json

4. Replay only after validation passes:

   rre play -config configs\config.json

5. Test a validation failure:
   - Stop mock-server
   - Open the matching file inside configs\sites\
   - Change or remove a Floor ID, Region ID, Reader ID or antenna port
   - Restart mock-server
   - Run validate again
   - Restore the original SiteGraph after testing

Also briefly explain:

GET /sites
→ lists available mocked sites

GET /sites/{siteId}
→ returns one complete SiteGraph

Keep the help concise.

Do not place the full README inside terminal help.

==================================================
README UPDATE
==================================================

Update README.md with:

1. What a SiteGraph is.
2. One JSON file represents one real Resonator site.
3. Difference between:
   - SiteInformation.site_json
   - RawReads.raw_payload
4. Real hierarchy:

Site
→ Floors
→ Regions
→ Readers
→ Antenna Ports

5. Recommended directory layout.
6. config.json format.
7. Why Site ID is required by the epic.
8. How all SiteGraphs are auto-discovered from configs/sites.
9. How root SiteGraph IDs are used as lookup keys.
10. How GET /sites works.
11. How GET /sites/{siteId} works.
12. How SQLite/site_id selects the correct target SiteGraph.
13. Validation rules.
14. Detailed pass output.
15. Detailed fail output.
16. QA testing workflow.
17. Config mode examples.
18. Direct argument examples.
19. Warning not to commit:
    - real customer SiteGraphs without approval
    - real Recorder SQLite databases
    - logs
    - runtime files
    - generated executables

==================================================
REAL RECORDER SQLITE SUPPORT
==================================================

Do not revert real Recorder schema support.

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

RecordingSession includes:

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

Do not reintroduce:

- test_description
- RawSiteJSON
- PascalCase SQL column names

Preserve nullable SQLite scanning.

==================================================
SECURITY REQUIREMENTS
==================================================

This repository uses Cycode security scanning.

Do not reintroduce previous findings.

Do not use:

- dynamic SQL
- fmt.Sprintf for SQL identifiers
- unrestricted arbitrary file paths
- unsafe parent traversal
- unrestricted arbitrary HTTP URLs
- removed SSRF protections
- ignored security findings

Use:

- hardcoded parameterized SQL
- read-only SQLite access
- filepath.Clean
- approved configs directory validation
- approved configs/sites directory validation
- `.json` extension validation
- existing safe URL validation patterns

For SiteGraph discovery:

- only allow the configured directory when it resolves inside the approved `configs/sites` area
- do not allow `..`
- do not recursively scan unexpected directories unless explicitly required
- reject symbolic/path escapes if supported by existing project security style
- return clear file-specific errors

Do not commit:

- rre.exe
- real SQLite files
- .sqlite-shm
- .sqlite-wal
- extracted .bin files
- replay_output.jsonl
- received_payloads.jsonl
- runtime logs
- temporary backup files

Update .gitignore where needed.

==================================================
FILES TO INSPECT
==================================================

Inspect existing code before editing.

Relevant files may include:

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

Create focused files if appropriate:

internal/site/validation_model.go
internal/site/parser.go
internal/site/parser_test.go
internal/site/summary.go
internal/mocktarget/site_store.go
internal/mocktarget/site_store_test.go

Use the exact module path from go.mod.

==================================================
TESTS
==================================================

Add or update tests for:

1. Discovering multiple SiteGraph JSON files from configs/sites.
2. Extracting root Site IDs.
3. Rejecting duplicate Site IDs.
4. Rejecting empty Site IDs.
5. Rejecting invalid JSON.
6. Rejecting unsafe SiteGraph directory paths.
7. Rejecting unsafe SiteGraph file paths.
8. GET /sites returns all available site summaries.
9. GET /sites/{siteId} returns the exact correct raw SiteGraph.
10. Unknown Site ID returns HTTP 404.
11. Multiple floors parse correctly.
12. Multiple regions under a floor parse correctly.
13. Multiple readers under a region parse correctly.
14. Reader antenna ports parse correctly.
15. Matching structures pass.
16. Missing floor fails.
17. Missing region under the correct floor fails.
18. Reader in the wrong region fails.
19. Missing reader fails.
20. Missing antenna port fails.
21. Extra target structures do not fail.
22. Configured Site ID mismatch fails.
23. SiteInformation Site ID mismatch fails.
24. Recorded SiteGraph root ID mismatch fails.
25. Validation category totals and matched counts are correct.
26. Clean terminal output does not contain [INFO], [OK], [ERROR] or [DEBUG].
27. Technical logs preserve log levels in log files.
28. Validate works through config mode.
29. Validate works through direct arguments.
30. Play stops before RawReads when validation fails.
31. Play continues when validation passes.
32. Summary remains working.
33. Version command remains working.
34. Help includes concise QA SiteGraph testing instructions.
35. Unknown full SiteGraph fields do not break parsing.
36. Complete raw SiteGraph JSON is preserved in GET /sites/{siteId}.

==================================================
IMPLEMENTATION PROCESS
==================================================

First inspect the current code.

Before editing, provide a brief implementation plan containing:

1. Existing relevant files.
2. Current incorrect assumptions.
3. Files to modify.
4. New files to create.
5. SiteGraph discovery/storage design.
6. Validation result design.
7. Terminal-versus-log separation.
8. Config and direct argument compatibility.
9. Help and README changes.

Then implement the changes.

Do not stop after only providing recommendations.

Make the actual code changes.

After implementation run:

gofmt on every changed Go file

go test ./...

go vet ./...

go build -o rre.exe ./cmd/rre

Do not stage or commit rre.exe.

Finally report:

1. Files changed.
2. Main design decisions.
3. Exact PowerShell commands for testing.
4. Expected mock-server startup output.
5. Expected GET /sites output.
6. Expected validation pass output.
7. Expected validation failure output.
8. Remaining assumptions requiring architect confirmation.