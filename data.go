[RRE][UI] Implement real-time playback progress tracking

Description:

Implement real-time playback progress tracking in the Replay Engine Status dashboard.

The dashboard should show the current progress of the active playback session using a progress bar and playback timer.

The progress information should update continuously while records are being injected by the Replay Engine.

Acceptance Criteria:

1. Progress Bar
- Display a visual progress bar in the Status tab.
- The progress bar represents the current playback progress.
- Progress updates continuously while playback is running.
- Progress reaches 100% when playback completes successfully.

2. Playback Timer
- Display playback time in HH:MM:SS / HH:MM:SS format.
- The first value represents the current elapsed playback time.
- The second value represents the total duration of the recording.
- Example: 00:38:42 / 01:25:17.
- The timer updates every second while playback is running.

3. Playback State
- Progress information starts updating when playback begins.
- Progress information stops updating when playback completes or is stopped.
- The final progress state is displayed correctly when playback completes.

4. Existing Functionality
- Existing playback functionality and Configuration/Status tab navigation must continue to work without regression.