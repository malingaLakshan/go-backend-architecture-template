// runValidate handles the validate command.
// Supports config-file mode and direct-argument mode.
func runValidate(flags *Flags) int {
	// Load values from the config file when -config is provided.
	if flags.Config != "" {
		cfg, err := config.Load(flags.Config)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			return 1
		}

		if cfg.RecordingFile == "" ||
			cfg.TargetURL == "" ||
			cfg.SiteID == "" {

			fmt.Fprintln(
				os.Stderr,
				"Error: config requires recording_file, target_url, and site_id",
			)
			return 1
		}

		flags.File = cfg.RecordingFile
		flags.TargetURL = cfg.TargetURL
		flags.SiteID = cfg.SiteID
	}

	// Validate required direct/config-loaded values.
	if flags.File == "" ||
		flags.TargetURL == "" ||
		flags.SiteID == "" {

		fmt.Fprintln(
			os.Stderr,
			"Error: -file, -target-url, and -site-id are required for validate",
		)
		return 1
	}

	log, closer := initLogger(flags.Config)
	defer closer.Close()

	configuredSiteID := strings.TrimSpace(flags.SiteID)

	log.Info(
		"Validation started for site %s",
		configuredSiteID,
	)

	// 1. Open the recording database.
	db, err := sqlite.OpenReadOnly(flags.File)
	if err != nil {
		log.Error("Failed to open recording database: %v", err)

		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "Validation failed.")
		fmt.Fprintf(
			os.Stderr,
			"Failed to open recording database: %v\n",
			err,
		)

		return 1
	}
	defer db.Close()

	// 2. Read the recorded SiteGraph.
	siteInfo, err := recording.GetSiteInfo(db, flags.SiteID)
	if err != nil {
		log.Error(
			"Failed to load recorded site information: %v",
			err,
		)

		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "Validation failed.")
		fmt.Fprintf(
			os.Stderr,
			"Failed to load recorded site information: %v\n",
			err,
		)

		return 1
	}

	recordedSite, err := site.ParseSiteGraph(siteInfo.SiteJSON)
	if err != nil {
		log.Error(
			"Failed to parse recorded SiteGraph: %v",
			err,
		)

		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "Validation failed.")
		fmt.Fprintf(
			os.Stderr,
			"Failed to parse recorded SiteGraph: %v\n",
			err,
		)

		return 1
	}

	// 3. Compare configured Site ID with the recorded SiteGraph ID.
	recordedSiteID := strings.TrimSpace(recordedSite.ID)

	if configuredSiteID != recordedSiteID {
		log.Error(
			"Configured Site ID %q does not match recorded Site ID %q",
			configuredSiteID,
			recordedSiteID,
		)

		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "Validation failed.")
		fmt.Fprintln(
			os.Stderr,
			"Configured Site ID does not match the recorded Site ID.",
		)
		fmt.Fprintf(
			os.Stderr,
			"  Configured Site ID: %s\n",
			configuredSiteID,
		)
		fmt.Fprintf(
			os.Stderr,
			"  Recorded Site ID:   %s\n",
			recordedSiteID,
		)

		return 1
	}

	// 4. Fetch the target SiteGraph.
	client := site.NewClient(flags.TargetURL)

	targetSite, err := client.FetchValidationSite(configuredSiteID)
	if err != nil {
		log.Error(
			"Failed to fetch target site: %v",
			err,
		)

		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "Validation failed.")
		fmt.Fprintf(
			os.Stderr,
			"No target site was found for Site ID %q.\n",
			configuredSiteID,
		)
		fmt.Fprintf(
			os.Stderr,
			"Details: %v\n",
			err,
		)

		return 1
	}

	// 5. Confirm that the returned target SiteGraph has the requested ID.
	targetSiteID := strings.TrimSpace(targetSite.ID)

	if targetSiteID != configuredSiteID {
		log.Error(
			"Target SiteGraph ID %q does not match requested Site ID %q",
			targetSiteID,
			configuredSiteID,
		)

		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "Validation failed.")
		fmt.Fprintln(
			os.Stderr,
			"Target SiteGraph ID does not match the requested Site ID.",
		)
		fmt.Fprintf(
			os.Stderr,
			"  Requested Site ID: %s\n",
			configuredSiteID,
		)
		fmt.Fprintf(
			os.Stderr,
			"  Target Site ID:    %s\n",
			targetSiteID,
		)

		return 1
	}

	// 6. Display summaries and validation results.
	fmt.Println()
	fmt.Println("Site Configuration Validation")
	fmt.Println()

	printSiteSummaryBlock("Recorded Site", recordedSite)
	printSiteSummaryBlock("Target Site", targetSite)

	fmt.Println("Validation Results")
	fmt.Println()

	// 7. Validate the complete site structure.
	result := site.ValidateStructure(recordedSite, targetSite)
	printValidationResults(result)

	if result.Passed {
		log.Success("Validation passed")
		return 0
	}

	log.Error(
		"Validation failed with %d mismatch(es)",
		len(result.Mismatches),
	)

	return 1
}