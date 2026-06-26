The mock server is still only printing summary logs like:

Received bundle: site_id=SITE-001 reader_id=READER-01 payload_size=244

That is not enough. I need to verify the exact raw HTTP request JSON body received by the mock server.

Please update the mock target server handler for POST /reader-bundles.

Requirements:

1. Read the full raw HTTP request body.
2. Print only the first received replay payload to the terminal with this heading:
    First received replay payload:
3. Save every received replay payload into this file:
    logs/received_payloads.jsonl
4. Each received payload should be written as one JSON line.
5. Create the logs folder if it does not exist.
6. Keep existing summary logs if needed.
7. Do not change replay timing.
8. Do not change InjectionTime pacing.
9. Do not change record ordering.
10. Do not change validation logic.
11. Do not change CLI command names.

Expected terminal output should include:

First received replay payload:

Then the actual JSON body received by the mock server.

Expected saved file:
logs/received_payloads.jsonl

Important:

* Do not hardcode sample values.
* Save the actual request body received from Replay Engine.
* The replay payload should match the required top-level structure with ProtoReaderBundle.
* Do not print every payload to terminal because it will be noisy.
* Only print the first payload, but save all received payloads to the jsonl file.
* After changes, tell me which file and function were updated.