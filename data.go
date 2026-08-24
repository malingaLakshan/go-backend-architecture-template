cd /opt/zebra/alt-resonate-locate/tnv/sim

~/resonate-sim-venv/bin/python - <<'PY'
import sys
sys.path.insert(0, "/opt/zebra/alt-resonate-locate/tnv/sim")

import paho.mqtt.client as mqtt
import messages_pb2

SITE_ID = "3b96f652-8200-3920-8a2c-0486c358964e"
RAW_TOPIC = f"resonate/locate/{SITE_ID}/rawrfid"

message_count = 0
total_reads = 0


def print_fields(protobuf_message):
    populated = {
        field.name: value
        for field, value in protobuf_message.ListFields()
    }

    for field in protobuf_message.DESCRIPTOR.fields:
        if field.name in populated:
            print(
                f"  Field {field.number} - "
                f"{field.name}: {populated[field.name]}"
            )
        else:
            print(
                f"  Field {field.number} - "
                f"{field.name}: <NOT SENT OR DEFAULT>"
            )


def on_connect(client, userdata, flags, reason_code, properties):
    print(f"Connected to MQTT: {reason_code}")
    print(f"Subscribed to: {RAW_TOPIC}")
    client.subscribe(RAW_TOPIC)


def on_message(client, userdata, message):
    global message_count, total_reads

    bundle = messages_pb2.ProtoReaderBundle()
    bundle.ParseFromString(message.payload)

    message_count += 1
    total_reads += len(bundle.reads)

    print("\n" + "=" * 75)
    print(f"MQTT MESSAGE #{message_count}")
    print(f"Reader ID: {bundle.reader_id}")
    print(f"Site ID: {bundle.site_id}")
    print(f"Reads in bundle: {len(bundle.reads)}")

    for index, read in enumerate(bundle.reads, start=1):
        print(f"\nREAD #{index}")
        print_fields(read)

    print("\nRUNNING TOTAL")
    print("-" * 75)
    print("MQTT messages received:", message_count)
    print("Individual reads received:", total_reads)


client = mqtt.Client(
    callback_api_version=mqtt.CallbackAPIVersion.VERSION2,
    client_id="continuous-rawrfid-checker"
)

client.on_connect = on_connect
client.on_message = on_message
client.connect("127.0.0.1", 1883, 60)

try:
    client.loop_forever()
except KeyboardInterrupt:
    print("\nStopped by user.")
    print("Total MQTT messages:", message_count)
    print("Total individual reads:", total_reads)
finally:
    client.disconnect()
PY