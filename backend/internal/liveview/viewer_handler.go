package liveview

import (
	"log"

	"github.com/gofiber/contrib/websocket"
)

func HandleAdminView(c *websocket.Conn) {
	deviceID, ok := c.Locals("deviceID").(string)
	if !ok || deviceID == "" {
		return
	}

	stream := GlobalHub.GetOrCreate(deviceID)

	viewerID := c.RemoteAddr().String()

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

		stream.mu.RLock()
		agent := stream.Agent
		stream.mu.RUnlock()

		if agent != nil {
			select {
			case stream.ControlChan <- msg:
			default:
			}
		}
	}
}
