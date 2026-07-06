Copilot Agent, do minimal changes only.

Very important:
Do not change existing function names.
Do not change function signatures.
Do not change return types.
Do not change command behavior.
Do not refactor service architecture.

Only fix Cycode SAST issues with the smallest possible changes.

Issue 1: internal/replay/injector.go
Cycode reports SSRF in HTTP request.

Keep:

func NewInjector(targetURL string) *Injector

Keep:

func (inj *Injector) Send(payload *ProtoReaderBundleWrapper) error

Inside Send, before http.NewRequest, validate and build the endpoint using a helper.

Do not use:

fmt.Sprintf("%s/reader-bundles", inj.TargetURL)

Instead:

* parse inj.TargetURL using net/url
* allow only http and https
* reject empty host
* reject URL with username/password
* build final endpoint using url.URL
* hardcode path as /reader-bundles
* pass only the validated endpoint string to http.NewRequest

Keep the existing behavior and timeout.

Issue 2: internal/recording/repository.go
Cycode reports SQL injection.

Do not change function signatures.

Keep existing functions:

CountTable(...)

GetUniqueCount(...)

But inside them, remove SQL identifier fmt.Sprintf.

For CountTable, use a switch with hardcoded SQL strings.

Example:

case “RawReads”:
query = “SELECT COUNT(*) FROM RawReads”

case “RecordingSession”:
query = “SELECT COUNT(*) FROM RecordingSession”

Return error for unsupported table.

For GetUniqueCount, use hardcoded table + column combinations only.

Example:

case table == “RawReads” && column == “TagID”:
query = “SELECT COUNT(DISTINCT TagID) FROM RawReads WHERE RecordingSessionID = ?”

case table == “RawReads” && column == “ReaderID”:
query = “SELECT COUNT(DISTINCT ReaderID) FROM RawReads WHERE RecordingSessionID = ?”

Return error for unsupported combination.

No fmt.Sprintf for SQL identifiers.

After changes:
Run gofmt.
Run go test ./....

Show only changed lines and confirm function signatures were not changed.