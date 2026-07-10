Now real Recorder SQLite schema support is restored. Next, fix the full config-based RRE CLI flow and validation behavior.

Context:
Project folder is resonate-replay-engine.
Executable is rre.exe.
CLI name in help should be rre.
Real Recorder SQLite schema is already supported using snake_case columns.
Do not revert repository/model/schema back to old camelCase sample schema.

Please fix all related files carefully.

1. CLI config support

Update internal/cli/args.go.

Flags struct must include:

type Flags struct {
    Command   string
    File      string
    Out       string
    TargetURL string
    SiteID    string
    Port      int
    Config    string
}

Add -config flag support for:
- summary
- validate
- play
- mock-server

Direct flag commands must still work:

./rre.exe summary -file data/recording_001.sqlite

./rre.exe validate -file data/recording_001.sqlite -target-url http://localhost:8080 -site-id b3489888-aacf-4451-893c-d7d994240f93

./rre.exe play -file data/recording_001.sqlite -target-url http://localhost:8080 -site-id b3489888-aacf-4451-893c-d7d994240f93

./rre.exe mock-server -port 8080

Config commands must work:

./rre.exe summary -config configs/pass_config.json
./rre.exe validate -config configs/pass_config.json
./rre.exe play -config configs/pass_config.json
./rre.exe mock-server -config configs/pass_config.json

./rre.exe mock-server -config configs/fail_config.json
./rre.exe validate -config configs/fail_config.json

-config must appear in command help:
./rre.exe summary -help
./rre.exe validate -help
./rre.exe play -help
./rre.exe mock-server -help

2. Config model

Create or fix internal/config/model.go.

Use this config structure:

type RunConfig struct {
    RecordingFile string `json:"recording_file"`
    TargetURL     string `json:"target_url"`
    SiteID        string `json:"site_id"`
    MockMode      string `json:"mock_mode"`
    MockSiteFile  string `json:"mock_site_file"`
}

Load(path string) should:
- read the JSON file
- unmarshal into RunConfig
- return helpful errors

3. summary with config

Update runSummary.

Direct:
./rre.exe summary -file data/recording_001.sqlite

Config:
./rre.exe summary -config configs/pass_config.json

Behavior:
- If -config is provided, load config file.
- Use cfg.RecordingFile as flags.File.
- summary does not need target_url, site_id, mock_mode, or mock_site_file.
- If both -file and -config are missing, return:
  [ERROR] -file or -config is required for summary

Do not require site_id for summary.

4. validate with config

Update runValidate.

Direct:
./rre.exe validate -file data/recording_001.sqlite -target-url http://localhost:8080 -site-id b3489888-aacf-4451-893c-d7d994240f93

Config:
./rre.exe validate -config configs/pass_config.json

Behavior:
- If -config is provided, load config file.
- Set flags.File = cfg.RecordingFile.
- Set flags.TargetURL = cfg.TargetURL.
- Set flags.SiteID = cfg.SiteID.
- If required values are missing, return:
  [ERROR] -file, -target-url, and -site-id are required for validate

Validate flow:
- Open SQLite recording file.
- Load recorded site config from SiteInformation.site_json using site_id.
- json.Unmarshal(siteInfo.SiteJSON, &recordedConfig)
- Fetch target site config from GET /sites/{siteId}.
- Validate recorded config vs target config.
- Print recorded site config summary.
- Print target site config summary.
- If validation fails, print all mismatch errors and return exit code 1.
- If validation passes, return exit code 0.

Do not use RawSiteJSON anywhere.
Error text should say SiteJSON.

5. play with config

Update runPlay.

Direct:
./rre.exe play -file data/recording_001.sqlite -target-url http://localhost:8080 -site-id b3489888-aacf-4451-893c-d7d994240f93

Config:
./rre.exe play -config configs/pass_config.json

Behavior:
- If -config is provided, load config file.
- Set flags.File = cfg.RecordingFile.
- Set flags.TargetURL = cfg.TargetURL.
- Set flags.SiteID = cfg.SiteID.
- If required values are missing, return:
  [ERROR] -file, -target-url, and -site-id are required for play

Play flow:
- Validate recorded site config against target site config first.
- Abort replay if validation fails.
- Load RawReads from real Recorder SQLite schema.
- Replay RawReads to POST /reader-bundles.
- Use injection_time_utc pacing.
- Write sender-side output to logs/replay_output.jsonl.

6. mock-server with config

Update runMockServer and runMockServerWithConfig.

Direct:
./rre.exe mock-server -port 8080

