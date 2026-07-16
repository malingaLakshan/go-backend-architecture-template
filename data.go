// Package mocktarget provides a mock Resonate HTTP server for testing replay.
package mocktarget

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"resonate-replay-engine-mvp-cli/internal/site"
)

// SiteSummary holds display-only information for one loaded SiteGraph.
type SiteSummary struct {
	SiteID       string
	Name         string
	Floors       int
	Regions      int
	Readers      int
	AntennaPorts int
}

// SiteStore holds all loaded SiteGraph raw JSON indexed by root Site ID.
//
// The original JSON bytes are preserved exactly as read from the source files.
// This allows GET /sites/{siteId} to return the complete SiteGraph without
// losing unknown fields through JSON re-serialization.
type SiteStore struct {
	raw       map[string][]byte
	summaries []SiteSummary
}

const approvedSiteGraphDirectory = "configs/sites"

// LoadSiteStore loads all SiteGraph JSON files from the approved directory.
//
// Security rules:
//
//   - the configured directory must remain inside configs/sites
//   - path traversal is rejected
//   - symbolic-link directory escapes are rejected
//   - only regular .json files are processed
//   - file names containing path separators are rejected
//   - symbolic-link files are rejected
//   - invalid JSON is rejected
//   - missing root Site IDs are rejected
//   - duplicate root Site IDs are rejected
//
// The root JSON "id" field is the source of truth. The file name is not used
// as the Site ID.
func LoadSiteStore(directory string) (*SiteStore, error) {
	cleanDir, err := validateSiteGraphDirectory(directory)
	if err != nil {
		return nil, err
	}

	return loadSiteStoreFromDirectory(cleanDir)
}

