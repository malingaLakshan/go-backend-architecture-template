cd /opt/zebra/alt-resonate-locate/tnv/sim

~/resonate-sim-venv/bin/python - <<'PY'
import sys
sys.path.insert(0, "/opt/zebra/alt-resonate-locate/tnv/sim")

import paho.mqtt.client as mqtt
import messages_pb2
from google.protobuf import message_factory

SITE_ID = "3b96f652-8200-3920-8a2c-0486c358964e"
LOCATION_TOPIC = f"resonate/locate/{SITE_ID}/locationUpdate"

message_count = 0
total_locations = 0


def find_location_bundle():
    for descriptor in messages_pb2.DESCRIPTOR.message_types_by_name.values():
        for field in descriptor.fields:
            if (
                field.label == field.LABEL_REPEATED
                and field.message_type is not None
            ):
                child_fields = {
                    child.name for child in field.message_type.fields
                }

                if {"item_id", "x", "y", "reason"}.issubset(child_fields):
                    return (
                        message_factory.GetMessageClass(descriptor),
                        field.name
                    )

    raise RuntimeError("Location bundle type was not found")


LocationBundle, locations_field = find_location_bundle()


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
    print(f"Subscribed to: {LOCATION_TOPIC}")
    print(f"Detected Protobuf type: {LocationBundle.DESCRIPTOR.name}")
    client.subscribe(LOCATION_TOPIC)


def on_message(client, userdata, message):
    global message_count, total_locations

    bundle = LocationBundle()
    bundle.ParseFromString(message.payload)

    locations = getattr(bundle, locations_field)

    message_count += 1
    total_locations += len(locations)

    print("\n" + "=" * 75)
    print(f"MQTT MESSAGE #{message_count}")
    print(f"Topic: {message.topic}")
    print(f"Locations in message: {len(locations)}")

    for index, location in enumerate(locations, start=1):
        print(f"\nLOCATION #{index}")
        print_fields(location)

    print("\nRUNNING TOTAL")
    print("-" * 75)
    print("MQTT messages received:", message_count)
    print("Individual locations received:", total_locations)


client = mqtt.Client(
    callback_api_version=mqtt.CallbackAPIVersion.VERSION2,
    client_id="continuous-location-checker"
)

client.on_connect = on_connect
client.on_message = on_message
client.connect("127.0.0.1", 1883, 60)

try:
    client.loop_forever()
except KeyboardInterrupt:
    print("\nStopped by user.")
    print("Total MQTT messages:", message_count)
    print("Total individual locations:", total_locations)
finally:
    client.disconnect()
PY