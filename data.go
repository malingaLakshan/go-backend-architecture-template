Please finalize the real Recorder SQLite schema integration for resonate-replay-engine.

Important:
Do not guess column names.
Do not use old camelCase DB columns.
Do not use test_description. Real column is description.
Do not create temporary duplicate files.
Do not commit logs, sqlite runtime side files, exe files, or generated output files.
Avoid Cycode/SAST issues: do not build SQL using fmt.Sprintf or string concatenation for table/column names. Use hardcoded SQL queries only.

Real SQLite schema from Recorder DB:

RecordingSession columns:
- recording_session_id TEXT
- test_name TEXT
- environment TEXT
- tester_name TEXT
- description TEXT
- start_time_utc TEXT
- end_time_utc TEXT
- resonate_build_number TEXT
- firmware_build_number TEXT
- reader_apps_build_number TEXT
- resonate_site_id TEXT
- state TEXT

RawReads columns:
- read_id TEXT
- recording_session_id TEXT
- tag_id TEXT
- reader_id TEXT
- antenna_id INTEGER
- antenna_type INTEGER
- source_timestamp_utc TEXT
- injection_time_utc TEXT
- confidence INTEGER
- rssi REAL
- tag_x INTEGER/REAL
- tag_y INTEGER/REAL
- floor_id INTEGER
- raw_payload BLOB

SiteInformation columns:
- site_information_id TEXT
- recording_session_id TEXT
- site_id TEXT
- site_name TEXT
- site_json BLOB

ResonateEvents columns:
- event_id TEXT
- recording_session_id TEXT
- tag_id TEXT
- event_type TEXT
- event_reason TEXT
- source_timestamp_utc TEXT
- injection_time_utc TEXT
- floor INTEGER
- x INTEGER
- y INTEGER
- z INTEGER
- region INTEGER
- event_details TEXT
- raw_payload BLOB

MLT_SOW_Locations columns:
- location_id TEXT
- recording_session_id TEXT
- tag_id TEXT
- source_timestamp_utc TEXT
- injection_time_utc TEXT
- floor INTEGER
- x INTEGER
- y INTEGER
- z INTEGER
- region INTEGER
- state TEXT
- confidence REAL
- raw_payload BLOB

Snapshots columns:
- snapshot_id TEXT
- recording_session_id TEXT
- timestamp_utc TEXT
- snapshot_name TEXT

SnapshotTagLocations columns:
- snapshot_tag_location_id TEXT
- snapshot_id TEXT
- tag_id TEXT
- x INTEGER
- y INTEGER
- z INTEGER
- floor INTEGER
- region INTEGER
- state TEXT

Task 1: Fix recording models

Update internal/recording/model.go to match real schema.

RecordingSession should map:
- RecordingSessionID <= recording_session_id
- TestName <= test_name
- Environment <= environment
- TesterName <= tester_name
- Description <= description
- StartTimeUTC <= start_time_utc parsed into time.Time
- EndTimeUTC <= end_time_utc parsed into time.Time
- ResonateBuildNumber <= resonate_build_number
- FirmwareBuildNumber <= firmware_build_number
- ReaderAppsBuildNumber <= reader_apps_build_number
- ResonateSiteID <= resonate_site_id
- State <= state

RawRead should map:
- ReadID string <= read_id
- RecordingSessionID string <= recording_session_id
- TagID string <= tag_id
- ReaderID string <= reader_id
- AntennaID int <= antenna_id
- AntennaTypeID int <= antenna_type
- SourceTimestampUtc string <= source_timestamp_utc
- InjectionTimeUtc string <= injection_time_utc
- Confidence int <= confidence
- RSSI float64 <= rssi
- TagX float64 <= tag_x
- TagY float64 <= tag_y
- FloorID int <= floor_id
- RawPayload []byte <= raw_payload
- Timestamp time.Time parsed from source_timestamp_utc
- InjectionTime time.Time parsed from injection_time_utc

