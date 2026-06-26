Now update the replay request payload format to match the required Resonate Replay Engine output JSON contract.

Required replay output JSON shape:

{
“ProtoReaderBundle”: {
“reader_id”: 42,
“reads”: [
{
“timestamp_ns”: 42,
“confidence”: 42,
“antenna_id”: 42,
“antenna_type”: 2,
“x”: 42,
“y”: 42,
“item_id”: “Sample text”,
“floor_id”: 42
}
],
“site_id”: “Sample text”,
“sent_timestamp_ms”: 42
}
}

Important:
The values above are only example values. Do not hardcode 42 or “Sample text” in replay logic.

Replay rules:

1. Read records from the SQLite RawReads table.
2. Records must be replayed strictly ordered by InjectionTime.
3. Each replayed HTTP request body must use the JSON wrapper key ProtoReaderBundle.
4. ProtoReaderBundle must contain:
    * reader_id
    * reads
    * site_id
    * sent_timestamp_ms
5. reads must be an array of ProtoRead objects.
6. Each ProtoRead object must contain:
    * timestamp_ns
    * confidence
    * antenna_id
    * antenna_type
    * x
    * y
    * item_id
    * floor_id
7. If the SQLite row already contains payload JSON, inspect it and convert/wrap it into the required ProtoReaderBundle structure if needed.
8. Do not send reads as an array of strings.
9. Do not change the database schema unless absolutely necessary.
10. Do not change site validation logic in this task.
11. Do not change CLI command names in this task.
12. Keep the existing pacing behavior based on InjectionTime.

Endpoint:
Continue sending replay payloads to the existing target HTTP endpoint used by the project, such as /reader-bundles, unless the current code already uses a different endpoint intentionally.

Add clear unit/helper functions if needed:

* Build replay payload from SQLite raw read row
* Validate replay payload shape before sending
* Marshal replay payload to JSON

Add one debug-safe log line that shows replay payload was prepared, but do not print full payload for every record unless debug mode already exists.

After changes, show:

* changed files
* sample generated replay payload
* command used to test replay