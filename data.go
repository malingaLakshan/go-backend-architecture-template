You are working inside this Go project:

C:\projects\ALTRFIDTools\resonate-replay-engine

Work only inside `resonate-replay-engine`.

Before changing anything, inspect the current project implementation carefully.

Do not remove or break any existing working functionality.

Current functionality that must remain working:

- help
- version
- generate-sample
- summary
- mock-server
- validate
- play
- config-file execution
- direct command-line argument execution
- real Recorder SQLite schema support
- SiteInformation.site_json validation
- RawReads.raw_payload playback
- injection_time_utc pacing
- versioning
- Cycode security fixes
- JSONL replay/mock output

Do not commit, push, or create a PR.

==================================================
DOMAIN AND CURRENT REQUIREMENT
==================================================

A SiteGraph is the complete configuration of one real Resonator site.

We currently have approximately four real SiteGraph JSON files.

Each SiteGraph represents a different target site and has a unique root `id`.

The real target API behaves like:

GET /sites/{siteId}

and returns the complete SiteGraph belonging to that Site ID.

The mock server must behave like a realistic Resonator target.

It must load all configured SiteGraph JSON files at startup and return the correct SiteGraph based on the requested Site ID.

Example:

GET /sites/site-a-id
→ returns SiteGraph A

GET /sites/site-b-id
→ returns SiteGraph B

GET /sites/unknown-id
→ returns HTTP 404 with a clear site-not-found response

The Recorder SQLite contains:

SiteInformation.site_json
→ the complete site configuration recorded with that recording

RawReads.raw_payload
→ the actual data used for playback

Validation must compare:

Recorded SiteInformation.site_json

against:

Target SiteGraph returned by GET /sites/{siteId}

Playback must start only after validation passes.

==================================================
REAL SITEGRAPH HIERARCHY
==================================================

The confirmed hierarchy is:

Site
  -> Floors
      -> Regions
          -> Readers
              -> Antennas

Confirmed details:

1. One site may have multiple floors.
2. Each floor contains regions.
3. Regions are directly inside the floor for the current Recorder output.
4. Each region contains readers.
5. Each reader contains antennas.
6. Reader identity uses the string field:

"id"

7. Antennas do not have a separate ID.
8. Antennas are identified by the integer field:

"port"

9. One reader may contain antenna ports 1 through 8.

Fields such as names, coordinates, bounds, networking, rotation, timeouts, physicality, behaviors and inventory type are not required for minimum structural validation.

Do not create hundreds of Go fields for every SiteGraph property.

==================================================
RECOMMENDED PROJECT STRUCTURE
==================================================

Use one central config and one separate full JSON file per real site.

Expected structure:

resonate-replay-engine/
├── configs/
│   ├── config.json
│   └── sites/
│       ├── bentonville.json
│       ├── site-b.json
│       ├── site-c.json
│       └── site-d.json
│
├── data/
├── logs/
├── internal/
└── README.md

Each file inside `configs/sites/` must contain one complete real SiteGraph.

Do not combine all full SiteGraphs into one giant file.

The central config should look like:

{
  "recording_file": "data/recording_001.sqlite",
  "target_url": "http://localhost:8080",
  "site_id": "b3489888-aacf-4451-893c-d7d994240f93",
  "mock_port": 8080,
  "site_graph_files": [
    "configs/sites/bentonville.json",
    "configs/sites/site-b.json",
    "configs/sites/site-c.json",
    "configs/sites/site-d.json"
  ]
}

Use actual existing SiteGraph filenames if they already exist.

Do not insert real customer SiteGraphs into source code.

Do not change the Site ID inside config.json automatically based on filename.

The root `id` inside each SiteGraph is the source of truth for that site.

==================================================
SITEGRAPH LOADING
==================================================

When mock-server starts using:

.\rre.exe mock-server -config configs\config.json

it must:

1. Load config.json.
2. Read every path in site_graph_files.
3. Safely read each SiteGraph JSON file.
4. Parse the root `id`.
5. Reject a SiteGraph with:
   - invalid JSON
   - empty root id
   - duplicate root id
6. Store each complete raw SiteGraph by Site ID.

Conceptually:

map[string][]byte

where:

