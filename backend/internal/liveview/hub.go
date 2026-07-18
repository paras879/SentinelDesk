package liveview

import (
	"sync"

	"github.com/gofiber/contrib/websocket"
)

type StreamState struct {
	mu          sync.RWMutex
	Agent       *websocket.Conn
	Viewers     map[string]*websocket.Conn
	ControlChan chan []byte
	LatestFrame []byte
}

type Hub struct {
	mu      sync.RWMutex
	Streams map[string]*StreamState
}

var GlobalHub = &Hub{
	Streams: make(map[string]*StreamState),
}

func (h *Hub) GetOrCreate(deviceID string) *StreamState {
	h.mu.Lock()
	defer h.mu.Unlock()
	if s, ok := h.Streams[deviceID]; ok {
		return s
	}
	s := &StreamState{
		Viewers:     make(map[string]*websocket.Conn),
		ControlChan: make(chan []byte, 64),
	}
	h.Streams[deviceID] = s
	return s
}

func (h *Hub) Get(deviceID string) *StreamState {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.Streams[deviceID]
}

func (h *Hub) Remove(deviceID string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.Streams, deviceID)
}
