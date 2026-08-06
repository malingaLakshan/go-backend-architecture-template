You are working in an existing production-quality Go project.

IMPORTANT CONTEXT

This repository already contains a fully working Replay Engine backend and CLI.

The CLI implementation is complete and production-ready.

The objective of this task is ONLY to begin implementing a new React UI for the existing backend.

The UI will consume backend APIs in future subtasks. Do NOT redesign or replace the backend.

The attached UI mockup images are the design reference for this task.

A future sprint will introduce a Playback Dashboard (progress, metrics, activity log, stop playback, etc.). Please design the architecture so it can easily support those future features, but DO NOT implement them now.

======================================================================
IMPORTANT CONSTRAINTS
======================================================================

DO NOT modify any existing Go backend functionality.

Specifically, DO NOT change:

- Replay logic
- Playback timing/pacing
- Validation logic
- Recording logic
- SQLite logic
- Site loading logic
- Mock target behavior
- CLI commands
- Existing APIs
- Config loading
- Logger
- Existing tests
- Existing package names
- Existing folder structure outside the UI

Treat the current backend as production-ready.

If backend functionality appears to be missing, prepare only the frontend service layer or interface. Do NOT implement or modify backend code.

Only create or modify files required for the React UI.

======================================================================
TASK
======================================================================

Implement ONLY the first Jira subtask:

[RRE][IMP][UI] Create application shell and base layout

Do NOT continue into later subtasks.

======================================================================
REQUIREMENTS
======================================================================

Create a new React + TypeScript application using Vite inside:

resonate-replay-engine/ui

Use a clean, professional, scalable architecture suitable for enterprise applications.

Recommended structure:

ui/
    src/
        app/
        features/
            replay/
                components/
                hooks/
                services/
                types/
                utils/
        shared/
            components/
            styles/
        layouts/
        pages/
        App.tsx
        main.tsx

Do not over-engineer the project.

Use simple React patterns.

Do not introduce Redux.

Use React state only.

======================================================================
UI REQUIREMENTS
======================================================================

Build the Configuration screen matching the attached mockup.

Include:

• Application title/header

• Configuration section

• Target Resonate URL textbox

• Connect button

• Target Site dropdown

• Select Recording button

• Recording information panel

• Validation status panel

• Play button

For this task:

- Use placeholder/local state only.
- No backend calls.
- No HTTP requests.
- No API implementation.
- No playback logic.
- No validation logic.

The Play button must remain disabled.

Use realistic placeholder values only where necessary.

======================================================================
QUALITY
======================================================================

Use:

- Functional components
- TypeScript
- Reusable components
- Accessible HTML
- Responsive layout
- Clean folder separation
- Clear naming
- Professional code style

Avoid unnecessary dependencies.

Avoid unnecessary abstractions.

======================================================================
FUTURE PREPARATION
======================================================================

Although only implementing the Configuration page, prepare the architecture so future stories can easily add:

- Playback Dashboard
- Progress updates
- Activity Log
- Stop Playback
- Live metrics
- API integration

Do NOT implement those features now.

Only prepare the structure.

======================================================================
DOCUMENTATION
======================================================================

Create a short README inside the ui folder explaining:

- prerequisites
- npm install
- npm run dev
- npm run build

======================================================================
BEFORE CODING
======================================================================

First inspect the existing repository.

Briefly explain:

- proposed architecture
- files to be created
- assumptions

Wait for my confirmation before generating or modifying any files.