package site

// SiteConfig represents a Resonate site configuration.
// This structure matches the real Recorder site_json format.
type SiteConfig struct {
	SiteID      string       `json:"id"`
	SiteName    string       `json:"name"`
	Expirations Expirations  `json:"expirations,omitempty"`
	Location    *Location    `json:"location,omitempty"`
	Networking  *Networking  `json:"networking,omitempty"`
	Readers     []Reader     `json:"readers,omitempty"`
	Floors      []Floor      `json:"floors,omitempty"`
	Regions     []Region     `json:"regions,omitempty"`
	Antennas    []Antenna    `json:"antennas,omitempty"`
}

// Expirations represents site-level timeout settings.
type Expirations struct {
	Departed                       string `json:"departed,omitempty"`
	Ghost                          string `json:"ghost,omitempty"`
	Missing                        string `json:"missing,omitempty"`
	ItemDepartedDebounceWindow     string `json:"itemDepartedDebounceWindow,omitempty"`
	AssetDepartedDebounceWindow    string `json:"assetDepartedDebounceWindow,omitempty"`
	DeviceDepartedDebounceWindow   string `json:"deviceDepartedDebounceWindow,omitempty"`
}

// Location represents site location details.
type Location struct {
	Address     string      `json:"address,omitempty"`
	Coordinates Coordinates `json:"coordinates,omitempty"`
	Origin      *Origin     `json:"origin,omitempty"`
	GeoPosition *GeoPosition `json:"geoPosition,omitempty"`
	Other       any         `json:"other,omitempty"`
}

// Coordinates represents location coordinate metadata.
type Coordinates struct {
	X float64 `json:"x,omitempty"`
	Y float64 `json:"y,omitempty"`
}

// Origin represents origin position data.
type Origin struct {
	Position Position `json:"position,omitempty"`
}

// Position represents x/y position.
type Position struct {
	X float64 `json:"x,omitempty"`
	Y float64 `json:"y,omitempty"`
}

// GeoPosition represents geo position.
type GeoPosition struct {
	Latitude  float64 `json:"latitude,omitempty"`
	Longitude float64 `json:"longitude,omitempty"`
}

// Networking represents site networking settings.
type Networking struct {
	NTPAddress                 string `json:"ntpAddress,omitempty"`
	SyslogAddress              string `json:"syslogAddress,omitempty"`
	TFTPAddress                string `json:"tftpAddress,omitempty"`
	DNSAddress                 string `json:"dnsAddress,omitempty"`
	FullyQualifiedDomainSuffix string `json:"fullyQualifiedDomainSuffix,omitempty"`
	CMM                        *CMM   `json:"cmm,omitempty"`
}

// CMM represents CMM networking configuration.
type CMM struct {
	Address  string `json:"address,omitempty"`
	Protocol string `json:"protocol,omitempty"`
	Port     int    `json:"port,omitempty"`
}

// Reader represents a physical RFID reader.
type Reader struct {
	ReaderID   string    `json:"id"`
	ReaderName string    `json:"name,omitempty"`
	ReaderType string    `json:"type,omitempty"`
	IPAddress  string    `json:"ipAddress,omitempty"`
	FloorID    string    `json:"floorId,omitempty"`
	X          float64   `json:"x,omitempty"`
	Y          float64   `json:"y,omitempty"`
	Antennas   []Antenna `json:"antennas,omitempty"`
}

// Antenna represents a reader antenna.
type Antenna struct {
	AntennaID   int     `json:"antenna_id"`
	AntennaType int     `json:"antenna_type,omitempty"`
	ReaderID    string  `json:"reader_id,omitempty"`
	X           float64 `json:"x,omitempty"`
	Y           float64 `json:"y,omitempty"`
}

// Floor represents a site floor.
// Real Recorder site_json uses UUID string IDs.
type Floor struct {
	FloorID   string   `json:"id"`
	FloorName string   `json:"name,omitempty"`
	Number    int      `json:"number,omitempty"`
	Width     float64  `json:"width,omitempty"`
	Height    float64  `json:"height,omitempty"`
	Regions   []Region `json:"regions,omitempty"`
	Bounds    []Point  `json:"bounds,omitempty"`
}

// Region represents a site region/zone.
type Region struct {
	RegionID      string   `json:"id"`
	RegionName    string   `json:"name,omitempty"`
	RegionType    string   `json:"type,omitempty"`
	Physicality   string   `json:"physicality,omitempty"`
	Behaviors     []string `json:"behaviors,omitempty"`
	InventoryType string   `json:"inventoryType,omitempty"`
	Timeouts      Timeouts `json:"timeouts,omitempty"`
	Regions       []Region `json:"regions,omitempty"`
	Bounds        []Point  `json:"bounds,omitempty"`
	FloorID        string   `json:"floor_id,omitempty"`
}

// Timeouts represents region timeout settings.
type Timeouts struct {
	AssetToGhost                    string `json:"assetToGhost,omitempty"`
	AssetToMissing                  string `json:"assetToMissing,omitempty"`
	AssetToDeparted                 string `json:"assetToDeparted,omitempty"`
	ItemToMissing                   string `json:"itemToMissing,omitempty"`
	ItemToGhost                     string `json:"itemToGhost,omitempty"`
	ItemToDeparted                  string `json:"itemToDeparted,omitempty"`
	DeviceToDeparted                string `json:"deviceToDeparted,omitempty"`
	PersonToDeparted                string `json:"personToDeparted,omitempty"`
	PointOfSaleToMissing            string `json:"pointOfSaleToMissing,omitempty"`
	PointOfSaleWindow               string `json:"pointOfSaleWindow,omitempty"`
	PointOfSaleZoneTimeTillMissing  string `json:"pointOfSaleZoneTimeTillMissing,omitempty"`
	PointOfSaleZoneWindow           string `json:"pointOfSaleZoneWindow,omitempty"`
	RFIDExitWindow                  string `json:"rfidExitWindow,omitempty"`
	RFIDEnterWindow                 string `json:"rfidEnterWindow,omitempty"`
}

// Point represents x/y boundary points.
type Point struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}