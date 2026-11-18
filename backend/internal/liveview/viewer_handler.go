package liveview

import (
	"context"
	"encoding/json"
	"log"

	"github.com/gofiber/contrib/websocket"
)

func HandleAdminView(c *websocket.Conn) {
	deviceID, ok := c.Locals("deviceID").(string)
	if !ok || deviceID == "" {
		return
	}

	viewerID := c.RemoteAddr().String()

	if GlobalRedisHub != nil {
		GlobalRedisHub.RegisterLocalViewer(deviceID, viewerID, c)
		defer GlobalRedisHub.UnregisterLocalViewer(deviceID, viewerID)
	}

	stream := GlobalHub.GetOrCreate(deviceID)

	stream.mu.Lock()
	stream.Viewers[viewerID] = c
	agentConnected := stream.Agent != nil
	latestFrame := stream.LatestFrame
	stream.mu.Unlock()

	if latestFrame != nil {
		c.WriteMessage(websocket.BinaryMessage, latestFrame)
	}

	if agentConnected {
		c.WriteJSON(map[string]string{"type": "status", "agent_connected": "true"})
	} else {
		c.WriteJSON(map[string]string{"type": "status", "agent_connected": "false"})
	}

	defer func() {
		stream.mu.Lock()
		delete(stream.Viewers, viewerID)
		agentStillConnected := stream.Agent != nil
		viewersLeft := len(stream.Viewers)
		stream.mu.Unlock()

		if !agentStillConnected && viewersLeft == 0 {
			GlobalHub.Remove(deviceID)
		}
	}()

	for {
		_, msg, err := c.ReadMessage()
		if err != nil {
			log.Println("Viewer disconnected:", deviceID, err)
			break
		}

		var pkt struct {
			Type string `json:"type"`
		}
		if err := json.Unmarshal(msg, &pkt); err != nil || pkt.Type == "" {
			continue
		}

		log.Printf("[Control] device=%s viewer=%s type=%s", deviceID, viewerID, pkt.Type)

		if GlobalRedisHub != nil {
			_ = GlobalRedisHub.PublishControl(context.Background(), deviceID, msg)
		} else {
			stream.mu.RLock()
			agent := stream.Agent
			stream.mu.RUnlock()

			if agent != nil {
				select {
				case stream.ControlChan <- msg:
				default:
					log.Printf("[Control] device=%s control channel full, dropping event", deviceID)
				}
			}
		}
	}
}
