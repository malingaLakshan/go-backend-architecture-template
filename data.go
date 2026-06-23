Please understand this project first. Do not modify code yet.

Project:
This is a Go CLI MVP for a Replay Engine inside the replay-engine-mvp-cli folder.

Domain:
A Recorder tool records RFID/location reader data into a SQLite database. The Replay Engine reads that recorded SQLite file and replays the same raw reader data into a target system. For this MVP, the target system is a mock HTTP server.

Main purpose:
The Replay Engine helps replay the same recorded RFID test multiple times so different target system versions/configurations can be tested with the same input data.

Important concepts:

* SQLite recording file = recorded test data
* SiteInformation = recorded site configuration
* RawReads = recorded RFID reader records
* InjectionTime = original timing used for replay pacing
* Pacing = replaying records with the same time gaps as the original recording
* Mock target server = fake target system for MVP testing

MVP features:

1. Generate sample SQLite recording data.
2. Insert sample SiteInformation.
3. Insert sample RawReads with x/y movement.
4. Show recording summary.
5. Validate recorded site against target site.
6. Replay raw records ordered by InjectionTime.
7. Preserve timing gaps using InjectionTime.
8. Send replay payloads to mock target server.
9. Maintain console logs and file logs.
10. Support safe Ctrl+C abort.

Important rules:

* Do not send all records immediately.
* First record can be sent immediately.
* Every next record must wait based on the InjectionTime gap.
* If validation fails, replay must not start.
* If one injection fails during replay, log it and continue.
* No pause feature. Only stop/abort.
* Logs should go to console and logs/rre.log.
* Generated SQLite files and log files should not be committed.

Please analyze the current project and explain:

1. Current folder structure.
2. Main CLI commands.
3. Package responsibilities.
4. Replay flow.
5. Validation flow.
6. Logging flow.
7. Whether this project matches the MVP requirements.
8. Any issues or improvements.

Do not change files yet. First give me your analysis.