key = root SiteGraph ID
value = complete original JSON bytes

Do not unmarshal the complete SiteGraph into an incomplete Go struct and then marshal it again.

The complete raw JSON must be preserved.

GET /sites/{siteId} must return the original complete SiteGraph JSON bytes.

==================================================
LIST AVAILABLE SITES
==================================================

Add this mock-server endpoint:

GET /sites

This must list all currently loaded mocked sites.

Do not add a separate `rre sites` command in this first implementation.

The API endpoint is the core design.

Recommended response:

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

The list should contain:

- Site ID
- Site name if available
- floor count
- region count
- reader count
- antenna-port count

The counts must include all floors, regions, readers and antennas in the real hierarchy.

GET /sites/{siteId} must continue to return the complete matching SiteGraph.

Update mock-server routing to support both:

GET /sites
GET /sites/{siteId}
POST /reader-bundles

==================================================
STRUCTURAL VALIDATION MODEL
==================================================

Create a focused internal normalized validation model similar to:

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

Follow current project conventions and exact module path from go.mod.

Suggested files:

internal/site/validation_model.go
internal/site/parser.go
internal/site/summary.go

The full SiteGraph parser must:

- read root `id`
- read every floor under `floors`
- read every floor `id`
- read every region under floor `regions`
- read every region `id`
- read every reader under region `readers`
- read every reader `id`
- read every antenna `port` under reader `antennas`
- process every site structure, not only the first one
- ignore unrelated fields safely
- return descriptive errors for malformed required data
- not silently produce false zero values

==================================================
VALIDATION RULES
==================================================

Before playback, validate:

1. Site ID matches.
2. Every recorded Floor ID exists in the target.
3. Every recorded Region ID exists under the correct Floor ID.
4. Every recorded Reader ID exists under the correct Region ID.
5. Every recorded antenna port exists under the correct Reader ID.

The target may contain additional structures.

Only recorded-required structures must exist in the target.

Preserve hierarchy.

Do not compare only global counts.

Use structural keys where helpful:

floorID

floorID + regionID

floorID + regionID + readerID

floorID + regionID + readerID + antennaPort

Validation errors must be clear, for example:

Configuration mismatch: target site is missing Floor ID <id>

Configuration mismatch: target site is missing Region ID <id> under Floor ID <floor-id>

Configuration mismatch: target site is missing Reader ID <id> under Region ID <region-id>

Configuration mismatch: target site is missing antenna port 4 under Reader ID <reader-id>

The play command must abort before reading RawReads when validation fails.

==================================================
VALIDATION TERMINAL OUTPUT
==================================================

Successful validation must show useful category details, not only “Validation passed”.

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
  Floors: 2
  Regions: 8
  Readers: 6
  Antenna Ports: 48

Validation Results

✓ Site ID matched
✓ 2 of 2 required floors matched
✓ 8 of 8 required regions matched
✓ 6 of 6 required readers matched
✓ 48 of 48 required antenna ports matched

Validation passed.
The recorded site is structurally compatible with the target site.

Failure output should show both the summary and exact mismatches:

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

Use a structured validation result model rather than only storing plain error strings.

Possible categories:

- Site ID
- Floors
- Regions
- Readers
- Antenna Ports

==================================================
TERMINAL OUTPUT AND LOGGING
==================================================

Remove log-style prefixes from normal terminal output.

Do not show terminal messages like:

[INFO]
[DEBUG]
[OK]
[ERROR]

The terminal should show clean user-facing messages.

Example:

Loaded 4 SiteGraphs.

Mock server is running at:
http://localhost:8080

Validation passed.

Instead of:

[INFO] Loaded 4 SiteGraphs
[OK] Validation passed

Keep proper log levels inside log files only.

Log files may still contain:

INFO
DEBUG
WARN
ERROR

Separate:

terminal presentation
from
technical file logging

Do not remove useful technical logging from log files.

==================================================
MOCK-SERVER TERMINAL OUTPUT
==================================================

Recommended startup output:

Resonate Replay Engine Mock Server

Configuration: configs\config.json
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

Do not print the full SiteGraph in the terminal.

==================================================
CONFIG AND DIRECT ARGUMENT SUPPORT
==================================================

