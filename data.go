Create a short Word document titled:

“JSON Event Stream Investigation – Current Status”

Use simple English and keep it to 1–2 pages. Cover only the Event stream.

Include these sections:

1. Event Stream Workflow

The application service writes asset updates into the MongoDB assets collection.

The datafeed-service watches those MongoDB changes and generates business events such as:

* ARRIVAL
* POSITION_CHANGE
* REGION_CHANGE
* FLOOR_CHANGE
* DEPARTURE
* EXIT
* GHOST
* MISSING

The workflow is:

Application Service → MongoDB Assets → Datafeed Service → JSON Events → Configured Output

2. Datafeed and Destination Topic

There is no default MQTT topic for Events.

A datafeed must be configured. The datafeed defines:

* Which events should be included
* Where the generated JSON should be sent
* The destination MQTT broker and topic

We configured this test destination topic:

resonate/locate/3b96f652-8200-3920-8a2c-0486c358964e/events/json

3. What We Tested

The feed was created and became ACTIVE.

The simulator was running and MongoDB asset data was changing.

The datafeed-service generated valid JSON Event arrays containing fields such as:

* id
* type
* productId
* site
* floor
* confidence
* region
* timestamp
* events
* x, y and z

However, subscribing to the configured MQTT topic returned no messages.

4. What We Found

The datafeed logs showed:

MOCK TRANSPORT: MQTT

This means the service generates the JSON Events but only prints them in the container logs. It does not actually publish them to the MQTT broker.

Therefore:

* Event generation is working.
* JSON Event payload creation is working.
* The configured feed is ACTIVE.
* Real MQTT publishing is not working.
* The Events can currently be viewed only through the datafeed-service logs.

5. Current State

The Event stream logic works inside datafeed-service, but there is currently no usable MQTT Event stream because the deployed service uses mock transport.

The dashboard WebSocket and /api/assets response contain current location or asset-state data. They are not the business Event stream produced by datafeed-service.

6. Information Required From the Service Team

Ask the following questions:

* Is mock transport intentional in this deployed version?
* How can we enable real MQTT transport?
* Is a configuration setting, environment variable, or different Docker image required?
* Is there a recommended standard MQTT topic for Events?
* If MQTT is unavailable, is another supported Event output available, such as HTTP, WebSocket or file output?
* How should JSON Events containing multiple values in events[] be stored in the ResonateEvents table?

Finish with this conclusion:

“Business Event generation has been verified successfully inside datafeed-service. The remaining blocker is Event delivery because the current service uses mock MQTT transport instead of publishing the generated JSON to the configured MQTT topic.”