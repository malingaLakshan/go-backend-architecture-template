Implement ONLY the minimum playback status integration needed to fix the current UI state and enable the live Status dashboard.

Current issue:
- Playback starts successfully.
- UI switches to Status.
- After replay completes, Configuration still shows the Play button as "Playing..." and disabled.
- Status page currently has 0 records, 0:00 / 0:00, no progress, and empty activity log.

Goal:
Expose real playback status from the existing replay process and wire it to the current React UI.

Requirements:

1. Backend playback status
Add a minimal read-only status endpoint, for example:

GET /api/v1/playback/status

Return the current replay state and metrics from the EXISTING replay execution.

Suggested response fields:
- state: idle | starting | playing | completed | failed | aborted
- processedRecords
- totalRecords
- elapsedMs
- totalDurationMs
- recordingFile
- targetSiteName
- message (optional)

2. Do NOT create a second replay engine or duplicate replay logic.
Observe/update status around the existing replay flow only.

3. UI polling
While state is starting/playing:
- poll the playback status endpoint about once per second;
- stop polling when state becomes completed, failed or aborted.

4. Fix the stuck Play button
- When backend reports completed/failed/aborted, update the UI workflow state.
- The button must no longer remain "Playing..." after replay finishes.
- Restore configuration controls after final state.
- Preserve validation result where appropriate so the user can start another run.

5. Status dashboard
Wire real values to:
- Replaying: selected recording
- Target Site
- Records Injected
- Total Records
- elapsed time / total duration
- progress bar

Progress:
processedRecords / totalRecords * 100

Do not use fake/mock dashboard values.

6. Timing
Do NOT invent timing in React.
Use backend/existing replay timing information for elapsed and total duration.
Do NOT modify pacing or scheduling calculations.

7. Activity Log
For now, only show minimal lifecycle messages if already available from replay status, such as:
- Playback started
- Playback completed
- Playback failed/aborted

Do NOT implement full log streaming yet.

8. State handling
Ensure:
- successful Play -> auto-switch to Status
- user can manually switch between Configuration and Status
- playing -> completed/failed/aborted transitions correctly
- stale status is cleared before a new replay starts

STRICT CONSTRAINTS:
- Do NOT modify replay timing/pacing logic.
- Do NOT modify payload generation, HTTP injection, replay summary, recording parsing, validation, SQLite, mock target, config loading, logger behavior, or existing CLI behavior.
- Do NOT refactor unrelated backend code.
- Do NOT implement Stop/Abort backend behavior yet.
- Keep changes minimal and quota-efficient.

Before editing:
- inspect only the current playback-start API, replay status object/state, and relevant React playback/status files;
- identify the smallest files to modify;
- show the proposed file list;
- wait for my confirmation.

After implementation:
- run relevant Go tests;
- run npm run build;
- report only modified files and manual verification steps.