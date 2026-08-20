[RRE][UI] Align Replay Engine UI with Command Center theming and shared UI components

Description:

Align the Replay Engine UI with the existing Command Center design system and shared UI components.

Update the Replay Engine UI styling and components to provide a consistent look and feel with the Command Center while preserving all existing Replay Engine functionality.

Acceptance Criteria:

1. Command Center Theming
- Apply the existing Command Center theme, colors, typography, spacing, and visual standards to the Replay Engine UI.
- Maintain a consistent appearance across Configuration and Status views.

2. Shared UI Components
- Reuse existing Command Center/shared UI components where applicable.
- Avoid duplicating components or styles that already exist in the shared component library.

3. Replay Engine UI
- Align buttons, inputs, dropdowns, tabs, cards, status indicators, dialogs, progress elements, and other applicable controls with the Command Center design.
- Preserve the existing responsive and centered layout behavior.

4. Existing Functionality
- Theming/component changes must not change Replay Engine business logic or workflow behavior.
- Existing target connection, validation, playback, status, progress, metrics, activity log, and Stop/Abort functionality must continue to work without regression.

5. Maintainability
- Avoid unnecessary custom styling where an existing shared component or theme can be reused.
- Keep the implementation consistent with existing Command Center frontend standards.