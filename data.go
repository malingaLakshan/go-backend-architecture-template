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

expected_fields = {
    "item_id", "x", "y", "z", "floor", "floor_id",
    "region", "confidence", "state", "timestamp_ns", "reason"
}


def all_descriptors():
    pending = list(
        messages_pb2.DESCRIPTOR.message_types_by_name.values()
    )

    while pending:
        descriptor = pending.pop(0)
        yield descriptor
        pending.extend(descriptor.nested_types)


def find_location_bundle():
    best_match = None

    for descriptor in all_descriptors():
        for field in descriptor.fields:
            if (
                field.label == field.LABEL_REPEATED
                and field.message_type is not None
            ):
                child_names = {
                    child.name for child in field.message_type.fields
                }

                score = len(expected_fields.intersection(child_names))

                if best_match is None or score > best_match[0]:
                    best_match = (
                        score,
                        descriptor,
                        field.name,
                        child_names
                    )

    if best_match is None or best_match[0] < 3:
        raise RuntimeError("Location bundle type was not found")

    score, descriptor, field_name, child_names = best_match

    print(f"Detected bundle: {descriptor.name}")
    print(f"Locations field: {field_name}")
    print(f"Location fields: {sorted(child_names)}")

    return (
        message_factory.GetMessageClass(descriptor),
        field_name
    )


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
    print(f"Locations received: {len(locations)}")

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