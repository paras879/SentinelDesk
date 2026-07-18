package process

import (
	"fmt"

	"sentineldesk/agent/internal/api"
	"sentineldesk/agent/internal/deviceid"
)

type SendProcessRequest struct {
	DeviceID  string       `json:"device_id"`
	Processes []ProcessInfo `json:"processes"`
}

func SendProcesses(processes []ProcessInfo) error {

	client := api.NewClient()

	req := SendProcessRequest{
		DeviceID:  deviceid.Get(),
		Processes: processes,
	}

	resp, err := client.R().
		SetBody(req).
		Post("/api/v1/devices/processes")

	if err != nil {
		return err
	}

	fmt.Println("ProcessInventory:", resp.Status())

	return nil
}
