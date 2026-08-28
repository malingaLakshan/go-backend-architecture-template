Hi Andrew,

Thank you for updating the requirements. We understand that the Recorder and Validator will read the transport, brokerUrl, and topic from the resonateEventStream configuration and subscribe to the configured MQTT topic.

One integration responsibility still needs clarification. During our testing, Events were published only after we manually created and started a datafeed with a configured MQTT destination topic. The topic we used was configured by us for that datafeed; it is not a default Resonate Event topic.

Could you please confirm whether creating, configuring, starting, and managing the datafeed should be treated as an external prerequisite managed by the Resonate deployment/operator, or whether any of this responsibility belongs to the Recorder or Validator?

If it is an external prerequisite, should this also be stated in the acceptance criteria or integration documentation?