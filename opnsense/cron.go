package opnsense

import (
	"fmt"
	"strings"

	"github.com/go-kit/log/level"
)

type cronSearchResponse struct {
	Rows []struct {
		UUID        string `json:"uuid"`
		Enabled     string `json:"enabled"`
		Minutes     string `json:"minutes"`
		Hours       string `json:"hours"`
		Days        string `json:"days"`
		Months      string `json:"months"`
		Weekdays    string `json:"weekdays"`
		Description string `json:"description"`
		Command     string `json:"command"`
		Origin      string `json:"origin"`
	} `json:"rows"`
	RowCount int `json:"rowCount"`
	Total    int `json:"total"`
	Current  int `json:"current"`
}

type CronStatus int

const (
	CronStatusDisabled CronStatus = iota
	CronStatusEnabled
)

type Cron struct {
	UUID        string
	Schedule    string
	Description string
	Command     string
	Origin      string
	Status      CronStatus
}

type CronTable struct {
	Cron         []Cron
	TotalEntries int
}

const fetchCronPayload = `{"current":1,"rowCount":-1,"sort":{},"searchPhrase":"","resolve":"no"}`

func (c *Client) FetchCronTable() (CronTable, *APICallError) {
	var resp cronSearchResponse
	var cronTable CronTable

	path, ok := c.endpoints["cronJobs"]
	if !ok {
		return cronTable, &APICallError{
			Endpoint:   "cronJobs",
			Message:    "endpoint not found",
			StatusCode: 0,
		}
	}

	if err := c.do("POST", path, strings.NewReader(fetchCronPayload), &resp); err != nil {
		return cronTable, err
	}

	for _, cron := range resp.Rows {
		cronTable.TotalEntries++

		intStatus, err := parseStringToInt(cron.Enabled, path)
		if err != nil {
			level.Warn(c.log).
				Log("msg", "unable to parse cron entry status", "err", err)
			continue
		}

		cronTable.Cron = append(cronTable.Cron, Cron{
			UUID:        cron.UUID,
			Status:      CronStatus(intStatus),
			Description: cron.Description,
			Schedule:    fmt.Sprintf("%s %s %s %s %s", cron.Minutes, cron.Hours, cron.Days, cron.Months, cron.Weekdays),
			Command:     cron.Command,
			Origin:      cron.Origin,
		})
	}

	return cronTable, nil
}
