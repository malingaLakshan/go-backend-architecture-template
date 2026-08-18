sudo docker logs --since 5m --timestamps alt-sl-datafeed-service 2>&1 | tail -150

sudo docker ps --format 'table {{.Names}}\t{{.Status}}\t{{.Ports}}' | grep -Ei 'mqtt|mosquitto'