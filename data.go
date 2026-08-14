cd /data/sim

./venv/bin/python3 - <<'PY'
import json
import base64

import paho.mqtt.client as mqtt
import messages_pb2
from google.protobuf import text_format

SITE_ID = "3b96f652-8200-3920-8a2c-0486c358964e"

RAW_TOPIC = f"resonate/locate/{SITE_ID}/rawrfid"
LOCATION_TOPIC = f"resonate/locate/{SITE_ID}/locationUpdate"

DETAIL_LIMIT = 3

raw_message_count = 0
location_message_count = 0


# Find available Protobuf message types
proto_types = messages_pb2.DESCRIPTOR.message_types_by_name

location_type_names = [
    name for name in proto_types
    if "location" in name.lower()
    or "position" in name.lower()
    or "coordinate" in name.lower()
]

print("AVAILABLE PROTOBUF TYPES:")
for name in proto_types:
    print(" -", name)

print("\nLOCATION PROTOBUF CANDIDATES:")
if location_type_names:
    for name in location_type_names:
        print(" -", name)
else:
    print(" - No location type found in messages_pb2.py")


def collect_repeated_counts(message, prefix="", totals=None):
    """Find the number of items inside every repeated Protobuf field."""

    if totals is None:
        totals = {}

    for field, value in message.ListFields():
        path = f"{prefix}.{field.name}" if prefix else field.name

        if field.label == field.LABEL_REPEATED:
            totals[path] = totals.get(path, 0) + len(value)

            if field.type == field.TYPE_MESSAGE:
                for item in value:
                    collect_repeated_counts(
                        item,
                        path + "[]",
                        totals
                    )

        elif field.type == field.TYPE_MESSAGE:
            collect_repeated_counts(value, path, totals)

    return totals


def print_repeated_counts(message):
    counts = collect_repeated_counts(message)

    print("\nREPEATED FIELD COUNTS:")

    if not counts:
        print(" - No repeated fields found")
        return

    for path, count in counts.items():
        print(f" - {path}: {count}")

    read_fields = []

    for path, count in counts.items():
        field_name = path.split(".")[-1].replace("[]", "")
        normalized = field_name.lower().replace("_", "")

        if (
            normalized in [
                "read",
                "reads",
                "rawread",
                "rawreads",
                "tagread",
                "tagreads",
                "rfidread",
                "rfidreads"
            ]
            or normalized.endswith("reads")
        ):
            read_fields.append((path, count))

    if read_fields:
        print("\nPOSSIBLE READ COUNT:")

        for path, count in read_fields:
            print(f" - {path}: {count}")
    else:
        print("\nCould not automatically identify the read field.")
        print("Check the repeated-field counts shown above.")


def try_decode_location(payload):
    """Try location-related Protobuf types available in messages_pb2."""

    results = []

    for type_name in location_type_names:
        message_class = getattr(messages_pb2, type_name, None)

        if message_class is None:
            continue

        try:
            location_message = message_class()
            location_message.ParseFromString(payload)

            populated_fields = len(location_message.ListFields())

            if populated_fields > 0:
                results.append(
                    (
                        populated_fields,
                        type_name,
                        location_message
                    )
                )

        except Exception:
            continue

    if not results:
        return False

    results.sort(key=lambda item: item[0], reverse=True)

    field_count, type_name, location_message = results[0]

    print("FORMAT: Protobuf")
    print("POSSIBLE PROTOBUF TYPE:", type_name)
    print("POPULATED FIELDS:", field_count)

    print("\nDECODED LOCATION DETAILS:")
    print(text_format.MessageToString(location_message))

    print_repeated_counts(location_message)

    return True


def on_connect(client, userdata, flags, rc):
    print("\nConnected to MQTT:", rc)

    client.subscribe(RAW_TOPIC, qos=0)
    client.subscribe(LOCATION_TOPIC, qos=0)

    print("\nSubscribed to:")
    print(" -", RAW_TOPIC)
    print(" -", LOCATION_TOPIC)
    print("\nWaiting for messages...")


def on_message(client, userdata, message):
    global raw_message_count
    global location_message_count

    payload = message.payload

    # -------------------------------------------------
    # RAW RFID
    # -------------------------------------------------
    if message.topic == RAW_TOPIC:
        raw_message_count += 1

        if raw_message_count <= DETAIL_LIMIT:
            print("\n" + "=" * 70)
            print("RAW RFID MESSAGE:", raw_message_count)
            print("TOPIC:", message.topic)
            print("QOS:", message.qos)
            print("PAYLOAD SIZE:", len(payload), "bytes")

            try:
                bundle = messages_pb2.ProtoReaderBundle()
                bundle.ParseFromString(payload)

                print("FORMAT: Protobuf")
                print("PROTOBUF TYPE: ProtoReaderBundle")

                print("\nDECODED RAW RFID DETAILS:")
                print(text_format.MessageToString(bundle))

                print_repeated_counts(bundle)

            except Exception as error:
                print("RAW RFID DECODE ERROR:", error)

        elif raw_message_count % 100 == 0:
            print(
                f"Received {raw_message_count} Raw RFID MQTT messages"
            )

    # -------------------------------------------------
    # LOCATION UPDATE
    # -------------------------------------------------
    elif message.topic == LOCATION_TOPIC:
        location_message_count += 1

        if location_message_count <= DETAIL_LIMIT:
            print("\n" + "=" * 70)
            print("LOCATION MESSAGE:", location_message_count)
            print("TOPIC:", message.topic)
            print("QOS:", message.qos)
            print("PAYLOAD SIZE:", len(payload), "bytes")

            decoded = try_decode_location(payload)

            if not decoded:
                print("FORMAT: Binary")
                print(
                    "No matching location Protobuf type was found "
                    "inside messages_pb2.py"
                )

                filename = (
                    f"locationUpdate_{location_message_count}.bin"
                )

                with open(filename, "wb") as file:
                    file.write(payload)

                print("Saved binary payload:", filename)
                print("HEX:", payload.hex())
                print(
                    "BASE64:",
                    base64.b64encode(payload).decode()
                )

        elif location_message_count % 10 == 0:
            print(
                f"Received {location_message_count} Location MQTT messages"
            )


client = mqtt.Client(client_id="both-topic-payload-checker")

client.on_connect = on_connect
client.on_message = on_message

client.connect("127.0.0.1", 1883, 60)

try:
    client.loop_forever()

except KeyboardInterrupt:
    print("\nStopped by user")
    print("TOTAL RAW MQTT MESSAGES:", raw_message_count)
    print("TOTAL LOCATION MESSAGES:", location_message_count)

finally:
    client.disconnect()
PY