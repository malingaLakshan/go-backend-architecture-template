sudo docker logs --since 10m alt-sl-datafeed-service 2>&1 \
| grep -iE 'mqtt|broker|connect|publish|error|failed|exception|lab-json-events' \
| tail -100.  sudo docker logs --since 10m mqtt 2>&1 | tail -100




