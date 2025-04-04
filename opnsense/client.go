package opnsense

import (
	"compress/gzip"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
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
	httpClient            *http.Client
	longRunningHttpClient *http.Client // A separate client is used for long running, expensive requests
	gatewayLossRegex      *regexp.Regexp
	gatewayRTTRegex       *regexp.Regexp
	log                   *slog.Logger
	headers               map[string]string
	endpoints             map[EndpointName]EndpointPath
	baseURL               string
	key                   string
	secret                string
	sslInsecure           bool
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
			"protocolStatistics":      "api/diagnostics/interface/getProtocolStatistics",
			"pfStatisticsByInterface": "api/diagnostics/firewall/pf_statistics/interfaces",
			"arp":                     "api/diagnostics/interface/search_arp",
			"dhcpv4":                  "api/dhcpv4/leases/searchLease",
			"openVPNInstances":        "api/openvpn/instances/search",
			"gatewaysStatus":          "api/routing/settings/searchGateway",
			"unboundDNSStatus":        "api/unbound/diagnostics/stats",
			"cronJobs":                "api/cron/settings/searchJobs",
			"wireguardClients":        "api/wireguard/service/show",
			"healthCheck":             "api/core/system/status",
			"firmware":                "api/core/firmware/status",
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
				IdleConnTimeout:       90 * time.Second,
				TLSHandshakeTimeout:   3 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
				ForceAttemptHTTP2:     true,
				MaxIdleConnsPerHost:   runtime.GOMAXPROCS(0) + 1,
			},
		},
		longRunningHttpClient: &http.Client{
			Timeout: 2 * time.Minute,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: cfg.Insecure,
					RootCAs:            sslPool,
				},
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
	return c.doWithClient(method, path, body, responseStruct, MaxRetries, c.httpClient)
}

func (c *Client) doLongRunning(method string, path EndpointPath, body io.Reader, responseStruct any) *APICallError {
	return c.doWithClient(method, path, body, responseStruct, 1, c.longRunningHttpClient)
}

func (c *Client) doWithClient(method string, path EndpointPath, body io.Reader, responseStruct any, maxRetries int, client *http.Client) *APICallError {
	url := fmt.Sprintf("%s/%s", c.baseURL, string(path))

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return &APICallError{
			Endpoint:   string(path),
			Message:    err.Error(),
			StatusCode: 0,
		}
	}

	req.SetBasicAuth(c.key, c.secret)

	for k, v := range c.headers {
		req.Header.Add(k, v)
	}

	if method == "POST" {
		req.Header.Add("Content-Type", "application/json;charset=utf-8")
	}

	c.log.Debug("fetching data", "component", "opnsense-client", "url", url, "method", method)

	// Retry the request up to MaxRetries times
	for i := 0; i < maxRetries; i++ {
		resp, err := client.Do(req)
		if err != nil {
			c.log.Error("failed to send request; retrying",
				"component", "opnsense-client",
				"err", err.Error())
			time.Sleep(25 * time.Millisecond)
			continue
		}

		var reader io.ReadCloser
		switch resp.Header.Get("Content-Encoding") {
		case "gzip":
			reader, err = gzip.NewReader(resp.Body)
			if err != nil {
				return &APICallError{
					Endpoint:   string(path),
					Message:    fmt.Sprintf("failed to decompress gzip response body: %s", err.Error()),
					StatusCode: resp.StatusCode,
				}
			}
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

		reader.Close()

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			err := json.Unmarshal(body, &responseStruct)
			if err != nil {
				return &APICallError{
					Endpoint:   string(path),
					Message:    fmt.Sprintf("failed to unmarshal response body: %s", err.Error()),
					StatusCode: resp.StatusCode,
				}
			}
			return nil
		} else {
			return &APICallError{
				Endpoint:   string(path),
				Message:    string(body),
				StatusCode: resp.StatusCode,
			}
		}

	}
	return &APICallError{
		Endpoint:   string(path),
		Message:    fmt.Sprintf("max retries of %d times reached", MaxRetries),
		StatusCode: 0,
	}
}
