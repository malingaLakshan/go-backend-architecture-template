Implement ONLY this Jira subtask:

[RRE][IMP][UI] Implement playback initialization and input locking

Context:
- Target connection, site selection, recording selection, and automatic validation are already working.
- The existing Go replay engine is already implemented and timing/pacing accuracy is verified.
- Do NOT change replay internals, pacing logic, timing calculations, payload generation, HTTP injection behavior, or replay summary logic.
- The UI must reuse the existing replay service exactly as it works today.

Requirements:

1. Play button
- Enable only when validation has successfully passed.
- Keep disabled for validation failure, missing inputs, or validation in progress.

2. Start playback
- Clicking Play must call a thin backend API that invokes the existing replay service.
- Do not duplicate or rewrite replay logic.
- Do not invoke CLI commands from React.
- Reuse the same backend functionality currently used by `rre play`.

3. Playback should start asynchronously.
- Do not keep the HTTP request open for the full recording duration.
- Return a minimal playback identifier/status after successful start.

4. When playback starts, lock configuration inputs:
- Target URL
- Connect button
- Target Site
- Recording selection
- Play button

5. UI states:
- starting
- playing
- start failed

Show a simple clear status only.
Do NOT implement the full playback dashboard, progress bar, activity log, metrics, timer, or stop functionality in this subtask.

6. If playback startup fails:
- show a clear error;
- restore the configuration controls;
- keep existing validation state where appropriate.

7. Backend/API
Add only the minimum endpoint required to start playback.
Use stable JSON responses, existing logger, request context, and existing services.

Prepare the API response so future playback-status/dashboard functionality can extend it without changing the existing replay engine.

8. Documentation
Update API documentation and README with the new Play flow.

STRICT CONSTRAINTS:
- Do NOT modify internal replay timing/pacing implementation.
- Do NOT modify `CalculateDelay` or the absolute scheduling logic.
- Do NOT modify payload building, HTTP injection, replay summary, recording parsing, validation, SQLite, mock target, config loading, logger behavior, or existing CLI command behavior.
- Do NOT refactor working backend packages.
- Do NOT change the existing `rre play` behavior.
- Do NOT implement later dashboard/stop/progress subtasks.
- Keep changes minimal and production quality.

Before editing:
- inspect how `rre play` currently creates and runs the replay service;
- identify exactly which existing functions/services will be reused;
- list the files you intend to create/modify;
- WAIT for my confirmation before changing code.

After implementation:
- run relevant Go tests;
- run `npm run build`;
- confirm no replay/pacing implementation files were modified unless strictly required for compilation;
- list all modified files;
- provide the manual test steps.