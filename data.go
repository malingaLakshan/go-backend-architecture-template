Implement ONLY this Jira subtask:

[RRE][IMP][UI] Implement recording selection and automatic validation

Context

- The React UI already connects to the existing backend.
- Target URL connection and Target Site selection are already working.
- Replay Engine CLI/backend business logic already exists and must be reused.
- Do NOT duplicate recording parsing or validation logic.

Requirements

1. Allow selecting a recording SQLite database.
- Accept only `.sqlite` files.
- Show the selected file name in the UI.
- Do not hardcode file paths.

2. Add the minimum backend API required for file selection and validation.
Reuse existing recording parsing and validation services.
Do NOT copy business logic.

3. After BOTH are available:
- Connected Target Site
- Selected recording

automatically trigger the existing validation logic.

No Validate button.

4. Validation must use the existing backend implementation already used by the CLI.

Do not implement validation in React.

5. Update the Validation panel.

Display:
- Loading while validating
- Validation success
- Validation failure
- Existing validation message/details returned by the backend

6. Behaviour

Changing:
- Target URL
- Target Site
- Recording

must clear the previous validation result and automatically revalidate when all required inputs are available.

7. Play button

- Keep disabled until validation succeeds.
- Do NOT implement playback in this subtask.

8. API

Reuse existing services.

Expose only the minimal endpoints required for:
- recording selection
- automatic validation

Use the existing logger, request context and reasonable timeouts.

9. Documentation

Update README.md with the new workflow.

Local Mock / QA flow

1. `rre serve`
2. `rre ui`
3. Open `http://localhost:9090`
4. Enter target URL (e.g. `http://localhost:8080`)
5. Connect
6. Select Target Site
7. Select Recording (.sqlite)
8. Validation starts automatically
9. Play becomes enabled only after successful validation

Strict constraints

- Reuse the existing recording parser, SQLite reader and validation services.
- Do NOT duplicate business logic.
- Do NOT modify replay, pacing, serve, summary, sites, config loading, logger or mock-target behaviour.
- Do NOT implement playback.
- Do NOT add future Jira functionality.
- Keep changes minimal and production quality.

Before editing:
- inspect the existing recording parsing and validation flow;
- list only the files that will be modified;
- wait for my confirmation before changing files.

After implementation:
- run relevant Go tests;
- run `npm run build`;
- report all modified files;
- describe how to manually verify the feature.