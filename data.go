cd /opt/zebra/alt-resonate-locate/svc

sudo docker cp "$(sudo docker compose ps -q api)":/app /tmp/smartlens-api

strings /tmp/smartlens-api | grep -iE 'trifecta|sds/sites|sitegraph|/sites' | sort -u | head -100