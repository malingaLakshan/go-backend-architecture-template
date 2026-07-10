We need to restore the real Recorder SQLite schema support in the RRE CLI before doing config/README work.

The project was already adapted to real Recorder SQLite output, but after folder rename/rebase some repository/model/schema code reverted to the old sample camelCase schema. Please fix all recording-related models, repositories, summary, sample schema, and tests to consistently use the real Recorder SQLite schema.

Important:
The real Recorder DB uses snake_case columns, not camelCase.
Do not use old columns like RecordingSessionID, SiteID, RawSiteJSON, ReaderID, TagID, InjectionTime, Timestamp, RawPayload.

Use the real schema below.

1. RecordingSession table

Use real snake_case columns.

Expected columns:
- recording_session_id
- test_name
- environment
- tester_name
- test_description
- site_id
- start_time_utc
- end_time_utc
- resonate_build_number
- firmware_build_number
- reader_apps_build_number
- resonate_site_id
- state

Update RecordingSession model to include these fields:

type RecordingSession struct {
    RecordingSessionID string
    TestName string
    Environment string
    TesterName string
    TestDescription string
    SiteID string
    StartTimeUTC time.Time
    EndTimeUTC time.Time
    ResonateBuildNumber string
    FirmwareBuildNumber string
    ReaderAppsBuildNumber string
    ResonateSiteID string
    State string
}

Update repository functions:
- GetSession
- GetFirstSession

They must query snake_case columns:

SELECT recording_session_id,
       test_name,
       environment,
       tester_name,
       test_description,
       site_id,
       start_time_utc,
       end_time_utc,
       resonate_build_number,
       firmware_build_number,
       reader_apps_build_number,
       resonate_site_id,
       state
FROM RecordingSession
WHERE recording_session_id = ?

And GetFirstSession should use:
FROM RecordingSession
LIMIT 1

2. SiteInformation table

Use real Recorder schema.

Expected columns:
- site_information_id
- recording_session_id
- site_id
- site_name
- site_json

Update SiteInformation model:

type SiteInformation struct {
    SiteInformationID string
    RecordingSessionID string
    SiteID string
    SiteName string
    SiteJSON []byte
}

Update GetSiteInfo to use:

SELECT site_information_id,
       recording_session_id,
       site_id,
       site_name,
       site_json
FROM SiteInformation
WHERE site_id = ?

Do not use RawSiteJSON anymore.
Any error message saying RawSiteJSON must be changed to SiteJSON.

3. RawReads table

Use real Recorder schema.

Expected columns:
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

Update RawRead model:

type RawRead struct {
    ReadID string
    RecordingSessionID string
    TagID string
    Z float64
    ReaderID string
    AntennaID int
    AntennaTypeID int
    SourceTimestampUtc string
    InjectionTimeUtc string
    Confidence int
    RSSI float64
    TagX float64
    TagY float64
    FloorID int
    RawPayload []byte

    Timestamp time.Time
    InjectionTime time.Time
}

Update GetRawReads to query:

SELECT read_id,
       recording_session_id,
       tag_id,
       z,
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

Scan into the RawRead fields above.
Parse SourceTimestampUtc into Timestamp.
Parse InjectionTimeUtc into InjectionTime.

4. Time range

Update GetRawReadTimeRange to use:

SELECT MIN(injection_time_utc),
       MAX(injection_time_utc)
FROM RawReads
WHERE recording_session_id = ?

5. Unique counts

Update GetUniqueCount to use only hardcoded safe SQL queries.

Support both old caller names and new snake_case names to avoid breaking existing code:
- "TagID" and "tag_id" should both count DISTINCT tag_id
- "ReaderID" and "reader_id" should both count DISTINCT reader_id

Use recording_session_id in WHERE clause.

Example:
SELECT COUNT(DISTINCT tag_id)
FROM RawReads
WHERE recording_session_id = ?

SELECT COUNT(DISTINCT reader_id)
FROM RawReads
WHERE recording_session_id = ?

Do not build SQL dynamically using fmt.Sprintf.

6. Summary

Update summary logic so it works with the real Recorder schema.

summary should print:
- total records
- unique readers
- unique tags
- first injection time
- last injection time
- total duration

It must use:
- RawReads
- tag_id
- reader_id
- injection_time_utc
- recording_session_id

7. Payload builder

Update payload.go for real Recorder RawRead.

RawPayload is []byte, not string.

Fix:
- if rawRead.RawPayload == "" is wrong
- use len(rawRead.RawPayload) == 0
- json.Unmarshal(rawRead.RawPayload, &payloadData)

If raw_payload is binary/empty/unusable, BuildPayload should build from structured RawRead columns:
- reader_id
- tag_id
- antenna_id
- antenna_type
- confidence
- tag_x
- tag_y
- floor_id
- injection timestamp

Add a clear error only if both raw_payload is empty AND structured fields are not enough.

8. Site JSON usage in commands

In runValidate and runPlay:
Use:

json.Unmarshal(siteInfo.SiteJSON, &recordedConfig)

Do not use:
siteInfo.RawSiteJSON
[]byte(siteInfo.RawSiteJSON)

Change error text:
"Failed to parse recorded SiteJSON"

9. SQLite sample schema and sample data

Update internal/sqlite/schema.go and sample_data.go to match the real Recorder schema, not the old camelCase schema.

RecordingSession should be created with snake_case columns.
SiteInformation should use site_json.
RawReads should use real snake_case columns.

10. Tests

Update tests to use real Recorder schema.
Remove old camelCase column references from tests.

Search and remove/fix all old names:
- RecordingSessionID
- SiteID
- RawSiteJSON
- ReaderID
- TagID
- InjectionTime
- Timestamp
- RawPayload

They should not appear inside SQL queries anymore, except maybe as Go struct field names.

11. Build

After changes:
- run gofmt
- run go test ./...
- run go build -o rre.exe ./cmd/rre

Manual checks:
./rre.exe summary -file data/recording_001.sqlite
./rre.exe validate -file data/recording_001.sqlite -target-url http://localhost:8080 -site-id b3489888-aacf-4451-893c-d7d994240f93

Please only focus on restoring real Recorder SQLite schema support first. Do not update README/config flow yet.