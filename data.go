if configuredSiteID != recordedSiteID {
	log.Error(
		"Configured Site ID %q does not match recorded Site ID %q",
		configuredSiteID,
		recordedSiteID,
	)

	fmt.Println()
	fmt.Println("Validation Results")
	fmt.Println()
	fmt.Println("✗ Site ID mismatch")
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