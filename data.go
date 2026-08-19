[RRE][UI] Implement real-time playback activity log

Description:

Implement the playback Activity Log in the Replay Engine Status dashboard.

The Activity Log should display important playback events in chronological order so the user can clearly understand what is happening during the active replay session.

Acceptance Criteria:

1. Activity Log
- Display an Activity Log section in the Status tab.
- New playback events should appear in chronological order.
- The log area should support scrolling when the number of entries increases.

2. Log Entries
- Each entry should include a timestamp and clear message.
- Show important events such as:
  - Playback started
  - Target connection/status messages
  - Playback warnings or failures
  - Playback completed successfully
  - Playback stopped/aborted by the user

3. Real-Time Updates
- New activity entries should appear while playback is running.
- The user should not need to refresh the page to see new events.

4. Final State
- The final completion, failure, or aborted message should remain visible after playback ends.
- Existing activity entries should remain available for the completed playback session.

5. Existing Functionality
- Existing playback, progress tracking, live metrics, timer, Stop/Abort control, and Configuration/Status tab behavior must continue to work without regression.