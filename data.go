We need to verify and fix only the replay HTTP request JSON body format.

Important:
Do not change InjectionTime pacing logic.
Do not change replay timing behavior.
Do not change record ordering.
Do not change validation logic.
Do not change CLI command names.
Do not change dashboard logic.

Current replay timing rule must remain:

* Read RawReads ordered by InjectionTime
* Send first record immediately
* For each next record, wait the time difference between current InjectionTime and previous InjectionTime
* Then send the next replay payload

Task:
Check the actual JSON body sent to the target /reader-bundles endpoint.

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

* Do not hardcode 42 or “Sample text”.
* Values must come from the SQLite RawReads payload/recorded data.
* The HTTP request body must have top-level key ProtoReaderBundle.
* reads must be an array of read objects, not strings.
* Keep existing InjectionTime-based scheduling exactly as it is.

For verification:

* Update mock target server debug/test output so it can show the first received replay JSON body clearly.
* It is okay to print only the first received payload or save received payloads into a debug file like logs/received_payloads.jsonl.
* Keep normal logs professional and not too noisy.

After changes, show:

1. Which function builds the replay payload
2. The exact first received JSON body from mock server
3. Confirmation that InjectionTime pacing was not changed
4. Changed files