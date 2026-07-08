Please first understand the project before making any changes.

Important:
Do not modify any files yet.
Do not change existing functionality.
Do not refactor the project.
Do not change command names.
Do not remove validation.
Do not change replay pacing.
Do not change injector behavior.
Do not change existing working real SQLite support.

Project context:
This is the RRE CLI MVP.
It reads Recorder SQLite files, validates recorded site configuration against target site configuration, and replays RawReads with InjectionTime-based pacing.

Current working flow:
- RRE reads recorded site config from SQLite SiteInformation.site_json
- RRE reads RawReads from the real Recorder SQLite schema
- validate command compares recorded site config with target site config
- play command validates first, then replays RawReads
- mock-server currently acts as target Resonate API
- validation and play are currently working with the real Recorder SQLite file

New QA requirement:
QA needs to test validation pass and validation fail scenarios using different Recorder SQLite files.

Problem:
QA cannot manually change the target URL site configuration every time to match each recorded SQLite file.

Requirement to understand:
We need config-based mock target support, so mock-server can return target site configuration from:
1. sqlite mode:
   Load SiteInformation.site_json from the given recording SQLite file.
   This should make validation pass when using the same SQLite file.
2. file mode:
   Load target site config from a provided JSON file.
   This allows QA to intentionally provide wrong/mismatched site config and verify validation fails.

Expected future commands:
rre mock-server -config configs/pass_config.json
rre validate -config configs/pass_config.json
rre play -config configs/pass_config.json

Old commands should continue to work:
rre mock-server -port 8080
rre validate -file ... -target-url ... -site-id ...
rre play -file ... -target-url ... -site-id ...

Please inspect the project and explain:
1. Which files are involved.
2. Current flow for mock-server, validate, and play.
3. Minimal safe implementation plan.
4. Any risks.
5. Exact files you would change.

Do not write code yet.