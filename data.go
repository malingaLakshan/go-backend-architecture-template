sudo docker exec mongodb mongosh \
"mongodb://admin:admin@localhost:27017/location_engine?authSource=admin" \
--quiet --eval '
printjson(
  db.feed_configs.find(
    {},
    {
      _id: 0,
      feedId: 1,
      siteId: 1,
      status: 1,
      eventFilters: 1,
      filterMode: 1,
      destination: 1
    }
  ).toArray()
)'