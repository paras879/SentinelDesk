package systeminfo

import (
	"fmt"

	"sentineldesk/agent/internal/api"
	"sentineldesk/agent/internal/config"
)

func SendSystemInfo(info *SystemInfo) error {
	endpoint := config.Get().ServerURL + "/api/v1/devices/system-info"

	client := api.NewClient()

	resp, err := client.R().
		SetBody(info).
		Post("/api/v1/devices/system-info")

	if err != nil {
		return fmt.Errorf("cannot reach backend at %s — %w", endpoint, err)
	}

	if resp.StatusCode() < 200 || resp.StatusCode() >= 300 {
		return fmt.Errorf("system info rejected [%s] %s — endpoint: %s", resp.Status(), string(resp.Body()), endpoint)
	}

	return nil
}
