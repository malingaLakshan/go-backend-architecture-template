cd /data/sim

./venv/bin/python3 - <<'PY'
import math
import os
import struct

import paho.mqtt.client as mqtt

SITE_ID = "3b96f652-8200-3920-8a2c-0486c358964e"
TOPIC = f"resonate/locate/{SITE_ID}/locationUpdate"

message_count = 0


def read_varint(data, position):
    value = 0
    shift = 0

    while position < len(data) and shift < 70:
        byte = data[position]
        position += 1
        value |= (byte & 0x7F) << shift

        if not byte & 0x80:
            return value, position

        shift += 7

    raise ValueError("Invalid varint")


def readable_text(data):
    try:
        text = data.decode("utf-8")
    except UnicodeDecodeError:
        return None

    if text and all(
        char.isprintable() or char in "\r\n\t"
        for char in text
    ):
        return text

    return None


def format_float(value):
    if not math.isfinite(value):
        return str(value)

    return f"{value:.9g}"


def decode_fields(data, level=0):
    lines = []
    position = 0
    indent = "  " * level

    while position < len(data):
        key, position = read_varint(data, position)

        field_number = key >> 3
        wire_type = key & 7

        if field_number == 0:
            raise ValueError("Invalid field number")

        # Varint: integer, boolean, enum or timestamp
        if wire_type == 0:
            value, position = read_varint(data, position)

            lines.append(
                f"{indent}Field {field_number} "
                f"(integer/varint): {value}"
            )

        # Fixed 64-bit value
        elif wire_type == 1:
            end = position + 8

            if end > len(data):
                raise ValueError("Truncated fixed64")

            raw = data[position:end]
            position = end

            integer_value = struct.unpack("<Q", raw)[0]
            double_value = struct.unpack("<d", raw)[0]

            lines.append(
                f"{indent}Field {field_number} (fixed64): "
                f"integer={integer_value}, "
                f"double={format_float(double_value)}"
            )

        # String, bytes or embedded Protobuf message
        elif wire_type == 2:
            size, position = read_varint(data, position)
            end = position + size

            if end > len(data):
                raise ValueError("Truncated field")

            raw = data[position:end]
            position = end

            text = readable_text(raw)

            if text is not None:
                lines.append(
                    f"{indent}Field {field_number} "
                    f"(string): {text!r}"
                )

            else:
                try:
                    nested = decode_fields(raw, level + 1)

                    if not nested:
                        raise ValueError("Empty embedded value")

                except Exception:
                    value = raw.hex()

                    if len(value) > 160:
                        value = value[:160] + "..."

                    lines.append(
                        f"{indent}Field {field_number} "
                        f"(bytes, {len(raw)} bytes): 0x{value}"
                    )

                else:
                    lines.append(
                        f"{indent}Field {field_number} "
                        "(embedded message) {"
                    )

                    lines.extend(nested)
                    lines.append(f"{indent}}}")

        # Fixed 32-bit value
        elif wire_type == 5:
            end = position + 4

            if end > len(data):
                raise ValueError("Truncated fixed32")

            raw = data[position:end]
            position = end

            integer_value = struct.unpack("<I", raw)[0]
            float_value = struct.unpack("<f", raw)[0]

            lines.append(
                f"{indent}Field {field_number} (fixed32): "
                f"integer={integer_value}, "
                f"float={format_float(float_value)}"
            )

        else:
            raise ValueError(
                f"Unsupported wire type: {wire_type}"
            )

    return lines


def on_connect(
    client,
    userdata,
    flags,
    reason_code,
    properties=None
):
    print("Connected:", reason_code)

    client.subscribe(TOPIC, qos=0)

    print("Subscribed to:", TOPIC)
    print("Waiting continuously for Location messages...")
    print("Press Ctrl+C whenever you want to stop.")


def on_message(client, userdata, mqtt_message):
    global message_count

    message_count += 1

    print("\n" + "=" * 75)
    print(f"LOCATION MQTT MESSAGE #{message_count}")
    print("=" * 75)

    print("Topic:", mqtt_message.topic)
    print("QoS:", mqtt_message.qos)
    print(
        "Payload size:",
        len(mqtt_message.payload),
        "bytes"
    )

    print("\nALL LOCATION FIELDS AND VALUES")
    print("-" * 75)

    try:
        fields = decode_fields(mqtt_message.payload)

        for field in fields:
            print(field)

    except Exception as error:
        print("Protobuf wire decoding failed:", error)


client_id = f"location-checker-{os.getpid()}"

if hasattr(mqtt, "CallbackAPIVersion"):
    client = mqtt.Client(
        callback_api_version=mqtt.CallbackAPIVersion.VERSION2,
        client_id=client_id
    )
else:
    client = mqtt.Client(client_id=client_id)

client.on_connect = on_connect
client.on_message = on_message

client.connect("127.0.0.1", 1883, 60)

try:
    client.loop_forever()

except KeyboardInterrupt:
    print("\nStopped by user.")
    print(
        "Total Location messages received:",
        message_count
    )

finally:
    client.disconnect()
PY