SiteInformation should map:
- SiteInformationID string <= site_information_id
- RecordingSessionID string <= recording_session_id
- SiteID string <= site_id
- SiteName string <= site_name
- SiteJSON []byte <= site_json

Task 2: Fix repository queries

Update internal/recording/repository.go.

All SQL must use the real snake_case columns.

GetSession must use:

SELECT recording_session_id,
       test_name,
       environment,
       tester_name,
       description,
       start_time_utc,
       end_time_utc,
       resonate_build_number,
       firmware_build_number,
       reader_apps_build_number,
       resonate_site_id,
       state
FROM RecordingSession
WHERE recording_session_id = ?

Do not use RecordingSessionID.
Do not use TestDescription.
Do not use test_description.
Do not use SiteID.
Use description and resonate_site_id.

GetFirstSession must use the same real columns and:
FROM RecordingSession
LIMIT 1

GetSiteInfo must use:

SELECT site_information_id,
       recording_session_id,
       site_id,
       site_name,
       site_json
FROM SiteInformation
WHERE site_id = ?

GetRawReads must use:

SELECT read_id,
       recording_session_id,
       tag_id,
       reader_id,
       antenna_id,
       antenna_type,
       source_timestamp_utc,
       injection_time_utc,
       confidence,
       rssi,
       tag_x,
       tag_y,
       floor_id,
       raw_payload
FROM RawReads
WHERE recording_session_id = ?
ORDER BY injection_time_utc ASC, read_id ASC

GetRawReadTimeRange must use:

SELECT MIN(injection_time_utc),
       MAX(injection_time_utc)
FROM RawReads
WHERE recording_session_id = ?

GetUniqueCount must not build SQL dynamically.
Use hardcoded allowed switch cases only:

- table == "RawReads" and column == "tag_id"
  SELECT COUNT(DISTINCT tag_id) FROM RawReads WHERE recording_session_id = ?

- table == "RawReads" and column == "reader_id"
  SELECT COUNT(DISTINCT reader_id) FROM RawReads WHERE recording_session_id = ?

Also support legacy caller input safely:
- "TagID" should map to tag_id query
- "ReaderID" should map to reader_id query

But still do not use dynamic SQL.

Task 3: Fix summary

Update internal/recording/summary.go.

Summary must use real columns and real functions.
Do not use old column names.

Summary must show:
- Recording Session ID
- Test Name
- Environment
- Tester Name
- Description
- Start Time
- End Time
- Resonate Build Number
- Firmware Build Number
- Reader Apps Build Number
- Resonate Site ID
- State
- Total RawReads
- Unique Tags
- Unique Readers
- First Injection Time
- Last Injection Time

summary command must work with direct file:

./rre.exe summary -file data/recording_001.sqlite

summary command must also work with config:

./rre.exe summary -config configs/pass_config.json

For summary with config:
- only recording_file is required
- do not require target_url
- do not require site_id

Task 4: Fix config support

Update internal/cli/args.go and internal/cli/commands.go.

Flags must include Config string.

-config must work for:
- summary
- validate
- play
- mock-server

Commands:

./rre.exe summary -config configs/pass_config.json
./rre.exe validate -config configs/pass_config.json
./rre.exe play -config configs/pass_config.json
./rre.exe mock-server -config configs/pass_config.json

Config file structure:

{
  "recording_file": "data/recording_001.sqlite",
  "target_url": "http://localhost:8080",
  "site_id": "b3489888-aacf-4451-893c-d7d994240f93",
  "mock_mode": "sqlite",
  "mock_site_file": ""
}

For validate/play:
- require recording_file, target_url, site_id

For summary:
- require only recording_file

For mock-server:
- if mock_mode is sqlite, require recording_file and site_id
- if mock_mode is file, require mock_site_file

Task 5: Fix SiteJSON usage

Do not use RawSiteJSON anywhere.

For validate and play:
- load recorded site info using GetSiteInfo
- unmarshal siteInfo.SiteJSON into site.SiteConfig
- error message should say SiteJSON, not RawSiteJSON

