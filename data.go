python3 - <<'PY'
import os
import struct
import paho.mqtt.client as mqtt

BROKER = "127.0.0.1"
PORT = 1883

# Listen for Location messages from any site.
TOPICS = [
    ("resonate/locate/+/locationUpdate", 0),
    ("resonate/locate/+/mlt_sow_locations", 0),
]

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

mqtt_message_count = 0
individual_location_count = 0


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
            end = position + 8

            if end > len(data):
                raise ValueError("Incomplete fixed64 field")

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
                raise ValueError("Incomplete fixed32 field")

            value = data[position:end]
            position = end

        else:
            raise ValueError(
                "Unsupported Protobuf wire type: {}".format(wire_type)
            )

        fields.append((number, wire_type, value))

    return fields


def readable_text(data):
    try:
        text = data.decode("utf-8")

        if text and all(
            character.isprintable() or character in "\r\n\t"
            for character in text
        ):
            return text
    except UnicodeDecodeError:
        pass

    return None


def format_value(number, wire_type, value, location=False):
    if wire_type == 0:
        return str(value)

    if wire_type == 1:
        return str(struct.unpack("<Q", value)[0])

    if wire_type == 2:
        text = readable_text(value)

        if text is not None:
            return text

        return "0x" + value.hex()

    if wire_type == 5:
        float_value = struct.unpack("<f", value)[0]

        if location and number == 8:
            return str(float_value)

        integer_value = struct.unpack("<I", value)[0]

        return "{} (uint32={})".format(
            float_value,
            integer_value,
        )

    return str(value)


def print_location(location_data, location_number):
    fields = parse_fields(location_data)
    fields_by_number = {}

    for number, wire_type, value in fields:
        fields_by_number.setdefault(number, []).append(
            (wire_type, value)
        )

    print("\nLOCATION #{}".format(location_number))

    for number, name in LOCATION_NAMES.items():
        values = fields_by_number.get(number)

        if not values:
            print(
                "  Field {} - {}: <NOT SENT OR DEFAULT>".format(
                    number,
                    name,
                )
            )
            continue

        for wire_type, value in values:
            formatted_value = format_value(
                number,
                wire_type,
                value,
                location=True,
            )

            if number in (9, 11):
                formatted_value += " (enum number)"

            print(
                "  Field {} - {}: {}".format(
                    number,
                    name,
                    formatted_value,
                )
            )

    # Print any additional fields not currently identified.
    for number in sorted(fields_by_number):
        if number in LOCATION_NAMES:
            continue

        for wire_type, value in fields_by_number[number]:
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


def on_connect(client, userdata, flags, reason_code, properties=None):
    print("Connected to MQTT:", reason_code)

    client.subscribe(TOPICS)

    print("Subscribed to:")
    print("  resonate/locate/+/locationUpdate")
    print("  resonate/locate/+/mlt_sow_locations")
    print("\nWaiting continuously.")
    print("Now run sim.py in another terminal.")
    print("Press Ctrl+C to stop.")


def on_message(client, userdata, message):
    global mqtt_message_count
    global individual_location_count

    try:
        bundle_fields = parse_fields(message.payload)
    except Exception as error:
        print("\nProtobuf decoding failed:", error)
        return

    mqtt_message_count += 1

    locations = [
        value
        for number, wire_type, value in bundle_fields
        if number == 2 and wire_type == 2
    ]

    print("\n" + "=" * 78)
    print("LOCATION MQTT MESSAGE #{}".format(mqtt_message_count))
    print("=" * 78)
    print("Topic:", message.topic)
    print("QoS:", message.qos)
    print("Payload size:", len(message.payload), "bytes")
    print("Protobuf type: Location Bundle")

    print("\nALL BUNDLE FIELDS")
    print("-" * 78)

    present_bundle_fields = set()

    for number, wire_type, value in bundle_fields:
        present_bundle_fields.add(number)

        if number == 2 and wire_type == 2:
            continue

        name = BUNDLE_NAMES.get(
            number,
            "unknown_field_{}".format(number),
        )

        print(
            "  Field {} - {}: {}".format(
                number,
                name,
                format_value(number, wire_type, value),
            )
        )

    print(
        "  Field 2 - locations: {} individual locations".format(
            len(locations)
        )
    )

    print("\nALL INDIVIDUAL LOCATIONS")
    print("-" * 78)

    for index, location_data in enumerate(locations, start=1):
        print_location(location_data, index)

    individual_location_count += len(locations)

    print("\nRUNNING TOTAL")
    print("-" * 78)
    print("MQTT messages received:", mqtt_message_count)
    print(
        "Individual locations received:",
        individual_location_count,
    )


client_id = "location-protobuf-check-{}".format(os.getpid())

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
    print("MQTT messages received:", mqtt_message_count)
    print(
        "Individual locations received:",
        individual_location_count,
    )
finally:
    client.disconnect()
PY