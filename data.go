~/resonate-sim-venv/bin/python -c '
import paho.mqtt.client as mqtt
c=mqtt.Client(mqtt.CallbackAPIVersion.VERSION1)
c.on_connect=lambda c,u,f,rc: c.subscribe("resonate/locate/3b96f652-8200-3920-8a2c-0486c358964e/events/json")
c.on_message=lambda c,u,m: print(m.topic, m.payload.decode(), flush=True)
c.connect("127.0.0.1",1883,60)
c.loop_forever()
'