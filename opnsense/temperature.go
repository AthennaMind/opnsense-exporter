package opnsense

import (
	"strconv"
)

// temperatureResponse is the bare JSON array returned by
// api/diagnostics/system/systemTemperature. The temperature value
// arrives as a quoted string and may be empty when the sensor
// has no reading.
type temperatureResponse []struct {
	Device      string  `json:"device"`
	DeviceSeq   flexInt `json:"device_seq"`
	Temperature string  `json:"temperature"`
	Type        string  `json:"type"`
}

type TemperatureSensor struct {
	Device    string
	DeviceSeq string
	Type      string
	Celsius   float64
}

type Temperatures struct {
	Sensors []TemperatureSensor
}

func parseTemperatures(resp temperatureResponse) Temperatures {
	var data Temperatures

	for _, row := range resp {
		celsius, err := strconv.ParseFloat(row.Temperature, 64)
		if err != nil {
			// A sensor without a reading must not become a 0 degree
			// sample, so the series is skipped instead.
			continue
		}

		data.Sensors = append(data.Sensors, TemperatureSensor{
			Device:    row.Device,
			DeviceSeq: strconv.Itoa(int(row.DeviceSeq)),
			Type:      row.Type,
			Celsius:   celsius,
		})
	}

	return data
}

func (c *Client) FetchTemperatures() (Temperatures, *APICallError) {
	var resp temperatureResponse
	var data Temperatures

	path, ok := c.endpoints["systemTemperature"]
	if !ok {
		return data, &APICallError{
			Endpoint:   "systemTemperature",
			Message:    "endpoint not found",
			StatusCode: 0,
		}
	}

	if err := c.do("GET", path, nil, &resp); err != nil {
		// Releases without this endpoint answer 404.
		// Treat that as a box with no sensors instead of an error.
		if err.StatusCode == 404 {
			return data, nil
		}
		return data, err
	}

	data = parseTemperatures(resp)

	if skipped := len(resp) - len(data.Sensors); skipped > 0 {
		c.log.Debug("temperature sensors without a reading skipped",
			"component", "opnsense-client",
			"skipped", skipped,
		)
	}

	return data, nil
}
