Real Resonate MQTT Investigation and Recorder Connection

Date: 12 August 2026

1. Purpose

The purpose of this investigation was to find how the Recorder can receive data from the real Resonate Locate instance.

The Recorder must connect directly to the Resonate MQTT broker. It should not use the Visualizer WebSocket or REST API.

2. Confirmed Test Environment

* Resonate server:

117108-trirh901.117asd.zebra.lan

* MQTT port:

1883

* Site ID:

3b96f652-8200-3920-8a2c-0486c358964e

* MQTT broker: EMQX 5.1
* Supplied simulator location:

/data/sim/sim.py

3. Confirmed Data Flow

Real RFID Readers or Simulator
              ↓
         MQTT Broker
          ↙       ↘
     Resonate     Recorder
         ↓           ↓
    Locations      SQLite

The simulator is used only to generate test Raw Reads. In the real environment, RFID readers will publish the Raw Reads.

MQTT can deliver the same message to multiple subscribers. Therefore, Resonate and Recorder can receive the same Raw Read messages.

4. Confirmed MQTT Topics

Raw Read stream

resonate/locate/3b96f652-8200-3920-8a2c-0486c358964e/rawrfid

* QoS: 0
* Payload: Protobuf/binary
* Existing Resonate subscriber: smartlens-ingest-1

This is the main topic the Recorder must subscribe to for Raw Reads.

Location stream

resonate/locate/3b96f652-8200-3920-8a2c-0486c358964e/locationUpdate

* QoS: 1
* Existing subscriber: location-processor

This MQTT topic was identified from the active broker subscriptions. Its MQTT payload has not yet been captured and decoded.

Event stream

The Event MQTT topic and payload have not yet been confirmed. It should not be added to the Recorder configuration until it is identified.

5. Recorder MQTT URLs

Use the URL that matches where the Recorder runs.

Recorder location	MQTT broker URL
Office machine connected to the Zebra network	tcp://117108-trirh901.117asd.zebra.lan:1883
Container inside the same Resonate Docker network	tcp://mqtt:1883
Directly on the Resonate Linux server	tcp://localhost:1883

The confirmed internal Resonate services use:

mqtt:1883

Before using the office-machine URL, confirm that port 1883 is reachable.

6. Recorder config.json

The following configuration contains the confirmed server, site and topic values:

{
  "mqtt": {
    "brokerUrl": "tcp://117108-trirh901.117asd.zebra.lan:1883",
    "siteId": "3b96f652-8200-3920-8a2c-0486c358964e",
    "rawReadsTopic": "resonate/locate/3b96f652-8200-3920-8a2c-0486c358964e/rawrfid",
    "rawReadsQoS": 0,
    "locationTopic": "resonate/locate/3b96f652-8200-3920-8a2c-0486c358964e/locationUpdate",
    "locationQoS": 1
  }
}

The JSON field names must be matched with the existing Recorder configuration model.

The Recorder must use its own unique MQTT client ID. It must not use the existing smartlens-ingest-1 or location-processor client IDs.

No MQTT username, password or TLS configuration was confirmed during this investigation. These values must not be invented.

7. How Recorder Receives Raw Reads

The Recorder must:

1. Connect to the MQTT broker.
2. Subscribe to the confirmed rawrfid topic using QoS 0.
3. Receive the Protobuf payload bytes.
4. Decode the payload using the matching Resonate Raw RFID Protobuf schema.
5. Extract the individual Raw Reads.
6. Save them in the SQLite RawReads table.

One MQTT message can contain a bundle of multiple Raw Reads.

8. Simple Test Commands

Check the MQTT broker

Run on the Resonate server:

sudo docker exec mqtt emqx ctl status

Expected result:

Node 'emqx@172.18.0.8'. 5.1 is started

View existing MQTT subscriptions

sudo docker exec mqtt emqx ctl subscriptions list

This should show the rawrfid and locationUpdate subscriptions.

Record the Raw Read counter

sudo docker exec mqtt emqx ctl clients show smartlens-ingest-1

Record the current delivered_msgs value.

During the investigation, the observed baseline was:

connected=true
delivered_msgs=20041

Run the supplied simulator

cd /data/sim
./venv/bin/python3 sim.py

Let it run for approximately 10 seconds and stop it using:

Ctrl + C

Run the counter command again:

sudo docker exec mqtt emqx ctl clients show smartlens-ingest-1

If delivered_msgs is greater than the first value, it proves:

Simulator published Raw Reads
→ MQTT broker received them
→ MQTT delivered them to Resonate

Test broker access from the Windows office machine

Run in PowerShell:

Test-NetConnection 117108-trirh901.117asd.zebra.lan -Port 1883

Expected result:

TcpTestSucceeded : True

If it returns False, the Recorder cannot connect from that machine until network access to port 1883 is allowed.

9. Confirmed Findings

* The Resonate MQTT broker is running on port 1883.
* The simulator successfully published Raw Read bundles.
* The exact Raw Read MQTT topic was identified.
* Resonate is subscribed to the Raw Read topic.
* After the simulator warm-up, Resonate produced live Location results.
* The exact Location MQTT topic was identified.
* The Visualizer displays Locations through WebSocket /ws.
* WebSocket /ws is only for the Visualizer and should not be used by Recorder.
* The Event MQTT topic is still pending.
* The Recorder must subscribe directly to MQTT and save the decoded data into SQLite.

10. Final Conclusion

The Recorder should access real Raw Reads using:

Broker:
tcp://117108-trirh901.117asd.zebra.lan:1883
Topic:
resonate/locate/3b96f652-8200-3920-8a2c-0486c358964e/rawrfid
QoS:
0

This is the confirmed connection path for receiving Raw Reads from the real Resonate test instance. The final implementation still requires the matching Protobuf schema to decode the received Raw Read bundles.