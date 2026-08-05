Task: Make the -config flag optional for all Replay Engine CLI commands without changing any existing functionality.

Requirements:
- If the user does not provide -config, automatically load the default configuration from configs/config.json.
- If the user provides -config <path>, continue using the specified configuration file.
- Apply this behavior consistently to all commands that currently use configuration, including:
  - serve
  - validate
  - play
  - kill
  - summary
  - and any other command that accepts -config.
- Preserve all existing functionality, logging, validation, error handling, and command behavior.
- Do not modify any business logic or backend functionality.
- Reuse the existing configuration-loading logic and avoid code duplication.
- Update the CLI help/usage text so that -config is shown as optional instead of required.
- Ensure both of the following work correctly:

rre serve
rre validate
rre play
rre kill
rre summary

and

rre serve -config configs/config.json
rre validate -config configs/config.json
rre play -config configs/config.json
rre kill -config configs/config.json
rre summary -config configs/config.json

Keep the implementation minimal, maintainable, and backward compatible. Before making changes, analyze the current configuration-loading flow and modify only the necessary files.