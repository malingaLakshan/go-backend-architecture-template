There is an incomplete implementation.

validate -config works, but mock-server -config does not work.

Current error:
flag provided but not defined: -config
Usage of mock-server:
  -port int

Please fix this only.

Requirements:
1. Add -config flag to the mock-server FlagSet in internal/cli/args.go.
2. Ensure Flags struct has Config string field if not already.
3. In runMockServer, if flags.Config is provided:
   - load config JSON
   - if mock_mode == "sqlite":
       read SiteInformation.site_json from recording_file
       unmarshal into site.SiteConfig
       start mock server with that config
   - if mock_mode == "file":
       read mock_site_file
       unmarshal into site.SiteConfig
       start mock server with that config
4. If mock_mode is invalid, return clear error.
5. Keep old command working:
   rre mock-server -port 8080
6. Do not change validate logic.
7. Do not change play logic.
8. Do not remove validation.
9. Add startup log:
   Mock target site config source: sqlite <recording_file>
   or
   Mock target site config source: file <mock_site_file>

Only touch:
internal/cli/args.go
internal/cli/commands.go
internal/mocktarget/server.go
internal/mocktarget/handler.go if needed.