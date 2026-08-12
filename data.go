{
  "buildNumbers": {
    "resonateUrl": "mock://build/resonate",
    "firmwareUrl": "mock://build/firmware",
    "readerAppsUrl": "mock://build/reader-apps"
  },
  "sites": {
    "listUrl": "http://117108-trirh901.117asd.zebra.lan:3389/api/locations",
    "detailsUrl": "http://117108-trirh901.117asd.zebra.lan:3389/api/locations/{siteId}"
  },
  "streams": {
    "rawReadsUrl": "resonate/locate/3b96f652-8200-3920-8a2c-0486c358964e/rawrfid",
    "locationsUrl": "resonate/locate/3b96f652-8200-3920-8a2c-0486c358964e/locationUpdate",
    "eventsUrl": ""
  },
  "snapshots": {
    "allTagLocationsUrl": "http://117108-trirh901.117asd.zebra.lan:3389/api/assets/{siteId}",
    "autoFinalSnapshot": true
  },
  "siteId": "3b96f652-8200-3920-8a2c-0486c358964e",
  "mock": {
    "enabled": false,
    "siteGraphDir": "./testdata",
    "messagesPerSecond": 10
  },
  "mqtt": {
    "brokerUrl": "tcp://117108-trirh901.117asd.zebra.lan:1883"
  },
  "storage": {
    "workingDirectory": "./data",
    "batchSize": 100,
    "flushIntervalMilliseconds": 1000
  },
  "control": {
    "host": "127.0.0.1",
    "port": 8787
  },
  "logging": {
    "directory": "./logs",
    "prefix": "rr",
    "minLevel": "debug"
  }
}