// GetRecordedSiteInfo retrieves the SiteInformation stored in the recording
// without filtering by the configured Site ID.
func GetRecordedSiteInfo(db *sql.DB) (*SiteInformation, error) {
	row := db.QueryRow(`
		SELECT site_information_id,
		       recording_session_id,
		       site_id,
		       site_name,
		       site_json
		FROM SiteInformation
		LIMIT 1
	`)

	var si SiteInformation

	err := row.Scan(
		&si.SiteInformationID,
		&si.RecordingSessionID,
		&si.SiteID,
		&si.SiteName,
		&si.SiteJSON,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get recorded site information: %w",
			err,
		)
	}

	return &si, nil
}