The config file is the main QA workflow.

Do not remove direct argument support.

Both modes must resolve into the same internal command logic.

Config mode:

.\rre.exe mock-server -config configs\config.json

.\rre.exe validate -config configs\config.json

.\rre.exe play -config configs\config.json

.\rre.exe summary -config configs\config.json

Direct mode:

.\rre.exe validate `
  -file data\recording_001.sqlite `
  -target-url http://localhost:8080 `
  -site-id <site-id>

.\rre.exe play `
  -file data\recording_001.sqlite `
  -target-url http://localhost:8080 `
  -site-id <site-id>

.\rre.exe summary `
  -file data\recording_001.sqlite

Keep existing direct mock-server flags working where currently supported.

RunConfig remains central.

Update internal/config/model.go to support:

- recording_file
- target_url
- site_id
- mock_port
- site_graph_files

Preserve compatibility only where reasonable and clearly documented.

==================================================
RECORDED SQLITE SITE SELECTION
==================================================

Validation should automatically select the correct target SiteGraph using the recorded Site ID.

Flow:

1. Open Recorder SQLite.
2. Read SiteInformation for the configured or recorded Site ID.
3. Read the recorded site_json.
4. Determine the recorded Site ID.
5. Call:

GET {target_url}/sites/{recordedSiteId}

6. Mock server searches the loaded SiteGraph map.
7. Mock server returns the matching SiteGraph.
8. RRE compares recorded and target structures.

QA should not manually select one of the four SiteGraph files for every validation.

The Site ID selects the correct target SiteGraph.

If no matching mocked SiteGraph exists, return a clear failure:

Target site configuration was not found for Site ID <id>

==================================================
PLAY COMMAND
==================================================

Preserve this exact order:

1. Resolve config or arguments.
2. Open Recorder SQLite.
3. Read SiteInformation.site_json.
4. Determine Site ID.
5. Fetch target SiteGraph.
6. Normalize recorded and target structures.
7. Validate.
8. Abort immediately if validation fails.
9. Only after validation succeeds:
   - read RawReads
   - order by injection_time_utc
   - use raw_payload or existing structured fallback
   - send to POST /reader-bundles

Do not break replay timing, payload construction or output logging.

Remember:

SiteInformation.site_json
→ validation

RawReads.raw_payload
→ playback

==================================================
HELP COMMAND IS IMPORTANT FOR QA
==================================================

Update the output of:

.\rre.exe help

QA will use this as the main testing guide.

Keep it short, readable and practical.

Add a section similar to:

QA SiteGraph Testing

1. Start the mock target:

   rre mock-server -config configs\config.json

2. View available mocked sites:

   Invoke-RestMethod http://localhost:8080/sites |
     ConvertTo-Json -Depth 10

3. Validate a recording:

   rre validate -config configs\config.json

4. Replay after validation passes:

   rre play -config configs\config.json

5. To test failure:
   - Open the required file inside configs\sites\
   - Change or remove a Floor ID, Region ID, Reader ID or antenna port
   - Restart mock-server
   - Run validate again
   - Restore the original SiteGraph after testing

Also show direct argument examples.

Mention:

GET /sites
→ lists all mocked sites

GET /sites/{siteId}
→ returns one complete SiteGraph

Keep the help concise enough for terminal use.

Do not turn help into the full README.

==================================================
README UPDATE
==================================================

Update README.md more fully.

Include:

1. Difference between:
   - SiteInformation.site_json
   - RawReads.raw_payload

2. SiteGraph meaning:
   - one JSON file represents one real Resonator site

3. SiteGraph hierarchy:

Site
→ Floors
→ Regions
→ Readers
→ Antenna Ports

4. Recommended directory structure:

configs/
├── config.json
└── sites/
    ├── bentonville.json
    ├── site-b.json
    ├── site-c.json
    └── site-d.json

5. config.json format.

6. How mock-server loads all sites.

7. How GET /sites works.

8. How GET /sites/{siteId} works.

9. How the recorded SQLite Site ID selects the correct target SiteGraph.

10. Validation rules.

11. Detailed pass and fail outputs.

12. QA test workflow.

13. Config and direct argument examples.

14. Warning:
    - do not commit customer/private SiteGraphs without approval
    - do not commit real Recorder SQLite databases
    - do not commit logs or generated files

Keep README accurate and practical.

==================================================
REAL RECORDER SQLITE SCHEMA
==================================================

Do not revert real Recorder SQLite support.

SiteInformation:

- site_information_id
- recording_session_id
- site_id
- site_name
- site_json

RawReads:

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
- old PascalCase SQL columns

Preserve nullable SQLite scanning.

==================================================
SECURITY REQUIREMENTS
==================================================

This repository uses Cycode security scanning.

Do not reintroduce previous issues.

Do not use:

- dynamic SQL
- fmt.Sprintf for SQL table or column identifiers
- unsafe unrestricted file paths
- path traversal
- unrestricted arbitrary HTTP URLs
- removed SSRF protections
- automatic ignoring of security findings

Use:

- hardcoded parameterized SQL
- read-only SQLite access
- filepath.Clean
- approved configs directory checks
- approved `configs/sites` directory checks
- JSON extension validation
- safe URL validation already present in the project

For SiteGraph files:

- only allow files under the approved `configs/sites` directory
- reject parent traversal
- reject non-JSON files
- return clear errors
- do not follow arbitrary external paths

Do not commit:

- rre.exe
- real company SQLite databases
- .sqlite-shm
- .sqlite-wal
- .bin extracts
- replay_output.jsonl
- received_payloads.jsonl
- runtime logs
- temporary backup files

Do not modify sibling projects:

../resonate-recorder
../resonate-analyzer

==================================================
FILES TO INSPECT
==================================================

Inspect the current implementation first.

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

Create focused files where helpful:

internal/site/validation_model.go
internal/site/parser.go
internal/site/parser_test.go
internal/site/summary.go
internal/mocktarget/site_store.go
internal/mocktarget/site_store_test.go

Use the exact Go module path from go.mod.

==================================================
TESTS
==================================================

Add or update tests for:

1. Loading multiple SiteGraph files.
2. Extracting the root Site ID.
3. Rejecting duplicate Site IDs.
4. Rejecting missing Site IDs.
5. Rejecting invalid SiteGraph JSON.
6. Rejecting unsafe SiteGraph paths.
7. GET /sites returns all loaded site summaries.
8. GET /sites/{siteId} returns the correct complete raw SiteGraph.
9. Unknown Site ID returns 404.
10. Multiple floors parse correctly.
11. Regions under the correct floor parse correctly.
12. Readers under the correct region parse correctly.
13. Antenna ports under the correct reader parse correctly.
14. Matching recorded and target structures pass.
15. Missing floor fails.
16. Missing region in the correct floor fails.
17. Reader in the wrong region fails.
18. Missing reader fails.
19. Missing antenna port fails.
20. Extra target structures do not fail.
21. Validation category totals and matched counts are correct.
22. Terminal formatter produces clean output without [INFO]/[OK]/[ERROR].
23. Technical logger keeps log levels in log files.
24. Validate works using config mode.
25. Validate works using direct arguments.
26. Play stops before RawReads when validation fails.
27. Play continues when validation succeeds.
28. Summary remains working.
29. Version command remains working.
30. Help contains the concise QA workflow.
31. Unknown full SiteGraph fields do not break parsing.
32. Full raw JSON is preserved when returned by GET /sites/{siteId}.

==================================================
IMPLEMENTATION PROCESS
==================================================

First inspect the current code.

Before editing, provide a short plan containing:

1. Existing relevant files.
2. Current incorrect assumptions.
3. Files to modify.
4. New files to add.
5. SiteGraph storage architecture.
6. Validation result design.
7. Terminal versus file-log separation.
8. Config and direct argument compatibility.
9. Help and README update plan.

Then implement the changes.

Do not stop after only explaining.

Make the actual code changes.

After implementation run:

gofmt on all changed Go files

go test ./...

go vet ./...

go build -o rre.exe ./cmd/rre

Do not stage or commit rre.exe.

Finally report:

1. Files changed.
2. Main design decisions.
3. Exact PowerShell commands to test.
4. Expected GET /sites output.
5. Expected validation-pass output.
6. Expected validation-fail output.
7. Any assumptions still requiring architect confirmation.