package software

import (
	"fmt"

	"sentineldesk/agent/internal/api"
	"sentineldesk/agent/internal/deviceid"
)

type SendSoftwareRequest struct {
	DeviceID string     `json:"device_id"`
	Software []Software `json:"software"`
}

func SendSoftware(software []Software) error {

	client := api.NewClient()

	req := SendSoftwareRequest{
		DeviceID: deviceid.Get(),
		Software: software,
	}

	resp, err := client.R().
		SetBody(req).
		Post("/api/v1/devices/software")

	if err != nil {
		return err
	}

	fmt.Println("SoftwareInventory:", resp.Status())

	return nil
}
