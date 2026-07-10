func Load(path string) (*RunConfig, error) {
	data, configName, err := readAllowedConfigFile(path)
	if err != nil {
		return nil, err
	}

	var cfg RunConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file %s: %w", configName, err)
	}

	return &cfg, nil
}

func readAllowedConfigFile(path string) ([]byte, string, error) {
	cleanPath := filepath.ToSlash(filepath.Clean(strings.TrimSpace(path)))

	switch cleanPath {
	case "configs/pass_config.json":
		data, err := os.ReadFile("configs/pass_config.json")
		return data, cleanPath, err

	case "configs/fail_config.json":
		data, err := os.ReadFile("configs/fail_config.json")
		return data, cleanPath, err

	case "configs/wrong_site_config.json":
		data, err := os.ReadFile("configs/wrong_site_config.json")
		return data, cleanPath, err

	default:
		return nil, "", fmt.Errorf("unsupported config file: %s", path)
	}
}