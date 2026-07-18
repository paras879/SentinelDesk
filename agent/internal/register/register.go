package register

import (
	"fmt"
	"log"

	"sentineldesk/agent/internal/api"
	"sentineldesk/agent/internal/config"
	"sentineldesk/agent/internal/deviceid"
	"sentineldesk/agent/internal/system"
)

func logRegistrationURL() {
	cfg := config.Get()
	registrationURL := cfg.ServerURL + "/api/v1/devices/register"
	log.Printf("Config File Path: config.yaml")
	log.Printf("Loaded Server URL: %s", cfg.ServerURL)
	log.Printf("Final Registration URL: %s", registrationURL)
}

type RegisterRequest struct {
	DeviceID        string                `json:"device_id"`
	DeviceName      string                `json:"device_name"`
	Hostname        string                `json:"hostname"`
	Username        string                `json:"username"`
	OS              string                `json:"os"`
	OSVersion       string                `json:"os_version"`
	IPAddress       string                `json:"ip_address"`
	MACAddress      string                `json:"mac_address"`
	ConnectedSubnet string                `json:"connected_subnet"`
	NetworkAdapters []system.NetworkAdapter `json:"network_adapters"`
	DefaultGateway  string                `json:"default_gateway"`
	NetworkGroupID  string                `json:"network_group_id"`
}

func RegisterDevice() error {

	logRegistrationURL()

	info, err := system.GetSystemInfo()
	if err != nil {
		return err
	}

	client := api.NewClient()

	req := RegisterRequest{
		DeviceID:        deviceid.Get(),
		DeviceName:      info.DeviceName,
		Hostname:        info.Hostname,
		Username:        info.Username,
		OS:              info.OS,
		OSVersion:       info.OSVersion,
		IPAddress:       info.IPAddress,
		MACAddress:      info.MACAddress,
		ConnectedSubnet: info.ConnectedSubnet,
		NetworkAdapters: info.NetworkAdapters,
		DefaultGateway:  info.DefaultGateway,
		NetworkGroupID:  info.NetworkGroupID,
	}

	resp, err := client.R().
		SetBody(req).
		Post("/api/v1/devices/register")

	if err != nil {
		return fmt.Errorf("cannot reach backend at %s — %v", config.Get().ServerURL, err)
	}

	switch resp.StatusCode() {

	case 200:
		fmt.Println("✓ Device already registered, record updated.")

	case 201:
		fmt.Println("✓ Device registered successfully.")

	default:
		return fmt.Errorf("registration failed: %s — %s", resp.Status(), string(resp.Body()))
	}

	return nil
}
