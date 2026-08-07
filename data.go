Approved. Implement the proposed plan now.

Keep the scope exactly as discussed:
- create only the required API/site integration files and frontend connection changes;
- reuse existing backend site-fetching logic;
- do not modify existing replay, pacing, recording, validation, SQLite, config, mock-target, logger, or CLI business behavior;
- no hardcoded target URL or site data;
- keep Play disabled and do not implement later subtasks.

After implementation:
- run Go tests for the affected packages;
- run npm run build;
- list the files created/modified and any startup command required for the new API.