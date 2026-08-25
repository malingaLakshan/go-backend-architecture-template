KT Guide — Using the Resonate JSON Event Stream

1. Purpose

Resonate does not provide a default MQTT topic for Events.

To receive Events:

1. Create a datafeed.
2. Configure an MQTT destination topic in that datafeed.
3. Start the datafeed.
4. Run the simulator to generate RFID movement.
5. Subscribe to the configured topic and receive JSON Events.

The datafeed monitors asset changes in MongoDB, generates business Events, and publishes them as JSON arrays.

⸻

2. Test environment

VM:

117108-trirh902.117asd.zebra.lan

Site ID:

3b96f652-8200-3920-8a2c-0486c358964e

Feed ID:

lab-json-events

Configured Event topic:

resonate/locate/3b96f652-8200-3920-8a2c-0486c358964e/events/json

MQTT port:

1883

Important:

* The Event topic is configurable; it is not a default Resonate topic.
* The topic does not need to be created separately in MQTT.
* It becomes available when the datafeed publishes its first message.

⸻

One-time setup

These steps are required only once for each Linux user account on the VM.

3. Connect to the new VM

Open a terminal and run:

ssh ML9352@117108-trirh902.117asd.zebra.lan

Enter the VM password.

4. Create a Python environment

python3 -m venv ~/resonate-sim-venv

Install the required packages:

~/resonate-sim-venv/bin/python -m pip install protobuf==6.31.1 paho-mqtt

Verify the packages:

~/resonate-sim-venv/bin/python -c "import google.protobuf; import paho.mqtt.client; print('Dependencies OK')"

Expected result:

Dependencies OK

These packages fixed the previous simulator errors:

No module named 'google'
No module named 'paho'

Do not use the existing venv folder under the simulator directory. Its Python link may point to another user account and can return Permission denied.

Use this environment instead:

~/resonate-sim-venv

⸻

Demonstration steps

Use three separate terminals.

Terminal 1 — Check Resonate services

Connect to the VM:

ssh ML9352@117108-trirh902.117asd.zebra.lan

Go to the Docker Compose directory:

cd /opt/zebra/alt-resonate-locate/svc

Check the services:

sudo docker compose ps

Confirm that the important services are running:

datafeed
location
mongodb
mqtt
processing-engine
api

⸻

Terminal 2 — Subscribe to the Event topic first

Open another terminal and connect to the VM:

ssh ML9352@117108-trirh902.117asd.zebra.lan

The mosquitto_sub command is not installed on this VM. Use the following Python subscriber instead:

~/resonate-sim-venv/bin/python -c 'import json; import paho.mqtt.client as mqtt; topic="resonate/locate/3b96f652-8200-3920-8a2c-0486c358964e/events/json"; client=mqtt.Client(mqtt.CallbackAPIVersion.VERSION2); client.on_connect=lambda c,u,f,r,p: (print("Connected to MQTT. Waiting for Events...",flush=True),c.subscribe(topic)); client.on_message=lambda c,u,m: print(json.dumps(json.loads(m.payload.decode("utf-8")),indent=2),flush=True); client.connect("127.0.0.1",1883,60); client.loop_forever()'

Keep this terminal open.

It may initially show only:

Connected to MQTT. Waiting for Events...

That is normal. Events will appear after the simulator creates asset movement.

Stop the subscriber using:

Ctrl+C

⸻

Terminal 3 — Check or create the datafeed

Open another terminal and connect to the VM:

ssh ML9352@117108-trirh902.117asd.zebra.lan

Check whether the feed already exists

curl -sS "http://localhost:3000/feed/siteId('3b96f652-8200-3920-8a2c-0486c358964e')"

Look for:

lab-json-events

If it already exists, do not create it again. Continue to the “Start the datafeed” step.

Create the datafeed only if it does not exist

curl -sS -w '\nHTTP %{http_code}\n' -X POST \
"http://localhost:3000/feed/siteId('3b96f652-8200-3920-8a2c-0486c358964e')" \
-H "Content-Type: application/json" \
-d '{
  "feedId": "lab-json-events",
  "eventFilters": [],
  "filterMode": "STRICT",
  "destination": {
    "protocol": "MQTT",
    "broker": "mqtt:1883",
    "topic": "resonate/locate/3b96f652-8200-3920-8a2c-0486c358964e/events/json"
  }
}'

Expected result:

HTTP 200

The newly created feed will normally have the status:

PAUSED

If the result is:

HTTP 409

the feed already exists. Do not create another one. Start the existing feed.

Start the datafeed

curl -sS -w '\nHTTP %{http_code}\n' -X POST \
"http://localhost:3000/feed/siteId('3b96f652-8200-3920-8a2c-0486c358964e')/lab-json-events/Start"

Expected result:

HTTP 200

⸻

Terminal 4 — Run the simulator

Open another terminal and connect to the VM:

ssh ML9352@117108-trirh902.117asd.zebra.lan

Go to the simulator directory:

cd /opt/zebra/alt-resonate-locate/tnv/sim

