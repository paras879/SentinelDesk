package liveview

import (
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/kbinani/screenshot"

	"sentineldesk/agent/internal/config"
	"sentineldesk/agent/internal/deviceid"
	"sentineldesk/agent/internal/remote"
)

func serverToWSURL(serverURL string) string {
	if strings.HasPrefix(serverURL, "https://") {
		return "wss://" + strings.TrimPrefix(serverURL, "https://")
	}
	return "ws://" + strings.TrimPrefix(serverURL, "http://")
}

type Streamer struct {
	capturer *Capturer
}

func NewStreamer() *Streamer {
	return &Streamer{
		capturer: NewCapturer(1280, 720),
	}
}

func (s *Streamer) Start() {
	for {
		if err := s.stream(); err != nil {
			log.Println("LiveView stream error:", err)
		}
		time.Sleep(3 * time.Second)
	}
}

func (s *Streamer) stream() error {
	cfg := config.Get()

	deviceID := deviceid.Get()

	wsURL := serverToWSURL(cfg.ServerURL) + "/ws/live/stream?device_id=" + deviceID

	header := map[string][]string{
		"X-Agent-Key": {cfg.AgentKey},
	}

	conn, _, err := websocket.DefaultDialer.Dial(wsURL, header)
	if err != nil {
		return err
	}
	defer conn.Close()

	log.Println("LiveView stream connected")

	s.capturer.HideDashboard()
	defer s.capturer.ShowDashboard()

	// Send screen info as first message
	bounds := screenshot.GetDisplayBounds(0)
	screenInfo, _ := json.Marshal(map[string]interface{}{
		"type":   "screen_info",
		"width":  bounds.Dx(),
		"height": bounds.Dy(),
	})
	conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
	conn.WriteMessage(websocket.TextMessage, screenInfo)
	conn.SetWriteDeadline(time.Time{})

	// Reader goroutine for remote control commands
	go func() {
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				return
			}
			remote.HandleMessage(msg)
		}
	}()

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		frame, err := s.capturer.Capture()
		if err != nil {
			log.Println("LiveView capture error:", err)
			continue
		}
		if frame == nil {
			continue
		}

		if err := conn.SetWriteDeadline(time.Now().Add(5 * time.Second)); err != nil {
			return err
		}
		if err := conn.WriteMessage(websocket.BinaryMessage, frame); err != nil {
			return err
		}
	}

	return nil
}
