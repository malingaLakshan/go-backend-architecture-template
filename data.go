1. Dashboard Transition
- When the user clicks Play and playback starts successfully, the UI automatically switches from the Configuration tab to the Status tab.
- The user can manually select the Configuration or Status tab.
- The automatic transition must not break the existing playback-start workflow.

2. Recording Information
- The Status tab displays the name of the recording file currently being replayed.
- The displayed recording filename matches the file selected in the Configuration tab.

3. Target Site Information
- The Status tab displays the selected Target Site name.
- The displayed Target Site matches the site selected/validated in the Configuration tab.

4. Existing Workflow
- Existing target connection, recording selection, site validation, and Play functionality continue to work correctly.
- No regression is introduced to the existing Configuration UI.