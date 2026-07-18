package liveview

import (
	"log"

	"github.com/gofiber/contrib/websocket"
)

func HandleAgentStream(c *websocket.Conn) {
	deviceID, ok := c.Locals("deviceID").(string)
	if !ok || deviceID == "" {
		return
	}

	stream := GlobalHub.GetOrCreate(deviceID)

	stream.mu.Lock()
	if stream.Agent != nil {
		stream.Agent.Close()
	}
	stream.Agent = c
	viewers := make([]*websocket.Conn, 0, len(stream.Viewers))
	for _, v := range stream.Viewers {
		viewers = append(viewers, v)
	}
	stream.mu.Unlock()

	for _, v := range viewers {
		v.WriteJSON(map[string]string{"type": "status", "agent_connected": "true"})
	}

	done := make(chan struct{})
	defer close(done)

	go func() {
		for {
			select {
			case <-done:
				return
			case msg := <-stream.ControlChan:
				stream.mu.RLock()
				agent := stream.Agent
				stream.mu.RUnlock()
				if agent != nil {
					agent.WriteMessage(websocket.TextMessage, msg)
				}
			}
		}
	}()

	defer func() {
		stream.mu.Lock()
		if stream.Agent == c {
			stream.Agent = nil
		}
		viewersLeft := make([]*websocket.Conn, 0, len(stream.Viewers))
		for _, v := range stream.Viewers {
			viewersLeft = append(viewersLeft, v)
		}
		stream.mu.Unlock()

		for _, v := range viewersLeft {
			v.WriteJSON(map[string]string{"type": "status", "agent_connected": "false"})
		}

		stream.mu.Lock()
		if stream.Agent == nil && len(stream.Viewers) == 0 {
			stream.mu.Unlock()
			GlobalHub.Remove(deviceID)
		} else {
			stream.mu.Unlock()
		}
	}()

	for {
		msgType, data, err := c.ReadMessage()
		if err != nil {
			log.Println("Agent stream disconnected:", deviceID, err)
			break
		}
		if msgType != websocket.BinaryMessage {
			continue
		}

		stream.mu.RLock()
		var failed []string
		for viewerID, viewer := range stream.Viewers {
			if err := viewer.WriteMessage(websocket.BinaryMessage, data); err != nil {
				viewer.Close()
				failed = append(failed, viewerID)
			}
		}
		stream.mu.RUnlock()

		if len(failed) > 0 {
			stream.mu.Lock()
			for _, id := range failed {
				delete(stream.Viewers, id)
			}
			stream.mu.Unlock()
		}
	}
}
