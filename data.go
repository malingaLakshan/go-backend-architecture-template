// runStop handles the stop command.
//
// It automatically discovers config.json when -config is not supplied,
// reads the configured mock-server port, finds the Windows process listening
// on that port, verifies that it belongs to rre.exe, and terminates its
// process tree.
//
// All technical messages use the existing common Replay Engine logger.
// No PID file or separate stop log file is created.
func runStop(flags *Flags) int {
	// Resolve either the explicit -config path or the standard config.json.
	configPath, err := resolveConfigPath(flags.Config)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Stop failed.")
		fmt.Fprintf(
			os.Stderr,
			"Configuration error: %v\n",
			err,
		)

		return 1
	}

	// Save the resolved path so both configuration loading and common
	// logger initialization use the same config.json.
	flags.Config = configPath

	log, closer := initLogger(configPath)
	defer closer.Close()

	cfg, err := config.Load(configPath)
	if err != nil {
		log.Error(
			"Failed to load configuration for stop command: %v",
			err,
		)

		fmt.Fprintln(os.Stderr, "Stop failed.")
		fmt.Fprintf(
			os.Stderr,
			"Failed to load configuration: %v\n",
			err,
		)

		return 1
	}

	// The explicit -port value starts with the parser default.
	// The configured mock_port takes precedence so `rre stop` matches
	// the port used by `rre serve`.
	port := flags.Port

	if cfg.MockPort > 0 {
		port = cfg.MockPort
	}

	if port <= 0 {
		port = 8080
	}

	log.Info(
		"Stop command started for port %d",
		port,
	)

	pid, err := findListeningPID(port)
	if err != nil {
		log.Error(
			"Failed to inspect port %d: %v",
			port,
			err,
		)

		fmt.Fprintln(os.Stderr, "Stop failed.")
		fmt.Fprintf(
			os.Stderr,
			"Failed to inspect port %d: %v\n",
			port,
			err,
		)

		return 1
	}

	if pid == 0 {
		log.Warn(
			"No process is listening on port %d",
			port,
		)

		fmt.Printf(
			"No Replay Engine server is running on port %d.\n",
			port,
		)

		return 0
	}

	imageName, err := processImageName(pid)
	if err != nil {
		log.Error(
			"Failed to identify process %d on port %d: %v",
			pid,
			port,
			err,
		)

		fmt.Fprintln(os.Stderr, "Stop failed.")
		fmt.Fprintf(
			os.Stderr,
			"Failed to identify the process on port %d: %v\n",
			port,
			err,
		)

		return 1
	}

	// Safety check: never terminate another application that happens
	// to be using the configured port.
	if !isReplayEngineProcess(imageName) {
		log.Error(
			"Refusing to terminate PID %d on port %d because the executable is %q",
			pid,
			port,
			imageName,
		)

		fmt.Fprintln(os.Stderr, "Stop refused.")
		fmt.Fprintf(
			os.Stderr,
			"Port %d is used by %s (PID %d), not rre.exe.\n",
			port,
			imageName,
			pid,
		)

		return 1
	}

	log.Info(
		"Replay Engine server found on port %d with PID %d",
		port,
		pid,
	)

	if err := terminateProcessTree(pid); err != nil {
		log.Error(
			"Failed to stop Replay Engine process %d: %v",
			pid,
			err,
		)

		fmt.Fprintln(os.Stderr, "Stop failed.")
		fmt.Fprintf(
			os.Stderr,
			"Failed to stop Replay Engine server on port %d: %v\n",
			port,
			err,
		)

		return 1
	}

	log.Success(
		"Replay Engine server stopped successfully on port %d",
		port,
	)

	fmt.Printf(
		"Replay Engine server stopped successfully on port %d.\n",
		port,
	)

	return 0
}