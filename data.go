// runValidate handles the validate command.
// Supports config-file mode and direct-argument mode.
func runValidate(flags *Flags) int {
	// Load command values from config when -config is supplied.
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

	// Validate required values.
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

	// 1. Open the recording database.
	db, err := sqlite.OpenReadOnly(flags.File)
	if err != nil {
		log.Error("Failed to open recording database: %v", err)

		fmt.Println()
		fmt.Println("Validation failed.")
		fmt.Printf("Failed to open recording database: %v\n", err)

		return 1
	}
	defer db.Close()

	// 2. Read the recorded SiteGraph.
	siteInfo, err := recording.GetSiteInfo(db, flags.SiteID)
	if err != nil {
		log.Error("Failed to load recorded site information: %v", err)

		fmt.Println()
		fmt.Println("Validation failed.")
		fmt.Printf(
			"Failed to load recorded site information: %v\n",
			err,
		)

		return 1
	}

	recordedSite, err := site.ParseSiteGraph(siteInfo.SiteJSON)
	if err != nil {
		log.Error("Failed to parse recorded SiteGraph: %v", err)

		fmt.Println()
		fmt.Println("Validation failed.")
		fmt.Printf(
			"Failed to parse recorded SiteGraph: %v\n",
			err,
		)

		return 1
	}

	// Read only the root Site ID from the recorded SiteGraph JSON.
	// This avoids depending on a model field such as recordedSite.ID.
	var recordedSiteMetadata struct {
		ID string `json:"id"`
	}

	if err := json.Unmarshal(
		[]byte(siteInfo.SiteJSON),
		&recordedSiteMetadata,
	); err != nil {
		log.Error(
			"Failed to read recorded Site ID from SiteGraph: %v",
			err,
		)

		fmt.Println()
		fmt.Println("Validation failed.")
		fmt.Printf(
			"Failed to read the recorded Site ID: %v\n",
			err,
		)

		return 1
	}

	configuredSiteID := strings.TrimSpace(flags.SiteID)
	recordedSiteID := strings.TrimSpace(recordedSiteMetadata.ID)

	// Stop immediately when the configured Site ID is incorrect.
	if configuredSiteID != recordedSiteID {
		log.Error(
			"Configured Site ID %q does not match recorded Site ID %q",
			configuredSiteID,
			recordedSiteID,
		)

		fmt.Println()
		fmt.Println("Site Configuration Validation")
		fmt.Println()
		fmt.Println("Validation Results")
		fmt.Println()
		fmt.Println("X Site ID mismatch")
		fmt.Println()
		fmt.Printf("  Configured Site ID: %s\n", configuredSiteID)
		fmt.Printf("  Recorded Site ID:   %s\n", recordedSiteID)
		fmt.Println()
		fmt.Println("Validation failed.")
		fmt.Println(
			"The configured Site ID does not match the recorded Site ID.",
		)

		return 1
	}

	// 3. Fetch the target SiteGraph.
	client := site.NewClient(flags.TargetURL)

	targetSite, err := client.FetchValidationSite(flags.SiteID)
	if err != nil {
		log.Error(
			"Failed to fetch target site for Site ID %s: %v",
			flags.SiteID,
			err,
		)

		fmt.Println()
		fmt.Println("Site Configuration Validation")
		fmt.Println()
		fmt.Println("Validation Results")
		fmt.Println()
		fmt.Printf(
			"X No target site found for Site ID: %s\n",
			flags.SiteID,
		)
		fmt.Println()
		fmt.Println("Validation failed.")

		return 1
	}

	// 4. Display site summaries.
	fmt.Println()
	fmt.Println("Site Configuration Validation")
	fmt.Println()

	printSiteSummaryBlock("Recorded Site", recordedSite)
	printSiteSummaryBlock("Target Site", targetSite)

	fmt.Println("Validation Results")
	fmt.Println()

	// 5. Validate the complete structure.
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