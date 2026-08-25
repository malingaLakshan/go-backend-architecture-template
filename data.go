for p in \
"/trifecta/v1/sds/sites" \
"/trifecta/v1/sds/sites/" \
"/api/trifecta/v1/sds/sites" \
"/v1/sds/sites" \
"/sites"
do
  code=$(curl -s -o /tmp/site-response.txt -w "%{http_code}" "http://localhost:8080$p")
  echo "$code  $p"
done