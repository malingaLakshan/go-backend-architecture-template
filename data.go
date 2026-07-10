## Summary
Updated the Replay Engine to support the real Recorder SQLite schema and config-based testing flow.

## Changes
- Updated Replay Engine code to work with the renamed `resonate-replay-engine` folder structure.
- Added config-based command support for validate, play, and mock-server flows.
- Synced RecordingSession, SiteInformation, and RawReads handling with the real Recorder SQLite schema.
- Fixed RawReads nullable column scanning for fields like RSSI, antenna, confidence, tag position, and floor.
- Updated site config validation to compare site-related IDs and detect mismatches.
- Added pass/fail config files for QA testing.
- Updated terminal output to show recorded and target site config summaries.
- Fixed README/help usage for config-based testing.

## Testing
- Ran validate with pass config and confirmed validation passes.
- Ran play with pass config and confirmed RawReads are loaded and replayed.
- Ran mock server with config and confirmed received payloads.
- Ran wrong site config and confirmed validation fails as expected.

## Notes
Generated logs and local database runtime files should not be committed.