curl -sS -w '\nHTTP %{http_code}\n' -X POST \
"http://localhost:3000/feed/siteId('3b96f652-8200-3920-8a2c-0486c358964e')" \
-H "Content-Type: application/json" \
-d '{"feedId":"lab-json-events","eventFilters":[],"filterMode":"STRICT","destination":{"protocol":"MQTT","broker":"mqtt:1883","topic":"resonate/locate/3b96f652-8200-3920-8a2c-0486c358964e/events/json"}}'