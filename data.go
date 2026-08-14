cd /data/sim

./venv/bin/python3 - <<'PY'
import json
import base64
import paho.mqtt.client as mqtt

TOPIC = "resonate/locate/3b96f652-8200-3920-8a2c-0486c358964e/locationUpdate"

def on_connect(client, userdata, flags, rc):
    print("Connected:", rc)
    client.subscribe(TOPIC, qos=1)
    print("Subscribed to:", TOPIC)
    print("Waiting for one location message...")

def on_message(client, userdata, message):
    payload = message.payload

    print("\nTOPIC:", message.topic)
    print("QOS:", message.qos)
    print("PAYLOAD SIZE:", len(payload), "bytes")

    try:
        text = payload.decode("utf-8")

        try:
            data = json.loads(text)
            print("FORMAT: JSON")
            print(json.dumps(data, indent=2))
        except json.JSONDecodeError:
            print("FORMAT: UTF-8 text, not JSON")
            print(repr(text))

    except UnicodeDecodeError:
        print("FORMAT: Binary")
        print("HEX:", payload.hex())
        print("BASE64:", base64.b64encode(payload).decode())

    client.disconnect()

client = mqtt.Client(client_id="location-payload-checker")
client.on_connect = on_connect
client.on_message = on_message
client.connect("127.0.0.1", 1883, 60)
client.loop_forever()
PY