package opnsense

type servicesSearchResponse struct {
	Total    int `json:"total"`
	RowCount int `json:"rowCount"`
	Current  int `json:"current"`
	Rows     []struct {
		ID          string `json:"id"`
		Locked      int    `json:"locked"`
		Running     int    `json:"running"`
		Description string `json:"description"`
		Name        string `json:"name"`
	} `json:"rows"`
}

type ServiceStatus int

const (
	ServiceStatusStopped ServiceStatus = iota
	ServiceStatusRunning
	ServiceStatusUnknown
)

type Service struct {
	Status      ServiceStatus
	Description string
	Name        string
}

type Services struct {
	TotalRunning int
	TotalStopped int
	Services     []Service
}

func (c *Client) FetchServices() (Services, *APICallError) {
	var resp servicesSearchResponse
	var services Services

	url, ok := c.endpoints["services"]

	if !ok {
		return services, &APICallError{
			Endpoint:   "services",
			Message:    "endpoint not found",
			StatusCode: 0,
		}
	}
	err := c.do("GET", url, nil, &resp)

	if err != nil {
		return services, err
	}
	for _, service := range resp.Rows {

		switch service.Running {
		case 0:
			services.TotalStopped++
		case 1:
			services.TotalRunning++
		}

		s := Service{
			Status:      ServiceStatus(service.Running),
			Description: service.Description,
			Name:        service.Name,
		}
		services.Services = append(services.Services, s)
	}

	return services, nil
}
