package opnsense

import (
	"fmt"
	"strings"
)

type firmwareStatusResponse struct {
	LastCheck      string `json:"last_check"`
	NeedsReboot    string `json:"needs_reboot"`
	OsVersion      string `json:"os_version"`
	ProductID      string `json:"product_id"`
	ProductVersion string `json:"product_version"`
	ProductAbi     string `json:"product_abi"`
	NewPackages    []struct {
		Name       string `json:"name"`
		Repository string `json:"repository"`
		Version    string `json:"version"`
	} `json:"new_packages"`
	UpgradePackages []struct {
		Name           string `json:"name"`
		Repository     string `json:"repository"`
		CurrentVersion string `json:"current_version"`
		NewVersion     string `json:"new_version"`
		Size           string `json:"size,omitempty"`
	} `json:"upgrade_packages"`
	Product struct {
		ProductCheck struct {
			UpgradeNeedsReboot string `json:"upgrade_needs_reboot"`
		} `json:"product_check"`
	} `json:"product"`
}

type FirmwareStatus struct {
	LastCheck          string
	NeedsReboot        int
	NewPackages        int
	OsVersion          string
	ProductABI         string
	ProductId          string
	ProductVersion     string
	UpgradePackages    int
	UpgradeNeedsReboot int
}

func (c *Client) FetchFirmwareStatus() (FirmwareStatus, *APICallError) {
	var resp firmwareStatusResponse
	var data FirmwareStatus

	url, ok := c.endpoints["firmware"]

	if !ok {
		return data, &APICallError{
			Endpoint:   "firmware",
			Message:    "Missing endpoint 'firmwareStatus'",
			StatusCode: 0,
		}
	}

	if err := c.do("GET", url, nil, &resp); err != nil {
		return data, err
	}

	data.LastCheck = resp.LastCheck
	data.OsVersion = resp.OsVersion
	data.ProductABI = resp.ProductAbi
	data.ProductId = resp.ProductID
	data.ProductVersion = resp.ProductVersion
	data.NewPackages = len(resp.NewPackages)
	data.UpgradePackages = len(resp.UpgradePackages)
	data.NeedsReboot = -1
	data.UpgradeNeedsReboot = -1

	// The following fields won't be set if there was not a firmware update check
	// since the last reboot.
	if resp.LastCheck == "" {
		return data, nil
	}

	if tNeedsReboot, err := parseStringToInt(resp.NeedsReboot, url); err != nil {
		c.log.Warn("firmware: failed to parse NeedsReboot", "details", err)
	} else {
		data.NeedsReboot = tNeedsReboot
	}

	if tUpgradeNeedsReboot, err := parseStringToInt(resp.Product.ProductCheck.UpgradeNeedsReboot, url); err != nil {
		c.log.Warn("firmware: failed to parse UpgradeNeedsReboot", "details", err)
	} else {
		data.UpgradeNeedsReboot = tUpgradeNeedsReboot
	}

	return data, nil
}

// Calling this function causes the OPNsense instance to check for firmware updates.
// This is a costly operation and should not be done on every scrape.
// This takes a long time (30s+) to complete and should not be done on every scrape.
func (c *Client) TriggerFirmwareStatusUpdate() error {
	url, ok := c.endpoints["firmware"]
	if !ok {
		return &APICallError{
			Endpoint:   "firmware",
			Message:    "Missing endpoint 'firmwareStatus'",
			StatusCode: 0,
		}
	}

	var resp any
	if err := c.doLongRunning("POST", url, strings.NewReader("{}"), &resp); err != nil {
		return &APICallError{
			Endpoint:   "firmware",
			Message:    fmt.Sprintf("failed while triggering firmware status update: %v", err),
			StatusCode: 0,
		}
	}

	return nil
}
