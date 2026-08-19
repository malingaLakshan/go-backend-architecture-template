EVENT_TOPIC='resonate/locate/3b96f652-8200-3920-8a2c-0486c358964e/events/json'

sudo docker exec mqtt emqx ctl trace start topic "$EVENT_TOPIC" /tmp/event-topic-check.log debug
sleep 30
sudo docker exec mqtt emqx ctl trace stop topic "$EVENT_TOPIC"
sudo docker exec mqtt sh -c 'wc -c /tmp/event-topic-check.log; tail -40 /tmp/event-topic-check.log'