// RawRead represents a single recorded RFID raw read.
type RawRead struct {
	ReadID             string
	RecordingSessionID string
	TagID              string
	ReaderID           string
	AntennaID          int
	AntennaTypeID      int
	SourceTimestampUtc string
	InjectionTimeUtc   string
	Confidence         int
	RSSI               float64
	TagX               float64
	TagY               float64
	FloorID            int
	RawPayload         []byte

	Timestamp     time.Time
	InjectionTime time.Time
}