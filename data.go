package mocktarget

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"replay-engine-mvp-cli/internal/site"
)

const (
	defaultPayloadFilePath = "logs/received_payloads.jsonl"

	realSiteID   = "b3489888-aacf-4451-893c-d7d994240f93"
	realSiteName = "Bentonville"

	mockSiteJSONEnv  = "RRE_MOCK_SITE_JSON"
	mockSiteJSONFile = "data/mock_site.json"
)

// Handler handles mock Resonate target requests.
type Handler struct {
	mu            sync.Mutex
	receivedCount int
	payloadFile   *os.File
	payloadsFile  string
}

// NewHandler creates a new mock target handler.
// Optional path can be passed from server.go, otherwise default log path is used.
func NewHandler(paths ...string) (*Handler, error) {
	payloadPath := defaultPayloadFilePath
	if len(paths) > 0 && strings.TrimSpace(paths[0]) != "" {
		payloadPath = paths[0]
	}

	if err := os.MkdirAll(filepath.Dir(payloadPath), 0755); err != nil {
		return nil, fmt.Errorf("failed to create payload log directory: %w", err)
	}

	file, err := os.Create(payloadPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create payload log file: %w", err)
	}

	return &Handler{
		payloadFile:  file,
		payloadsFile: payloadPath,
	}, nil
}

