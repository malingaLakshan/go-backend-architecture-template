Implement a new CLI command:

rre ui

Purpose:
Start the Replay Engine browser UI and its HTTP API together for QA/development use.

Expected workflow:

1. `rre serve`
   - starts the existing mock Resonate target
   - example target URL: http://localhost:8080

2. `rre ui`
   - starts the Replay Engine UI/API server
   - default port: 9090

3. QA opens:
   http://localhost:9090

4. QA enters:
   http://localhost:8080

5. Connect retrieves the available sites from the existing mock target.

Requirements:

- Reuse the existing API server implementation already added for `/api/v1/...`.
- Serve the built React production files from `ui/dist`.
- Serve both:
  - `/api/v1/...` API routes
  - React UI/static assets
  from the same HTTP server/port.
- Visiting `/` must load the React application.
- React client-side routes, if any are added later, should safely fall back to `index.html` without intercepting `/api/*`.
- Default UI port should be 9090.
- Support an optional port argument only if consistent with the existing CLI style.
- If `ui/dist` is missing, return a clear startup error telling the developer to run the React production build.
- Use the existing common logger.
- Handle graceful shutdown using context/OS signals where practical.

Strict constraints:
- Do NOT change the behavior of existing commands such as `serve`, `play`, `validate`, `summary`, or `kill`.
- Do NOT modify replay, pacing, validation, recording, SQLite, mock-target, config-loading, or existing business logic.
- Do NOT require Node.js or `npm run dev` for QA runtime.
- Do NOT automatically start the mock target from `rre ui`.
- Keep `rre serve` and `rre ui` as separate responsibilities.
- Keep changes minimal and production quality.

Development/QA target flow must remain:

`rre serve` -> mock Resonate target
`rre ui`    -> Replay Engine UI + API

Before editing:
- inspect the existing CLI command registration and current UI API server;
- list the exact files you intend to modify;
- explain briefly how `/api` and React static-file routing will coexist;
- WAIT for my confirmation before changing files.