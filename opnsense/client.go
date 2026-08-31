package opnsense

import (
	"bytes"
	"compress/gzip"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"regexp"
	"runtime"
	"time"

	"github.com/AthennaMind/opnsense-exporter/internal/options"
)

// MaxRetries is the maximum number of retries
// when a request to the OPNsense API fails
const MaxRetries = 3

// EndpointName is the custom type for name of an endpoint definition
type EndpointName string

// EndpointPath is the custom type for url path of an endpoint definition
type EndpointPath string

// Client is an OPNsense API client
type Client struct {
	httpClient       *http.Client
	gatewayLossRegex *regexp.Regexp
	gatewayRTTRegex  *regexp.Regexp
	log              *slog.Logger
	headers          map[string]string
	endpoints        map[EndpointName]EndpointPath
	baseURL          string
	key              string
	secret           string
	sslInsecure      bool
}

// NewClient creates a new OPNsense API Client
func NewClient(cfg options.OPNSenseConfig, userAgentVersion string, log *slog.Logger) (Client, error) {
	sslPool, err := x509.SystemCertPool()
	if err != nil {
		return Client{}, errors.Join(fmt.Errorf("failed to load system cert pool"), err)
	}

	gatewayLossRegex, err := regexp.Compile(`\d\.\d %`)
	if err != nil {
		return Client{}, errors.Join(fmt.Errorf("failed to build regex for gatewayLoss calculation"), err)
	}

	gatewayRTTRegex, err := regexp.Compile(`\d+\.\d+ ms`)
	if err != nil {
		return Client{}, errors.Join(fmt.Errorf("failed to build regex for gatewayRTT calculation"), err)
	}
	client := Client{
		log:              log,
		baseURL:          fmt.Sprintf("%s://%s", cfg.Protocol, cfg.Host),
		key:              cfg.APIKey,
		secret:           cfg.APISecret,
		gatewayLossRegex: gatewayLossRegex,
		gatewayRTTRegex:  gatewayRTTRegex,
		endpoints: map[EndpointName]EndpointPath{
			"services":                "api/core/service/search",
			"interfaces":              "api/diagnostics/traffic/interface",
			"protocolStatistics":      "api/diagnostics/interface/get_protocol_statistics",
			"pfStatisticsByInterface": "api/diagnostics/firewall/pf_statistics/interfaces",
			"arp":                     "api/diagnostics/interface/search_arp",
			"dhcpv4":                  "api/dhcpv4/leases/searchLease",
			"openVPNInstances":        "api/openvpn/instances/search",
			"openVPNSessions":         "api/openvpn/service/search_sessions",
			"gatewaysStatus":          "api/routing/settings/searchGateway",
			"unboundDNSStatus":        "api/unbound/diagnostics/stats",
			"cronJobs":                "api/cron/settings/searchJobs",
			"wireguardClients":        "api/wireguard/service/show",
			"ipsecPhase1":             "api/ipsec/sessions/search_phase1",
			"ipsecPhase2":             "api/ipsec/sessions/search_phase2",
			"healthCheck":             "api/core/system/status",
			"systemTemperature":       "api/diagnostics/system/systemTemperature",
			"firmware":                "api/core/firmware/status",
			"keaDhcpv4":               "api/kea/leases4/search",
			"keaDhcpv6":               "api/kea/leases6/search",
		},
		headers: map[string]string{
			"Accept":          "application/json",
			"User-Agent":      fmt.Sprintf("prometheus-opnsense-exporter/%s", userAgentVersion),
			"Accept-Encoding": "gzip, deflate, br",
		},
		sslInsecure: cfg.Insecure,
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: cfg.Insecure,
					RootCAs:            sslPool,
				},
				// Bound the TCP dial. Without this a dial to an unreachable
				// host is canceled only by the client timeout and the
				// half-open socket can linger in SYN-SENT far longer.
				DialContext: (&net.Dialer{
					Timeout:   5 * time.Second,
					KeepAlive: 30 * time.Second,
				}).DialContext,
				IdleConnTimeout:       90 * time.Second,
				TLSHandshakeTimeout:   3 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
				ForceAttemptHTTP2:     true,
				MaxIdleConnsPerHost:   runtime.GOMAXPROCS(0) + 1,
			},
		},
	}

	return client, nil
}

