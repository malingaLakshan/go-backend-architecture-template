Real Resonate MQTT Integration

For: Simulation Engine and Recorder Engine Developers

1. Confirmed Connection Details

Server:
117108-trirh901.117asd.zebra.lan
Site name:
Simulator Site 5260
Site ID:
3b96f652-8200-3920-8a2c-0486c358964e
MQTT port:
1883

Use the correct broker URL based on where the engine runs:

From an office/development machine:
tcp://117108-trirh901.117asd.zebra.lan:1883
Inside the Resonate Docker network:
tcp://mqtt:1883

Before connecting from an office machine, check access:

Test-NetConnection 117108-trirh901.117asd.zebra.lan -Port 1883

Continue only if:

TcpTestSucceeded : True

No MQTT username, password, or TLS configuration was confirmed.

2. Simulation Engine Configuration

This config connects the Simulation Engine to the real Resonate MQTT broker.

{
  "mqtt": {
    "brokerUrl": "tcp://117108-trirh901.117asd.zebra.lan:1883",
    "siteId": "3b96f652-8200-3920-8a2c-0486c358964e",
    "publishTopic": "resonate/locate/3b96f652-8200-3920-8a2c-0486c358964e/rawrfid",
    "qos": 0
  }
}

If the Simulation Engine runs inside the Resonate Docker network, change only:

"brokerUrl": "tcp://mqtt:1883"

The JSON key names can be mapped to the existing Simulation Engine configuration model. The values above must remain unchanged.

3. How the Simulation Engine Must Publish

The Simulation Engine must:

1. Create a ProtoReaderBundle.
2. Set the confirmed site ID.
3. Add the simulated Raw Reads.
4. Serialize the bundle as Protobuf binary.
5. Publish the binary bytes to the rawrfid topic using QoS 0.

The supplied simulator confirms this publishing method:

payload = m.SerializeToString()
mqtt_client.publish(mqtt_topic, payload, qos=0)

The payload must be published as raw Protobuf bytes.

Do not publish it as:

* JSON
* A normal string
* Base64 text

Reader, antenna, floor, and item values must match the selected site configuration. Do not use random identifiers.

4. Confirmed Raw Read Protobuf Model

Generated model on the Resonate server:

/data/sim/messages_pb2.py

Root message:

smartlens.messages.input.ProtoReaderBundle

The following schema was reconstructed from the generated Protobuf descriptor. It describes the payload structure, not actual payload values.

syntax = "proto3";
package smartlens.messages.input;
message ProtoRead {
  uint64 timestamp_ns = 1;
  uint32 confidence = 2;
  uint32 antenna_id = 3;
  uint32 antenna_type = 4;
  int32 x = 5;
  int32 y = 6;
  string item_id = 7;
  uint32 floor_id = 8;
}
message ProtoReaderBundle {
  uint32 reader_id = 1;
  repeated ProtoRead reads = 2;
  string site_id = 3;
  uint64 sent_timestamp_ms = 4;
}

One ProtoReaderBundle can contain multiple ProtoRead records.

5. Recorder Engine Configuration

{
  "mqtt": {
    "brokerUrl": "tcp://117108-trirh901.117asd.zebra.lan:1883",
    "siteId": "3b96f652-8200-3920-8a2c-0486c358964e",
    "subscriptions": [
      {
        "topic": "resonate/locate/3b96f652-8200-3920-8a2c-0486c358964e/rawrfid",
        "qos": 0,
        "payloadType": "ProtoReaderBundle"
      },
      {
        "topic": "resonate/locate/3b96f652-8200-3920-8a2c-0486c358964e/locationUpdate",
        "qos": 1,
        "payloadType": "Pending confirmation"
      }
    ]
  }
}

If the Recorder runs inside the Resonate Docker network, change only:

"brokerUrl": "tcp://mqtt:1883"

Each engine must use its own unique MQTT client ID. Do not reuse these existing Resonate client IDs:

smartlens-ingest-1
location-processor

6. How the Recorder Must Receive Raw Reads

The Recorder must:

1. Connect to the configured MQTT broker.
2. Subscribe to the rawrfid topic using QoS 0.
3. Receive the binary payload.
4. Deserialize it as ProtoReaderBundle.
5. Extract all records from the repeated reads field.
6. Save the extracted Raw Reads to SQLite.

The Raw Read payload must not be parsed as JSON.

The Recorder may subscribe to the locationUpdate topic using QoS 1, but its MQTT payload format must be confirmed before decoding it.

7. Site and Asset Endpoints

Get the available site

GET http://117108-trirh901.117asd.zebra.lan:3389/api/locations

Real response received:

[
  {
    "friendlyName": "Simulator Site 5260",
    "uniqueIdentifier": "3b96f652-8200-3920-8a2c-0486c358964e"
  }
]

Get site assets

GET http://117108-trirh901.117asd.zebra.lan:3389/api/assets/3b96f652-8200-3920-8a2c-0486c358964e

This endpoint returned 200 OK with JSON. The response was approximately 8.8 MB, so the full body is not included.

The Simulation Engine should use the site asset information when mapping valid reader, antenna, and floor values.

8. Observed Location JSON

After Raw Reads were published, Resonate generated Locations. This real object was observed through the Visualizer WebSocket:

{
  "id": "3055873B3B61BEDEE00AB6CE",
  "site": "3b96f652-8200-3920-8a2c-0486c358964e",
  "r": "sim-5260-world",
  "s": "VALID",
  "x": 5982.0,
  "y": 5316.0
}

This is Visualizer WebSocket JSON. It does not confirm that the MQTT locationUpdate payload uses the same JSON format.

9. Still Pending

* Capture one serialized Raw Read MQTT payload sample.
* Confirm the MQTT locationUpdate payload format.
* Identify the Event MQTT topic and payload format.
* Confirm external access to MQTT port 1883 from the developers’ machines.