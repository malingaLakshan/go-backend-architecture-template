Implement ONLY this Jira subtask:

[RRE][IMP][UI] Implement target connection and site retrieval

Context:
- The React UI shell already exists under resonate-replay-engine/ui.
- The existing Go Replay Engine CLI and backend are working and must remain unchanged.
- During development, `rre serve` is started separately to run the existing local mock Resonate target.
- The UI connects to whatever Target Resonate URL the user enters.
- `http://localhost:8080` is only a development example. Never hardcode it as application behavior.
- In future, the same UI must work with an already-running real Resonate target without code changes.

Requirements:

1. Inspect the existing Go site client/service and REUSE the current site-fetching logic. Do not duplicate site business logic.

2. Add a thin HTTP API layer under `internal/api/`.

Implement:

GET /api/v1/health

GET /api/v1/sites?targetUrl=<user-supplied-url>

3. `/api/v1/sites` must:
- require and validate `targetUrl`;
- use the supplied URL dynamically;
- support the HTTP/HTTPS schemes already supported by the existing backend;
- call the existing site client/service;
- return the actual sites exposed by that target;
- use request context and a reasonable timeout;
- use the existing common logger;
- return stable, safe JSON errors.

Success shape:

{
  "data": [
    {
      "id": "...",
      "name": "..."
    }
  ]
}

Error shape:

{
  "error": {
    "code": "TARGET_UNAVAILABLE",
    "message": "Unable to connect to the target Resonate instance."
  }
}

4. Add only the minimum API server startup/entry point required for the React UI. Existing CLI commands and behavior must continue unchanged.

5. Connect the existing React UI:
- use the Target URL currently entered by the user;
- Connect calls `/api/v1/sites?targetUrl=<encoded-url>`;
- show loading, success, and error states;
- populate Target Site dropdown with the real returned sites;
- clear stale sites/selection when URL changes or reconnect starts;
- prevent duplicate Connect requests while loading;
- keep Play disabled;
- do NOT implement recording selection, validation, or playback yet.

6. Keep frontend HTTP access in the existing API/service layer (`replayApi.ts` or equivalent), not directly scattered through UI components.

7. Configure Vite `/api` proxy if required for development. Avoid unrestricted CORS.

8. Add concise `docs/API.md` documentation for the health and sites endpoints.

9. Add focused tests for missing/invalid target URL, unavailable target, and successful site retrieval where practical.

STRICT:
- Do NOT modify replay, pacing, recording, validation, SQLite, config loading, mock-target behavior, logger behavior, or existing CLI business logic.
- Do NOT hardcode site data or target URLs.
- Do NOT automatically start the mock target from Connect.
- Do NOT implement later UI subtasks.
- Keep changes minimal and production quality.

Expected local development flow:
1. Run `rre serve` separately.
2. Start the Replay Engine UI/API.
3. Enter `http://localhost:8080` in the UI.
4. Click Connect.
5. Load the real SiteGraphs/sites exposed by the running mock target.

Before editing:
- inspect the existing site-fetching implementation and API startup options;
- identify the exact existing function/service you will reuse;
- list only the files you intend to create/modify;
- WAIT for my confirmation before changing code.