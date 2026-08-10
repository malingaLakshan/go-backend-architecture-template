Implement ONLY the minimum UI/state foundation needed for:

[RRE][IMP][UI] Add workflow state management and error handling

and prepare the UI structure for the upcoming real-time playback dashboard.

IMPORTANT:
Copilot quota is limited. Keep the implementation very focused.
Do not inspect or refactor unrelated backend code.
Do not implement new backend APIs in this task.

Current state:
- Configuration UI already works.
- Target connection works.
- Recording selection works.
- Automatic validation works.
- Play starts the existing replay successfully.
- Existing replay/pacing/timing/backend business logic must remain unchanged.

GOAL

Add a clean two-tab UI:

Configuration | Status

1. CONFIGURATION TAB
- Move/reuse the CURRENT working configuration UI inside the Configuration tab.
- Do not redesign or rewrite its working behavior.
- Preserve:
  - Target URL
  - Connect
  - Target Site
  - Recording selection
  - Validation
  - Play
- During active playback, configuration controls remain disabled/read-only.
- User must still be able to manually view the Configuration tab.

2. STATUS TAB
Create the Status dashboard layout matching the existing Jira mockup.

Display sections for:
- Replaying: selected recording file name
- Target Site: selected site name
- Playback Progress bar
- elapsed time / total duration
- Records Injected
- Total Records
- Activity Log area
- Stop button placeholder

IMPORTANT:
- Do NOT hardcode example values such as morning.sqlite, Store #4271, 28,392, etc.
- Use existing selected recording/site values where already available.
- For live values not yet available, use typed empty/default state only, not fake data.

3. TAB BEHAVIOR
- Both Configuration and Status tabs must always be manually clickable.
- After Play STARTS successfully, automatically switch to Status.
- Do NOT wait until playback completes.
- During playback the user may manually switch between Configuration and Status.
- Status remains visible after completion/failure for future final-state handling.

4. CENTRAL WORKFLOW STATE
Introduce a simple centralized workflow state using the existing React state/hook architecture.
Do NOT add Redux/Zustand.

Use a clear state model such as:

idle
connecting
connected
validating
ready
starting
playing
completed
error

Reuse existing state where possible instead of duplicating booleans.

5. PLAYBACK STATUS TYPE
Add a typed model/interface for future live dashboard data, for example:

state
recordingFile
targetSiteName
processedRecords
totalRecords
elapsedMs
totalDurationMs
message/error

Do NOT implement polling or backend status retrieval yet.

6. ERROR HANDLING FOUNDATION
Preserve current errors and ensure workflow state does not become stuck after:
- connection failure
- validation failure
- playback start failure

Do not add new backend error logic.

7. UI REQUIREMENTS
- Preserve current full-screen responsive layout.
- Keep content centered as already implemented.
- Match the Jira Status layout as closely as practical using the existing styling approach.
- No unnecessary new dependencies.

STRICT CONSTRAINTS

- Frontend/state structure only unless a tiny existing API-call adjustment is required for tab switching.
- Do NOT modify replay, pacing, timing, recording, validation, SQLite, mock-target, config loading, logger, CLI behavior, or backend business logic.
- Do NOT implement playback polling.
- Do NOT implement real progress calculation yet.
- Do NOT implement Stop/Abort backend behavior yet.
- Do NOT implement activity-log backend streaming yet.
- Do NOT add hardcoded dashboard data.
- Keep changes minimal and quota-efficient.

Before editing:
1. Inspect ONLY the relevant existing React files.
2. List the exact files you plan to modify/create.
3. Explain the state/tab approach briefly.
4. WAIT for my confirmation.

After implementation:
- run npm run build;
- list modified files;
- explain how to manually verify:
  - Configuration tab still works
  - Status tab can be opened manually
  - successful Play automatically switches to Status
  - Configuration can still be viewed during playback.