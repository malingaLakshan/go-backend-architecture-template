Approved. Implement exactly the proposed changes.

Requirements

1. Add a new CLI command:

   rre ui

2. `rre ui` must:
- start the existing Replay Engine UI/API server;
- default to port 9090;
- serve the built React app from `ui/dist`;
- expose the existing `/api/v1/*` endpoints from the same server;
- open the UI at http://localhost:9090;
- not start the mock target automatically.

3. Keep `rre serve` unchanged.
It continues to start the local mock Resonate target (e.g. localhost:8080).

4. Update `rre help`.

Keep the current format.

Recommended workflow:

1. rre summary
2. rre serve
3. rre ui
4. rre sites
5. rre validate
6. rre play
7. rre kill

Available commands:

rre summary    Summarize a recorded site configuration
rre serve      Start the target server
rre ui         Start the Replay Engine UI
rre sites      List all available sites
rre validate   Validate recorded site configurations against target site
rre play       Replay RawReads into target Resonate
rre kill       Kill the server

Keep "Recommended workflow" and "Available commands" as separate sections exactly like today. Do not add unnecessary text.

5. Update README.md with a short "Local Mock / QA" section:

1. Run `rre serve`
2. Run `rre ui`
3. Open `http://localhost:9090`
4. Enter `http://localhost:8080`
5. Click Connect

Also mention:
- `rre serve` is only required for the local mock target.
- For a real deployment, run `rre ui` and enter the real Resonate target URL.

6. If `ui/dist` does not exist, `rre ui` must return a clear error telling the user to build the React UI first.

7. Reuse the existing UI/API server. Keep the implementation minimal.

Strict constraints

- Do NOT change the behavior of `serve`, `summary`, `sites`, `validate`, `play`, or `kill`.
- Do NOT modify replay, pacing, recording, validation, SQLite, config loading, mock target, logger, or existing business logic.
- Do NOT require `npm run dev` or Node.js at QA runtime.
- Do NOT add unrelated commands or features.

After implementation:
- run the relevant Go tests;
- run `npm run build`;
- confirm `ui/dist` exists;
- show all modified files;
- show the exact command to rebuild `rre.exe`.