Polish the current React UI shell only.

Changes:
- Increase the application width and height so it uses the browser space better and matches the provided mockup.
- Remove the fake minimize, maximize, and close buttons from the header.
- Keep a clean application header with only the title “Resonate Replay Engine”.
- Improve spacing and alignment for the form sections.
- Keep the Play button disabled.
- Do not implement backend calls, validation, playback, or real SQLite loading.
- For the Select Recording button, add a temporary disabled or non-functional placeholder state with clear text such as “Recording selection will be implemented in the next subtask”, so users are not confused.
- Modify only files inside ui/.
- Run npm run build and report the changed files.