// Endpoints returns a map of all the endpoints
// that are called by the client.
func (c *Client) Endpoints() map[EndpointName]EndpointPath {
	return c.endpoints
}

// do sends a request to the OPNsense API.
// The response is unmarshalled
// into the responseStruc
func (c *Client) do(method string, path EndpointPath, body io.Reader, responseStruct any) *APICallError {
	url := fmt.Sprintf("%s/%s", c.baseURL, string(path))

	// Buffer the payload so every retry sends a fresh body.
	// A plain io.Reader is consumed by the first attempt.
	var payload []byte
	if body != nil {
		var err error
		payload, err = io.ReadAll(body)
		if err != nil {
			return &APICallError{
				Endpoint:   string(path),
				Message:    err.Error(),
				StatusCode: 0,
			}
		}
	}

	c.log.Debug("fetching data", "component", "opnsense-client", "url", url, "method", method)

	// Retry the request up to MaxRetries times
	for i := 0; i < MaxRetries; i++ {
		req, err := c.newRequest(method, url, payload)
		if err != nil {
			return &APICallError{
				Endpoint:   string(path),
				Message:    err.Error(),
				StatusCode: 0,
			}
		}

		resp, err := c.httpClient.Do(req)
		if err != nil {
			c.log.Error("failed to send request; retrying",
				"component", "opnsense-client",
				"err", err.Error())
			time.Sleep(25 * time.Millisecond)
			continue
		}

		return c.handleResponse(resp, path, url, responseStruct)
	}
	return &APICallError{
		Endpoint:   string(path),
		Message:    fmt.Sprintf("max retries of %d times reached", MaxRetries),
		StatusCode: 0,
	}
}

// newRequest builds a request with auth and headers set.
func (c *Client) newRequest(method, url string, payload []byte) (*http.Request, error) {
	var body io.Reader
	if payload != nil {
		body = bytes.NewReader(payload)
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(c.key, c.secret)

	for k, v := range c.headers {
		req.Header.Add(k, v)
	}

	if method == "POST" {
		req.Header.Add("Content-Type", "application/json;charset=utf-8")
	}

	return req, nil
}

// handleResponse reads the response and unmarshals it into responseStruct.
// It always drains and closes the response body so the
// connection is returned to the pool instead of leaked.
func (c *Client) handleResponse(resp *http.Response, path EndpointPath, url string, responseStruct any) *APICallError {
	defer func() {
		_, _ = io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}()

	var reader io.Reader
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		gzReader, err := gzip.NewReader(resp.Body)
		if err != nil {
			return &APICallError{
				Endpoint:   string(path),
				Message:    fmt.Sprintf("failed to decompress gzip response body: %s", err.Error()),
				StatusCode: resp.StatusCode,
			}
		}
		defer gzReader.Close()
		reader = gzReader
	default:
		reader = resp.Body
	}

	body, err := io.ReadAll(reader)
	if err != nil {
		return &APICallError{
			Endpoint:   string(path),
			Message:    fmt.Sprintf("failed to read response body: %s", err.Error()),
			StatusCode: resp.StatusCode,
		}
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return &APICallError{
			Endpoint:   string(path),
			Message:    string(body),
			StatusCode: resp.StatusCode,
		}
	}

	if err := json.Unmarshal(body, &responseStruct); err != nil {
		return &APICallError{
			Endpoint:   string(path),
			Message:    fmt.Sprintf("failed to unmarshal response body: %s", err.Error()),
			StatusCode: resp.StatusCode,
		}
	}

	c.log.Debug("returned data", "component", "opnsense-client", "url", url, "data", string(body))

	return nil
}
