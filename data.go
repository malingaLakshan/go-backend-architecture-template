I need to add versioning support to the Resonate Replay Engine CLI (`rre`) in this Go project.

Context:
- This is inside a monorepo: `ALTRFIDTools/resonate-replay-engine`
- Binary name is `rre`
- Current CLI has commands like:
  - help
  - generate-sample
  - summary
  - mock-server
  - validate
  - play
- I need to add version support based on the project versioning guide.
- Version should be injected at build time using Go ldflags.
- Local/default version should be `dev`.

Required behavior:
1. Add support for:
   - `rre version`
   - `rre --version`
   - `rre -v`

2. Output should be:
   `rre version <version>`

   Example:
   `rre version v0.1.0`

3. Default local build output should be:
   `rre version dev`

4. Build command with injected version should work:
   `go build -ldflags="-X main.version=v0.1.0" -o rre.exe ./cmd/rre`

5. Do not break existing commands:
   - `rre help`
   - `rre summary`
   - `rre validate`
   - `rre play`
   - `rre mock-server`
   - config based commands must still work

Implementation preference:
- Add `var version = "dev"` in `cmd/rre/main.go` package main.
- In `main.go`, handle root-level version flags before normal CLI command parsing:
  - if first arg is `version`, `--version`, or `-v`, print version and exit.
- Also add `version` as a normal command in CLI help if the existing dispatcher requires it.
- Update `internal/cli/args.go` help text to include:
  - `version          Show Replay Engine version`
  - examples for `rre version`, `rre --version`, and `rre -v`
  - build example using ldflags.

Important:
- Do not modify unrelated modules like `resonate-recorder` or `resonate-analyzer`.
- Do not add generated files to Git:
  - `rre.exe`
  - `logs/*.jsonl`
  - `data/*.sqlite-shm`
  - `data/*.sqlite-wal`
  - `data/*.bin`
- Do not reintroduce Cycode violations.
- Do not use user-controlled paths directly in `os.ReadFile` or HTTP request calls.
- Keep changes small and focused only on versioning.

After implementation, run or make sure these commands pass:
- `gofmt -w cmd/rre/main.go internal/cli/args.go`
- `go test ./...`
- `go build -ldflags="-X main.version=v0.1.0" -o rre.exe ./cmd/rre`
- `.\rre.exe version`
- `.\rre.exe --version`
- `.\rre.exe -v`
- `.\rre.exe help`

Expected version command result:
`rre version v0.1.0`

Please implement this cleanly and show me the changed files.