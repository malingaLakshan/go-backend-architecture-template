import time
from pathlib import Path

import paho.mqtt.client as mqtt
import messages_pb2

SITE_ID = "3b96f652-8200-3920-8a2c-0486c358964e"

RAW_TOPIC = f"resonate/locate/{SITE_ID}/rawrfid"
LOCATION_TOPIC = f"resonate/locate/{SITE_ID}/locationUpdate"

stats = {
    "raw": 0,
    "raw_decoded": 0,
    "raw_errors": 0,
    "location": 0
}

raw_sample_saved = False
location_sample_saved = False


def on_connect(client, userdata, flags, rc):
    print("Connected to MQTT:", rc)

    client.subscribe(RAW_TOPIC, qos=0)
    client.subscribe(LOCATION_TOPIC, qos=0)

    print("Subscribed to:")
    print(" -", RAW_TOPIC)
    print(" -", LOCATION_TOPIC)
    print("Waiting for messages...")


def on_message(client, userdata, message):
    global raw_sample_saved, location_sample_saved

    if message.topic == RAW_TOPIC:
        stats["raw"] += 1

        try:
            bundle = messages_pb2.ProtoReaderBundle()
            bundle.ParseFromString(message.payload)
            stats["raw_decoded"] += 1

            if not raw_sample_saved:
                Path("rawrfid_first.bin").write_bytes(message.payload)
                Path("rawrfid_first_decoded.txt").write_text(
                    str(bundle),
                    encoding="utf-8"
                )
                raw_sample_saved = True
                print("\nSaved first decoded Raw RFID sample.")

        except Exception as error:
            stats["raw_errors"] += 1
            print("Raw RFID decode error:", error)

    elif message.topic == LOCATION_TOPIC:
        stats["location"] += 1

        if not location_sample_saved:
            Path("locationUpdate_first.bin").write_bytes(message.payload)
            location_sample_saved = True
            print(
                "\nSaved first Location Update sample:",
                len(message.payload),
                "bytes"
            )


client = mqtt.Client(client_id="resonate-topic-checker")
client.on_connect = on_connect
client.on_message = on_message

client.connect("127.0.0.1", 1883, 60)
client.loop_start()

previous_raw = 0
previous_location = 0

try:
    while True:
        time.sleep(5)

        new_raw = stats["raw"] - previous_raw
        new_location = stats["location"] - previous_location

        print(
            f"RAW MQTT messages: {stats['raw']} (+{new_raw}) | "
            f"Decoded: {stats['raw_decoded']} | "
            f"Decode errors: {stats['raw_errors']} | "
            f"LOCATION messages: {stats['location']} (+{new_location})"
        )

        previous_raw = stats["raw"]
        previous_location = stats["location"]

except KeyboardInterrupt:
    print("\nStopping topic checker...")

finally:
    client.loop_stop()
    client.disconnect()