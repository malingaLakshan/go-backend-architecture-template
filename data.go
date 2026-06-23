Please add an interactive CLI dashboard/menu to the existing Replay Engine CLI MVP.

Do not rewrite the full project. Keep all existing commands working.

Project folder:
replay-engine-mvp-cli

Required new command:
go run ./cmd/rre dashboard

Also, if possible, when user runs:
go run ./cmd/rre

show the dashboard by default.

Files to change/add:

1. cmd/rre/main.go

* Register or route the new dashboard command.
* If no command is provided, open dashboard mode.

2. internal/cli/commands.go

* Add dashboard command handling.
* Keep existing commands unchanged:
    * generate-sample
    * summary
    * mock-server
    * validate
    * play

3. internal/cli/dashboard.go

* Create this new file.
* Implement interactive terminal menu.

4. internal/cli/args.go

* Add default dashboard values if needed:
    * recording file: data/sample_recording.sqlite
    * target URL: http://localhost:9090
    * site ID: SITE-001
    * mock server port: 9090

5. internal/mocktarget/server.go

* Only change this if needed to allow mock server to start from dashboard mode in a goroutine.
* Add safe handling so dashboard does not start multiple servers on the same port.

Do not change these unless absolutely necessary:

* internal/replay/service.go
* internal/replay/pacing.go
* internal/replay/injector.go
* internal/recording/repository.go
* internal/site/validator.go
* internal/logger/logger.go

Dashboard menu:

Replay Engine CLI MVP

1. Generate sample recording
2. Show recording summary
3. Start mock target server
4. Validate site
5. Play replay
6. Run full demo flow
7. Exit

Expected behavior:

1. Generate sample recording

* Use default file path: data/sample_recording.sqlite
* Generate 50 sample raw records.
* Print success message.

2. Show recording summary

* Read default SQLite file.
* Print recording details and estimated replay duration.

3. Start mock target server

* Start server on port 9090.
* In dashboard mode, start it in background using goroutine.
* If already running, show “Mock server already running”.

4. Validate site

* Validate SITE-001 against http://localhost:9090.
* Show pass/fail clearly.

5. Play replay

* Run validation first.
* If validation fails, do not replay.
* If validation passes, replay records using InjectionTime pacing.
* Show terminal progress and final summary.

6. Run full demo flow

* Generate sample recording.
* Show summary.
* Start mock server if not running.
* Validate SITE-001.
* Play replay.
* Print final summary.

7. Exit

* Exit cleanly.

Important rules:

* Use simple terminal input/output.
* Prefer Go standard library only.
* Do not add heavy CLI UI libraries.
* Keep existing CLI commands working.
* Do not commit generated data/sample_recording.sqlite or logs/rre.log.
* Do not rewrite unrelated files.
* Reuse existing services/functions as much as possible.
* First explain your planned changes file by file, then provide code changes.