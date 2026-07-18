package windowsservice

import (
	"fmt"

	"sentineldesk/agent/internal/api"
	"sentineldesk/agent/internal/deviceid"
)

type SendServicesRequest struct {
	DeviceID string        `json:"device_id"`
	Services []ServiceInfo `json:"services"`
}

func SendServices(services []ServiceInfo) error {

	client := api.NewClient()

	req := SendServicesRequest{
		DeviceID: deviceid.Get(),
		Services: services,
	}

	resp, err := client.R().
		SetBody(req).
		Post("/api/v1/devices/services")

	if err != nil {
		return err
	}

	fmt.Println("WindowsServices:", resp.Status())

	return nil
}
