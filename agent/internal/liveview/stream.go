package liveview

import (
	"encoding/json"
	"log"
	"math"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
	"github.com/kbinani/screenshot"

	"sentineldesk/agent/internal/config"
	"sentineldesk/agent/internal/deviceid"
)

func serverToWSURL(serverURL string) string {
	if strings.HasPrefix(serverURL, "https://") {
		return "wss://" + strings.TrimPrefix(serverURL, "https://")
	}
	return "ws://" + strings.TrimPrefix(serverURL, "http://")
}

type framePacket struct {
	data []byte
}

type Streamer struct {
	capturer       *Capturer
	frameChan      chan *framePacket
	quality        int32
	targetFPS      int32
	scaleWidth     int32
	scaleHeight    int32
	sendTimeEWMA   int64
	encodeTimeEWMA int64
	lastFrameSize  int64
	mu             sync.Mutex
}

func NewStreamer() *Streamer {
	return &Streamer{
		capturer:    NewCapturer(1920, 1080),
		frameChan:   make(chan *framePacket, 1),
		quality:     85,
		targetFPS:   30,
		scaleWidth:  1920,
		scaleHeight: 1080,
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

	bounds := screenshot.GetDisplayBounds(0)
	screenInfo, _ := json.Marshal(map[string]interface{}{
		"type":   "screen_info",
		"width":  bounds.Dx(),
		"height": bounds.Dy(),
	})
	conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
	conn.WriteMessage(websocket.TextMessage, screenInfo)
	conn.SetWriteDeadline(time.Time{})

	go func() {
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				return
			}
			// Remote control disabled
		}
	}()

	stop := make(chan struct{})

	go s.captureLoop(stop)
	go s.sendLoop(conn, stop)

	<-stop
	return nil
}

func (s *Streamer) captureLoop(stop chan struct{}) {

	for {
		select {
		case <-stop:
			return
		default:
		}

		start := time.Now()

		q := int(atomic.LoadInt32(&s.quality))
		w := int(atomic.LoadInt32(&s.scaleWidth))
		h := int(atomic.LoadInt32(&s.scaleHeight))

		s.mu.Lock()
		if s.capturer.width != w || s.capturer.height != h {
			s.capturer.Resize(w, h)
		}
		s.mu.Unlock()

		frame, err := s.capturer.Capture(q)
		if err != nil {
			log.Println("LiveView capture error:", err)
			time.Sleep(100 * time.Millisecond)
			continue
		}

		if frame != nil {
			pkt := &framePacket{data: frame}
			atomic.StoreInt64(&s.lastFrameSize, int64(len(frame)))

			select {
			case <-s.frameChan:
			default:
			}
			s.frameChan <- pkt

			elapsed := time.Since(start).Microseconds()
			s.updateEncodeEWMA(elapsed)
		}

		target := float64(atomic.LoadInt32(&s.targetFPS))
		if target < 1 {
			target = 10
		}
		interval := time.Duration(float64(time.Second) / target)
		elapsed := time.Since(start)
		if elapsed < interval {
			time.Sleep(interval - elapsed)
		}
	}
}

func (s *Streamer) sendLoop(conn *websocket.Conn, stop chan struct{}) {
	defer close(stop)

	adaptTicker := time.NewTicker(2 * time.Second)
	defer adaptTicker.Stop()

	keepaliveTicker := time.NewTicker(30 * time.Second)
	defer keepaliveTicker.Stop()

	var sent int
	for {
		select {
		case <-stop:
			return
		case pkt := <-s.frameChan:
			start := time.Now()
			conn.SetWriteDeadline(time.Now().Add(2 * time.Second))
			if err := conn.WriteMessage(websocket.BinaryMessage, pkt.data); err != nil {
				return
			}
			elapsed := time.Since(start).Microseconds()
			s.updateSendEWMA(elapsed)
			sent++
			if sent%30 == 1 {
				log.Printf("[Send] frame size=%d sendUs=%d total=%d", len(pkt.data), elapsed, sent)
			}
		case <-adaptTicker.C:
			s.adapt()
		case <-keepaliveTicker.C:
			conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
			if err := conn.WriteControl(websocket.PingMessage, nil, time.Now().Add(5*time.Second)); err != nil {
				return
			}
		}
	}
}

func (s *Streamer) updateEncodeEWMA(elapsed int64) {
	old := atomic.LoadInt64(&s.encodeTimeEWMA)
	if old == 0 {
		atomic.StoreInt64(&s.encodeTimeEWMA, elapsed)
		return
	}
	atomic.StoreInt64(&s.encodeTimeEWMA, int64(0.3*float64(elapsed)+0.7*float64(old)))
}

func (s *Streamer) updateSendEWMA(elapsed int64) {
	old := atomic.LoadInt64(&s.sendTimeEWMA)
	if old == 0 {
		atomic.StoreInt64(&s.sendTimeEWMA, elapsed)
		return
	}
	atomic.StoreInt64(&s.sendTimeEWMA, int64(0.3*float64(elapsed)+0.7*float64(old)))
}

func (s *Streamer) adapt() {
	sendUs := atomic.LoadInt64(&s.sendTimeEWMA)
	encodeUs := atomic.LoadInt64(&s.encodeTimeEWMA)

	totalUs := sendUs + encodeUs
	if totalUs < 1 {
		return
	}

	// Target: total time (encode + send) should be ~50% of frame budget
	// 20 FPS → 50ms budget → 25ms target = 25000µs
	idealBudget := 25000.0

	ratio := float64(idealBudget) / float64(totalUs)

	currentFPS := atomic.LoadInt32(&s.targetFPS)
	currentQuality := atomic.LoadInt32(&s.quality)

	if ratio > 1.5 {
		newFPS := int32(math.Min(60, float64(currentFPS)*1.2))
		atomic.StoreInt32(&s.targetFPS, newFPS)

		newQ := int32(math.Min(100, float64(currentQuality)*1.05))
		atomic.StoreInt32(&s.quality, newQ)

		nw := int32(math.Min(2560, float64(atomic.LoadInt32(&s.scaleWidth))*1.2))
		nh := nw * 9 / 16
		atomic.StoreInt32(&s.scaleWidth, nw)
		atomic.StoreInt32(&s.scaleHeight, nh)
	} else if ratio < 0.5 {
		newQ := int32(math.Max(70, float64(currentQuality)*0.9))
		atomic.StoreInt32(&s.quality, newQ)

		if currentQuality <= 75 {
			newFPS := int32(math.Max(15, float64(currentFPS)*0.8))
			atomic.StoreInt32(&s.targetFPS, newFPS)

			nw := int32(math.Max(1280, float64(atomic.LoadInt32(&s.scaleWidth))*0.8))
			nh := nw * 9 / 16
			atomic.StoreInt32(&s.scaleWidth, nw)
			atomic.StoreInt32(&s.scaleHeight, nh)
		}
	}
}
