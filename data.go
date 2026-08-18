strings /tmp/datafeed-service-app \
| grep -iE 'mqtt://|tcp://|broker|mqtt.*publish|publish.*mqtt|eventFilters|filterMode' \
| head -120