// loadSiteStoreFromDirectory performs SiteGraph loading from a supplied
// directory.
//
// LoadSiteStore calls this only after approved-directory validation.
// Unit tests may call this helper with a temporary test directory.
func loadSiteStoreFromDirectory(directory string) (*SiteStore, error) {
	if strings.TrimSpace(directory) == "" {
		return nil, fmt.Errorf("SiteGraph directory is required")
	}

	cleanDir := filepath.Clean(directory)

	absoluteDir, err := filepath.Abs(cleanDir)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to resolve SiteGraph directory %q: %w",
			cleanDir,
			err,
		)
	}

	directoryInfo, err := os.Stat(absoluteDir)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to inspect SiteGraph directory %q: %w",
			cleanDir,
			err,
		)
	}

	if !directoryInfo.IsDir() {
		return nil, fmt.Errorf(
			"SiteGraph path is not a directory: %s",
			cleanDir,
		)
	}

	// Create a filesystem rooted at the already resolved directory.
	// All subsequent reads use names relative to this root and cannot use an
	// absolute path.
	siteGraphFS := os.DirFS(absoluteDir)

	entries, err := fs.ReadDir(siteGraphFS, ".")
	if err != nil {
		return nil, fmt.Errorf(
			"failed to read SiteGraph directory %q: %w",
			cleanDir,
			err,
		)
	}

	store := &SiteStore{
		raw:       make(map[string][]byte),
		summaries: make([]SiteSummary, 0),
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()

		if !strings.EqualFold(filepath.Ext(name), ".json") {
			continue
		}

		if err := validateSiteGraphFileName(name); err != nil {
			return nil, err
		}

		// Do not follow symbolic-link files. A SiteGraph must be a regular
		// JSON file directly within the selected directory.
		if entry.Type()&os.ModeSymlink != 0 {
			return nil, fmt.Errorf(
				"symbolic-link SiteGraph files are not allowed: %s",
				name,
			)
		}

		entryInfo, infoErr := entry.Info()
		if infoErr != nil {
			return nil, fmt.Errorf(
				"failed to inspect SiteGraph file %s: %w",
				name,
				infoErr,
			)
		}

		if !entryInfo.Mode().IsRegular() {
			return nil, fmt.Errorf(
				"SiteGraph file must be a regular file: %s",
				name,
			)
		}

		// Security: read through a filesystem rooted at the validated
		// directory. No user-controlled full filesystem path is constructed.
		data, readErr := fs.ReadFile(siteGraphFS, name)
		if readErr != nil {
			return nil, fmt.Errorf(
				"failed to read SiteGraph file %s: %w",
				name,
				readErr,
			)
		}

		if len(data) == 0 {
			return nil, fmt.Errorf(
				"SiteGraph file %s is empty",
				name,
			)
		}

		// Extract only root metadata here. ParseSiteGraph performs the
		// structural parsing needed for summary counts.
		var minimalSite struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		}

		if jsonErr := json.Unmarshal(data, &minimalSite); jsonErr != nil {
			return nil, fmt.Errorf(
				"invalid JSON in SiteGraph file %s: %w",
				name,
				jsonErr,
			)
		}

		minimalSite.ID = strings.TrimSpace(minimalSite.ID)

		if minimalSite.ID == "" {
			return nil, fmt.Errorf(
				"SiteGraph file %s is missing required root 'id' field",
				name,
			)
		}

		if _, exists := store.raw[minimalSite.ID]; exists {
			return nil, fmt.Errorf(
				"duplicate Site ID %q found in SiteGraph file %s",
				minimalSite.ID,
				name,
			)
		}

		// Parse only the structural hierarchy required for summary counts:
		//
		// Site
		// └── Floors
		//     ├── Regions recursively
		//     └── Readers
		//         └── Antenna ports
		validationSite, parseErr := site.ParseSiteGraph(data)
		if parseErr != nil {
			return nil, fmt.Errorf(
				"failed to parse SiteGraph hierarchy in %s: %w",
				name,
				parseErr,
			)
		}

		if validationSite.SiteID != minimalSite.ID {
			return nil, fmt.Errorf(
				"SiteGraph root ID mismatch in %s: metadata ID %q, parsed ID %q",
				name,
				minimalSite.ID,
				validationSite.SiteID,
			)
		}

		counts := site.CountValidationSite(validationSite)

		// Preserve the complete original JSON bytes.
		//
		// Copying prevents an accidental mutation if implementation details
		// change later.
		rawCopy := append([]byte(nil), data...)
		store.raw[minimalSite.ID] = rawCopy

		store.summaries = append(
			store.summaries,
			SiteSummary{
				SiteID:       minimalSite.ID,
				Name:         strings.TrimSpace(minimalSite.Name),
				Floors:       counts.Floors,
				Regions:      counts.Regions,
				Readers:      counts.Readers,
				AntennaPorts: counts.AntennaPorts,
			},
		)
	}

	return store, nil
}

// Get returns the full original SiteGraph JSON for a Site ID.
//
// The returned byte slice is a copy so callers cannot mutate the bytes stored
// inside SiteStore.
func (s *SiteStore) Get(siteID string) ([]byte, bool) {
	if s == nil {
		return nil, false
	}

	data, exists := s.raw[siteID]
	if !exists {
		return nil, false
	}

	return append([]byte(nil), data...), true
}

// Summaries returns display summaries for all loaded sites.
//
// A copy is returned so callers cannot modify SiteStore's internal slice.
func (s *SiteStore) Summaries() []SiteSummary {
	if s == nil {
		return nil
	}

	return append([]SiteSummary(nil), s.summaries...)
}

// Count returns the number of loaded SiteGraphs.
func (s *SiteStore) Count() int {
	if s == nil {
		return 0
	}

	return len(s.raw)
}

