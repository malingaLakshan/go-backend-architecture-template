1. Create RRE UI application shell and base layout

Summary:
RRE UI - Create application shell and base layout

Description:
Create the initial RRE graphical application and implement the single-window layout shown in the story mock-up.

Include:

* Target Resonate URL field
* Connect button
* Target Site dropdown
* Select Recording button
* Selected recording information area
* Site validation status area
* Play button
* Correct initial enabled/disabled states

Acceptance criteria:

* Application launches successfully with a clean window.
* All required controls are visible.
* Target Site dropdown is disabled until connection succeeds.
* Play button is disabled by default.
* Layout follows the attached mock-up.

Estimate: 1 day

⸻

2. Connect to Target Resonate and load available sites

Summary:
RRE UI - Implement target connection and site retrieval

Description:
Connect the UI to the selected Target Resonate instance. When the user clicks Connect, call the existing backend client to retrieve available sites using the target URL.

Acceptance criteria:

* Target URL is required before connecting.
* Connect calls the target /sites endpoint.
* Successful response populates the Target Site dropdown.
* Site name and identifier are retained for the selected site.
* Loading, successful connection, no-sites, and connection-failure states are displayed.
* Changing the URL clears previous sites, validation results, and disables Play.

Estimate: 1.5 days

⸻

3. Select SQLite recording and display recording information

Summary:
RRE UI - Implement SQLite recording file selection

Description:
Add the native recording file picker and integrate it with the existing SQLite reading functionality.

Acceptance criteria:

* Select Recording opens the native file browser.
* File picker filters for .sqlite files.
* Selected database is opened in read-only mode.
* UI displays the selected filename.
* UI reads and displays the original recorded Site Name from SiteInformation.
* Invalid files, inaccessible files, and missing SiteInformation data show clear errors.
* Selecting another recording clears the previous validation result and disables Play.

Estimate: 1.5 days

⸻

4. Integrate automatic site configuration validation

Summary:
RRE UI - Integrate automatic site validation

Description:
Connect the UI to the existing RRE site-validation logic. Validation must run automatically when both a recording and Target Site have been selected.

Acceptance criteria:

* Validation starts automatically when both required selections exist.
* UI shows a validating/loading state while validation runs.
* Existing Go validation logic is reused.
* The selected target site is fetched using /sites/{siteId}.
* Successful validation displays a green indicator and Ready to Play.
* Failed validation displays a red indicator and the specific backend error.
* Changing the recording, target site, or URL invalidates the previous result.
* Play is enabled only for the currently validated recording and target combination.

Estimate: 2 days

⸻

5. Implement playback initialization and configuration locking

Summary:
RRE UI - Implement playback initialization and input locking

Description:
Connect the Play button to the existing RRE playback engine so that pacing and data injection begin using the selected recording and Target Site.

Acceptance criteria:

* Play is disabled initially, during validation, and after failed validation.
* Play becomes enabled only after successful validation.
* Clicking Play starts the existing backend playback process.
* Multiple playback starts or double-click execution are prevented.
* After playback starts, URL, Connect, Target Site, and recording controls become read-only or disabled.
* UI displays a clear starting/running state.
* A playback-start failure displays the error and allows a safe retry.

Estimate: 1.5 days

⸻

6. Add workflow state management, error handling, and logging

Summary:
RRE UI - Add workflow state management and error handling

Description:
Make the complete UI workflow reliable and ensure backend failures are clearly presented without leaving the UI in an incorrect state.

Acceptance criteria:

* UI handles connecting, connected, file-loaded, validating, ready, failed, starting, and running states.
* Stale validation responses cannot enable Play after an input changes.
* Network, SQLite, validation, and playback errors are user-friendly.
* Technical details are written through the existing common logger.
* UI cannot enter an invalid combination of enabled controls.

Estimate: 1 day

⸻

7. Add tests, documentation, and acceptance-criteria verification

Summary:
RRE UI - Add automated tests and verify story acceptance criteria

Description:
Add tests for the new UI workflow and perform complete story-level verification.

Acceptance criteria:

* Tests cover connection success and failure.
* Tests cover valid and invalid SQLite recordings.
* Tests cover automatic validation pass and failure.
* Tests verify Play button enablement rules.
* Tests verify configuration locking after playback starts.
* Existing RRE tests continue to pass.
* Build, CI, security, and quality checks pass.
* README contains UI build, launch, and usage instructions.
* All four sections of TRIF-6996 are manually verified against the mock-up.