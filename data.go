Hi Andrew and David,

Thank you for the clarification regarding the Event topic. We understand that the Recorder and Validator should accept the topic name through an input parameter or configuration file. Once Andrew updates the requirements, we can make the necessary changes to the Recorder and Validator on our side.

Andrew, I also tested the Site List path you provided using the services available on the new VM. I tried the following URLs:

* http://117108-trirh902.117asd.zebra.lan:8080/trifecta/v1/sds/sites
* http://117108-trirh902.117asd.zebra.lan:8080/trifecta/v1/sds/sites/
* http://117108-trirh902.117asd.zebra.lan:8080/api/trifecta/v1/sds/sites
* http://117108-trirh902.117asd.zebra.lan:8080/v1/sds/sites
* http://117108-trirh902.117asd.zebra.lan:8080/sites
* http://117108-trirh902.117asd.zebra.lan:3389/trifecta/v1/sds/sites

All returned HTTP 404. Could you please provide the correct base URL, including the host and port, for /trifecta/v1/sds/sites and confirm whether authentication is required?

The following single-site SiteGraph endpoint works when opened directly:

http://117108-trirh902.117asd.zebra.lan:3389/api/sitegraph/<site-id>

However, requests from our application are blocked by CORS. Could you please confirm whether we should call this endpoint through our backend or whether CORS will be enabled on the service?