Example:
json.Unmarshal(siteInfo.SiteJSON, &recordedConfig)

Task 6: Fix validator

Update internal/site/validator.go.

Validation must check actual IDs, not only counts.

Validate:
- site ID equality
- reader ID equality
- antenna ID equality
- floor ID equality
- region ID equality

If IDs differ, validation must fail even when counts are same.

Error messages should be clear:
- Site ID mismatch: recorded=..., target=...
- Reader ID missing in target: ...
- Reader ID missing in recorded: ...
- Antenna ID missing in target: ...
- Antenna ID missing in recorded: ...
- Floor ID missing in target: ...
- Floor ID missing in recorded: ...
- Region ID missing in target: ...
- Region ID missing in recorded: ...

Task 7: Fix mock server config mode

mock-server should support:

./rre.exe mock-server -config configs/pass_config.json
./rre.exe mock-server -config configs/fail_config.json

pass_config.json:
- mock_mode sqlite
- load site config from Recording DB SiteInformation.site_json
- serve exact recorded site config

fail_config.json:
- mock_mode file
- load mock_site_file
- serve wrong config

wrong_site_config.json must be valid JSON and must match SiteConfig structure.
Do not put strings inside readers array.
readers, floors, regions, antennas must be arrays of objects.

Task 8: Fix replay payload

RawPayload is []byte/BLOB.
Do not compare RawPayload with "".

Use:
if len(rawRead.RawPayload) == 0

Unmarshal:
json.Unmarshal(rawRead.RawPayload, &payloadData)

ReadID is string.
Use %s for ReadID errors, not %d.

Task 9: .gitignore safety

Update .gitignore to avoid committing generated/runtime files:

rre.exe
logs/*.jsonl
*.sqlite-shm
*.sqlite-wal
*.tmp
*.backup
*.bin

Do not ignore required config JSON files.
Do not ignore configs/pass_config.json.
Do not ignore configs/fail_config.json.
Do not ignore configs/wrong_site_config.json.

Do not commit:
- logs/replay_output.jsonl
- logs/received_payloads.jsonl
- sqlite -shm files
- sqlite -wal files
- rre.exe

Task 10: README and help

Update README and CLI help with working commands:

Build:
go build -o rre.exe ./cmd/rre

Summary:
./rre.exe summary -file data/recording_001.sqlite
./rre.exe summary -config configs/pass_config.json

Start mock server pass mode:
./rre.exe mock-server -config configs/pass_config.json

Validate pass:
./rre.exe validate -config configs/pass_config.json

Play:
./rre.exe play -config configs/pass_config.json

Start mock server fail mode:
./rre.exe mock-server -config configs/fail_config.json

Validate fail:
./rre.exe validate -config configs/fail_config.json

Task 11: Tests

Update tests so they match the real schema.

Tests must not expect old camelCase columns.

Add/adjust tests for:
- GetSession uses description column
- GetFirstSession uses description column
- GetRawReads uses real RawReads columns
- GetSiteInfo uses site_json
- Summary works with real schema
- Summary works with config
- Validate pass config
- Validate fail config
- Validator fails when counts match but IDs differ
- Wrong site config JSON unmarshals successfully

Task 12: Final verification

Run:

gofmt -w ./internal/cli ./internal/config ./internal/mocktarget ./internal/recording ./internal/replay ./internal/site ./internal/sqlite
go test ./...
go build -o rre.exe ./cmd/rre

Then manually test:

./rre.exe summary -config configs/pass_config.json

Terminal 1:
./rre.exe mock-server -config configs/pass_config.json

Terminal 2:
./rre.exe validate -config configs/pass_config.json

Terminal 2:
./rre.exe play -config configs/pass_config.json

Stop server.

Terminal 1:
./rre.exe mock-server -config configs/fail_config.json

Terminal 2:
./rre.exe validate -config configs/fail_config.json

Expected:
- pass config validation passes
- fail config validation fails with clear mismatch errors
- no Cycode SAST issue
- no generated logs/exe/sqlite runtime files staged