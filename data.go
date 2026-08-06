Implement ONLY this Jira subtask:

[RRE][IMP][UI] Implement target connection and site retrieval

Context:
- The React UI shell already exists under resonate-replay-engine/ui.
- The existing Go Replay Engine CLI and backend functionality are working and must remain unchanged.
- The UI must now connect to the real existing backend logic. Do not use hardcoded or mock site data.

Requirements:

1. Inspect the existing Go packages and reuse the current site-fetching logic used by the Replay Engine.
2. Add a thin HTTP API layer under a new package such as:

   internal/api/

3. Implement these endpoints:

   GET /api/v1/health

   GET /api/v1/sites?targetUrl=http://localhost:8080

4. The sites endpoint must:
   - validate the target URL;
   - call the existing backend site client/service;
   - return the actual available sites from the supplied target;
   - not duplicate existing site-fetching business logic;
   - use the existing common logger;
   - use request context and a reasonable timeout.

5. Use a stable JSON response structure.

Successful sites response example:

{
  "data": [
    {
      "id": "site-id",
      "name": "site-name"
    }
  ]
}

Error response example:

{
  "error": {
    "code": "TARGET_UNAVAILABLE",
    "message": "Unable to connect to the target Resonate instance."
  }
}

6. Add a new command or entry point for starting the UI API only if required, but do not change the behavior of existing CLI commands.

7. Update the React UI:
   - implement the real API call in replayApi.ts;
   - use the entered Target Resonate URL;
   - call the sites endpoint when Connect is clicked;
   - show loading, connected, and error states;
   - populate the Target Site dropdown using the real response;
   - keep the Play button disabled;
   - do not implement file selection, validation, or playback yet.

8. During React development, configure Vite proxying for /api if needed. Avoid unrestricted CORS.

9. Add concise API documentation for these endpoints in:

   docs/API.md

10. Add focused Go handler tests and frontend tests only where practical.

Strict constraints:
- Do not modify replay, pacing, recording, validation, SQLite, config-loading, mock-target, logger, or existing CLI business logic.
- Reuse existing backend services instead of copying logic.
- Do not add hardcoded site data.
- Do not implement later Jira subtasks.
- Keep changes minimal and production quality.

Before editing:
- Inspect the existing site client, CLI startup flow, logger, and module structure.
- Briefly show the files you intend to create or modify.
- Wait for my confirmation before changing files.