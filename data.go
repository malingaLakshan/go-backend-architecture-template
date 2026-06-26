The previous change did not produce the expected output.

Problem:
After running replay, I still cannot see the actual received JSON payload in terminal or in a separate payload file. I only see normal logs/summary.

Please fix this specifically in the mock target server POST /reader-bundles handler.

Required behavior:

1. When mock server receives POST /reader-bundles, read the raw HTTP request body using io.ReadAll(r.Body).
2. Save the exact raw request body into:
    logs/received_payloads.jsonl
3. Each HTTP request body must be written as one line in that file.
4. Create the logs folder if it does not exist.
5. Print only the first received raw request body to terminal with this heading:
    First received replay payload:
6. The terminal output must show the real JSON body, not only site_id, reader_id, or payload_size.
7. Keep the existing summary log also if needed.

Important:

* Do not write normal logger messages into received_payloads.jsonl.
* received_payloads.jsonl must contain only replay HTTP request JSON bodies.
* Do not change InjectionTime pacing.
* Do not change replay timing.
* Do not change record ordering.
* Do not change validation logic.
* Do not change CLI command names.

Implementation hint:
In the POST /reader-bundles handler, do something similar to:

* bodyBytes, err := io.ReadAll(r.Body)
* save bodyBytes to logs/received_payloads.jsonl
* print the first payload to terminal
* then unmarshal/parse the same bodyBytes if the handler still needs to extract site_id or reader_id

After changes, show me:

1. The file/function changed
2. The code section that reads r.Body
3. The code section that writes to logs/received_payloads.jsonl
4. The terminal output showing First received replay payload: