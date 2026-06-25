Please remove all emojis from CLI terminal output and logs in the replay-engine-mvp-cli project.

Reason:
This is an internal/company CLI tool, so terminal output should be professional, plain-text, and compatible with all terminals.

Files to check:

* internal/cli/commands.go
* internal/cli/dashboard.go
* any other file that prints CLI messages using fmt.Print, fmt.Println, or fmt.Printf

Required changes:
Replace emoji-based messages with plain text status prefixes.

Use this style:

* ✅ Validation passed. Ready to replay. → [OK] Validation passed. Ready to replay.
* ❌ Validation failed. ... → [ERROR] Validation failed: ...
* ⚠️ ... → [WARN] ...
* ℹ️ ... → [INFO] ...
* 🚀 Replay started → [INFO] Replay started.
* 🎉 Replay completed → [OK] Replay completed.
* 🛑 Replay aborted → [WARN] Replay aborted by user.

Rules:

* Do not change business logic.
* Do not refactor unrelated code.
* Do not change function names.
* Do not change CLI command behavior.
* Only update user-facing terminal output and any log messages that contain emojis.
* Keep output clear, simple, and professional.
* After changes, show me the list of changed files.