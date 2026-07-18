package software

import (
	"fmt"

	"sentineldesk/agent/internal/api"
	"sentineldesk/agent/internal/config"
	"sentineldesk/agent/internal/deviceid"
)

type SendSoftwareRequest struct {
	DeviceID string     `json:"device_id"`
	Software []Software `json:"software"`
}

func SendSoftware(software []Software) error {
	endpoint := config.Get().ServerURL + "/api/v1/devices/software"

	client := api.NewClient()

	req := SendSoftwareRequest{
		DeviceID: deviceid.Get(),
		Software: software,
	}

	resp, err := client.R().
		SetBody(req).
		Post("/api/v1/devices/software")

	if err != nil {
		return fmt.Errorf("cannot reach backend at %s — %w", endpoint, err)
	}

	if resp.StatusCode() < 200 || resp.StatusCode() >= 300 {
		return fmt.Errorf("software inventory rejected [%s] %s — endpoint: %s", resp.Status(), string(resp.Body()), endpoint)
	}

	return nil
}
