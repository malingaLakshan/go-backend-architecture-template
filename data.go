cd ~/sim

python3 - <<'PY'
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
    7: "floor_id",
    8: "confidence",
    10: "timestamp_ns",
}

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

    raise ValueError("Incomplete varint")


def readable_text(data):
    try:
        value = data.decode("utf-8")
        if value and all(c.isprintable() for c in value):
            return value
    except UnicodeDecodeError:
        pass

    return None


def decode(data, names, depth=0):
    position = 0
    indent = "  " * depth

    while position < len(data):
        key, position = read_varint(data, position)
        number = key >> 3
        wire_type = key & 7
        name = names.get(number, f"unknown_field_{number}")

        if wire_type == 0:
            value, position = read_varint(data, position)
            print(f"{indent}Field {number} - {name}: {value}")

        elif wire_type == 1:
            raw = data[position:position + 8]
            position += 8

            integer_value = struct.unpack("<Q", raw)[0]
            double_value = struct.unpack("<d", raw)[0]

            print(
                f"{indent}Field {number} - {name}: "
                f"{integer_value} (double={double_value})"
            )

        elif wire_type == 2:
            size, position = read_varint(data, position)
            raw = data[position:position + size]
            position += size

            # Top-level field 2 contains individual Location messages.
            if depth == 0 and number == 2:
                print(f"{indent}Field {number} - {name} {{")
                decode(raw, LOCATION_NAMES, depth + 1)
                print(f"{indent}}}")
                continue

            text = readable_text(raw)

            if text is not None:
                print(f"{indent}Field {number} - {name}: {text}")
            else:
                try:
                    print(f"{indent}Field {number} - {name} {{")
                    decode(raw, {}, depth + 1)
                    print(f"{indent}}}")
                except Exception:
                    print(
                        f"{indent}Field {number} - {name}: "
                        f"0x{raw.hex()}"
                    )

        elif wire_type == 5:
            raw = data[position:position + 4]
            position += 4

            integer_value = struct.unpack("<I", raw)[0]
            float_value = struct.unpack("<f", raw)[0]

            print(
                f"{indent}Field {number} - {name}: "
                f"{float_value} "
                f"(uint32={integer_value})"
            )

        else:
            raise ValueError(f"Unsupported wire type: {wire_type}")


def on_connect(client, userdata, flags, reason_code, properties=None):
    print("Connected to MQTT:", reason_code)
    client.subscribe(TOPIC, qos=1)
    print("Subscribed to:", TOPIC)
    print("Waiting continuously. Press Ctrl+C to stop.\n")


def on_message(client, userdata, message):
    global message_count
    message_count += 1

    print("\n" + "=" * 75)
    print(f"LOCATION MQTT MESSAGE #{message_count}")
    print("=" * 75)
    print("Topic:", message.topic)
    print("QoS:", message.qos)
    print("Payload size:", len(message.payload), "bytes")
    print("\nALL LOCATION FIELDS AND VALUES")
    print("-" * 75)

    try:
        decode(message.payload, BUNDLE_NAMES)
    except Exception as error:
        print("Decode error:", error)


try:
    client = mqtt.Client(mqtt.CallbackAPIVersion.VERSION2)
except (AttributeError, TypeError):
    client = mqtt.Client()

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