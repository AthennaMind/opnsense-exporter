package opnsense

type servicesSearchResponse struct {
	Rows []struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		Locked      int    `json:"locked"`
		Running     int    `json:"running"`
	} `json:"rows"`
	Total    int `json:"total"`
	RowCount int `json:"rowCount"`
	Current  int `json:"current"`
}

type ServiceStatus int

const (
	ServiceStatusStopped ServiceStatus = iota
	ServiceStatusRunning
	ServiceStatusUnknown
)

type Service struct {
	Description string
	Name        string
	Status      ServiceStatus
}

type Services struct {
	Services     []Service
	TotalRunning int
	TotalStopped int
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
