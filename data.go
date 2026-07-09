Please improve terminal output for config-based QA testing.

Important:
Do not change validation logic.
Do not change replay logic.
Do not remove any existing functionality.
Do not print full sensitive site_json by default.
Only add clear terminal summaries.

Problem:
In pass_config.json, mock_site_file is empty because mock_mode is sqlite.
This is correct, but it is not clear in terminal where readers, antennas, floors, and regions come from.

Requirement:
When running mock-server with -config, print a clear config summary.

Example command:
rre mock-server -config configs/pass_config.json

Expected terminal output for sqlite mode:

[INFO] Loaded run config: configs/pass_config.json
[INFO] Mock mode: sqlite
[INFO] Recording file: data/recording_001.sqlite
[INFO] Site config source: SQLite SiteInformation.site_json
[INFO] Serving site config:
       Site ID: b3489888-aacf-4451-893c-d7d994240f93
       Site Name: Bentonville
       Readers: 3
       Antennas: 5
       Floors: 1
       Regions: 3

Example command:
rre mock-server -config configs/fail_config.json

Expected terminal output for file mode:

[INFO] Loaded run config: configs/fail_config.json
[INFO] Mock mode: file
[INFO] Mock site file: configs/wrong_site_config.json
[INFO] Site config source: JSON file
[INFO] Serving site config:
       Site ID: WRONG-SITE-ID
       Site Name: Wrong Site
       Readers: 0
       Antennas: 0
       Floors: 0
       Regions: 0

Also improve validate -config terminal output.

When running:
rre validate -config configs/pass_config.json

Print:

[INFO] Loaded run config: configs/pass_config.json
[INFO] Recording file: data/recording_001.sqlite
[INFO] Target URL: http://localhost:8080
[INFO] Site ID: b3489888-aacf-4451-893c-d7d994240f93
[INFO] Recorded site config:
       Site ID: ...
       Site Name: ...
       Readers: ...
       Antennas: ...
       Floors: ...
       Regions: ...
[INFO] Target site config:
       Site ID: ...
       Site Name: ...
       Readers: ...
       Antennas: ...
       Floors: ...
       Regions: ...

Then existing validation result:
[OK] Validation passed: recorded site matches target site

For fail config, expected output should show different recorded vs target summary, then:
[ERROR] Validation failed with X error(s)

Implementation guidance:
1. Add a small helper function to print site config summary.
   Example name:
   printSiteConfigSummary(title string, cfg *site.SiteConfig)

2. Summary should print:
   - Site ID
   - Site Name
   - Readers count
   - Antennas count
   - Floors count
   - Regions count

3. Reuse this helper in:
   - runMockServerWithConfig after loading siteConfig
   - runValidate after loading recordedConfig and targetConfig
   - optional: runPlay before validation result, because play validates first

4. Do not print full readers/antenna details by default.
   Only counts and site id/name.

5. Do not change validation behavior.
6. Do not change mock server behavior.
7. Do not change replay pacing or payload sending.

Files likely involved:
- internal/cli/commands.go

Only touch other files if absolutely necessary.