Config:
./rre.exe mock-server -config configs/pass_config.json
./rre.exe mock-server -config configs/fail_config.json

Behavior:
- If no -config, start default mock server on selected port.
- If -config exists, load config and use mock_mode.

mock_mode = sqlite:
- Requires recording_file and site_id.
- Open recording_file.
- Load SiteInformation.site_json by site_id.
- Unmarshal into site.SiteConfig.
- Start mock server using this exact SiteConfig.
- This is the pass test mode.

mock_mode = file:
- Requires mock_site_file.
- Read mock_site_file.
- Unmarshal into site.SiteConfig.
- Start mock server using this SiteConfig.
- This is the fail test mode.

Errors:
- If mock_mode is not sqlite or file:
  [ERROR] config mock_mode must be "sqlite" or "file"
- If sqlite mode and recording_file missing:
  [ERROR] config missing required field for sqlite mode: recording_file
- If sqlite mode and site_id missing:
  [ERROR] config missing required field for sqlite mode: site_id
- If file mode and mock_site_file missing:
  [ERROR] config missing required field for file mode: mock_site_file
- If mock_site_file JSON is invalid, print clear JSON parse error.

mock-server startup output should show:
- loaded config path
- mock mode
- recording file or mock site file
- site config source
- site ID
- site name
- reader count
- antenna count
- floor count
- region count
- listening URL
- endpoints
- received payload log path

7. mocktarget server/handler

Update internal/mocktarget/server.go and handler.go.

Support:

func StartServer(port int) error
func StartServerWithConfig(port int, cfg *site.SiteConfig) error

Handler support:
func NewHandler(paths ...string) (*Handler, error)
func NewHandlerWithConfig(cfg *site.SiteConfig, paths ...string) (*Handler, error)

GET /sites/{siteId}:
- If handler has siteConfig:
  - If requested siteId equals siteConfig.SiteID, return that config as JSON.
  - If requested siteId does not match, return 404.
- If handler has no siteConfig:
  - Return default fallback config only for the fallback test site.

POST /reader-bundles:
- Read request body.
- Append each payload line to logs/received_payloads.jsonl.
- Print first payload as pretty JSON.
- Print compact receive line:
  [INFO] Received #N | site_id=... | reader_id=... | reads=... | size=... bytes
- Return:
  {"status":"accepted"}

8. Config files

Create or restore these files:

configs/pass_config.json
configs/fail_config.json
configs/wrong_site_config.json

pass_config.json:

{
  "recording_file": "data/recording_001.sqlite",
  "target_url": "http://localhost:8080",
  "site_id": "b3489888-aacf-4451-893c-d7d994240f93",
  "mock_mode": "sqlite",
  "mock_site_file": ""
}

fail_config.json:

{
  "recording_file": "data/recording_001.sqlite",
  "target_url": "http://localhost:8080",
  "site_id": "b3489888-aacf-4451-893c-d7d994240f93",
  "mock_mode": "file",
  "mock_site_file": "configs/wrong_site_config.json"
}

wrong_site_config.json must be valid JSON and must match SiteConfig structure.
Do not use string values for readers/floors/regions/antennas.
Use arrays of objects.

Use this:

{
  "id": "WRONG-SITE-ID",
  "name": "Wrong Site",
  "readers": [
    {
      "id": "WRONG-READER-01",
      "name": "Wrong Reader",
      "type": "RFID",
      "ipAddress": "192.168.1.200",
      "floorId": "WRONG-FLOOR-01",
      "x": 1,
      "y": 1,
      "antennas": [
        {
          "antenna_id": 99,
          "antenna_type": 2,
          "reader_id": "WRONG-READER-01",
          "x": 1,
          "y": 1
        }
      ]
    }
  ],
  "floors": [
    {
      "id": "WRONG-FLOOR-01",
      "name": "Wrong Floor",
      "number": 1,
      "width": 100,
      "height": 100,
      "regions": [
        {
          "id": "WRONG-REGION-01",
          "name": "Wrong Region",
          "type": "WRONG_TYPE",
          "physicality": "VIRTUAL",
          "inventoryType": "OTHER"
        }
      ]
    }
  ],
  "regions": [
    {
      "id": "WRONG-REGION-01",
      "name": "Wrong Region",
      "type": "WRONG_TYPE",
      "physicality": "VIRTUAL",
      "inventoryType": "OTHER"
    }
  ],
  "antennas": [
    {
      "antenna_id": 99,
      "antenna_type": 2,
      "reader_id": "WRONG-READER-01",
      "x": 1,
      "y": 1
    }
  ]
}

