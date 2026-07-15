Fix only the current SiteGraph parsing and validation bug. Do not redesign unrelated parts of the project.

Confirmed real structure:

Site
  -> Floors
      -> Regions
          -> nested Regions recursively at any depth
      -> Readers
          -> Antennas

Important:

- Readers are directly inside each floor.
- Readers are not inside regions.
- Regions can contain nested regions recursively.
- A region can contain multiple child regions.
- Antennas are inside readers.
- Antennas are identified by integer field `port`.

Current wrong output:

Regions: 1
Readers: 0
Antenna Ports: 0

But the real JSON contains many nested regions, floor-level readers and antennas.

Inspect only the relevant files, likely:

internal/site/model.go
internal/site/parser.go
internal/site/summary.go
internal/site/validator.go
internal/site/*_test.go
internal/mocktarget/site_store.go
internal/mocktarget/server.go

Fix the normalized model to follow:

type ValidationSite struct {
    SiteID string
    Floors []ValidationFloor
}

type ValidationFloor struct {
    ID      string
    Regions []ValidationRegion
    Readers []ValidationReader
}

type ValidationRegion struct {
    ID      string
    Regions []ValidationRegion
}

type ValidationReader struct {
    ID           string
    AntennaPorts []int
}

Requirements:

1. Parse all floors.
2. Recursively parse every nested region at any depth.
3. Preserve each region’s parent-child relationship.
4. Parse readers from floor.readers.
5. Parse antennas from reader.antennas.
6. Use antenna.port.
7. Count all nested regions recursively.
8. Count all floor readers.
9. Count all antenna ports.
10. Validate:
   - floor under site
   - top-level region under correct floor
   - child region under correct parent region
   - reader under correct floor
   - antenna port under correct reader
11. Do not flatten region relationships in a way that allows false matches.
12. Update only the necessary tests.
13. Do not change config flow, replay flow, versioning, README, help, logging, Git files or unrelated code.
14. Do not commit anything.

After changes run:

gofmt -w .
go test ./...
go vet ./...
go build -o rre.exe ./cmd/rre

Then report only:
- files changed
- tests run
- any remaining compile/test issue