## Resonate Replay Engine v0.3.0

### Added

- Integrated the shared ALTRFIDTools common logging module.
- Added logger configuration through the Replay Engine configuration file.
- Added logging for summary, validation, playback, and mock-server workflows.
- Added mock-server lifecycle and error logging.
- Added support for appending multiple command executions to the configured Replay Engine log file.

### Changed

- Replaced the Replay Engine-specific logger implementation with the shared logger.
- Updated Replay Engine module dependencies and configuration model.
- Updated documentation for architecture, commands, validation, and logging.

### Validation

- All automated tests passed.
- Static analysis completed.
- Release executable built successfully.
- `serve`, `validate`, and `play` manually verified.