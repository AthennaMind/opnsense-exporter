package opnsense

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

	tNeedsReboot, err := parseStringToInt(resp.NeedsReboot, url)
	if err != nil {
		c.log.Warn("firmware: failed to parse NeedsRebot", "details", err)
		data.NeedsReboot = -1
	}
	data.NeedsReboot = tNeedsReboot

	data.NewPackages = len(resp.NewPackages)
	data.UpgradePackages = len(resp.UpgradePackages)

	tUpgradeNeedsReboot, err := parseStringToInt(resp.Product.ProductCheck.UpgradeNeedsReboot, url)
	if err != nil {
		c.log.Warn("firmware: failed to parse UpgradeNeedsReboot", "details", err)
		data.UpgradeNeedsReboot = -1
	}
	data.UpgradeNeedsReboot = tUpgradeNeedsReboot

	return data, nil
}
