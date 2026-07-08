Please update the README.md for the new RRE config-based QA testing support.

Important:
Do not change any Go code.
Only update documentation files, mainly README.md.
Keep the wording simple and clear.
Do not remove existing usage unless it is wrong.
Keep old CLI flag usage documented because old commands still work.

Project context:
RRE CLI can:
- generate sample recording
- show summary
- validate recorded site config against target site config
- replay RawReads using InjectionTime-based pacing
- start mock-server as a target Resonate API

Recently added:
Config-based command support for:
- mock-server
- validate
- play

New config file format:

{
  "recording_file": "data/recording_001.sqlite",
  "target_url": "http://localhost:8080",
  "site_id": "b3489888-aacf-4451-893c-d7d994240f93",
  "mock_mode": "sqlite",
  "mock_site_file": ""
}

Fields:
- recording_file: path to Recorder SQLite file
- target_url: target Resonate/mock target base URL
- site_id: site ID to validate/replay
- mock_mode: only used by mock-server config mode
  - sqlite: mock-server loads SiteInformation.site_json from recording_file
  - file: mock-server loads target site config from mock_site_file
- mock_site_file: path to JSON site config file, only required when mock_mode is file

Document old commands:

rre mock-server -port 8080
rre validate -file data/recording_001.sqlite -target-url http://localhost:8080 -site-id <site-id>
rre play -file data/recording_001.sqlite -target-url http://localhost:8080 -site-id <site-id>

Document new config commands:

rre mock-server -config configs/pass_config.json
rre validate -config configs/pass_config.json
rre play -config configs/pass_config.json

Add QA testing section:

1. Validation pass scenario:
Create configs/pass_config.json:

{
  "recording_file": "data/recording_001.sqlite",
  "target_url": "http://localhost:8080",
  "site_id": "b3489888-aacf-4451-893c-d7d994240f93",
  "mock_mode": "sqlite",
  "mock_site_file": ""
}

Terminal 1:
rre mock-server -config configs/pass_config.json

Terminal 2:
rre validate -config configs/pass_config.json

Expected:
Validation passed: recorded site matches target site

Optional:
rre play -config configs/pass_config.json

2. Validation fail scenario:
Create configs/wrong_site_config.json with mismatched site details, for example:

{
  "id": "WRONG-SITE-ID",
  "name": "Wrong Site",
  "readers": [],
  "floors": [],
  "regions": [],
  "antennas": []
}

Create configs/fail_config.json:

{
  "recording_file": "data/recording_001.sqlite",
  "target_url": "http://localhost:8080",
  "site_id": "b3489888-aacf-4451-893c-d7d994240f93",
  "mock_mode": "file",
  "mock_site_file": "configs/wrong_site_config.json"
}

Terminal 1:
rre mock-server -config configs/fail_config.json

Terminal 2:
rre validate -config configs/fail_config.json

Expected:
Validation failed

Explain clearly:
- Validation is not removed or bypassed.
- validate/play still compare recorded SQLite SiteInformation.site_json against target config fetched from target_url /sites/{siteId}.
- mock-server config mode only controls what target config is returned for QA testing.
- sqlite mode is used to create matching target config and test pass case.
- file mode is used to provide wrong target config and test fail case.

Add troubleshooting:
- If fail_config still passes, make sure mock-server was restarted with fail_config.
- Correct command is:
  rre mock-server -config configs/fail_config.json
  not:
  rre mock-server --configs/fail_config.json
- Stop old mock-server before switching configs.
- Check target config manually:
  PowerShell:
  Invoke-RestMethod http://localhost:8080/sites/<site-id> | ConvertTo-Json -Depth 10

Mention not to commit real Recorder SQLite files or sensitive site data unless approved by the team.

Please update README.md only.