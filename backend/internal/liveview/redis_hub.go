package liveview

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/redis/go-redis/v9"
)

// RedisHub manages multi-node pub/sub for live stream frames and remote control packets.
// When deployed across multiple backend instances, agents and viewers can be connected
// to different instances and still communicate seamlessly.
type RedisHub struct {
	client *redis.Client

	mu           sync.RWMutex
	localAgents  map[string]*websocket.Conn
	localViewers map[string]map[string]*websocket.Conn

	subMu        sync.Mutex
	frameSubs    map[string]context.CancelFunc
	controlSubs  map[string]context.CancelFunc
}

var GlobalRedisHub *RedisHub

// NewRedisHub initializes a connection to Redis and tests it with PING.
func NewRedisHub(redisURL string) (*RedisHub, error) {
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(opt)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &RedisHub{
		client:       client,
		localAgents:  make(map[string]*websocket.Conn),
		localViewers: make(map[string]map[string]*websocket.Conn),
		frameSubs:    make(map[string]context.CancelFunc),
		controlSubs:  make(map[string]context.CancelFunc),
	}, nil
}

// PublishFrame publishes a JPEG binary frame to the Redis channel for the device.
func (h *RedisHub) PublishFrame(ctx context.Context, deviceID string, frame []byte) error {
	channel := "liveview:frame:" + deviceID
	return h.client.Publish(ctx, channel, frame).Err()
}

// PublishControl publishes a remote control JSON packet to the Redis channel for the device.
func (h *RedisHub) PublishControl(ctx context.Context, deviceID string, msg []byte) error {
	channel := "liveview:control:" + deviceID
	return h.client.Publish(ctx, channel, msg).Err()
}

// RegisterLocalAgent registers an agent connection connected to this local instance.
func (h *RedisHub) RegisterLocalAgent(deviceID string, c *websocket.Conn) {
	h.mu.Lock()
	if old, exists := h.localAgents[deviceID]; exists && old != c {
		old.Close()
	}
	h.localAgents[deviceID] = c
	h.mu.Unlock()

	h.ensureControlSubscription(deviceID)
}

// UnregisterLocalAgent unregisters an agent connection from this local instance.
func (h *RedisHub) UnregisterLocalAgent(deviceID string, c *websocket.Conn) {
	h.mu.Lock()
	if current, exists := h.localAgents[deviceID]; exists && current == c {
		delete(h.localAgents, deviceID)
	}
	h.mu.Unlock()

	h.subMu.Lock()
	if cancel, ok := h.controlSubs[deviceID]; ok {
		cancel()
		delete(h.controlSubs, deviceID)
	}
	h.subMu.Unlock()
}

// RegisterLocalViewer registers a viewer WebSocket connected to this local instance.
func (h *RedisHub) RegisterLocalViewer(deviceID string, viewerID string, c *websocket.Conn) {
	h.mu.Lock()
	if _, ok := h.localViewers[deviceID]; !ok {
		h.localViewers[deviceID] = make(map[string]*websocket.Conn)
	}
	h.localViewers[deviceID][viewerID] = c
	h.mu.Unlock()

	h.ensureFrameSubscription(deviceID)
}

// UnregisterLocalViewer removes a viewer WebSocket connected to this local instance.
func (h *RedisHub) UnregisterLocalViewer(deviceID string, viewerID string) {
	h.mu.Lock()
	if viewers, ok := h.localViewers[deviceID]; ok {
		delete(viewers, viewerID)
		if len(viewers) == 0 {
			delete(h.localViewers, deviceID)
		}
	}
	remainingViewers := len(h.localViewers[deviceID])
	h.mu.Unlock()

	if remainingViewers == 0 {
		h.subMu.Lock()
		if cancel, ok := h.frameSubs[deviceID]; ok {
			cancel()
			delete(h.frameSubs, deviceID)
		}
		h.subMu.Unlock()
	}
}

// ensureFrameSubscription subscribes to frame updates for this device if not already subscribed.
func (h *RedisHub) ensureFrameSubscription(deviceID string) {
	h.subMu.Lock()
	if _, exists := h.frameSubs[deviceID]; exists {
		h.subMu.Unlock()
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	h.frameSubs[deviceID] = cancel
	h.subMu.Unlock()

	go func() {
		channel := "liveview:frame:" + deviceID
		pubsub := h.client.Subscribe(ctx, channel)
		defer pubsub.Close()

		ch := pubsub.Channel()
		for {
			select {
			case <-ctx.Done():
				return
			case msg, ok := <-ch:
				if !ok {
					return
				}
				frameData := []byte(msg.Payload)

				h.mu.RLock()
				viewersMap := h.localViewers[deviceID]
				viewers := make([]*websocket.Conn, 0, len(viewersMap))
				for _, v := range viewersMap {
					viewers = append(viewers, v)
				}
				h.mu.RUnlock()

				var failed []*websocket.Conn
				for _, v := range viewers {
					if err := v.WriteMessage(websocket.BinaryMessage, frameData); err != nil {
						failed = append(failed, v)
					}
				}

				if len(failed) > 0 {
					h.mu.Lock()
					if currentMap, ok := h.localViewers[deviceID]; ok {
						for id, conn := range currentMap {
							for _, f := range failed {
								if conn == f {
									conn.Close()
									delete(currentMap, id)
								}
							}
						}
					}
					h.mu.Unlock()
				}
			}
		}
	}()
}

// ensureControlSubscription subscribes to control packets from viewers across the cluster.
func (h *RedisHub) ensureControlSubscription(deviceID string) {
	h.subMu.Lock()
	if _, exists := h.controlSubs[deviceID]; exists {
		h.subMu.Unlock()
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	h.controlSubs[deviceID] = cancel
	h.subMu.Unlock()

	go func() {
		channel := "liveview:control:" + deviceID
		pubsub := h.client.Subscribe(ctx, channel)
		defer pubsub.Close()

		ch := pubsub.Channel()
		for {
			select {
			case <-ctx.Done():
				return
			case msg, ok := <-ch:
				if !ok {
					return
				}
				controlData := []byte(msg.Payload)

				h.mu.RLock()
				agent := h.localAgents[deviceID]
				h.mu.RUnlock()

				if agent != nil {
					if err := agent.WriteMessage(websocket.TextMessage, controlData); err != nil {
						log.Printf("[RedisHub] Failed to forward control packet to agent %s: %v", deviceID, err)
					}
				}
			}
		}
	}()
}
