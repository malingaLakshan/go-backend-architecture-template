Title: [R&D] Investigate and document the Resonate Locate instance

Description:

Conduct a hands-on R&D investigation of the provided Resonate Locate instance to understand its current deployment, components, configuration, interfaces, data flow, and available capabilities.

The investigation should determine what the instance currently contains, what activities can be performed using it, and how it could support the Resonate tool suite, including Recorder, Replay Engine, Validator, Analyzer, and Simulator.

RawReads acquisition and capture should be investigated as one part of the research, but it is not the only objective.

Scope:

* Verify access to the provided instance and dashboard.
* Inspect the host environment and existing Docker Compose deployment.
* Identify running containers, services, images, ports, networks, volumes, and dependencies.
* Understand the purpose of each available component.
* Review available configuration files without exposing sensitive values.
* Identify available user interfaces, APIs, MQTT connections, databases, logs, and other interfaces.
* Review the provided simulator and understand how it generates and sends RFID data.
* Investigate whether the instance receives data from real readers, simulated readers, or both.
* Investigate how RawReads enter the system and whether Recorder can capture them.
* Explore how the instance could be used for Replay Engine, validation, simulation, analysis, and end-to-end testing.
* Perform safe and non-disruptive practical tests where permitted.
* Determine whether any additional deployment or configuration changes are required.
* Identify technical blockers, limitations, dependencies, access issues, and risks.
* Produce an R&D findings document containing evidence, conclusions, and recommended next steps.

Acceptance Criteria:

* The current deployment and its components are documented.
* The purpose of the identified services and components is explained.
* A high-level architecture and data-flow diagram is included.
* Available interfaces, APIs, MQTT connections, storage, configuration, and logs are documented where accessible.
* Possible uses of the instance for the Resonate tool suite are identified.
* RawReads generation, ingestion, storage, and capture findings are documented.
* Safe practical experiments and their results are recorded.
* Required deployment or configuration changes are clearly stated.
* Blockers, limitations, risks, and unanswered questions are documented.
* Recommended next steps are provided.
* The completed R&D findings document is attached or linked to the Jira ticket.

Out of Scope:

* Permanent configuration changes.
* Restarting, upgrading, or redeploying shared services without approval.
* Implementing product changes identified during the investigation.