// HandleGetSite handles GET /sites/{siteId}.
func (h *Handler) HandleGetSite(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	siteID := strings.TrimPrefix(r.URL.Path, "/sites/")
	siteID = strings.Trim(siteID, "/")

	if siteID == "" {
		http.Error(w, "site id is required", http.StatusBadRequest)
		return
	}

	config := getMockSiteConfig(siteID)
	if config == nil {
		http.Error(w, fmt.Sprintf("site not found: %s", siteID), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(config); err != nil {
		http.Error(w, "failed to encode site config", http.StatusInternalServerError)
		return
	}
}

// HandleReaderBundle handles POST /reader-bundles requests.
func (h *Handler) HandleReaderBundle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	h.mu.Lock()
	h.receivedCount++
	count := h.receivedCount

	if h.payloadFile != nil {
		fmt.Fprintln(h.payloadFile, string(body))
	}
	h.mu.Unlock()

	var wrapper struct {
		ProtoReaderBundle struct {
			ReaderID interface{}   `json:"reader_id"`
			SiteID   string        `json:"site_id"`
			Reads    []interface{} `json:"reads"`
		} `json:"ProtoReaderBundle"`
	}

	_ = json.Unmarshal(body, &wrapper)

	if count == 1 {
		var prettyJSON map[string]interface{}
		_ = json.Unmarshal(body, &prettyJSON)

		pretty, _ := json.MarshalIndent(prettyJSON, "", "  ")

		fmt.Println()
		fmt.Println("First received replay payload:")
		fmt.Println(string(pretty))
		fmt.Println()
	}

	fmt.Printf(
		"[INFO] Received #%d | site_id=%s | reader_id=%v | reads=%d | size=%d bytes\n",
		count,
		wrapper.ProtoReaderBundle.SiteID,
		wrapper.ProtoReaderBundle.ReaderID,
		len(wrapper.ProtoReaderBundle.Reads),
		len(body),
	)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(map[string]string{
		"status": "accepted",
	})
}

// PrintSummary prints the final received payload summary.
func (h *Handler) PrintSummary() {
	h.mu.Lock()
	defer h.mu.Unlock()

	fmt.Println()
	fmt.Println("====================================================")
	fmt.Println("                 Mock Server Summary")
	fmt.Println("====================================================")
	fmt.Printf("  Total Received:        %d\n", h.receivedCount)
	fmt.Printf("  Payloads File:         %s\n", h.payloadsFile)
	fmt.Println("====================================================")

	if h.payloadFile != nil {
		_ = h.payloadFile.Close()
	}
}

// getMockSiteConfig returns a mock site configuration.
//
// Best option:
// If data/mock_site.json exists, it returns the exact real Recorder site_json.
// Otherwise it returns a fallback real-style Bentonville config.
func getMockSiteConfig(siteID string) *site.SiteConfig {
	if config := loadSiteConfigFromFile(siteID); config != nil {
		return config
	}

	if siteID != realSiteID {
		return nil
	}

	floorID := "b2c98296-6c4e-4a52-8380-17a11ddb2b2c"

	regions := []site.Region{
		{
			RegionID:      "3caf648f-fce3-424e-bc08-e711451ddaab",
			RegionName:    "World",
			RegionType:    "EXCLUSION",
			Physicality:   "VIRTUAL",
			InventoryType: "OTHER",
		},
		{
			RegionID:      "a81bb44a-62a6-4075-8703-976c4dd252e6",
			RegionName:    "Truck1",
			RegionType:    "DEPARTURE",
			Physicality:   "PHYSICAL",
			InventoryType: "EXIT",
		},
		{
			RegionID:      "e9f2e1c7-9221-4e5a-8349-08c54e1afb4a",
			RegionName:    "Truck2",
			RegionType:    "DEPARTURE",
			Physicality:   "PHYSICAL",
			InventoryType: "EXIT",
		},
	}

	readers := []site.Reader{
		{
			ReaderID:   "READER-01",
			ReaderName: "Entrance Reader",
			ReaderType: "RFID",
			IPAddress:  "192.168.1.101",
			FloorID:    floorID,
			X:          5.0,
			Y:          10.0,
			Antennas: []site.Antenna{
				{AntennaID: 1, AntennaType: 2, ReaderID: "READER-01", X: 5.0, Y: 10.0},
				{AntennaID: 2, AntennaType: 2, ReaderID: "READER-01", X: 5.5, Y: 10.0},
			},
		},
		{
			ReaderID:   "READER-02",
			ReaderName: "Fitting Room Reader",
			ReaderType: "RFID",
			IPAddress:  "192.168.1.102",
			FloorID:    floorID,
			X:          15.0,
			Y:          20.0,
			Antennas: []site.Antenna{
				{AntennaID: 3, AntennaType: 2, ReaderID: "READER-02", X: 15.0, Y: 20.0},
				{AntennaID: 4, AntennaType: 1, ReaderID: "READER-02", X: 15.5, Y: 20.0},
			},
		},
		{
			ReaderID:   "READER-03",
			ReaderName: "POS Reader",
			ReaderType: "RFID",
			IPAddress:  "192.168.1.103",
			FloorID:    floorID,
			X:          25.0,
			Y:          5.0,
			Antennas: []site.Antenna{
				{AntennaID: 5, AntennaType: 2, ReaderID: "READER-03", X: 25.0, Y: 5.0},
			},
		},
	}

	antennas := []site.Antenna{
		{AntennaID: 1, AntennaType: 2, ReaderID: "READER-01", X: 5.0, Y: 10.0},
		{AntennaID: 2, AntennaType: 2, ReaderID: "READER-01", X: 5.5, Y: 10.0},
		{AntennaID: 3, AntennaType: 2, ReaderID: "READER-02", X: 15.0, Y: 20.0},
		{AntennaID: 4, AntennaType: 1, ReaderID: "READER-02", X: 15.5, Y: 20.0},
		{AntennaID: 5, AntennaType: 2, ReaderID: "READER-03", X: 25.0, Y: 5.0},
	}

	return &site.SiteConfig{
		SiteID:   realSiteID,
		SiteName: realSiteName,
		Readers:  readers,
		Floors: []site.Floor{
			{
				FloorID:   floorID,
				FloorName: "Floor001",
				Number:    0,
				Height:    420,
				Regions:   regions,
			},
		},
		Regions:  regions,
		Antennas: antennas,
	}
}

func loadSiteConfigFromFile(siteID string) *site.SiteConfig {
	path := strings.TrimSpace(os.Getenv(mockSiteJSONEnv))
	if path == "" {
		path = mockSiteJSONFile
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return nil
	}

	var config site.SiteConfig
	if err := json.Unmarshal(content, &config); err != nil {
		return nil
	}

	if config.SiteID != siteID {
		return nil
	}

	if len(config.Regions) == 0 {
		config.Regions = collectRegions(config.Floors)
	}

	if len(config.Antennas) == 0 {
		config.Antennas = collectAntennas(config.Readers)
	}

	return &config
}

func collectRegions(floors []site.Floor) []site.Region {
	var regions []site.Region

	for _, floor := range floors {
		regions = append(regions, floor.Regions...)
	}

	return regions
}

func collectAntennas(readers []site.Reader) []site.Antenna {
	var antennas []site.Antenna

	for _, reader := range readers {
		antennas = append(antennas, reader.Antennas...)
	}

	return antennas
}