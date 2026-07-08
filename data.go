func BuildPayload(
	rawRead *recording.RawRead,
	siteID string,
) (*ProtoReaderBundleWrapper, error) {

	// First try old/sample JSON RawPayload.
	// This keeps generated sample DB support.
	if len(rawRead.RawPayload) > 0 {
		var payloadData RawPayloadData

		if err := json.Unmarshal(rawRead.RawPayload, &payloadData); err == nil &&
			len(payloadData.Reads) > 0 {

			var reads []ProtoRead
			for _, readMap := range payloadData.Reads {
				pr := ProtoRead{
					TimestampNs: toInt64(readMap["timestamp_ns"]),
					Confidence:  toInt(readMap["confidence"]),
					AntennaID:   toInt(readMap["antenna_id"]),
					AntennaType: toInt(readMap["antenna_type"]),
					X:           toFloat64(readMap["x"]),
					Y:           toFloat64(readMap["y"]),
					ItemID:      toString(readMap["item_id"]),
					FloorID:     toInt(readMap["floor_id"]),
				}

				reads = append(reads, pr)
			}

			bundleSiteID := payloadData.SiteID
			if bundleSiteID == "" {
				bundleSiteID = siteID
			}

			sentTimestamp := payloadData.SentTimestampMs
			if sentTimestamp == 0 {
				sentTimestamp = rawRead.InjectionTime.UnixMilli()
			}

			return &ProtoReaderBundleWrapper{
				ProtoReaderBundle: ProtoReaderBundle{
					ReaderID:        payloadData.ReaderID,
					Reads:           reads,
					SiteID:          bundleSiteID,
					SentTimestampMs: sentTimestamp,
				},
			}, nil
		}
	}

	// Real Recorder DB path:
	// raw_payload is binary, so build payload from RawReads columns.
	read := ProtoRead{
		TimestampNs: rawRead.Timestamp.UnixNano(),
		Confidence:  rawRead.Confidence,
		AntennaID:   rawRead.AntennaID,
		AntennaType: rawRead.AntennaTypeID,
		X:           rawRead.TagX,
		Y:           rawRead.TagY,
		ItemID:      rawRead.TagID,
		FloorID:     rawRead.FloorID,
	}

	readerID := toInt(rawRead.ReaderID)

	return &ProtoReaderBundleWrapper{
		ProtoReaderBundle: ProtoReaderBundle{
			ReaderID:        readerID,
			Reads:           []ProtoRead{read},
			SiteID:          siteID,
			SentTimestampMs: rawRead.InjectionTime.UnixMilli(),
		},
	}, nil
}