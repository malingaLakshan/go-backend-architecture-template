func GetRawReadTimeRange(
	db *sql.DB,
	sessionID string,
) (first, last time.Time, err error) {
	var firstStr, lastStr string

	err = db.QueryRow(`
		SELECT MIN(injection_time_utc),
		       MAX(injection_time_utc)
		FROM RawReads
		WHERE recording_session_id = ?
	`, sessionID).Scan(&firstStr, &lastStr)
	if err != nil {
		return time.Time{}, time.Time{},
			fmt.Errorf("failed to get time range: %w", err)
	}

	first, _ = parseTime(firstStr)
	last, _ = parseTime(lastStr)

	return first, last, nil
}