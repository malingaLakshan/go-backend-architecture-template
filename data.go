Good plan. Please proceed, but follow these rules strictly:

1. Keep old commands working exactly as they are.
   Do not break:
   rre mock-server -port 8080
   rre validate -file ... -target-url ... -site-id ...
   rre play -file ... -target-url ... -site-id ...

2. Add -config support as an additional option only.

3. If -config is provided, config values can populate:
   recording_file -> file
   target_url -> target-url
   site_id -> site-id
   mock_mode -> mock server config mode
   mock_site_file -> mock site JSON file

4. For mock-server:
   mock_mode = sqlite
   means read SiteInformation.site_json from the configured recording_file and serve it from GET /sites/{siteId}.

5. mock_mode = file
   means read mock_site_file JSON and serve that from GET /sites/{siteId}.

6. Do not bypass validation.
   validate and play must still compare:
   recorded SQLite site_json
   vs
   target site config from target URL.

7. mock-server is only helping QA control what target config is returned.
   Validation logic must stay unchanged.

8. In config mode, validate required fields clearly.
   For validate/play config needs:
   recording_file, target_url, site_id.
   For mock-server sqlite mode config needs:
   recording_file, site_id, mock_mode.
   For mock-server file mode config needs:
   mock_site_file, site_id, mock_mode.

9. Do not silently ignore invalid config.
   Return clear error messages.

10. Please implement small changes only in:
   internal/config/model.go
   internal/cli/args.go
   internal/cli/commands.go
   internal/mocktarget/handler.go
   internal/mocktarget/server.go

11. After changes, run:
   gofmt
   go test ./...
   go build -o rre.exe ./cmd/rre