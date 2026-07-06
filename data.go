Please fix the GitHub PR security scan issues only.

Context:
This is a Go CLI MVP for the Resonate Replay Engine.
The failing scan is from Cycode SAST.

Detected issues:

1. SSRF risk in internal/replay/injector.go
    because TargetURL is used to build an HTTP request.
2. SQL injection risk in
    internal/recording/repository.go
    because table/column names are used with fmt.Sprintf.

Important:
Do not change the main replay behavior.
Do not change command names.
Do not change output format unless required.
Do not add large new architecture changes.
Keep the fix small, safe, and review-friendly.

Required fix 1:
For injector.go, validate and sanitize the target URL
before creating the POST request.

Allow only http and https schemes.
Reject empty host.
Reject unsupported schemes like file, ftp, gopher, etc.
Build the final /reader-bundles URL using net/url,
not simple string formatting.
Avoid path traversal or malformed URL issues.
Keep the existing 30 second HTTP timeout.

Required fix 2:
For repository.go, remove unsafe dynamic SQL.

For CountTable, do not accept arbitrary table names.
Use a whitelist of allowed table names, for example:
RecordingSession, SiteInformation, RawReads,
ResonateEvents, Snapshots, SnapshotTagLocations,
MLT_SO_Locations.

Only build the SQL query after checking the table name
against the whitelist.

For GetUniqueCount, do not accept arbitrary table or
column names.
Use a whitelist of allowed table + column combinations.
Only allow the combinations actually used by the summary
feature, such as RawReads.TagID and RawReads.ReaderID.

Do not use user-provided values directly in SQL identifiers.

After changes:
Run gofmt.
Run go test ./....
Show me the exact files changed and a short explanation.

Goal:
Make Cycode SAST pass while keeping the replay engine
behavior the same.

⸻

After Copilot changes, check these two files mainly:

internal/replay/injector.go
internal/recording/repository.go

For GitHub comment, after fixing you can write:

Fixed the Cycode SAST findings by validating target URLs before HTTP requests and restricting dynamic SQL identifiers to internal whitelists.