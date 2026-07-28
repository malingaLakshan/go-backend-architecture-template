// resolveConfigPath returns an explicitly supplied configuration path,
// or automatically discovers the Replay Engine config.json.
//
// Search order:
//  1. Explicit -config path
//  2. configs/config.json relative to the current working directory
//  3. config.json relative to the current working directory
//  4. configs/config.json relative to rre.exe
//  5. config.json relative to rre.exe
func resolveConfigPath(configPath string) (string, error) {
	configPath = strings.TrimSpace(configPath)

	if configPath != "" {
		info, err := os.Stat(configPath)
		if err != nil {
			return "", fmt.Errorf(
				"configuration file %q is not accessible: %w",
				configPath,
				err,
			)
		}

		if info.IsDir() {
			return "", fmt.Errorf(
				"configuration path %q is a directory",
				configPath,
			)
		}

		absolutePath, err := filepath.Abs(configPath)
		if err != nil {
			return configPath, nil
		}

		return absolutePath, nil
	}

	candidates := []string{
		filepath.Join("configs", "config.json"),
		"config.json",
	}

	executablePath, err := os.Executable()
	if err == nil {
		executableDirectory := filepath.Dir(executablePath)

		candidates = append(
			candidates,
			filepath.Join(
				executableDirectory,
				"configs",
				"config.json",
			),
			filepath.Join(
				executableDirectory,
				"config.json",
			),
		)
	}

	for _, candidate := range candidates {
		info, err := os.Stat(candidate)
		if err != nil || info.IsDir() {
			continue
		}

		absolutePath, err := filepath.Abs(candidate)
		if err != nil {
			return candidate, nil
		}

		return absolutePath, nil
	}

	return "", fmt.Errorf(
		"config.json was not found in the working directory, configs directory, or executable directory",
	)
}