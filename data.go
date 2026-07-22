// Package config provides the RunConfig model for config-file-based commands.
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// RunConfig holds configuration loaded from a JSON config file.
//
// Fields:
//
//	recording_file       - path to the Recorder SQLite file
//	target_url           - target Resonate HTTP base URL
//	site_id              - target site ID used for validation
//	mock_port            - port used by the mock server
//	site_graph_directory - directory containing SiteGraph JSON files
type RunConfig struct {
	RecordingFile      string `json:"recording_file"`
	TargetURL          string `json:"target_url"`
	SiteID             string `json:"site_id"`
	MockPort           int    `json:"mock_port,omitempty"`
	SiteGraphDirectory string `json:"site_graph_directory,omitempty"`
}

// Load reads and parses a RunConfig from the provided JSON file path.
//
// The configuration file may be stored in any directory.
func Load(path string) (*RunConfig, error) {
	data, configPath, err := readConfigFile(path)
	if err != nil {
		return nil, err
	}

	var cfg RunConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf(
			"failed to parse config file %s: %w",
			configPath,
			err,
		)
	}

	return &cfg, nil
}

// Save writes the RunConfig to the provided JSON file path.
//
// The file is written atomically by first creating a temporary file in the
// same directory and then renaming it to replace the original file.
func (cfg *RunConfig) Save(path string) error {
	if cfg == nil {
		return fmt.Errorf("config cannot be nil")
	}

	cleanPath, err := resolveConfigPath(path)
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	configDirectory := filepath.Dir(cleanPath)

	info, err := os.Stat(configDirectory)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf(
				"config directory does not exist: %s",
				configDirectory,
			)
		}

		return fmt.Errorf(
			"failed to access config directory %s: %w",
			configDirectory,
			err,
		)
	}

	if !info.IsDir() {
		return fmt.Errorf(
			"config directory path is not a directory: %s",
			configDirectory,
		)
	}

	tempFile, err := os.CreateTemp(configDirectory, "config-*.tmp")
	if err != nil {
		return fmt.Errorf("failed to create temporary config file: %w", err)
	}

	tempPath := tempFile.Name()
	removeTempFile := true

	defer func() {
		if removeTempFile {
			_ = os.Remove(tempPath)
		}
	}()

	if _, err := tempFile.Write(data); err != nil {
		_ = tempFile.Close()
		return fmt.Errorf("failed to write temporary config file: %w", err)
	}

	if err := tempFile.Sync(); err != nil {
		_ = tempFile.Close()
		return fmt.Errorf("failed to sync temporary config file: %w", err)
	}

	if err := tempFile.Close(); err != nil {
		return fmt.Errorf("failed to close temporary config file: %w", err)
	}

	// Windows does not allow os.Rename to replace an existing file.
	// Remove the existing target before renaming the temporary file.
	if _, err := os.Stat(cleanPath); err == nil {
		if err := os.Remove(cleanPath); err != nil {
			return fmt.Errorf(
				"failed to replace existing config file %s: %w",
				cleanPath,
				err,
			)
		}
	} else if !os.IsNotExist(err) {
		return fmt.Errorf(
			"failed to access existing config file %s: %w",
			cleanPath,
			err,
		)
	}

	if err := os.Rename(tempPath, cleanPath); err != nil {
		return fmt.Errorf(
			"failed to finalize config file %s: %w",
			cleanPath,
			err,
		)
	}

	removeTempFile = false
	return nil
}

// readConfigFile reads a configuration file from the provided path.
//
// Both relative and absolute paths are supported. The file must exist,
// must be a regular file, and must have a .json extension.
func readConfigFile(path string) ([]byte, string, error) {
	cleanPath, err := resolveConfigPath(path)
	if err != nil {
		return nil, "", err
	}

	info, err := os.Stat(cleanPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, "", fmt.Errorf(
				"config file not found: %s",
				cleanPath,
			)
		}

		return nil, "", fmt.Errorf(
			"failed to access config file %s: %w",
			cleanPath,
			err,
		)
	}

	if info.IsDir() {
		return nil, "", fmt.Errorf(
			"config path must point to a file: %s",
			cleanPath,
		)
	}

	if !info.Mode().IsRegular() {
		return nil, "", fmt.Errorf(
			"config path is not a regular file: %s",
			cleanPath,
		)
	}

	data, err := os.ReadFile(cleanPath)
	if err != nil {
		return nil, "", fmt.Errorf(
			"failed to read config file %s: %w",
			cleanPath,
			err,
		)
	}

	return data, cleanPath, nil
}

// resolveConfigPath validates and resolves a config file path.
//
// The path may point to any directory, but it must use the .json extension.
func resolveConfigPath(path string) (string, error) {
	trimmedPath := strings.TrimSpace(path)
	if trimmedPath == "" {
		return "", fmt.Errorf("config file path is required")
	}

	cleanPath := filepath.Clean(trimmedPath)

	if !strings.EqualFold(filepath.Ext(cleanPath), ".json") {
		return "", fmt.Errorf(
			"config file must have a .json extension: %s",
			path,
		)
	}

	absolutePath, err := filepath.Abs(cleanPath)
	if err != nil {
		return "", fmt.Errorf(
			"failed to resolve config file path %s: %w",
			path,
			err,
		)
	}

	return absolutePath, nil
}