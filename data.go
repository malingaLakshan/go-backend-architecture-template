cd /data/sim

./venv/bin/python3 - <<'PY'
import shutil
import subprocess

import paho.mqtt.client as mqtt
import messages_pb2
from google.protobuf import text_format

SITE_ID = "3b96f652-8200-3920-8a2c-0486c358964e"

RAW_TOPIC = f"resonate/locate/{SITE_ID}/rawrfid"
LOCATION_TOPIC = f"resonate/locate/{SITE_ID}/locationUpdate"

raw_count = 0
location_count = 0
has_protoc = shutil.which("protoc") is not None


def on_connect(client, userdata, flags, reason_code, properties=None):
    print("Connected to MQTT:", reason_code)

    client.subscribe([
        (RAW_TOPIC, 0),
        (LOCATION_TOPIC, 0)
    ])

    print("Subscribed to:")
    print(" -", RAW_TOPIC)
    print(" -", LOCATION_TOPIC)
    print("Waiting for all messages...")


def on_message(client, userdata, message):
    global raw_count, location_count

    print("\n" + "=" * 70)
    print("TOPIC:", message.topic)
    print("QOS:", message.qos)
    print("PAYLOAD SIZE:", len(message.payload), "bytes")

    # RAW RFID
    if message.topic == RAW_TOPIC:
        raw_count += 1
        print("RAW RFID MESSAGE:", raw_count)

        try:
            bundle = messages_pb2.ProtoReaderBundle()
            bundle.ParseFromString(message.payload)

            print("FORMAT: Protobuf")
            print("TYPE: ProtoReaderBundle")
            print("READS INSIDE THIS MESSAGE:", len(bundle.reads))

            print("\nHUMAN-READABLE PROTOBUF:")
            print(
                text_format.MessageToString(
                    bundle,
                    as_utf8=True
                )
            )

        except Exception as error:
            print("RAW RFID DECODE FAILED:", error)

    # LOCATION UPDATE
    elif message.topic == LOCATION_TOPIC:
        location_count += 1
        print("LOCATION MESSAGE:", location_count)

        if not has_protoc:
            print("Cannot decode location fields.")
            print("The protoc command is not installed.")
            return

        result = subprocess.run(
            ["protoc", "--decode_raw"],
            input=message.payload,
            capture_output=True
        )

        if result.returncode == 0:
            print("FORMAT: Possible Protobuf")
            print("\nHUMAN-READABLE PROTOBUF FIELD VALUES:")
            print(
                result.stdout.decode(
                    "utf-8",
                    errors="replace"
                )
            )
        else:
            print("LOCATION PAYLOAD COULD NOT BE DECODED AS RAW PROTOBUF.")
            print(
                result.stderr.decode(
                    "utf-8",
                    errors="replace"
                )
            )


if hasattr(mqtt, "CallbackAPIVersion"):
    client = mqtt.Client(
        callback_api_version=mqtt.CallbackAPIVersion.VERSION2,
        client_id="mqtt-protobuf-checker"
    )
else:
    client = mqtt.Client(client_id="mqtt-protobuf-checker")

client.on_connect = on_connect
client.on_message = on_message

client.connect("127.0.0.1", 1883, 60)

try:
    client.loop_forever()

except KeyboardInterrupt:
    print("\nStopped.")
    print("Total Raw RFID messages:", raw_count)
    print("Total Location messages:", location_count)

finally:
    client.disconnect()
PY