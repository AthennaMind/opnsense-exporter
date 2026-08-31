package opnsense

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/AthennaMind/opnsense-exporter/internal/options"
)

func TestParseTemperatures(t *testing.T) {
	payload := `[
		{"device": "hw.acpi.thermal.tz0", "device_seq": 0, "temperature": "45.5", "type_translated": "Zone", "type": "acpi_tz"},
		{"device": "dev.cpu", "device_seq": "1", "temperature": "62.0", "type_translated": "CPU", "type": "coretemp"},
		{"device": "dev.cpu", "device_seq": 2, "temperature": "", "type_translated": "CPU", "type": "coretemp"}
	]`

	var resp temperatureResponse
	if err := json.Unmarshal([]byte(payload), &resp); err != nil {
		t.Fatalf("failed to unmarshal temperature response: %v", err)
	}

	data := parseTemperatures(resp)

	if len(data.Sensors) != 2 {
		t.Fatalf("expected 2 sensors, the one without a reading skipped, got %d", len(data.Sensors))
	}

	first := data.Sensors[0]
	if first.Device != "hw.acpi.thermal.tz0" || first.DeviceSeq != "0" || first.Type != "acpi_tz" || first.Celsius != 45.5 {
		t.Errorf("unexpected first sensor: %+v", first)
	}

	second := data.Sensors[1]
	if second.Device != "dev.cpu" || second.DeviceSeq != "1" || second.Celsius != 62.0 {
		t.Errorf("unexpected second sensor: %+v", second)
	}
}

func TestParseTemperaturesEmpty(t *testing.T) {
	var resp temperatureResponse
	if err := json.Unmarshal([]byte(`[]`), &resp); err != nil {
		t.Fatalf("failed to unmarshal empty response: %v", err)
	}
	if data := parseTemperatures(resp); len(data.Sensors) != 0 {
		t.Errorf("expected no sensors, got %d", len(data.Sensors))
	}
}

func TestFetchTemperaturesNotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	}))
	defer server.Close()

	client, err := NewClient(options.OPNSenseConfig{
		Protocol:  "http",
		Host:      strings.TrimPrefix(server.URL, "http://"),
		APIKey:    "test",
		APISecret: "test",
	}, "test", slog.New(slog.NewTextHandler(io.Discard, nil)))
	if err != nil {
		t.Fatalf("failed to build client: %v", err)
	}

	data, apiErr := client.FetchTemperatures()
	if apiErr != nil {
		t.Fatalf("expected 404 to be tolerated, got error: %s", apiErr.Error())
	}
	if len(data.Sensors) != 0 {
		t.Errorf("expected no sensors on 404, got %d", len(data.Sensors))
	}
}
