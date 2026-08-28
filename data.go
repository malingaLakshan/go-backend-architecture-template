Hi Andrew,

Thank you. I tested the Site List endpoint using the new VM hostname:

https://117108-trirh902.117asd.zebra.lan/trifecta/v1/sds/sites

The browser returned ERR_CONNECTION_CLOSED. I also checked the VM and confirmed that no service is currently listening on host port 443.

Could you please confirm whether port 443 should be exposed on this lab VM or whether hostname.zebra.lan refers to a different host? If it is a different service, could you please provide the correct hostname and authentication method?

Regarding the single-site SiteGraph endpoint, our application is already calling it through our backend. The endpoint works when opened directly in the browser. However, when the request is made through our application’s backend flow, the browser reports a CORS error and the request fails.

Could you please confirm the expected CORS or proxy configuration for this flow and whether any specific authentication, headers or certificates are required?