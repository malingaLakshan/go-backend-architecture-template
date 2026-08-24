Good morning Steve,

We confirmed that the RawRFID and SOW/location streams are working and have default MQTT topics:

* resonate/locate/<site-id>/rawrfid
* resonate/locate/<site-id>/locationUpdate

We also successfully received JSON Events on the new VM. However, Events do not have a default topic—we had to create and start a datafeed with a configured destination topic.

Could you please confirm whether our Recorder and Validator should create and manage this datafeed, or whether Resonate will provide a preconfigured Event datafeed and topic?