9. Validator must check IDs, not only counts

Update internal/site/validator.go.

Validation must fail if IDs do not match, even if counts match.

Check:
- site ID equality
- reader IDs equality
- antenna IDs equality
- floor IDs equality
- region IDs equality

Rules:
- If recorded site ID != target site ID, fail.
- Every recorded reader ID must exist in target.
- Every target reader ID must exist in recorded.
- Every recorded antenna ID must exist in target.
- Every target antenna ID must exist in recorded.
- Every recorded floor ID must exist in target.
- Every target floor ID must exist in recorded.
- Every recorded region ID must exist in target.
- Every target region ID must exist in recorded.

Error examples:
- Site ID mismatch: recorded=..., target=...
- Reader ID missing in target: ...
- Reader ID missing in recorded: ...
- Antenna ID missing in target: ...
- Antenna ID missing in recorded: ...
- Floor ID missing in target: ...
- Floor ID missing in recorded: ...
- Region ID missing in target: ...
- Region ID missing in recorded: ...

If a config has zero readers/antennas but recorded config also has zero, that section can pass. But if IDs differ, it must fail.

10. SiteConfig model

Do not break real site_json unmarshalling.

Ensure internal/site/model.go supports:

type SiteConfig struct {
    SiteID   string    `json:"id"`
    SiteName string    `json:"name"`
    Readers  []Reader  `json:"readers,omitempty"`
    Floors   []Floor   `json:"floors,omitempty"`
    Regions  []Region  `json:"regions,omitempty"`
    Antennas []Antenna `json:"antennas,omitempty"`
}

Reader:
- id
- name
- type
- ipAddress
- floorId
- x
- y
- antennas

Antenna:
- antenna_id
- antenna_type
- reader_id
- x
- y

Floor:
- id
- name
- number
- width
- height
- regions

Region:
- id
- name
- type
- physicality
- inventoryType

Go should ignore extra JSON fields from real Recorder site_json.

11. .gitignore

Do not commit generated runtime files.

Update .gitignore if needed:

rre.exe
logs/*.jsonl
*.sqlite-shm
*.sqlite-wal
*.bin
*.tmp
*.backup

Do not delete the logs folder itself if needed; keep .gitkeep if required.

12. README and help text

Update README and internal CLI help.

Include direct usage:

./rre.exe summary -file data/recording_001.sqlite
./rre.exe mock-server -port 8080
./rre.exe validate -file data/recording_001.sqlite -target-url http://localhost:8080 -site-id b3489888-aacf-4451-893c-d7d994240f93
./rre.exe play -file data/recording_001.sqlite -target-url http://localhost:8080 -site-id b3489888-aacf-4451-893c-d7d994240f93

Include config usage:

./rre.exe summary -config configs/pass_config.json
./rre.exe mock-server -config configs/pass_config.json
./rre.exe validate -config configs/pass_config.json
./rre.exe play -config configs/pass_config.json

Fail validation demo:

Terminal 1:
./rre.exe mock-server -config configs/fail_config.json

Terminal 2:
./rre.exe validate -config configs/fail_config.json

Expected:
Validation fails with mismatch errors.

Pass validation demo:

Terminal 1:
./rre.exe mock-server -config configs/pass_config.json

Terminal 2:
./rre.exe validate -config configs/pass_config.json

Expected:
Validation passes.

13. Tests

Update or add tests for:
- summary direct file
- summary config file
- validate pass config
- validate fail config
- mock-server sqlite config mode
- mock-server file config mode
- validator fails when counts match but IDs differ
- wrong_site_config.json unmarshals into SiteConfig

14. Build/test

After changes:
- gofmt all changed Go files
- go test ./...
- go build -o rre.exe ./cmd/rre

Manual check commands:

./rre.exe summary -help
./rre.exe validate -help
./rre.exe play -help
./rre.exe mock-server -help

./rre.exe summary -config configs/pass_config.json

Terminal 1:
./rre.exe mock-server -config configs/pass_config.json

Terminal 2:
./rre.exe validate -config configs/pass_config.json

Then stop server.

Terminal 1:
./rre.exe mock-server -config configs/fail_config.json

Terminal 2:
./rre.exe validate -config configs/fail_config.json

Then stop server.

Terminal 1:
./rre.exe mock-server -config configs/pass_config.json

Terminal 2:
./rre.exe play -config configs/pass_config.json

Important:
Do not revert real Recorder SQLite repository/model/schema support.
Do not bring back old camelCase SQL columns.
Do not use RawSiteJSON.
Do not create temp duplicate files.
Do not commit logs or generated files.