cd /data/sim

./venv/bin/python3 - <<'PY'
import json
import paho.mqtt.client as mqtt
import messages_pb2

site_id = "3b96f652-8200-3920-8a2c-0486c358964e"

raw_topic = f"resonate/locate/{site_id}/rawrfid"
location_topic = f"resonate/locate/{site_id}/locationUpdate"

seen = set()

def on_connect(client, userdata, flags, rc):
    print("Connected to MQTT, result:", rc)
    client.subscribe([
        (raw_topic, 0),
        (location_topic, 1)
    ])
    print("Subscribed to:")
    print(raw_topic)
    print(location_topic)
    print("Waiting for messages...")

def on_message(client, userdata, message):
    if message.topic in seen:
        return

    seen.add(message.topic)

    print("\n================================")
    print("TOPIC:", message.topic)
    print("QOS:", message.qos)
    print("PAYLOAD SIZE:", len(message.payload), "bytes")

    if message.topic == raw_topic:
        try:
            bundle = messages_pb2.ProtoReaderBundle()
            bundle.ParseFromString(message.payload)

            print("FORMAT: ProtoReaderBundle")
            print("SITE ID:", bundle.site_id)
            print("READER ID:", bundle.reader_id)
            print("READ COUNT:", len(bundle.reads))
            print("SENT TIMESTAMP:", bundle.sent_timestamp_ms)
        except Exception as error:
            print("Protobuf decoding failed:", error)
            print("First bytes:", message.payload[:100].hex())

    if message.topic == location_topic:
        try:
            text = message.payload.decode("utf-8")
            try:
                print("FORMAT: JSON")
                print(json.dumps(json.loads(text), indent=2))
            except json.JSONDecodeError:
                print("FORMAT: UTF-8 text")
                print(text)
        except UnicodeDecodeError:
            print("FORMAT: Binary")
            print("First bytes:", message.payload[:100].hex())

    if raw_topic in seen and location_topic in seen:
        print("\nSUCCESS: Messages received from both topics.")
        client.disconnect()

client = mqtt.Client(client_id="resonate-topic-checker")
client.on_connect = on_connect
client.on_message = on_message
client.connect("127.0.0.1", 1883, 60)
client.loop_forever()
PY