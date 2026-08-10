Title: [RRE][R&D] Investigate RawReads acquisition from a real Resonate Locate instance

Description:

Conduct an R&D investigation using the provided Resonate Locate instance to understand its deployment, simulator, RFID data flow, and determine whether RawReads can be obtained and captured for Recorder and Replay Engine testing.

Scope:

* Verify SSH and dashboard access to the provided instance.
* Inspect the currently running Docker Compose services without modifying them.
* Identify container images, ports, networks, volumes, configuration files, and service dependencies.
* Review the provided sim.py simulator.
* Identify the MQTT broker, topic, protobuf message type, site details, reader details, and publish interval used by the simulator.
* Determine whether the instance receives real reader data, simulator-generated data, or both.
* Determine whether Recorder can capture these reads into the SQLite RawReads table.
* Confirm whether any additional deployment or configuration is required.
* Document technical blockers, findings, evidence, and recommended next steps.

Acceptance Criteria:

* Existing Resonate Locate services and their status are documented.
* A simplified RawReads data flow is documented.
* The purpose and behaviour of sim.py are understood.
* It is confirmed whether RawReads can be obtained and captured by Recorder, or the reason they cannot be captured is documented.
* The need for deployment or configuration changes is clearly stated.
* All identified blockers and recommended next actions are documented.

Safety note:

No containers, services, configurations, or simulator processes should be changed or started without confirming their impact on the shared environment.