package process

import (
	"fmt"

	"sentineldesk/agent/internal/api"
	"sentineldesk/agent/internal/config"
	"sentineldesk/agent/internal/deviceid"
)

type SendProcessRequest struct {
	DeviceID  string        `json:"device_id"`
	Processes []ProcessInfo `json:"processes"`
}

func SendProcesses(processes []ProcessInfo) error {
	endpoint := config.Get().ServerURL + "/api/v1/devices/processes"

	client := api.NewClient()

	req := SendProcessRequest{
		DeviceID:  deviceid.Get(),
		Processes: processes,
	}

	resp, err := client.R().
		SetBody(req).
		Post("/api/v1/devices/processes")

	if err != nil {
		return fmt.Errorf("cannot reach backend at %s — %w", endpoint, err)
	}

	if resp.StatusCode() < 200 || resp.StatusCode() >= 300 {
		return fmt.Errorf("process inventory rejected [%s] %s — endpoint: %s", resp.Status(), string(resp.Body()), endpoint)
	}

	return nil
}
