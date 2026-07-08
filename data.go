for rows.Next() {
	var r RawRead

	var rssi sql.NullFloat64
	var tagX sql.NullFloat64
	var tagY sql.NullFloat64
	var floorID sql.NullInt64
	var confidence sql.NullInt64
	var antennaID sql.NullInt64
	var antennaType sql.NullInt64

	err := rows.Scan(
		&r.ReadID,
		&r.RecordingSessionID,
		&r.TagID,
		&r.ReaderID,
		&antennaID,
		&antennaType,
		&r.SourceTimestampUtc,
		&r.InjectionTimeUtc,
		&confidence,
		&rssi,
		&tagX,
		&tagY,
		&floorID,
		&r.RawPayload,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to scan raw read: %w", err)
	}

	if antennaID.Valid {
		r.AntennaID = int(antennaID.Int64)
	}

	if antennaType.Valid {
		r.AntennaTypeID = int(antennaType.Int64)
	}

	if confidence.Valid {
		r.Confidence = int(confidence.Int64)
	}

	if rssi.Valid {
		r.RSSI = rssi.Float64
	}

	if tagX.Valid {
		r.TagX = tagX.Float64
	}

	if tagY.Valid {
		r.TagY = tagY.Float64
	}

	if floorID.Valid {
		r.FloorID = int(floorID.Int64)
	}

	r.Timestamp, _ = parseTime(r.SourceTimestampUtc)
	r.InjectionTime, _ = parseTime(r.InjectionTimeUtc)

	reads = append(reads, r)
}