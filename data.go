## Summary

Implements Replay Engine site validation against full Recorder and target SiteGraphs.

## Changes

- Loads multiple full SiteGraph JSON files from the configured sites directory
- Adds GET /sites and GET /sites/{siteId} mock endpoints
- Preserves complete original SiteGraph JSON responses
- Parses recursive Region hierarchies
- Parses Readers directly under Floors
- Parses antenna ports under Readers
- Validates recorded structures against the target hierarchy
- Allows additional structures in the target
- Prevents playback when validation fails
- Updates CLI help and documentation
- Aligns tests with the real SiteGraph hierarchy

## Validation rules

- Site ID must match
- Recorded Floors must exist in the target
- Regions must exist under the correct Floor and parent Region
- Readers must exist under the correct Floor
- Antenna ports must exist under the correct Reader

## Verification

- `go test ./...`
- `go vet ./...`
- `go build -o rre.exe ./cmd/rre`
- Manual mock-server test
- Manual GET /sites test
- Manual validation test
- Manual playback test

## Security

- Keeps parameterized, hardcoded SQL
- Keeps read-only SQLite access
- Keeps existing SSRF protections
- Keeps approved SiteGraph directory/path validation
- Does not include private SQLite recordings or SiteGraph files