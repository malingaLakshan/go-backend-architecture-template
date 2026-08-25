cd /opt/zebra/alt-resonate-locate/svc

sudo docker inspect "$(sudo docker compose ps -q api)" \
  --format '{{range .Config.Env}}{{println .}}{{end}}' |
  sort |
  grep -iE 'url|host|port|site|sds|trifecta'