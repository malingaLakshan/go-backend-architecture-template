// runValidate handles the validate command.
// Supports config-file mode and direct-argument mode.
func runValidate(flags *Flags) int {
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
				"Error: config requires: recording_file, target_url, site_id",
			)
			return 1
		}

		flags.File = cfg.RecordingFile
		flags.TargetURL = cfg.TargetURL
		flags.SiteID = cfg.SiteID
	}

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

	log.Info("Validation started for site %s", flags.SiteID)

	// 1. Read recorded SiteGraph.
	db, err := sqlite.OpenReadOnly(flags.File)
	if err != nil {
		log.Error("%v", err)

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

	siteInfo, err := recording.GetSiteInfo(db, flags.SiteID)
	if err != nil {
		log.Error(
			"Failed to load recorded site info: %v",
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

	// 2. Fetch target SiteGraph.
	client := site.NewClient(flags.TargetURL)

	targetSite, err := client.FetchValidationSite(flags.SiteID)
	if err != nil {
		log.Error(
			"Failed to fetch target site for Site ID %s: %v",
			flags.SiteID,
			err,
		)

		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "Validation failed.")
		fmt.Fprintf(
			os.Stderr,
			"No target site was found for Site ID: %s\n",
			flags.SiteID,
		)
		fmt.Fprintf(
			os.Stderr,
			"Details: %v\n",
			err,
		)

		return 1
	}

	// 3. Display summaries and validation results.
	fmt.Println()
	fmt.Println("Site Configuration Validation")
	fmt.Println()

	printSiteSummaryBlock("Recorded Site", recordedSite)
	printSiteSummaryBlock("Target Site", targetSite)

	fmt.Println("Validation Results")
	fmt.Println()

	// 4. Validate structure.
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