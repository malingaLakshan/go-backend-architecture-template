func Load(path string) (*RunConfig, error) {
	cleanPath, err := validateConfigPath(path)
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(cleanPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", cleanPath, err)
	}

	var cfg RunConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file %s: %w", cleanPath, err)
	}

	return &cfg, nil
}

func validateConfigPath(path string) (string, error) {
	if strings.TrimSpace(path) == "" {
		return "", fmt.Errorf("config path is required")
	}

	cleanPath := filepath.Clean(path)

	// Only allow JSON config files inside configs folder.
	if filepath.IsAbs(cleanPath) {
		return "", fmt.Errorf("absolute config paths are not allowed")
	}

	if filepath.Ext(cleanPath) != ".json" {
		return "", fmt.Errorf("config file must be a .json file")
	}

	configDir := "configs" + string(os.PathSeparator)

	if cleanPath != "configs" && !strings.HasPrefix(cleanPath, configDir) {
		return "", fmt.Errorf("config file must be inside configs directory")
	}

	if strings.Contains(cleanPath, "..") {
		return "", fmt.Errorf("config path cannot contain parent directory traversal")
	}

	return cleanPath, nil
}