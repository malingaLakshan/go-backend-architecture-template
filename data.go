package replay

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// Injector sends ProtoReaderBundle payloads to a target Resonate instance.
type Injector struct {
	TargetURL  string
	HTTPClient *http.Client
}

// buildReaderBundlesURL validates the targetURL and returns the safe
// /reader-bundles endpoint.
// Only http and https are allowed. The final path is hardcoded.
func buildReaderBundlesURL(targetURL string) (string, error) {
	if targetURL == "" {
		return "", fmt.Errorf("target URL must not be empty")
	}

	parsed, err := url.Parse(targetURL)
	if err != nil {
		return "", fmt.Errorf("invalid target URL: %w", err)
	}

	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return "", fmt.Errorf(
			"unsupported URL scheme %q: only http and https are allowed",
			parsed.Scheme,
		)
	}

	if parsed.Host == "" {
		return "", fmt.Errorf("target URL must include a host")
	}

	if parsed.User != nil {
		return "", fmt.Errorf("target URL must not contain credentials")
	}

	endpoint := &url.URL{
		Scheme: parsed.Scheme,
		Host:   parsed.Host,
		Path:   "/reader-bundles",
	}

	return endpoint.String(), nil
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