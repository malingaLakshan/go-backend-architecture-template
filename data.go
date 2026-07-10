Fix RawReads scanning for nullable SQLite columns.

Current play command fails with:
sql: Scan error on column index 9, name "rssi": converting NULL to float64 is unsupported

Real Recorder DB can contain NULL values for RawReads columns like:
- antenna_id
- antenna_type
- confidence
- rssi
- tag_x
- tag_y
- floor_id

Update internal/recording/repository.go GetRawReads() so nullable DB columns are scanned safely.

Use:
- sql.NullInt64 for antenna_id, antenna_type, confidence, floor_id
- sql.NullFloat64 for rssi, tag_x, tag_y

After scanning, assign defaults only when Valid is true:
- if antennaID.Valid { r.AntennaID = int(antennaID.Int64) }
- if antennaType.Valid { r.AntennaTypeID = int(antennaType.Int64) }
- if confidence.Valid { r.Confidence = int(confidence.Int64) }
- if rssi.Valid { r.RSSI = rssi.Float64 }
- if tagX.Valid { r.TagX = tagX.Float64 }
- if tagY.Valid { r.TagY = tagY.Float64 }
- if floorID.Valid { r.FloorID = int(floorID.Int64) }

Do not change the real Recorder column names.
Do not use dynamic SQL.
Do not change read_id from string.
Do not change raw_payload from []byte.

Then run:
gofmt -w internal/recording/repository.go
go test ./...
go build -o rre.exe ./cmd/rre