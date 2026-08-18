sudo docker cp alt-sl-datafeed-service:/app /tmp/datafeed-service-app

strings /tmp/datafeed-service-app \
| grep -iE 'mqtt://|tcp://|broker|mqtt.*publish|publish.*mqtt|eventFilters|filterMode' \
| head -120