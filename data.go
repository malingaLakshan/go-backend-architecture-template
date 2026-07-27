if flags.SiteID != siteInfo.SiteID {
	log.Error(
		"Configured Site ID %q does not match recorded Site ID %q",
		flags.SiteID,
		siteInfo.SiteID,
	)

	fmt.Println()
	fmt.Println("Site Configuration Validation")
	fmt.Println()
	fmt.Println("Validation Results")
	fmt.Println()
	fmt.Println("X Site ID mismatch")
	fmt.Println()
	fmt.Printf("  Configured Site ID: %s\n", flags.SiteID)
	fmt.Printf("  Recorded Site ID:   %s\n", siteInfo.SiteID)
	fmt.Println()
	fmt.Println("Validation failed.")

	return 1
}