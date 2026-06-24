The current Mermaid diagram is mostly correct, but please improve it to make the Replay Engine flow clearer.

Required changes:

1. Show validation before replay starts.
2. Show that Site Validator fetches target site configuration from Mock Target Server.
3. Show that Replay Engine reads raw records only after validation passes.
4. Show Pacing Logic between Replay Engine and HTTP Injector.
5. Show HTTP Injector sends replayed payloads to Mock Target Server.
6. Show Logger as a shared component receiving logs from CLI Dashboard, Site Validator, Replay Engine, HTTP Injector, and Mock Target Server.
7. Keep it high-level and professional.
8. Update replay-engine-mvp-cli/docs/ARCHITECTURE.md automatically.
9. Do not change Go source code.

Main flow should be:
User → CLI Dashboard / Commands → SQLite Recording File → Recording Reader → Site Validator → Replay Engine → Pacing Logic → HTTP Injector → Mock Target Server

Also show:
Site Validator → Mock Target Server for target site config
All main components → Logger