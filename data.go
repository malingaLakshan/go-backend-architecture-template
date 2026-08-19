[RRE][UI] Implement playback Stop/Abort control

Description:

Implement the Stop/Abort control for an active Replay Engine playback session.

The Status dashboard should allow the user to safely stop an active playback. The UI must ask for confirmation before stopping because playback cannot be resumed.

Acceptance Criteria:

1. Stop Button
- Display a prominent "Stop" button in the Status tab while playback is active.
- Do not add a Pause button.

2. Confirmation
- When Stop is clicked, display a confirmation dialog.
- The dialog should clearly inform the user that playback cannot be resumed.
- Playback must continue if the user cancels the confirmation.

3. Stop Playback
- If the user confirms, send the stop/abort request to the backend.
- The active playback must stop data injection safely.
- The playback state must update to aborted/stopped.

4. UI Reset
- After playback is stopped, unlock the Configuration controls.
- Allow the user to configure and start a new playback session.
- The final stopped/aborted state should be reflected correctly in the UI.

5. Existing Functionality
- Existing playback, progress tracking, metrics, timer, and tab navigation must continue to work without regression.