// Package logger configures structured logging using Go's standard log/slog.
// Logs are written to both stdout and a log file simultaneously.
package logger

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

// Setup initialises the global slog logger.
// It writes logs to both stdout (for development) and a file (for persistence).
func Setup(level string, filePath string) (*slog.Logger, error) {
	// Ensure the log directory exists.
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}

	// Open (or create) the log file in append mode.
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return nil, err
	}

	// Write to both stdout and the log file.
	multiWriter := io.MultiWriter(os.Stdout, file)

	// Parse the log level.
	logLevel := parseLevel(level)

	// Create a JSON handler for structured log output.
	handler := slog.NewJSONHandler(multiWriter, &slog.HandlerOptions{
		Level:     logLevel,
		AddSource: false,
	})

	logger := slog.New(handler)

	// Set as the default logger for the entire application.
	slog.SetDefault(logger)

	return logger, nil
}

// parseLevel converts a string level name to a slog.Level.
func parseLevel(level string) slog.Level {
	switch strings.ToUpper(level) {
	case "DEBUG":
		return slog.LevelDebug
	case "INFO":
		return slog.LevelInfo
	case "WARN":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}





package logger

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

func NewLogger() (*slog.Logger, error) {
	baseLogDir := "logs"
	logFileName := "app.log"

	cleanBaseDir := filepath.Clean(baseLogDir)
	cleanLogPath := filepath.Clean(filepath.Join(cleanBaseDir, logFileName))

	if !strings.HasPrefix(cleanLogPath, cleanBaseDir) {
		return nil, os.ErrPermission
	}

	if err := os.MkdirAll(cleanBaseDir, 0750); err != nil {
		return nil, err
	}

	file, err := os.OpenFile(cleanLogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		return nil, err
	}

	writer := io.MultiWriter(os.Stdout, file)

	handler := slog.NewJSONHandler(writer, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	return slog.New(handler), nil
}






