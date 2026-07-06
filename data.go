Copilot Agent, fix the remaining Cycode SAST failures exactly at:

internal/recording/repository.go line 181
internal/replay/injector.go line 92

Do not make large changes. Keep behavior same.

Issue 1: repository.go SQL injection

Cycode still detects unsafe SQL because the code still builds SQL using dynamic table/column values.

Fix this by removing generic dynamic SQL for the scanned functions.

For CountTable, do not use:

fmt.Sprintf("SELECT COUNT(*) FROM %s", table)

Instead, use a switch statement and hardcoded SQL per allowed table.

Example style:

case "RawReads": query = "SELECT COUNT(*) FROM RawReads"

Do the same for all allowed internal tables.

For GetUniqueCount, do not build SQL using dynamic table or column names.

Use hardcoded switch cases for allowed combinations only.

Example:

case table == "RawReads" && column == "TagID": query = "SELECT COUNT(DISTINCT TagID) FROM RawReads WHERE RecordingSessionID = ?"

case table == "RawReads" && column == "ReaderID": query = "SELECT COUNT(DISTINCT ReaderID) FROM RawReads WHERE RecordingSessionID = ?"

Return an error for unsupported table/column.

Goal: no fmt.Sprintf for SQL identifiers in these functions.

Issue 2: injector.go SSRF

Cycode still detects SSRF because request URL is created from TargetURL.

Fix by adding strict target URL validation before http.NewRequest.

Create a helper like:

buildReaderBundlesURL(targetURL string) (string, error)

Inside it:

* parse using url.ParseRequestURI or url.Parse
* allow only http and https
* reject empty host
* reject URLs with username/password
* reject unsupported schemes
* cleanly join path with /reader-bundles
* do not use fmt.Sprintf("%s/reader-bundles", targetURL)

Use url.URL to build the final URL.

Also add tests for:

* valid http://localhost:8080
* valid https://example.com/api
* invalid file:///etc/passwd
* invalid ftp://example.com
* invalid empty target URL
* invalid URL without host

After changes:

Run:

gofmt

go test ./...

Then show changed files and explain the fix.

Important:
Do not ignore or mark false positive.
Actually change code so Cycode SAST can pass.

⸻

Also, after Copilot fixes it, check there is no SQL line like this:

fmt.Sprintf("SELECT ... %s ...", table)

and no URL line like this:

fmt.Sprintf("%s/reader-bundles", inj.TargetURL)