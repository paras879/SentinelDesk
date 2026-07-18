package heartbeat

import (
	"fmt"

	"sentineldesk/agent/internal/api"
	"sentineldesk/agent/internal/deviceid"
	"sentineldesk/agent/internal/system"
)

type HeartbeatRequest struct {
	DeviceID        string `json:"device_id"`
	Hostname        string `json:"hostname"`
	IPAddress       string `json:"ip_address"`
	MACAddress      string `json:"mac_address"`
	ConnectedSubnet string `json:"connected_subnet"`
	DefaultGateway  string `json:"default_gateway"`
	NetworkGroupID  string `json:"network_group_id"`
}

func SendHeartbeat() error {

	info, err := system.GetSystemInfo()
	if err != nil {
		return err
	}

	client := api.NewClient()

	req := HeartbeatRequest{
		DeviceID:        deviceid.Get(),
		Hostname:        info.Hostname,
		IPAddress:       info.IPAddress,
		MACAddress:      info.MACAddress,
		ConnectedSubnet: info.ConnectedSubnet,
		DefaultGateway:  info.DefaultGateway,
		NetworkGroupID:  info.NetworkGroupID,
	}

	resp, err := client.R().
		SetBody(req).
		Post("/api/v1/devices/heartbeat")

	if err != nil {
		return err
	}

	if resp.StatusCode() < 200 || resp.StatusCode() >= 300 {
		return fmt.Errorf("heartbeat rejected: %s", resp.Status())
	}

	fmt.Println("Heartbeat:", resp.Status())

	return nil
}
