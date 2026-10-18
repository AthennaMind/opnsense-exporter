package opnsense

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/AthennaMind/opnsense-exporter/internal/options"
)

func TestFetchResourceUsage(t *testing.T) {
	endpoints := []struct {
		name    string
		path    string
		payload func(string, string) string
		fetch   func(*Client) (ResourceUsage, *APICallError)
	}{
		{
			name: "pf states",
			path: "/api/diagnostics/firewall/pf_states",
			payload: func(current, limit string) string {
				return fmt.Sprintf(`{"current":%s,"limit":%s}`, current, limit)
			},
			fetch: (*Client).FetchPFStates,
		},
		{
			name: "mbuf clusters",
			path: "/api/diagnostics/system/system_mbuf",
			payload: func(current, limit string) string {
				return fmt.Sprintf(`{"mbuf-statistics":{"cluster-total":%s,"cluster-max":%s,"mbuf-current":999}}`, current, limit)
			},
			fetch: (*Client).FetchMbufUsage,
		},
	}

	for _, endpoint := range endpoints {
		t.Run(endpoint.name, func(t *testing.T) {
			cases := []struct {
				name       string
				payload    string
				statusCode int
				want       ResourceUsage
				wantError  bool
			}{
				{name: "strings", payload: endpoint.payload(`"123"`, `"1000000"`), want: ResourceUsage{123, 1000000}},
				{name: "numbers", payload: endpoint.payload(`123`, `1000000`), want: ResourceUsage{123, 1000000}},
				{name: "mixed types", payload: endpoint.payload(`"123"`, `1000000`), want: ResourceUsage{123, 1000000}},
				{name: "zero usage", payload: endpoint.payload(`"0"`, `1000`), want: ResourceUsage{0, 1000}},
				{name: "large counts", payload: endpoint.payload(`4294967296`, `"8589934592"`), want: ResourceUsage{4294967296, 8589934592}},
				{name: "backend failure", payload: `{"result":"failed"}`, wantError: true},
				{name: "empty response", payload: `{}`, wantError: true},
				{name: "null response", payload: `null`, wantError: true},
				{name: "null current", payload: endpoint.payload(`null`, `1000`), wantError: true},
				{name: "null limit", payload: endpoint.payload(`123`, `null`), wantError: true},
				{name: "empty count", payload: endpoint.payload(`""`, `1000`), wantError: true},
				{name: "invalid count", payload: endpoint.payload(`"unknown"`, `1000`), wantError: true},
				{name: "negative count", payload: endpoint.payload(`-1`, `1000`), wantError: true},
				{name: "negative limit", payload: endpoint.payload(`123`, `-1`), wantError: true},
				{name: "fractional count", payload: endpoint.payload(`1.5`, `1000`), wantError: true},
				{name: "overflow", payload: endpoint.payload(`18446744073709551616`, `1000`), wantError: true},
				{name: "invalid JSON", payload: `{`, wantError: true},
				{name: "forbidden", payload: `{"message":"Forbidden"}`, statusCode: http.StatusForbidden, wantError: true},
				{name: "not found", payload: `{}`, statusCode: http.StatusNotFound, wantError: true},
				{name: "server error", payload: `{}`, statusCode: http.StatusInternalServerError, wantError: true},
			}
			for _, tc := range cases {
				t.Run(tc.name, func(t *testing.T) {
					server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						if r.Method != http.MethodGet || r.URL.Path != endpoint.path {
							t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
							http.NotFound(w, r)
							return
						}
						w.Header().Set("Content-Type", "application/json")
						if tc.statusCode != 0 {
							w.WriteHeader(tc.statusCode)
						}
						_, _ = io.WriteString(w, tc.payload)
					}))
					defer server.Close()
					client, err := NewClient(options.OPNSenseConfig{
						Protocol: "http",
						Host:     strings.TrimPrefix(server.URL, "http://"),
					}, "test", slog.New(slog.NewTextHandler(io.Discard, nil)))
					if err != nil {
						t.Fatal(err)
					}
					got, apiErr := endpoint.fetch(&client)
					if tc.wantError {
						if apiErr == nil {
							t.Fatalf("expected error, got %+v", got)
						}
						if apiErr.Endpoint != strings.TrimPrefix(endpoint.path, "/") {
							t.Errorf("unexpected error endpoint: %s", apiErr.Endpoint)
						}
						if tc.statusCode != 0 && apiErr.StatusCode != tc.statusCode {
							t.Errorf("expected status %d, got %d", tc.statusCode, apiErr.StatusCode)
						}
						if got != (ResourceUsage{}) {
							t.Errorf("unexpected partial data on failure: %+v", got)
						}
						return
					}
					if apiErr != nil {
						t.Fatal(apiErr)
					}
					if got != tc.want {
						t.Errorf("got %+v, want %+v", got, tc.want)
					}
				})
			}
		})
	}
}
