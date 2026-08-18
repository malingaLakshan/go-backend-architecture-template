curl -sS -w '\nHTTP %{http_code}\n' -X POST \
"http://localhost:3000/siteId('3b96f652-8200-3920-8a2c-0486c358964e')/Feeds/lab-json-events/Reset"