COPY THIS INTO GEMINI

Create a concise text file named:
Resonate_MQTT_Recorder_Field_Gap_Analysis.txt

This is an R&D report only. Do not change code or the database. Verify exact
SQLite and Protobuf types from the repository. Do not guess unknown fields.

1. RAW RFID
===========

Topic:
resonate/locate/3b96f652-8200-3920-8a2c-0486c358964e/rawrfid

Bundle fields:
reader_id, reads, site_id, sent_timestamp_ms

Individual read fields:
timestamp_ns, confidence, antenna_id, antenna_type, x, y, item_id, floor_id

RawReads database columns:
read_id, recording_session_id, tag_id, reader_id, antenna_id, antenna_type,
source_timestamp_utc, injection_time_utc, confidence, rssi, tag_x, tag_y,
floor_id, raw_payload

Mappings:
- tag_id <- item_id
- reader_id <- bundle reader_id
- antenna_id <- antenna_id
- antenna_type <- antenna_type
- source_timestamp_utc <- timestamp_ns (convert to UTC)
- confidence <- confidence
- tag_x <- x
- tag_y <- y
- floor_id <- floor_id
- read_id and recording_session_id are Recorder-generated
- raw_payload stores the original payload if confirmed by code

Highlight:
- MISSING: rssi is not in the Raw RFID Protobuf.
- NOT POPULATED: floor_id exists but was empty in live messages.
- UNCONFIRMED: injection_time_utc source.
- MQTT-ONLY: bundle site_id and sent_timestamp_ms have no obvious RawReads
  columns.

2. LOCATIONS
============

Topic:
resonate/locate/3b96f652-8200-3920-8a2c-0486c358964e/locationUpdate

Bundle fields:
site_id, locations, sent_timestamp_ms

Individual fields:
1=UNKNOWN, 2=item_id, 3=x, 4=y, 5=z, 6=floor, 7=region,
8=confidence, 9=state, 10=timestamp_ns, 11=reason

MLT_SOW_Locations database columns:
location_id, recording_session_id, tag_id, source_timestamp_utc,
injection_time_utc, floor, x, y, z, region, state, confidence, raw_payload

Mappings:
- tag_id <- field 2 item_id
- x <- field 3
- y <- field 4
- z <- field 5
- floor <- field 6
- region <- field 7
- confidence <- field 8
- state <- field 9
- source_timestamp_utc <- field 10 timestamp_ns (convert to UTC)
- location_id and recording_session_id are Recorder-generated
- raw_payload stores the original payload if confirmed by code

Highlight:
- MISSING DB COLUMN: field 11 reason.
- UNKNOWN: field 1 name, meaning and database mapping.
- UNCONFIRMED: state/reason enum text mappings.
- UNCONFIRMED: injection_time_utc source.
- MQTT-ONLY: bundle site_id and sent_timestamp_ms have no obvious location-table
  columns.

OUTPUT
======

Keep the report within about two pages and include only:

1. One short summary.
2. Raw RFID comparison table:
   MQTT field | Database field | Status
3. Location comparison table:
   MQTT field | Database field | Status
4. Missing/Unconfirmed Fields list.

Include every field listed above. Confirm exact types from the repository and
mark unsupported information as UNCONFIRMED. Do not add implementation steps,
architecture explanations or unnecessary background.

END