Run the simulator using the Python environment created earlier:

~/resonate-sim-venv/bin/python sim-mqtt.py

Expected output includes:

Simulator Online (MQTT mode)
Connecting to MQTT Broker at 127.0.0.1:1883
Publishing on Topic: resonate/locate/.../rawrfid

The simulator may first show:

WARMUP

Keep it running. Events may take some time to appear while assets are created and begin moving.

After the warm-up, the output changes to:

LIVE

Stop the simulator using:

Ctrl+C

Do not run two simulator instances at the same time.

⸻

5. Confirm the Event output

Return to the Event subscriber terminal.

You should receive JSON arrays similar to:

[
  {
    "id": "30D113479A249C969A24DE4",
    "type": "ITEM",
    "productId": "2864761863432",
    "site": "3b96f652-8200-3920-8a2c-0486c358964e",
    "floor": 0,
    "confidence": 100,
    "region": "region-id",
    "timestamp": "2026-08-21T06:17:29Z",
    "events": [
      "DEPARTURE",
      "EXIT",
      "POSITION_CHANGE",
      "REGION_CHANGE"
    ],
    "x": 6432.0,
    "y": 864.0,
    "z": 0.0
  }
]

Important payload details:

* Each MQTT message contains a JSON array.
* One array can contain multiple Event records.
* One record can contain multiple values inside the events array.

Events observed during testing include:

ARRIVAL
POSITION_CHANGE
REGION_CHANGE
FLOOR_CHANGE
DEPARTURE
EXIT

Common fields include:

id
type
productId
site
floor
confidence
region
timestamp
events
x
y
z

⸻

6. Confirm that real MQTT publishing is being used

Open another terminal:

ssh ML9352@117108-trirh902.117asd.zebra.lan

Run:

cd /opt/zebra/alt-resonate-locate/svc

Then:

sudo docker compose logs --since 5m datafeed 2>&1 | grep -E "Successfully published|MOCK TRANSPORT"

Correct result:

Successfully published ... to MQTT

Incorrect result:

MOCK TRANSPORT

The new VM should use real MQTT publishing. MOCK TRANSPORT was the problem with the old datafeed image.

⸻

7. Validator integration guidance

For the current implementation, the Validator should not hard-code the Event topic.

It should accept the following values through an input parameter or configuration file:

MQTT host
MQTT port
Event topic

Example configuration:

{
  "mqttHost": "117108-trirh902.117asd.zebra.lan",
  "mqttPort": 1883,
  "eventTopic": "resonate/locate/3b96f652-8200-3920-8a2c-0486c358964e/events/json"
}

Connection address depends on where the Validator runs:

* Validator running outside the VM:

117108-trirh902.117asd.zebra.lan:1883

* Validator running directly on the VM:

127.0.0.1:1883

* Validator running as a container on the same Docker network:

mqtt:1883

The Validator should:

1. Read the MQTT host, port and Event topic from configuration.
2. Connect to the MQTT broker.
3. Subscribe to the configured Event topic.
4. Decode each MQTT payload as a JSON array.
5. Iterate through every Event record in the array.
6. Validate the required fields and Event types.
7. Record any missing or invalid values.

For this KT, the datafeed is prepared manually in the test environment. The Validator only needs to consume the configured topic unless the updated requirement explicitly assigns datafeed creation and lifecycle management to the Validator.

⸻

8. What must be repeated after stopping work

You do not need to create the datafeed every day.

The datafeed configuration is stored on the Resonate VM, not on the laptop.

At the beginning of another test:

1. Check that the Docker services are running.
2. Check that lab-json-events still exists.
3. Start it if its status is PAUSED.
4. Start the Event subscriber.
5. Run the simulator.

The following must be started again after they are stopped:

* Event subscriber
* Simulator

If the VM or datafeed container was restarted, verify the feed status and start it again if required.

⸻

9. Common problems

mosquitto_sub: command not found

Use the Python subscriber provided in this guide.

No module named 'google'

Run:

~/resonate-sim-venv/bin/python -m pip install protobuf==6.31.1

No module named 'paho'

Run:

~/resonate-sim-venv/bin/python -m pip install paho-mqtt

Permission denied when using ./venv/bin/python

Do not use the simulator’s existing venv.

Use:

~/resonate-sim-venv/bin/python

Datafeed creation returns HTTP 409

The feed already exists. Do not recreate it. Start the existing feed.

Subscriber connects but displays nothing

Check the following:

1. Datafeed has been started.
2. Simulator is running.
3. Simulator has completed enough of its warm-up period.
4. Subscriber is using the exact configured Event topic.
5. Datafeed logs show Successfully published ... to MQTT.
6. Datafeed logs do not show MOCK TRANSPORT.

⸻

10. Final confirmed workflow

Simulator
→ RawRFID MQTT
→ Resonate location processing
→ application-service
→ MongoDB asset changes
→ datafeed-service
→ configured Event MQTT topic
→ Validator

This method has been tested successfully on the new Resonate VM, and real JSON Events were received through MQTT.