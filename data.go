cd /opt/zebra/alt-resonate-locate/svc
sudo docker compose logs --since 10m datafeed 2>&1 | tail -100