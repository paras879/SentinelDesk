package systeminfo

import (
	"fmt"

	"sentineldesk/agent/internal/api"
)

func SendSystemInfo(info *SystemInfo) error {

	client := api.NewClient()

	resp, err := client.R().
		SetBody(info).
		Post("/api/v1/devices/system-info")

	if err != nil {
		return err
	}

	fmt.Println("SystemInfo:", resp.Status())

	return nil
}
