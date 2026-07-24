func runMockServer(flags *Flags) int {
	log, closer := initLogger(flags.Config)
	defer closer.Close()

	log.Info("Mock server command started")

	if flags.Config == "" {
		log.Info("Starting mock server without config on port %d", flags.Port)

		if err := mocktarget.StartServer(flags.Port); err != nil {
			log.Error("Mock server failed: %v", err)
			return 1
		}

		log.Success("Mock server stopped")
		return 0
	}

	cfg, err := config.Load(flags.Config)
	if err != nil {
		log.Error("Failed to load configuration: %v", err)
		return 1
	}

	siteGraphDir := cfg.SiteGraphDirectory
	if siteGraphDir == "" {
		siteGraphDir = "configs/sites"
	}

	log.Info("Loading SiteGraphs from directory: %s", siteGraphDir)

	store, err := mocktarget.LoadSiteStore(siteGraphDir)
	if err != nil {
		log.Error("Failed to load SiteGraphs: %v", err)
		return 1
	}

	port := flags.Port
	if cfg.MockPort > 0 && flags.Port == 8080 {
		port = cfg.MockPort
	}

	log.Success("Mock server starting on port %d", port)

	if err := mocktarget.StartServerWithSiteStore(
		port,
		store,
		flags.Config,
		siteGraphDir,
	); err != nil {
		log.Error("Mock server stopped with error: %v", err)
		return 1
	}

	log.Success("Mock server stopped")
	return 0
}