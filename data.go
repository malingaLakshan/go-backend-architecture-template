package replay

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type Injector struct {
	TargetURL  string
	HTTPClient *http.Client
}

func NewInjector(targetURL string) *Injector {
	return &Injector{
		TargetURL: targetURL,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func buildReaderBundlesURL(targetURL string) (string, error) {
	if strings.TrimSpace(targetURL) == "" {
		return "", fmt.Errorf("target URL must not be empty")
	}

	parsed, err := url.Parse(targetURL)
	if err != nil {
		return "", fmt.Errorf("invalid target URL: %w", err)
	}

	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return "", fmt.Errorf("unsupported URL scheme: %s", parsed.Scheme)
	}

	if parsed.Host == "" {
		return "", fmt.Errorf("target URL must include a host")
	}

	if parsed.User != nil {
		return "", fmt.Errorf("target URL must not contain credentials")
	}

	host := parsed.Hostname()
	if !isAllowedTargetHost(host) {
		return "", fmt.Errorf("target host is not allowed: %s", host)
	}

	endpoint := url.URL{
		Scheme: parsed.Scheme,
		Host:   parsed.Host,
		Path:   "/reader-bundles",
	}

	return endpoint.String(), nil
}

func isAllowedTargetHost(host string) bool {
	host = strings.ToLower(strings.TrimSpace(host))

	if host == "localhost" || host == "127.0.0.1" || host == "::1" {
		return true
	}

	allowedHosts := os.Getenv("RRE_ALLOWED_TARGET_HOSTS")
	if allowedHosts == "" {
		return false
	}

	for _, allowed := range strings.Split(allowedHosts, ",") {
		allowed = strings.ToLower(strings.TrimSpace(allowed))
		if allowed == "" {
			continue
		}

		if host == allowed {
			return true
		}

		if strings.HasPrefix(allowed, ".") &&
			strings.HasSuffix(host, allowed) {
			return true
		}
	}

	ip := net.ParseIP(host)
	if ip != nil {
		for _, allowed := range strings.Split(allowedHosts, ",") {
			allowed = strings.TrimSpace(allowed)

			_, cidr, err := net.ParseCIDR(allowed)
			if err == nil && cidr.Contains(ip) {
				return true
			}
		}
	}

	return false
}

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