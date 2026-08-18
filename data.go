sudo docker ps --format 'table {{.Names}}\t{{.Status}}\t{{.Ports}}' | grep -E 'datafeed|3000'

sudo docker port alt-sl-datafeed-service

sudo docker inspect alt-sl-datafeed-service \
--format 'WorkingDir={{.Config.WorkingDir}} Cmd={{json .Config.Cmd}} Entrypoint={{json .Config.Entrypoint}} ExposedPorts={{json .Config.ExposedPorts}}'

sudo docker logs --tail 100 alt-sl-datafeed-service 2>&1