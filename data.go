Please analyze the current replay-engine-mvp-cli Go project and create/update architecture documentation automatically.

Do not change application logic.

Target file:
replay-engine-mvp-cli/docs/ARCHITECTURE.md

If the file does not exist, create it.
If it already exists, update it cleanly.

Project context:
This is a Go CLI MVP for a Replay Engine.

Current structure:

* cmd/rre/main.go = application entry point
* internal/cli = command parsing, dashboard menu, command orchestration
* internal/sqlite = SQLite connection, schema, sample data generation
* internal/recording = recording summary, site information, raw reads from SQLite
* internal/site = target site fetching and site validation
* internal/replay = replay service, InjectionTime pacing, status, payload injection
* internal/mocktarget = local mock target HTTP server
* internal/logger = console and file logging
* data/sample_recording.sqlite = generated sample SQLite recording file
* logs/rre.log = generated log file

Need:
Add a professional high-level architecture document with Mermaid diagrams.

In docs/ARCHITECTURE.md, include:

1. Title:
    Replay Engine CLI MVP - Architecture
2. Short overview:
    Explain that this tool reads recorded RFID/location data from SQLite and replays it to a target system using original InjectionTime pacing.
3. Diagram 1:
    High-level component architecture Mermaid diagram.

It should include:

* User
* CLI Dashboard / Commands
* SQLite Recording File
* Recording Reader
* Site Validator
* Replay Engine
* Pacing Logic
* HTTP Injector
* Mock Target Server
* Logger

4. Diagram 2:
    Replay flow sequence Mermaid diagram.

It should show:

* User starts dashboard/play
* CLI opens SQLite recording
* Reads SiteInformation
* Fetches target site configuration
* Validates recorded site vs target site
* Reads RawReads ordered by InjectionTime
* Calculates pacing delay
* Sends payload to mock target server
* Logs progress
* Completes or aborts safely

5. Add short section:
    Package responsibilities

Explain each package shortly:

* cmd/rre
* internal/cli
* internal/sqlite
* internal/recording
* internal/site
* internal/replay
* internal/mocktarget
* internal/logger

6. Add short section:
    Replay rules

Mention:

* Validation must pass before replay starts.
* Raw records must be ordered by InjectionTime.
* First record can be sent immediately.
* Next records must wait based on InjectionTime gap.
* Single injection failure should be logged and replay should continue.
* Ctrl+C should abort safely.
* Generated SQLite/log files should not be committed.

Important:

* Use Mermaid code blocks inside the markdown file.
* Keep the document clean, professional, and easy to understand.
* Do not include confidential company data.
* Use generic names only.
* Do not modify Go source code.
* After updating the file, show me the final file path and a short summary of what was added.