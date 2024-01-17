package opnsense

import (
	"compress/gzip"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"runtime"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
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
	log              log.Logger
	baseURL          string
	key              string
	secret           string
	sslInsecure      bool
	endpoints        map[EndpointName]EndpointPath
	httpClient       *http.Client
	headers          map[string]string
	gatewayLossRegex *regexp.Regexp
	gatewayRTTRegex  *regexp.Regexp
}

// NewClient creates a new OPNsense API Client
func NewClient(protocol, address, key, secret, userAgentVersion string, sslInsecure bool, log log.Logger) (Client, error) {

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
		baseURL:          fmt.Sprintf("%s://%s", protocol, address),
		key:              key,
		secret:           secret,
		gatewayLossRegex: gatewayLossRegex,
		gatewayRTTRegex:  gatewayRTTRegex,
		endpoints: map[EndpointName]EndpointPath{
			"services":           "api/core/service/search",
			"protocolStatistics": "api/diagnostics/interface/getProtocolStatistics",
			"arp":                "api/diagnostics/interface/search_arp",
			"dhcpv4":             "api/dhcpv4/leases/searchLease",
			"openVPNInstances":   "api/openvpn/instances/search",
			"interfaces":         "api/diagnostics/traffic/interface",
			"systemInfo":         "widgets/api/get.php?load=system%2Ctemperature",
			"gatewaysStatus":     "api/routes/gateway/status",
			"unboundDNSStatus":   "api/unbound/diagnostics/stats",
			"cronJobs":           "api/cron/settings/searchJobs",
		},
		headers: map[string]string{
			"Accept":          "application/json",
			"User-Agent":      fmt.Sprintf("prometheus-opnsense-exporter/%s", userAgentVersion),
			"Accept-Encoding": "gzip, deflate, br",
		},
		sslInsecure: sslInsecure,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: sslInsecure,
					RootCAs:            sslPool,
				},
				IdleConnTimeout:       90 * time.Second,
				TLSHandshakeTimeout:   1 * time.Second,
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

	level.Debug(c.log).
		Log("msg", "fetching data", "component", "opnsense-client", "url", url, "method", method)

	// Retry the request up to MaxRetries times
	for i := 0; i < MaxRetries; i++ {
		resp, err := c.httpClient.Do(req)
		if err != nil {
			level.Error(c.log).
				Log("msg", "failed to send request; retrying",
					"component", "opnsense-client",
					"err", err.Error())
			time.Sleep(25 * time.Millisecond)
			continue
		}

		var reader io.ReadCloser
		switch resp.Header.Get("Content-Encoding") {
		case "gzip":
			reader, err = gzip.NewReader(resp.Body)
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
				fmt.Println(url)
				fmt.Println(string(body))
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
