## Summary
Added versioning support to the Resonate Replay Engine CLI.

## Changes
- Added `rre version` command.
- Added `rre --version` and `rre -v` support.
- Added default local version as `dev`.
- Added build-time version injection using Go ldflags.
- Updated help/README with versioning usage and release build command.

## Build Example
```powershell
go build -ldflags="-X main.version=v0.1.0" -o rre.exe ./cmd/rre