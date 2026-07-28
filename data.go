// runStop handles the stop command.
//
// It locates the process listening on the configured mock-server port,
// verifies that the process belongs to the Replay Engine, and terminates
// the process tree.
//
// Logging is written through the existing common logger configuration.
func runStop(flags *Flags) int {
	log, closer := initLogger(flags.Config)
	defer closer.Close()

	port := flags.Port

	// Use the configured mock-server port when a configuration file is supplied.
	if flags.Config != "" {
		cfg, err := config.Load(flags.Config)
		if err != nil {
			log.Error("Failed to load configuration for stop command: %v", err)

			fmt.Fprintln(os.Stderr)
			fmt.Fprintln(os.Stderr, "Stop failed.")
			fmt.Fprintf(
				os.Stderr,
				"Failed to load configuration: %v\n",
				err,
			)

			return 1
		}

		if cfg.MockPort > 0 {
			port = cfg.MockPort
		}
	}

	// Preserve the current Replay Engine default when no port was supplied.
	if port <= 0 {
		port = 8080
	}

	log.Info("Stop command started for port %d", port)

	pid, err := findListeningPID(port)
	if err != nil {
		log.Error(
			"Failed to inspect port %d: %v",
			port,
			err,
		)

		fmt.Fprintln(os.Stderr)
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

		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "Stop failed.")
		fmt.Fprintf(
			os.Stderr,
			"Failed to identify process %d on port %d: %v\n",
			pid,
			port,
			err,
		)

		return 1
	}

	// Safety check: never terminate an unrelated program using the same port.
	if !isReplayEngineProcess(imageName) {
		log.Error(
			"Refusing to terminate process %d on port %d because its executable is %q",
			pid,
			port,
			imageName,
		)

		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "Stop refused.")
		fmt.Fprintf(
			os.Stderr,
			"Port %d is being used by %s (PID %d), not rre.exe.\n",
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

		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "Stop failed.")
		fmt.Fprintf(
			os.Stderr,
			"Failed to stop Replay Engine process %d: %v\n",
			pid,
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