## Replay Engine v0.4.0

This release adds Replay Engine server termination support and improves playback pacing accuracy.

### New Feature

- Added the `rre stop` command to terminate the Replay Engine server running on the configured port.
- Added Replay Engine process verification before termination.
- Added process-tree termination support.
- Integrated stop-command activity with the common logger.
- Added clear terminal messages for successful termination and cases where no Replay Engine server is running.

### Playback Pacing Improvements

- Updated replay pacing to follow an absolute timeline based on the first recorded `InjectionTime`.
- Prevented payload processing, JSON serialization, file I/O, HTTP request time, and scheduling overhead from accumulating between records.
- Preserved the existing cancellation, status, logging, and replay-output behavior.

### Validation Results

Playback pacing was validated using a recording containing 7,045 RawReads:

- Original recording duration: `11m42.9801641s`
- Replay duration: `11m42.9855934s`
- Timing difference: approximately `5.43ms`
- Successful records: `7045/7045`
- Failed records: `0`
- Replay status: `completed`

The Replay Engine now closely follows the original recorded timeline without cumulative timing drift.