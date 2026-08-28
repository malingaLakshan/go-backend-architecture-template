Hi Andrew,

Thank you. I tested the Site List endpoint using the new VM hostname:

https://117108-trirh902.117asd.zebra.lan/trifecta/v1/sds/sites

The browser returned ERR_CONNECTION_CLOSED. I also checked the VM, and no service is currently listening on host port 443.

Could you please confirm whether port 443 should be exposed on this lab VM, or whether hostname.zebra.lan refers to a different host? If it refers to another service, could you please provide the intended hostname and the required authentication method?

Regarding the single-site SiteGraph endpoint, our application is already calling it through our backend. The endpoint works when opened directly in the browser, but the request made through our backend fails.

Could you please confirm whether the backend request requires any authentication, specific headers, certificates, or additional service configuration?