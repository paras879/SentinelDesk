package liveview

import (
	"log"
	"time"

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
			stream.LatestFrame = nil
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

	c.SetReadDeadline(time.Now().Add(90 * time.Second))
	c.SetPongHandler(func(string) error {
		c.SetReadDeadline(time.Now().Add(90 * time.Second))
		return nil
	})

	var frameCount int
	for {
		msgType, data, err := c.ReadMessage()
		if err != nil {
			log.Println("Agent stream disconnected:", deviceID, err)
			break
		}
		if msgType != websocket.BinaryMessage {
			continue
		}

		frameCount++
		if frameCount%30 == 1 {
			log.Printf("[Backend] frame #%d from agent %s size=%d viewers=%d",
				frameCount, deviceID, len(data), len(stream.Viewers))
		}

		stream.mu.Lock()
		stream.LatestFrame = data
		var failed []string
		for viewerID, viewer := range stream.Viewers {
			select {
			case <-done:
				stream.mu.Unlock()
				return
			default:
			}
			if err := viewer.WriteMessage(websocket.BinaryMessage, data); err != nil {
				viewer.Close()
				failed = append(failed, viewerID)
			}
		}
		stream.mu.Unlock()

		if len(failed) > 0 {
			stream.mu.Lock()
			for _, id := range failed {
				delete(stream.Viewers, id)
			}
			stream.mu.Unlock()
		}
	}
}
