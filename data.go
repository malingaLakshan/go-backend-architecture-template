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

const (
	targetLocalhost = "localhost"
	targetLoopback  = "127.0.0.1"

	readerBundlesLocalhost = "http://localhost:8080/reader-bundles"
	readerBundlesLoopback  = "http://127.0.0.1:8080/reader-bundles"
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

// validateAllowedTarget validates the input target URL and returns
// an internal allowlist key. For this MVP, only the local mock server is allowed.
func validateAllowedTarget(targetURL string) (string, error) {
	if strings.TrimSpace(targetURL) == "" {
		return "", fmt.Errorf("target URL must not be empty")
	}

	parsed, err := url.Parse(targetURL)
	if err != nil {
		return "", fmt.Errorf("invalid target URL: %w", err)
	}

	if parsed.Scheme != "http" {
		return "", fmt.Errorf("unsupported target URL scheme: %s", parsed.Scheme)
	}

	if parsed.User != nil {
		return "", fmt.Errorf("target URL must not contain credentials")
	}

	if parsed.Host == "" {
		return "", fmt.Errorf("target URL must include a host")
	}

	if parsed.Path != "" && parsed.Path != "/" {
		return "", fmt.Errorf("target URL path is not allowed")
	}

	if parsed.RawQuery != "" || parsed.Fragment != "" {
		return "", fmt.Errorf("target URL query or fragment is not allowed")
	}

	host := strings.ToLower(parsed.Hostname())
	port := parsed.Port()

	if host == "localhost" && port == "8080" {
		return targetLocalhost, nil
	}

	if host == "127.0.0.1" && port == "8080" {
		return targetLoopback, nil
	}

	return "", fmt.Errorf("target URL is not allowlisted")
}

// buildReaderBundlesURL is kept for tests and validation.
// It returns only hardcoded allowlisted endpoints.
func buildReaderBundlesURL(targetURL string) (string, error) {
	target, err := validateAllowedTarget(targetURL)
	if err != nil {
		return "", err
	}

	if target == targetLocalhost {
		return readerBundlesLocalhost, nil
	}

	if target == targetLoopback {
		return readerBundlesLoopback, nil
	}

	return "", fmt.Errorf("target URL is not allowlisted")
}

// Send posts a ProtoReaderBundle payload to the target endpoint.
// Endpoint: POST /reader-bundles
func (inj *Injector) Send(payload *ProtoReaderBundleWrapper) error {
	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	target, err := validateAllowedTarget(inj.TargetURL)
	if err != nil {
		return fmt.Errorf("invalid target URL: %w", err)
	}

	if target == targetLocalhost {
		return inj.sendToLocalhost(jsonBytes)
	}

	if target == targetLoopback {
		return inj.sendToLoopback(jsonBytes)
	}

	return fmt.Errorf("target URL is not allowlisted")
}

func (inj *Injector) sendToLocalhost(jsonBytes []byte) error {
	req, err := http.NewRequest(
		http.MethodPost,
		readerBundlesLocalhost,
		bytes.NewReader(jsonBytes),
	)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := inj.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send payload: %w", err)
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

func (inj *Injector) sendToLoopback(jsonBytes []byte) error {
	req, err := http.NewRequest(
		http.MethodPost,
		readerBundlesLoopback,
		bytes.NewReader(jsonBytes),
	)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := inj.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send payload: %w", err)
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