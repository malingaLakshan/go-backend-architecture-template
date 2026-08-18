python3 - <<'PY'
import math
import os
import struct
import paho.mqtt.client as mqtt

BROKER = "127.0.0.1"
PORT = 1883
TOPIC = "resonate/locate/3b96f652-8200-3920-8a2c-0486c358964e/locationUpdate"

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

STRING_FIELDS = {2, 7}
message_count = 0


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
        if shift > 70:
            raise ValueError("Invalid varint")

    raise ValueError("Incomplete varint")


def parse_fields(data):
    fields = []
    position = 0

    while position < len(data):
        key, position = read_varint(data, position)
        field_number = key >> 3
        wire_type = key & 7

        if field_number == 0:
            raise ValueError("Invalid field number 0")

        if wire_type == 0:
            value, position = read_varint(data, position)

        elif wire_type == 1:
            end = position + 8
            if end > len(data):
                raise ValueError("Incomplete fixed64")
            value = data[position:end]
            position = end

        elif wire_type == 2:
            size, position = read_varint(data, position)
            end = position + size
            if end > len(data):
                raise ValueError("Incomplete length-delimited field")
            value = data[position:end]
            position = end

        elif wire_type == 5:
            end = position + 4
            if end > len(data):
                raise ValueError("Incomplete fixed32")
            value = data[position:end]
            position = end

        else:
            raise ValueError(
                "Unsupported Protobuf wire type: {}".format(wire_type)
            )

        fields.append((field_number, wire_type, value))

    return fields


def readable_text(data):
    try:
        value = data.decode("utf-8")
        if value and all(
            character.isprintable() or character in "\r\n\t"
            for character in value
        ):
            return value
    except UnicodeDecodeError:
        pass

    return None


def format_float(value):
    if not math.isfinite(value):
        return str(value)
    return "{:.9g}".format(value)


def format_value(field_number, wire_type, value, location=False):
    if wire_type == 0:
        if location and field_number in (9, 11):
            return "{} (enum number)".format(value)
        return str(value)

    if wire_type == 1:
        integer_value = struct.unpack("<Q", value)[0]
        double_value = struct.unpack("<d", value)[0]
        return "uint64={}, double={}".format(
            integer_value,
            format_float(double_value),
        )

    if wire_type == 2:
        text = readable_text(value)
        if text is not None:
            return text
        return "0x" + value.hex()

    if wire_type == 5:
        integer_value = struct.unpack("<I", value)[0]
        float_value = struct.unpack("<f", value)[0]

        if location and field_number == 8:
            return format_float(float_value)

        return "float={}, uint32={}".format(
            format_float(float_value),
            integer_value,
        )

    return str(value)


def missing_default(field_number):
    if field_number in STRING_FIELDS:
        return "<NOT SENT; default empty string>"
    return "<NOT SENT; protobuf default 0>"


def print_location(data, location_number):
    print("Field 2 - locations #{} {{".format(location_number))

    try:
        fields = parse_fields(data)
    except Exception as error:
        print("  Could not decode location:", error)
        print("  Raw bytes:", data.hex())
        print("}")
        return

    fields_by_number = {}
    for field_number, wire_type, value in fields:
        fields_by_number.setdefault(field_number, []).append(
            (wire_type, value)
        )

    # Print every known Location field, including fields omitted as zero.
    for field_number, field_name in LOCATION_NAMES.items():
        values = fields_by_number.get(field_number, [])

        if not values:
            print(
                "  Field {} - {}: {}".format(
                    field_number,
                    field_name,
                    missing_default(field_number),
                )
            )
            continue

        for wire_type, value in values:
            print(
                "  Field {} - {}: {}".format(
                    field_number,
                    field_name,
                    format_value(
                        field_number,
                        wire_type,
                        value,
                        location=True,
                    ),
                )
            )

    # Automatically print any future/additional fields.
    for field_number in sorted(fields_by_number):
        if field_number in LOCATION_NAMES:
            continue

        for wire_type, value in fields_by_number[field_number]:
            print(
                "  Field {} - unknown_field_{}: {}".format(
                    field_number,
                    field_number,
                    format_value(
                        field_number,
                        wire_type,
                        value,
                        location=True,
                    ),
                )
            )

    print("}")


def on_connect(client, userdata, flags, reason_code, properties=None):
    print("Connected to MQTT:", reason_code)
    client.subscribe(TOPIC, qos=1)
    print("Subscribed to:", TOPIC)
    print("Waiting continuously. Press Ctrl+C to stop.")


def on_message(client, userdata, message):
    global message_count
    message_count += 1

    print("\n" + "=" * 78)
    print("LOCATION MQTT MESSAGE #{}".format(message_count))
    print("=" * 78)
    print("Topic:", message.topic)
    print("QoS:", message.qos)
    print("Payload size:", len(message.payload), "bytes")
    print("\nALL LOCATION FIELDS AND VALUES")
    print("-" * 78)

    try:
        bundle_fields = parse_fields(message.payload)
    except Exception as error:
        print("Protobuf decoding failed:", error)
        return

    present_bundle_fields = set()
    location_number = 0

    for field_number, wire_type, value in bundle_fields:
        present_bundle_fields.add(field_number)
        field_name = BUNDLE_NAMES.get(
            field_number,
            "unknown_field_{}".format(field_number),
        )

        if field_number == 2 and wire_type == 2:
            location_number += 1
            print_location(value, location_number)
        else:
            print(
                "Field {} - {}: {}".format(
                    field_number,
                    field_name,
                    format_value(field_number, wire_type, value),
                )
            )

    for field_number, field_name in BUNDLE_NAMES.items():
        if field_number not in present_bundle_fields:
            print(
                "Field {} - {}: <NOT SENT>".format(
                    field_number,
                    field_name,
                )
            )


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
    print("Total Location messages received:", message_count)
finally:
    client.disconnect()
PY