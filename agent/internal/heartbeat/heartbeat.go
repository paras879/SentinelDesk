package heartbeat

import (
	"fmt"

	"sentineldesk/agent/internal/api"
	"sentineldesk/agent/internal/config"
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
	endpoint := config.Get().ServerURL + "/api/v1/devices/heartbeat"

	info, err := system.GetSystemInfo()
	if err != nil {
		return fmt.Errorf("failed to get system info: %w", err)
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
		return fmt.Errorf("cannot reach backend at %s — %w", endpoint, err)
	}

	if resp.StatusCode() < 200 || resp.StatusCode() >= 300 {
		return fmt.Errorf("heartbeat rejected [%s] %s — endpoint: %s", resp.Status(), string(resp.Body()), endpoint)
	}

	return nil
}
