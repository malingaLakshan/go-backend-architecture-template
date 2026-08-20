Hi all, I checked the JSON Event stream.

The datafeed-service creates business Events by watching asset changes in MongoDB. Before it can send those Events, we must create a datafeed. A datafeed is simply a configuration that defines which Events to generate and where to send the JSON output.

I created a datafeed and selected this test MQTT topic as its output address:

resonate/locate/3b96f652-8200-3920-8a2c-0486c358964e/events/json

The Events were generated successfully. However, the service itself displayed MOCK TRANSPORT: MQTT in its logs. We did not enable mock transport manually.

This means the service only printed the Event JSON in the logs and did not publish it to the MQTT broker.

Could you please confirm how to enable real MQTT publishing? Does it require a configuration change, environment variable, or different datafeed-service image?