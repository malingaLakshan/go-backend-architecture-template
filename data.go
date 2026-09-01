Thank you, Astha. This is helpful.

However, the events/json destination topic shown here was created by us only for the investigation. It is not an Event topic provided by the Resonate environment.

RawRFID and Location already have supported topics. Similarly, we need the Resonate team to configure the correct Event datafeed and provide us with the Event topic and endpoint details. Our Recorder and Validator will then subscribe to the provided topic.

The current test feed also produces delayed batches, and some messages fail because of the MQTT packet-size limit. Could we review these points during the meeting?