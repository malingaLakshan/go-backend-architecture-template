Investigation Conclusion and Current Status

The investigation covered the RawRFID, Location and Event streams from the real Resonate environment.

Work completed:

* Verified the Resonate services, MQTT broker and simulator setup on the VM.
* Fixed the simulator’s missing Python dependencies and successfully generated live RFID movement.
* Created subscriber utilities to inspect and decode the RawRFID and Location payloads.
* Compared the received stream fields with the fields expected by our Recorder database.
* Investigated the Event generation workflow, datafeed configuration and SiteGraph APIs.
* Coordinated with the Resonate team regarding the old datafeed image, Event topic handling and the latest test VM.

Stream findings:

1. RawRFID stream

* Default topic: resonate/locate/<site-id>/rawrfid
* Payload format: Protobuf
* Live RawRFID bundles were received and decoded successfully.

2. Location stream

* Default topic: resonate/locate/<site-id>/locationUpdate
* Payload format: Protobuf
* Live location updates were received and decoded successfully after RawRFID data was processed by Resonate.

3. Event stream

* Events do not have a default MQTT topic.
* The datafeed-service generates business Events by monitoring asset changes in MongoDB.
* A datafeed must be created to define the required Event types and the MQTT destination topic.
* We configured the following test destination topic:

resonate/locate/<site-id>/events/json

* On the original VM, the datafeed used MOCK TRANSPORT. It generated Event JSON only in the service logs and did not publish it to MQTT.
* The Resonate team provided a new VM containing the latest datafeed service.
* On the new VM, the datafeed was created and started successfully, and JSON Event arrays were received through MQTT.
* Events such as ARRIVAL, POSITION_CHANGE, REGION_CHANGE, FLOOR_CHANGE, DEPARTURE and EXIT were observed.

Current Event integration approach:

The Validator and Recorder should accept the configured Event topic through an input parameter or configuration file. David confirmed this approach, and Andrew will update the requirements. Once the requirement is updated, we can make the necessary changes on our side.

Related API findings:

* The single-site SiteGraph endpoint works when accessed directly:

http://117108-trirh902.117asd.zebra.lan:3389/api/sitegraph/<site-id>

* Calls from our application are currently blocked by CORS.
* The provided Site List path, /trifecta/v1/sds/sites, was tested with the available services and common route variations on ports 8080 and 3389. All tests returned HTTP 404.
* We are waiting for confirmation of the correct Site List base URL and the expected approach for accessing the SiteGraph endpoint from our application.

Final status:

The RawRFID, Location and Event streams have now been identified and verified. RawRFID and Location use default MQTT topics, while the Event stream requires a datafeed and a configured destination topic. The remaining work is to apply the updated Validator and Recorder requirements and resolve the Site List API and SiteGraph CORS details.