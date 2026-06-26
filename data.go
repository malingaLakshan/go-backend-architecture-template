We need to update the Replay Engine CLI identity and help output.

Context:
This project is the Resonate Replay Engine, so the command name must be rre, not rr. rr can be confused with Recorder-related tooling.

Scope:
Only update CLI argument parsing/help text, README usage examples, and documentation examples.
Do not change replay business logic, validation logic, SQLite logic, pacing logic, or mock server logic.

Required command name:

* Tool name: rre
* Windows executable name: rre.exe
* Keep Go entry folder as cmd/rre

Help behavior:
The following should show root help:

* go run ./cmd/rre help
* go run ./cmd/rre -help
* go run ./cmd/rre --help
* .\rre.exe help
* .\rre.exe -help
* .\rre.exe --help

Root help must clearly show:

Tool name:
rre - Resonate Replay Engine

Description:
Replay recorded RFID/location reader data from a SQLite recording file into a target Resonate HTTP instance.

Usage:
rre <command> [flags]

Windows PowerShell note:
When running a local executable in Windows PowerShell, use .\rre.exe instead of rre unless the executable folder is added to PATH.

Available commands:

* help - Show help information
* generate-sample - Generate a sample SQLite recording file
* summary - Show recording summary
* mock-server - Start mock target server
* validate - Validate recorded site configuration against target site
* play - Replay recorded raw reads to target Resonate instance
* dashboard - Open optional interactive dashboard

Examples must use rre, not rr.

Development examples:

* go run ./cmd/rre help
* go run ./cmd/rre generate-sample -out data/sample_recording.sqlite
* go run ./cmd/rre summary -file data/sample_recording.sqlite
* go run ./cmd/rre mock-server -port 8080
* go run ./cmd/rre validate -file data/sample_recording.sqlite -target-url http://localhost:8080 -site-id SITE-001
* go run ./cmd/rre play -file data/sample_recording.sqlite -target-url http://localhost:8080 -site-id SITE-001

Windows executable examples:

* .\rre.exe help
* .\rre.exe generate-sample -out data/sample_recording.sqlite
* .\rre.exe summary -file data/sample_recording.sqlite
* .\rre.exe mock-server -port 8080
* .\rre.exe validate -file data/sample_recording.sqlite -target-url http://localhost:8080 -site-id SITE-001
* .\rre.exe play -file data/sample_recording.sqlite -target-url http://localhost:8080 -site-id SITE-001

Build command:
go build -o rre.exe ./cmd/rre

Rules:

* Do not remove dashboard, but keep it as optional command only.
* When no command is provided, show root help clearly.
* Remove old rr command examples from help/README/docs.
* After changes, show changed files and final help output.