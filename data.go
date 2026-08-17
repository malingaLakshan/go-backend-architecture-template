cd /data/sim

./venv/bin/python3 - <<'PY'
import paho.mqtt.client as mqtt
import messages_pb2

SITE_ID = "3b96f652-8200-3920-8a2c-0486c358964e"
RAW_TOPIC = f"resonate/locate/{SITE_ID}/rawrfid"


def print_all_fields(message):
    present_fields = {
        field.name: value
        for field, value in message.ListFields()
    }

    for field in message.DESCRIPTOR.fields:
        if field.name in present_fields:
            value = present_fields[field.name]
            print(f"  Field {field.number} - {field.name}: {value}")
        else:
            print(
                f"  Field {field.number} - {field.name}: "
                "<NOT SENT OR DEFAULT>"
            )


def on_connect(client, userdata, flags, reason_code, properties=None):
    print("Connected:", reason_code)
    client.subscribe(RAW_TOPIC, qos=0)
    print("Subscribed to:", RAW_TOPIC)
    print("Waiting for one complete Raw RFID message...")


def on_message(client, userdata, mqtt_message):
    bundle = messages_pb2.ProtoReaderBundle()

    try:
        bundle.ParseFromString(mqtt_message.payload)
    except Exception as error:
        print("Protobuf decoding failed:", error)
        client.disconnect()
        return

    print("\n" + "=" * 75)
    print("MQTT INFORMATION")
    print("=" * 75)
    print("Topic:", mqtt_message.topic)
    print("QoS:", mqtt_message.qos)
    print("Payload size:", len(mqtt_message.payload), "bytes")
    print("Protobuf type: ProtoReaderBundle")

    print("\n" + "=" * 75)
    print("ALL BUNDLE FIELDS")
    print("=" * 75)

    present_bundle_fields = {
        field.name: value
        for field, value in bundle.ListFields()
    }

    for field in bundle.DESCRIPTOR.fields:
        if field.name == "reads":
            print(
                f"  Field {field.number} - reads: "
                f"{len(bundle.reads)} individual reads"
            )
        elif field.name in present_bundle_fields:
            print(
                f"  Field {field.number} - {field.name}: "
                f"{present_bundle_fields[field.name]}"
            )
        else:
            print(
                f"  Field {field.number} - {field.name}: "
                "<NOT SENT OR DEFAULT>"
            )

    print("\n" + "=" * 75)
    print("ALL INDIVIDUAL READS")
    print("=" * 75)

    for index, read in enumerate(bundle.reads, start=1):
        print(f"\nREAD #{index}")
        print("-" * 40)
        print_all_fields(read)

    print("\n" + "=" * 75)
    print("SUMMARY")
    print("=" * 75)
    print("Bundle fields defined:", len(bundle.DESCRIPTOR.fields))
    print("Individual reads received:", len(bundle.reads))
    print(
        "Fields defined for each read:",
        len(bundle.DESCRIPTOR.fields_by_name["reads"].message_type.fields)
    )

    client.disconnect()


if hasattr(mqtt, "CallbackAPIVersion"):
    client = mqtt.Client(
        callback_api_version=mqtt.CallbackAPIVersion.VERSION2,
        client_id="complete-rawrfid-checker"
    )
else:
    client = mqtt.Client(client_id="complete-rawrfid-checker")

client.on_connect = on_connect
client.on_message = on_message

client.connect("127.0.0.1", 1883, 60)
client.loop_forever()
PY