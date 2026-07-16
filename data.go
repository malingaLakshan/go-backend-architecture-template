## Resonate Replay Engine v0.1.2

This release improves SiteGraph handling, structural validation, and multi-site mock-server support for the Resonate Replay Engine.

### Highlights

- Added support for loading multiple SiteGraph JSON files
- Added automatic SiteGraph discovery from the configured sites directory
- Added `GET /sites` to list available mocked sites
- Added `GET /sites/{siteId}` to return the complete matching SiteGraph
- Preserved the original full SiteGraph JSON without re-serialization
- Added recursive nested-region parsing and validation
- Updated reader handling so readers are validated directly under floors
- Added antenna-port validation under the correct reader
- Added detailed validation results for missing floors, regions, readers, and antenna ports
- Ensured playback stops before RawReads are loaded when validation fails
- Improved QA-friendly terminal output
- Updated CLI help, configuration examples, and documentation
- Restricted SiteGraph file access to the approved directory

### Validation hierarchy

The Replay Engine now validates the confirmed SiteGraph structure:

```text
Site
└── Floors
    ├── Regions
    │   └── Nested Regions recursively
    └── Readers
        └── Antenna Ports