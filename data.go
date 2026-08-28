import json
from datetime import datetime

import paho.mqtt.client as mqtt

BROKER = "127.0.0.1"
PORT = 1883
TOPIC = "resonate/locate/3b96f652-8200-3920-8a2c-0486c358964e/events/json"

batch_count = 0
event_count = 0


def on_connect(client, userdata, flags, reason_code, properties):
    if reason_code == 0:
        print(f"Connected to MQTT broker: {BROKER}:{PORT}")
        print(f"Subscribed to: {TOPIC}")
        print("Waiting for Events...\n")
        client.subscribe(TOPIC)
    else:
        print(f"MQTT connection failed: {reason_code}")


def on_message(client, userdata, message):
    global batch_count, event_count

    try:
        payload = json.loads(message.payload.decode("utf-8"))
        records = payload if isinstance(payload, list) else [payload]

        batch_count += 1

        print("\n" + "=" * 70)
        print(
            f"BATCH {batch_count} | "
            f"Received: {datetime.now().isoformat(timespec='seconds')} | "
            f"Records: {len(records)}"
        )
        print("=" * 70)

        for record in records:
            event_count += 1

            print(f"\nEVENT RECORD {event_count}")
            print(f"ID:         {record.get('id')}")
            print(f"Type:       {record.get('type')}")
            print(f"Events:     {record.get('events')}")
            print(f"Site:       {record.get('site')}")
            print(f"Floor:      {record.get('floor')}")
            print(f"Region:     {record.get('region')}")
            print(f"Position:   x={record.get('x')}, y={record.get('y')}, z={record.get('z')}")
            print(f"Timestamp:  {record.get('timestamp')}")
            print("-" * 70)

    except Exception as error:
        print(f"Failed to process MQTT message: {error}")
        print(message.payload.decode("utf-8", errors="replace"))


client = mqtt.Client(mqtt.CallbackAPIVersion.VERSION2)
client.on_connect = on_connect
client.on_message = on_message

client.connect(BROKER, PORT, 60)
client.loop_forever()