// validateSiteGraphFileName ensures a directory entry is a simple JSON file
// name and not a filesystem path.
func validateSiteGraphFileName(name string) error {
	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("empty SiteGraph file name is not allowed")
	}

	if name == "." || name == ".." {
		return fmt.Errorf(
			"unsafe SiteGraph file name rejected: %s",
			name,
		)
	}

	// A direct directory entry must equal its base name.
	if filepath.Base(name) != name {
		return fmt.Errorf(
			"unsafe SiteGraph file name rejected: %s",
			name,
		)
	}

	if strings.ContainsAny(name, `/\`) {
		return fmt.Errorf(
			"SiteGraph file name must not contain path separators: %s",
			name,
		)
	}

	// fs.ValidPath accepts slash-separated relative paths. Since separators
	// were rejected above, this confirms the value is a valid simple name.
	if !fs.ValidPath(name) {
		return fmt.Errorf(
			"invalid SiteGraph file name rejected: %s",
			name,
		)
	}

	if !strings.EqualFold(filepath.Ext(name), ".json") {
		return fmt.Errorf(
			"SiteGraph file must use the .json extension: %s",
			name,
		)
	}

	return nil
}

// validateSiteGraphDirectory ensures the configured SiteGraph directory stays
// inside the approved configs/sites area.
//
// It performs both lexical containment checking and resolved symbolic-link
// containment checking.
func validateSiteGraphDirectory(directory string) (string, error) {
	if strings.TrimSpace(directory) == "" {
		return "", fmt.Errorf("site_graph_directory is required")
	}

	cleaned := filepath.Clean(directory)
	approvedRoot := filepath.Clean(approvedSiteGraphDirectory)

	absoluteRoot, err := filepath.Abs(approvedRoot)
	if err != nil {
		return "", fmt.Errorf(
			"failed to resolve approved SiteGraph directory: %w",
			err,
		)
	}

	absoluteDirectory, err := filepath.Abs(cleaned)
	if err != nil {
		return "", fmt.Errorf(
			"failed to resolve site_graph_directory %q: %w",
			directory,
			err,
		)
	}

	if err := ensurePathWithinRoot(
		absoluteRoot,
		absoluteDirectory,
	); err != nil {
		return "", fmt.Errorf(
			"site_graph_directory must be within %s: %w",
			filepath.ToSlash(approvedRoot),
			err,
		)
	}

	// Resolve symbolic links after confirming lexical containment.
	//
	// This prevents a directory such as configs/sites/external-link from
	// escaping the approved root through a symbolic link.
	resolvedRoot, err := filepath.EvalSymlinks(absoluteRoot)
	if err != nil {
		return "", fmt.Errorf(
			"failed to resolve approved SiteGraph directory %q: %w",
			approvedRoot,
			err,
		)
	}

	resolvedDirectory, err := filepath.EvalSymlinks(absoluteDirectory)
	if err != nil {
		return "", fmt.Errorf(
			"failed to resolve site_graph_directory %q: %w",
			directory,
			err,
		)
	}

	if err := ensurePathWithinRoot(
		resolvedRoot,
		resolvedDirectory,
	); err != nil {
		return "", fmt.Errorf(
			"resolved site_graph_directory escapes approved directory: %w",
			err,
		)
	}

	directoryInfo, err := os.Stat(resolvedDirectory)
	if err != nil {
		return "", fmt.Errorf(
			"failed to inspect site_graph_directory %q: %w",
			directory,
			err,
		)
	}

	if !directoryInfo.IsDir() {
		return "", fmt.Errorf(
			"site_graph_directory is not a directory: %s",
			directory,
		)
	}

	return resolvedDirectory, nil
}

// ensurePathWithinRoot checks that candidate is either the root itself or is
// located beneath root.
func ensurePathWithinRoot(root, candidate string) error {
	relativePath, err := filepath.Rel(root, candidate)
	if err != nil {
		return fmt.Errorf(
			"failed to compare path with approved root: %w",
			err,
		)
	}

	if relativePath == "." {
		return nil
	}

	if relativePath == ".." ||
		strings.HasPrefix(relativePath, ".."+string(filepath.Separator)) ||
		filepath.IsAbs(relativePath) {
		return fmt.Errorf(
			"path %q is outside approved root %q",
			candidate,
			root,
		)
	}

	return nil
}