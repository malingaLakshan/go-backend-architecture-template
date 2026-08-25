cd /opt/zebra/alt-resonate-locate/svc

sudo docker inspect "$(sudo docker compose ps -q api)" \
--format 'Path={{.Path}} Args={{json .Args}} WorkDir={{.Config.WorkingDir}}'