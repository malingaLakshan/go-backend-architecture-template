// resolveConfigPath returns the explicit config path or the standard
// Replay Engine config path.
//
// config.Load requires the path to remain relative and inside configs/.
func resolveConfigPath(configPath string) (string, error) {
	configPath = strings.TrimSpace(configPath)

	if configPath != "" {
		cleanPath := filepath.Clean(configPath)

		info, err := os.Stat(cleanPath)
		if err != nil {
			return "", fmt.Errorf(
				"configuration file %q is not accessible: %w",
				cleanPath,
				err,
			)
		}

		if info.IsDir() {
			return "", fmt.Errorf(
				"configuration path %q is a directory",
				cleanPath,
			)
		}

		return cleanPath, nil
	}

	defaultPath := filepath.Join("configs", "config.json")

	info, err := os.Stat(defaultPath)
	if err != nil {
		return "", fmt.Errorf(
			"default configuration file %q is not accessible: %w",
			defaultPath,
			err,
		)
	}

	if info.IsDir() {
		return "", fmt.Errorf(
			"default configuration path %q is a directory",
			defaultPath,
		)
	}

	return defaultPath, nil
}