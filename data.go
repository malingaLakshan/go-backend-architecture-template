python3 - <<'PY'
import os
import struct
import paho.mqtt.client as mqtt

BROKER = "127.0.0.1"
PORT = 1883
TOPIC = "resonate/locate/#"

BUNDLE_NAMES = {
    1: "site_id",
    2: "locations",
    3: "sent_timestamp_ms",
}

LOCATION_NAMES = {
    2: "item_id",
    3: "x",
    4: "y",
    5: "z",
    6: "floor",
    7: "region",
    8: "confidence",
    9: "state",
    10: "timestamp_ns",
    11: "reason",
}

location_count = 0
raw_count = 0


def read_varint(data, position):
    value = 0
    shift = 0

    while position < len(data):
        byte = data[position]
        position += 1
        value |= (byte & 0x7F) << shift

        if not byte & 0x80:
            return value, position

        shift += 7

    raise ValueError("Incomplete varint")


def parse_fields(data):
    fields = []
    position = 0

    while position < len(data):
        key, position = read_varint(data, position)
        number = key >> 3
        wire_type = key & 7

        if wire_type == 0:
            value, position = read_varint(data, position)

        elif wire_type == 1:
            value = data[position:position + 8]
            position += 8

        elif wire_type == 2:
            size, position = read_varint(data, position)
            value = data[position:position + size]
            position += size

        elif wire_type == 5:
            value = data[position:position + 4]
            position += 4

        else:
            raise ValueError(
                "Unsupported wire type {}".format(wire_type)
            )

        fields.append((number, wire_type, value))

    return fields


def format_value(number, wire_type, value, location=False):
    if wire_type == 0:
        return str(value)

    if wire_type == 1:
        return str(struct.unpack("<Q", value)[0])

    if wire_type == 5:
        float_value = struct.unpack("<f", value)[0]

        if location and number == 8:
            return str(float_value)

        return "{} (uint32={})".format(
            float_value,
            struct.unpack("<I", value)[0],
        )

    if wire_type == 2:
        try:
            text = value.decode("utf-8")
            if text and all(
                character.isprintable() for character in text
            ):
                return text
        except UnicodeDecodeError:
            pass

        return "0x" + value.hex()

    return str(value)


def print_location(data, index):
    print("Field 2 - locations #{} {{".format(index))

    fields = parse_fields(data)
    values_by_number = {}

    for number, wire_type, value in fields:
        values_by_number.setdefault(number, []).append(
            (wire_type, value)
        )

    # Display all currently identified fields.
    for number, name in LOCATION_NAMES.items():
        values = values_by_number.get(number)

        if not values:
            print(
                "  Field {} - {}: <NOT SENT; default 0>".format(
                    number,
                    name,
                )
            )
            continue

        for wire_type, value in values:
            formatted = format_value(
                number,
                wire_type,
                value,
                location=True,
            )

            if number in (9, 11):
                formatted += " (enum number)"

            print(
                "  Field {} - {}: {}".format(
                    number,
                    name,
                    formatted,
                )
            )

    # Display any additional fields automatically.
    for number in sorted(values_by_number):
        if number in LOCATION_NAMES:
            continue

        for wire_type, value in values_by_number[number]:
            print(
                "  Field {} - unknown_field_{}: {}".format(
                    number,
                    number,
                    format_value(
                        number,
                        wire_type,
                        value,
                        location=True,
                    ),
                )
            )

    print("}")


def decode_location_message(message):
    global location_count
    location_count += 1

    print("\n" + "=" * 78)
    print("LOCATION MQTT MESSAGE #{}".format(location_count))
    print("=" * 78)
    print("Topic:", message.topic)
    print("QoS:", message.qos)
    print("Payload size:", len(message.payload), "bytes")
    print("\nALL LOCATION FIELDS AND VALUES")
    print("-" * 78)

    fields = parse_fields(message.payload)
    location_index = 0

    for number, wire_type, value in fields:
        name = BUNDLE_NAMES.get(
            number,
            "unknown_field_{}".format(number),
        )

        if number == 2 and wire_type == 2:
            location_index += 1
            print_location(value, location_index)
        else:
            print(
                "Field {} - {}: {}".format(
                    number,
                    name,
                    format_value(number, wire_type, value),
                )
            )


def on_connect(client, userdata, flags, reason_code, properties=None):
    print("Connected to MQTT:", reason_code)
    client.subscribe(TOPIC, qos=0)
    print("Listening to:", TOPIC)
    print("Now run sim.py in another terminal.")
    print("Press Ctrl+C to stop.\n")


def on_message(client, userdata, message):
    global raw_count

    topic_lower = message.topic.lower()

    if topic_lower.endswith("/rawrfid"):
        raw_count += 1
        print(
            "[RAW RFID #{}] received: {} bytes from {}".format(
                raw_count,
                len(message.payload),
                message.topic,
            )
        )
        return

    if (
        topic_lower.endswith("/locationupdate")
        or topic_lower.endswith("/mlt_sow_locations")
    ):
        try:
            decode_location_message(message)
        except Exception as error:
            print("Location Protobuf decoding failed:", error)


client_id = "location-check-{}".format(os.getpid())

try:
    client = mqtt.Client(
        callback_api_version=mqtt.CallbackAPIVersion.VERSION2,
        client_id=client_id,
    )
except (AttributeError, TypeError):
    client = mqtt.Client(client_id=client_id)

client.on_connect = on_connect
client.on_message = on_message
client.connect(BROKER, PORT, 60)

try:
    client.loop_forever()
except KeyboardInterrupt:
    print("\nStopped by user.")
    print("Raw RFID messages:", raw_count)
    print("Location messages:", location_count)
finally:
    client.disconnect()
PY