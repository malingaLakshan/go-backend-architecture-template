We need to convert the Replay Engine MVP into a proper command-line tool.

Current issue:
The project currently has an interactive dashboard flow, but QA and the team need to test it by typing commands directly in the terminal. The tool should behave like a professional CLI tool.

Goal:
Make the main usage command-based, similar to:

* rr -help
* rr --help
* rr help
* rr generate-sample -out data/sample_recording.sqlite
* rr summary -file data/sample_recording.sqlite
* rr mock-server -port 8080
* rr validate -file data/sample_recording.sqlite -target-url http://localhost:8080 -site-id SITE-001
* rr play -file data/sample_recording.sqlite -target-url http://localhost:8080 -site-id SITE-001

Important:
Do not remove the dashboard completely. Keep it as an optional command:

* rr dashboard

But the default and primary behavior must be command-line usage.

Required changes:

1. Update CLI argument parsing so root help works:
    * rr -help
    * rr --help
    * rr help
2. Root help should display:
    * tool name
    * short description
    * available commands
    * example commands
3. Each command should support its own help:
    * rr play -help
    * rr validate -help
    * rr generate-sample -help
    * rr summary -help
    * rr mock-server -help
    * rr dashboard -help
4. The tool should not force the dashboard when no command is provided.
    Instead, show help and return a clear message.
5. Keep existing business logic:
    * sample generation
    * summary
    * mock server
    * validation
    * replay
    * dashboard
6. Do not refactor unrelated packages.
7. Do not change database schema.
8. Do not change replay logic.
9. Do not change validation logic.
10. Only update CLI parsing, help text, and command routing where needed.

Expected command behavior:

rr -help
Should print available commands and examples.

rr generate-sample -out data/sample_recording.sqlite
Should generate a sample SQLite recording file.

rr summary -file data/sample_recording.sqlite
Should print recording summary.

rr mock-server -port 8080
Should start mock target server.

rr validate -file data/sample_recording.sqlite -target-url http://localhost:8080 -site-id SITE-001
Should validate recorded site data against target site data.

rr play -file data/sample_recording.sqlite -target-url http://localhost:8080 -site-id SITE-001
Should replay raw records using InjectionTime pacing.

rr dashboard
Should open the existing interactive dashboard.

Also update README usage section with the new command-based examples.

After changes, show:

* changed files
* final command examples
* how to build Windows executable as rr.exe