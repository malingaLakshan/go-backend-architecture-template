cd /data/sim

./venv/bin/python3 - <<'PY'
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


def on_connect(client, userdata, flags, reason_code, properties=None):
    print("Connected:", reason_code)
    client.subscribe(RAW_TOPIC, qos=0)
    print("Subscribed to:", RAW_TOPIC)
    print("Continuously waiting for Raw RFID messages...")
    print("Press Ctrl+C whenever you want to stop.")


def on_message(client, userdata, mqtt_message):
    global message_count, total_reads

    bundle = messages_pb2.ProtoReaderBundle()

    try:
        bundle.ParseFromString(mqtt_message.payload)
    except Exception as error:
        print("Protobuf decoding failed:", error)
        return

    message_count += 1
    total_reads += len(bundle.reads)

    print("\n" + "=" * 75)
    print(f"RAW RFID MQTT MESSAGE #{message_count}")
    print("=" * 75)
    print("Topic:", mqtt_message.topic)
    print("QoS:", mqtt_message.qos)
    print("Payload size:", len(mqtt_message.payload), "bytes")
    print("Protobuf type: ProtoReaderBundle")

    print("\nALL BUNDLE FIELDS")
    print("-" * 75)

    populated_bundle = {
        field.name: value
        for field, value in bundle.ListFields()
    }

    for field in bundle.DESCRIPTOR.fields:
        if field.name == "reads":
            print(
                f"  Field {field.number} - reads: "
                f"{len(bundle.reads)} individual reads"
            )
        elif field.name in populated_bundle:
            print(
                f"  Field {field.number} - "
                f"{field.name}: {populated_bundle[field.name]}"
            )
        else:
            print(
                f"  Field {field.number} - "
                f"{field.name}: <NOT SENT OR DEFAULT>"
            )

    print("\nALL INDIVIDUAL READS")
    print("-" * 75)

    for index, read in enumerate(bundle.reads, start=1):
        print(f"\nREAD #{index}")
        print_fields(read)

    print("\nRUNNING TOTAL")
    print("-" * 75)
    print("MQTT messages received:", message_count)
    print("Individual reads received:", total_reads)


if hasattr(mqtt, "CallbackAPIVersion"):
    client = mqtt.Client(
        callback_api_version=mqtt.CallbackAPIVersion.VERSION2,
        client_id="continuous-rawrfid-checker"
    )
else:
    client = mqtt.Client(client_id="continuous-rawrfid-checker")

client.on_connect = on_connect
client.on_message = on_message

client.connect("127.0.0.1", 1883, 60)

try:
    client.loop_forever()

except KeyboardInterrupt:
    print("\n\nStopped by user.")
    print("Total MQTT messages:", message_count)
    print("Total individual reads:", total_reads)

finally:
    client.disconnect()
PY