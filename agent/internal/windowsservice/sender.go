package windowsservice

import (
	"fmt"

	"sentineldesk/agent/internal/api"
	"sentineldesk/agent/internal/config"
	"sentineldesk/agent/internal/deviceid"
)

type SendServicesRequest struct {
	DeviceID string        `json:"device_id"`
	Services []ServiceInfo `json:"services"`
}

func SendServices(services []ServiceInfo) error {
	endpoint := config.Get().ServerURL + "/api/v1/devices/services"

	client := api.NewClient()

	req := SendServicesRequest{
		DeviceID: deviceid.Get(),
		Services: services,
	}

	resp, err := client.R().
		SetBody(req).
		Post("/api/v1/devices/services")

	if err != nil {
		return fmt.Errorf("cannot reach backend at %s — %w", endpoint, err)
	}

	if resp.StatusCode() < 200 || resp.StatusCode() >= 300 {
		return fmt.Errorf("windows services rejected [%s] %s — endpoint: %s", resp.Status(), string(resp.Body()), endpoint)
	}

	return nil
}
