package replay

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Injector sends ProtoReaderBundle payloads to a target Resonate instance.
type Injector struct {
	TargetURL  string
	HTTPClient *http.Client
}

// NewInjector creates a new HTTP injector for the given target URL.
func NewInjector(targetURL string) *Injector {
	return &Injector{
		TargetURL: targetURL,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// buildReaderBundlesURL validates the target URL and returns a safe,
// allowlisted /reader-bundles endpoint.
//
// This is intentionally strict for the MVP/demo.
// It only allows the local mock target server.
func buildReaderBundlesURL(targetURL string) (string, error) {
	if strings.TrimSpace(targetURL) == "" {
		return "", fmt.Errorf("target URL must not be empty")
	}

	parsed, err := url.Parse(targetURL)
	if err != nil {
		return "", fmt.Errorf("invalid target URL: %w", err)
	}

	if parsed.User != nil {
		return "", fmt.Errorf("target URL must not contain credentials")
	}

	if parsed.Path != "" && parsed.Path != "/" {
		return "", fmt.Errorf("target URL path is not allowed")
	}

	host := strings.ToLower(parsed.Hostname())
	port := parsed.Port()

	switch {
	case parsed.Scheme == "http" &&
		host == "localhost" &&
		port == "8080":
		return "http://localhost:8080/reader-bundles", nil

	case parsed.Scheme == "http" &&
		host == "127.0.0.1" &&
		port == "8080":
		return "http://127.0.0.1:8080/reader-bundles", nil

	default:
		return "", fmt.Errorf("target URL is not allowed: %s", targetURL)
	}
}

// Send posts a ProtoReaderBundle payload to the target endpoint.
// Endpoint: POST /reader-bundles
func (inj *Injector) Send(payload *ProtoReaderBundleWrapper) error {
	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	endpoint, err := buildReaderBundlesURL(inj.TargetURL)
	if err != nil {
		return fmt.Errorf("invalid target URL: %w", err)
	}

	req, err := http.NewRequest(
		http.MethodPost,
		endpoint,
		bytes.NewReader(jsonBytes),
	)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := inj.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send payload to %s: %w", endpoint, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf(
			"target returned status %d: %s",
			resp.StatusCode,
			string(body),
		)
	}

	return nil
}