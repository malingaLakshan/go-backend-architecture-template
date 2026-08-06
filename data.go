Implement the first Replay Engine UI subtask:

[RRE][IMP][UI] Create application shell and base layout

Context:
- This repository already contains a working Go-based Replay Engine CLI and backend packages.
- The CLI must continue working unchanged.
- We are adding a React UI for the existing backend.
- The attached screenshots are the required visual reference for the current configuration screen.
- A playback dashboard will be implemented in a future story, so keep the structure extensible, but do not implement dashboard functionality now.

Requirements:
1. Create a React + TypeScript UI using Vite inside:

   resonate-replay-engine/ui

2. Use a clean, professional, feature-based architecture suitable for future backend API integration.

3. Create the application shell and base configuration-page layout matching the attached mockup as closely as reasonably possible.

4. The screen must include these UI sections:
   - Application header titled "Resonate Replay Engine"
   - Configuration tab or configuration section
   - Target Resonate URL input
   - Connect button
   - Target Site dropdown
   - Select Recording button
   - Selected recording information area
   - Site validation status area
   - Play button

5. For this subtask only:
   - Use local placeholder state.
   - Do not add real API calls.
   - Do not create or modify Go backend endpoints.
   - Do not duplicate validation or replay logic in React.
   - The Play button must be disabled by default.
   - Site dropdown and recording information may use empty placeholder states.
   - Include visual placeholders for validation success and failure states, but do not implement validation behavior yet.

6. Prepare the structure for later integration with the existing Go backend using separate API, service, hook, and type layers.

7. Use a simple structure similar to:

   ui/src/
     app/
     features/replay/components/
     features/replay/hooks/
     features/replay/services/
     features/replay/types/
     shared/components/
     shared/styles/

   Keep it practical and avoid unnecessary abstractions.

8. Do not add Redux or another external state-management library. Use React state only for this shell.

9. Use accessible HTML:
   - Proper labels for inputs
   - Keyboard-accessible buttons
   - Disabled states
   - Semantic headings
   - Clear focus styles

10. Keep styling clean, simple, responsive, and close to the attached mockup. Avoid excessive animations or decorative UI.

11. Do not modify existing Go CLI, replay, validation, recording, site, logger, config, or pacing logic.

12. Add basic component tests only if the repository already has an established React testing setup. Otherwise, do not add a large testing framework in this subtask.

13. Update or add a short UI README containing:
   - prerequisites;
   - npm install;
   - npm run dev;
   - npm run build.

Before editing:
- Inspect the repository structure.
- Confirm the proposed files and approach briefly.
- Modify only the files required for this subtask.

After implementation:
- Run npm install if necessary.
- Run npm run build.
- Show the exact files created or changed.
- Summarize any assumptions.
- Do not continue into target connection, file selection, validation integration, or playback integration.