func GetRawReads(db *sql.DB, sessionID string) ([]RawRead, error) {
	rows, err := db.Query(`
		SELECT read_id,
		       recording_session_id,
		       tag_id,
		       reader_id,
		       antenna_id,
		       antenna_type,
		       source_timestamp_utc,
		       injection_time_utc,
		       confidence,
		       rssi,
		       tag_x,
		       tag_y,
		       floor_id,
		       raw_payload
		FROM RawReads
		WHERE recording_session_id = ?
		ORDER BY injection_time_utc ASC, read_id ASC
	`, sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to query raw reads: %w", err)
	}
	defer rows.Close()

	var reads []RawRead

	for rows.Next() {
		var r RawRead

		err := rows.Scan(
			&r.ReadID,
			&r.RecordingSessionID,
			&r.TagID,
			&r.ReaderID,
			&r.AntennaID,
			&r.AntennaTypeID,
			&r.SourceTimestampUtc,
			&r.InjectionTimeUtc,
			&r.Confidence,
			&r.RSSI,
			&r.TagX,
			&r.TagY,
			&r.FloorID,
			&r.RawPayload,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan raw read: %w", err)
		}

		r.Timestamp, _ = parseTime(r.SourceTimestampUtc)
		r.InjectionTime, _ = parseTime(r.InjectionTimeUtc)

		reads = append(reads, r)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating raw reads: %w", err)
	}

	return reads, nil
}