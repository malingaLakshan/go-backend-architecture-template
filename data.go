Resonate Instance Investigation – Handover Notes

1. Investigation purpose

The purpose was to verify:

* How the provided simulator sends RFID data.
* Which MQTT topics are available.
* What payload types are used.
* How Resonate generates location data.
* How the Recorder should receive and store this data.
* Which HTTP APIs are available.

2. Test environment

Item	Value
Resonate host	117108-trirh901.117asd.zebra.lan
HTTP API port	3389
MQTT port	1883
MQTT broker	EMQX 5.5.1
Site name	Simulator Site 5260
Site ID	3b96f652-8200-3920-8a2c-0486c358964e

When running commands inside the Resonate machine, the MQTT broker can be accessed using:

127.0.0.1:1883

3. Confirmed data flow

sim.py
  ↓
Publishes raw RFID bundles
  ↓
rawrfid MQTT topic
  ↓
Resonate processes the raw reads
  ↓
Publishes calculated locations
  ↓
locationUpdate MQTT topic
  ↓
Recorder decodes the payloads
  ↓
SQLite tables

Important:

* MQTT carries binary messages.
* MQTT messages do not look like SQLite rows.
* The Recorder must decode and split the messages before saving them.
* One MQTT bundle can contain multiple RFID reads.
* Raw reads and location updates are not necessarily generated one-to-one.

4. SiteGraph API

The following endpoint was confirmed:

GET http://117108-trirh901.117asd.zebra.lan:3389/api/sitegraph/3b96f652-8200-3920-8a2c-0486c358964e

Result:

* HTTP status: 200 OK
* Response type: JSON
* Site name: Simulator Site 5260
* Correct Site ID field: uniqueIdentifier

Important Site ID finding

The correct Site ID is:

3b96f652-8200-3920-8a2c-0486c358964e

The following value is the MongoDB document ID and must not be used as the Site ID:

6a70ef3aa10865c7e0194ce5

The endpoint below returns 404:

GET http://117108-trirh901.117asd.zebra.lan:3389/api/sitegraph

Therefore, only the SiteGraph details endpoint was confirmed. An API for listing all sites was not found.

/api/locations is not a Site listing API.

5. Confirmed MQTT topics

Raw RFID topic

resonate/locate/3b96f652-8200-3920-8a2c-0486c358964e/rawrfid

Property	Result
Topic status	Confirmed working
Observed QoS	0
Transport format	Binary
Serialization	Protocol Buffers
Protobuf message type	ProtoReaderBundle
Publisher	Provided simulator

A sample was successfully decoded using:

bundle = messages_pb2.ProtoReaderBundle()
bundle.ParseFromString(message.payload)

The decoded sample included:

* Site ID
* Reader ID
* RFID reads
* Sent timestamp

Location update topic

resonate/locate/3b96f652-8200-3920-8a2c-0486c358964e/locationUpdate

Property	Result
Topic status	Confirmed working
Observed QoS	0
Sample payload size	329 bytes
Transport format	Binary
Exact message type	Not identified yet
Publisher	Resonate

The payload could not be decoded as JSON or normal UTF-8 text.

The HEX and Base64 values shown in the terminal are not the actual data type. They are only printable representations of the unknown binary payload.

The correct Protobuf definition or generated message class is still required to decode the location fields.

Event topic

The configuration currently contains:

mock://streams/events

A real Resonate event MQTT topic and its payload type were not confirmed during this investigation.

6. Simulator behaviour

The simulator was run from:

cd /data/sim
./venv/bin/python3 sim.py

The simulator has two main stages:

WARMUP stage

* Raw reads are published.
* Moving tag count may remain at zero.
* Location updates may not be generated immediately.

LIVE stage

* Moving tags become active.
* The simulator continues publishing raw RFID bundles.
* Resonate processes the reads and produces location updates.

One observed run showed approximately:

100 bundles per second
1,900–2,100 total RFID reads per second

This confirms that a bundle and an individual read are different things.

7. Why the subscriber showed only one message

The temporary subscriber was created only to capture a sample payload.

One version ignored later messages using:

if message.topic in seen:
    return

The location test stopped after the first message using:

client.disconnect()

Therefore, the one-message output does not mean that only one message was published. The script intentionally captured only the first message.

A continuous subscriber is required to measure the real message count.

8. MQTT subscription details

The topics can be checked from the Resonate server using:

sudo docker exec mqtt emqx ctl topics list

A Paho MQTT client can subscribe using:

import paho.mqtt.client as mqtt
raw_topic = "resonate/locate/3b96f652-8200-3920-8a2c-0486c358964e/rawrfid"
location_topic = "resonate/locate/3b96f652-8200-3920-8a2c-0486c358964e/locationUpdate"
def on_connect(client, userdata, flags, rc):
    client.subscribe(raw_topic, qos=0)
    client.subscribe(location_topic, qos=0)
def on_message(client, userdata, message):
    print(message.topic, len(message.payload), message.qos)
client = mqtt.Client()
client.on_connect = on_connect
client.on_message = on_message
client.connect("127.0.0.1", 1883, 60)
client.loop_forever()

The subscriber should be started before starting the simulator.

QoS 0 is best-effort delivery. It does not provide acknowledgement or retry, so the subscriber should remain connected while the simulator is running.

9. MQTT messages compared with Recorder tables

The MQTT payloads are the input to the Recorder. They are not already formatted database rows.

Expected mapping:

MQTT stream	Recorder table
rawrfid	RawReads
locationUpdate	MLT_SOW_Locations
Event stream	ResonateEvents

After decoding ProtoReaderBundle, the Recorder can create individual rows in the RawReads table.

After decoding the location message, the Recorder should create rows in the MLT_SOW_Locations table.

These database screenshots show the expected stored structure. They are different from the original MQTT binary payloads.

10. Recorder startup issue

The Recorder currently fails while retrieving site information:

start recording: retrieve site information: site ID ... not found

However:

* The configured uniqueIdentifier is correct.
* The same SiteGraph URL works in the browser.
* The MQTT topics also work independently.

This indicates a possible mismatch between the Recorder’s expected Site API response and the SiteGraph API exposed by this Resonate instance.

The Recorder fails before completing the recording startup process. Therefore, live MQTT-to-SQLite persistence has not yet been confirmed.

11. Final investigation status

Confirmed

* Resonate Docker services are running.
* SiteGraph details API works.
* Correct Site ID was identified.
* MQTT broker is available on port 1883.
* Simulator publishes a high volume of raw RFID bundles.
* rawrfid topic works.
* Raw payload is ProtoReaderBundle.
* Resonate consumes raw reads.
* locationUpdate topic works.
* Location payload is binary.
* MQTT messages are bundles and are different from database rows.

Still required

* Identify the exact location payload message type/schema.
* Confirm the real event topic and payload type.
* Confirm the correct Site listing API, if one exists.
* Resolve the Recorder site-information lookup issue.
* Run the Recorder successfully and verify live rows in SQLite.