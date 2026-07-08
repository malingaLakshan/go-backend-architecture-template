func runMockServer(flags *Flags) int {
	if flags.Config != "" {
		return runMockServerWithConfig(flags)
	}

	if err := mocktarget.StartServer(flags.Port); err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] %v\n", err)
		return 1
	}

	return 0
}