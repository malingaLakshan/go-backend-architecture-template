Hi Andrew,

To explain our current finding clearly, the Resonate environment provides two default MQTT topics:

* RawRFID: resonate/locate/<site-id>/rawrfid
* Location: resonate/locate/<site-id>/locationUpdate

There is no default MQTT topic for Events.

For our investigation, we manually created and started a datafeed and configured our own destination topic:

resonate/locate/<site-id>/events/json

We then subscribed to this topic and successfully received JSON Events. This was only a manual test setup; the Event topic was not provided by Resonate.

We understand that the updated Recorder and Validator requirements allow the broker URL and Event topic to be provided through configuration. For this integration, we need Resonate to provide and manage a running datafeed with a configured Event topic, so the Recorder and Validator can subscribe to it.

Could you please confirm this setup and provide the Event topic that should be used?