package cli

import (
	"io"
	"os"

	commonlogger "altrfidtools/common/logger"
	"altrfidtools/resonate-replay-engine/internal/config"
)

// initLogger creates the Replay Engine logger using the logging
// configuration from config.json.
//
// If the config or log file cannot be loaded, logging falls back
// to stderr so that command execution is not blocked.
func initLogger(configPath string) (*commonlogger.Logger, io.Closer) {
	cfg, err := config.Load(configPath)
	if err != nil {
		return commonlogger.New(os.Stderr), nopCloser{}
	}

	log, closer, err := commonlogger.NewFromConfig(cfg.Logging)
	if err != nil {
		return commonlogger.New(os.Stderr), nopCloser{}
	}

	return log, closer
}

// nopCloser is used when the fallback stderr logger is active.
type nopCloser struct{}

func (nopCloser) Close() error {
	return nil
}