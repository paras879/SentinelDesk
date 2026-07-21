package heartbeat

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gorilla/websocket"

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

func serverToWSURL(serverURL string) string {
	if strings.HasPrefix(serverURL, "https://") {
		return "wss://" + strings.TrimPrefix(serverURL, "https://")
	}
	return "ws://" + strings.TrimPrefix(serverURL, "http://")
}

func StartHeartbeatLoop(ctx context.Context, interval int) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		err := connectAndStreamHeartbeat(ctx, interval)
		if err != nil {
			log.Println("[ERROR] Heartbeat WS disconnected:", err)
		}

		select {
		case <-ctx.Done():
			return
		case <-time.After(5 * time.Second):
			// reconnect delay
		}
	}
}

func connectAndStreamHeartbeat(ctx context.Context, interval int) error {
	cfg := config.Get()
	wsURL := serverToWSURL(cfg.ServerURL) + "/api/v1/devices/ws/heartbeat"

	header := map[string][]string{
		"X-Agent-Key": {cfg.AgentKey},
	}

	conn, _, err := websocket.DefaultDialer.Dial(wsURL, header)
	if err != nil {
		return fmt.Errorf("dial failed: %w", err)
	}
	defer conn.Close()

	log.Println("[SUCCESS] Heartbeat WebSocket connected")

	// Start a goroutine to read messages (like PONGs or control messages) so the connection doesn't block
	go func() {
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				return
			}
		}
	}()

	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()

	// Send initial heartbeat
	if err := sendWSHeartbeat(conn); err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			if err := sendWSHeartbeat(conn); err != nil {
				return err
			}
		}
	}
}

func sendWSHeartbeat(conn *websocket.Conn) error {
	info, err := system.GetSystemInfo()
	if err != nil {
		return fmt.Errorf("failed to get system info: %w", err)
	}

	req := HeartbeatRequest{
		DeviceID:        deviceid.Get(),
		Hostname:        info.Hostname,
		IPAddress:       info.IPAddress,
		MACAddress:      info.MACAddress,
		ConnectedSubnet: info.ConnectedSubnet,
		DefaultGateway:  info.DefaultGateway,
		NetworkGroupID:  info.NetworkGroupID,
	}

	data, err := json.Marshal(req)
	if err != nil {
		return err
	}

	conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
	return conn.WriteMessage(websocket.TextMessage, data)
}
