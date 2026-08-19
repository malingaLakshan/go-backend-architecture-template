[RRE][UI] Implement live playback metrics

Description:

Implement live playback metrics in the Replay Engine Status dashboard.

The dashboard should display the number of records successfully injected during playback and the total number of records available in the selected recording.

The values must update based on the active playback session.

Acceptance Criteria:

1. Records Injected
- Display a "Records Injected" value in the Status tab.
- The value represents the number of RawReads successfully processed/sent so far.
- The value updates continuously while playback is running.

2. Total Records
- Display a "Total Records" value in the Status tab.
- The value represents the total number of RawReads in the selected .sqlite recording.
- The total value remains consistent throughout the playback session.

3. Playback Completion
- When playback completes successfully, Records Injected should match Total Records.
- The final values remain visible after playback completes.

4. Existing Functionality
- Existing playback, progress tracking, timer, and Configuration/Status tab behavior must continue to work without regression.