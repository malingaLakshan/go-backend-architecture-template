The previous implementation for making -config optional caused replay behavior to change. Revert any complex fallback or automatic-discovery logic and implement the requirement only through CLI flag defaults.

Requirements:
- For every command that accepts -config, set its default value in args.go to "configs/config.json" instead of an empty string.
- Keep explicit -config <path> overrides working.
- Do not change replay, pacing, validation, logging, configuration parsing, or any business logic.
- Do not add automatic file-search or config-discovery functions.
- Do not modify configs/config.json values.
- Keep commands.go behavior unchanged unless a minimal help-text update is required.
- Limit changes to internal/cli/args.go and directly related CLI tests/help text.
- Run all tests after the change and show the exact diff.