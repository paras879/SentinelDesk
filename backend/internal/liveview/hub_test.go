package liveview

import (
	"sync"
	"testing"
)

func newTestHub() *Hub {
	return &Hub{
		Streams: make(map[string]*StreamState),
	}
}

func TestGetOrCreate_CreatesNew(t *testing.T) {
	h := newTestHub()
	s := h.GetOrCreate("device-1")
	if s == nil {
		t.Fatal("GetOrCreate returned nil for new device")
	}
}

func TestGetOrCreate_Idempotent(t *testing.T) {
	h := newTestHub()
	s1 := h.GetOrCreate("device-1")
	s2 := h.GetOrCreate("device-1")
	if s1 != s2 {
		t.Fatal("GetOrCreate returned different pointers for the same deviceID")
	}
}

func TestGet_UnknownDevice_ReturnsNil(t *testing.T) {
	h := newTestHub()
	s := h.Get("nonexistent")
	if s != nil {
		t.Fatal("expected nil for unknown device, got non-nil")
	}
}

func TestGet_AfterGetOrCreate_ReturnsState(t *testing.T) {
	h := newTestHub()
	created := h.GetOrCreate("device-2")
	got := h.Get("device-2")
	if got == nil {
		t.Fatal("Get returned nil after GetOrCreate")
	}
	if got != created {
		t.Fatal("Get returned different pointer than GetOrCreate")
	}
}

func TestRemove_MakesGetReturnNil(t *testing.T) {
	h := newTestHub()
	h.GetOrCreate("device-3")
	h.Remove("device-3")
	if s := h.Get("device-3"); s != nil {
		t.Fatal("expected nil after Remove, got non-nil")
	}
}

func TestGetOrCreate_Concurrent(t *testing.T) {
	h := newTestHub()
	const goroutines = 10
	results := make([]*StreamState, goroutines)
	var wg sync.WaitGroup
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func(idx int) {
			defer wg.Done()
			results[idx] = h.GetOrCreate("shared-device")
		}(i)
	}
	wg.Wait()
	for i, s := range results {
		if s == nil {
			t.Errorf("goroutine %d got nil from GetOrCreate", i)
		}
	}
	// All must be the same pointer
	for i := 1; i < goroutines; i++ {
		if results[i] != results[0] {
			t.Errorf("goroutine %d returned different pointer than goroutine 0", i)
		}
	}
}
