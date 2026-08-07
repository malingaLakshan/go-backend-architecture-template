Fix the UI Play output with the smallest possible change.

Inspect the current UI Play handler and the existing MVP `runPlay` implementation first.

Requirements:
1. Make the UI Play action reuse the existing MVP play execution path and summary formatting.
2. Do not create another replay function or duplicate replay logic.
3. Preserve all existing playback behavior: validation, RawReads loading, InjectionTime pacing, HTTP sending, cancellation, JSONL output, success/failure counting, and logging.
4. On successful playback, the visible UI terminal must show exactly the useful MVP output:
   - `First replay output payload:` once with formatted JSON
   - `Replay Summary`
   - Total Records
   - Successful
   - Failed
   - Duration
   - Status
5. Do not display timestamped `[INFO]`, `[DEBUG]`, or other logger lines in the UI terminal. They must continue going only to the existing common log file.
6. Remove UI-only lines such as:
   - `Playback started: session=...`
   - `[INFO] Playback started...`
   - `[INFO] Playback finished...`
   - `Playback session ... finished: ...`
7. Keep concise user-facing error messages in the UI when playback fails; detailed errors remain in the log file.
8. Do not change unrelated backend logic, CLI behaviour, logging paths, validation, or pacing.
9. Update only the minimum required tests and run all tests.

If `runPlay` is only a CLI adapter, reuse its existing underlying play service and existing MVP output/summary formatter. Do not implement a new playback flow.

Before editing, briefly identify why the UI currently bypasses the MVP summary. Then implement and report only changed files and test results.