sudo docker exec mongodb mongosh \
"mongodb://admin:admin@localhost:27017/location_engine?authSource=admin" \
--quiet --eval '
db.feed_configs.find({}).forEach(d => printjson({
  feedId: d.feedId,
  keys: Object.keys(d),
  siteId: d.siteId,
  status: d.status,
  filters: d.filters,
  pMode: d.pMode,
  eventFilters: d.eventFilters,
  filterMode: d.filterMode,
  destination: d.destination
}))
'