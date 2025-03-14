package util

import (
	"sync"
	"testing"
	"time"

	"github.com/pureapi/pureapi-core/util/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegisterAndEmit(t *testing.T) {
	emitter := NewEventEmitter()
	ch := make(chan *types.Event, 1)

	// Register a listener that sends the event to the channel.
	emitter.RegisterListener("test", func(e *types.Event) {
		ch <- e
	})

	evt := types.NewEvent("test", "foo")
	emitter.Emit(evt)

	select {
	case received := <-ch:
		assert.Equal(t, "foo", received.Message)
	case <-time.After(500 * time.Millisecond):
		t.Error("timeout waiting for event callback")
	}
}

func TestMultipleListeners(t *testing.T) {
	emitter := NewEventEmitter()
	var wg sync.WaitGroup
	count := 0
	var mu sync.Mutex
	numListeners := 3

	// Register several listeners for the same event.
	for i := 0; i < numListeners; i++ {
		wg.Add(1)
		emitter.RegisterListener("mult", func(e *types.Event) {
			mu.Lock()
			count++
			mu.Unlock()
			wg.Done()
		})
	}

	emitter.Emit(types.NewEvent("mult", "bar"))

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		assert.Equal(t, numListeners, count)
	case <-time.After(500 * time.Millisecond):
		t.Error("timeout waiting for multiple listeners")
	}
}

func TestRemoveListener(t *testing.T) {
	emitter := NewEventEmitter()
	var wg sync.WaitGroup
	var mu sync.Mutex
	callCount := 0

	// Register two listeners for the event "rm".
	wg.Add(2)
	emitter.RegisterListener("rm", func(e *types.Event) {
		mu.Lock()
		callCount++
		mu.Unlock()
		wg.Done()
	})
	emitter.RegisterListener("rm", func(e *types.Event) {
		mu.Lock()
		callCount++
		mu.Unlock()
		wg.Done()
	})

	// Since the IDs are generated internally, inspect the listeners map.
	emitter.mu.RLock()
	listeners := emitter.listeners["rm"]
	require.Len(t, listeners, 2)
	idToRemove := listeners[0].id
	emitter.mu.RUnlock()

	// Remove the first listener.
	emitter.RemoveListener("rm", idToRemove)

	// After removal, verify only one listener remains.
	emitter.mu.RLock()
	remaining := emitter.listeners["rm"]
	emitter.mu.RUnlock()
	require.Len(t, remaining, 1)

	// Reset the wait group to expect one callback.
	wg.Add(-1) // We already had 2 added, so subtract one.
	// Reset call count.
	mu.Lock()
	callCount = 0
	mu.Unlock()

	// Emit event.
	emitter.Emit(types.NewEvent("rm", "remove"))

	select {
	case <-time.After(500 * time.Millisecond):
		// Give a little time for the callback.
	}

	mu.Lock()
	finalCount := callCount
	mu.Unlock()
	assert.Equal(t, 1, finalCount, "only one listener should be invoked after removal")
}

func TestEmitNoListeners(t *testing.T) {
	emitter := NewEventEmitter()
	// Emit an event for which no listener is registered. Should not panic.
	assert.NotPanics(t, func() {
		emitter.Emit(types.NewEvent("none", "no listeners"))
	})
}

func TestWithTimeoutOption(t *testing.T) {
	timeoutDuration := 1000 * time.Millisecond
	emitter := NewEventEmitter(WithTimeout(timeoutDuration))
	// Apply the timeout option.
	require.NotNil(t, emitter.timeout)
	assert.Equal(t, timeoutDuration, *emitter.timeout)

	ch := make(chan *types.Event, 1)
	// Register a listener that simulates some work.
	emitter.RegisterListener("timeout", func(e *types.Event) {
		time.Sleep(2000 * time.Millisecond)
		ch <- e
	})
	emitter.Emit(types.NewEvent("timeout", "with timeout"))

	select {
	case received := <-ch:
		assert.Equal(t, "with timeout", received.Message)
	case <-time.After(5000 * time.Millisecond):
		t.Error("timeout waiting for callback with timeout option")
	}
}

func TestConcurrentEmit(t *testing.T) {
	emitter := NewEventEmitter()
	var wg sync.WaitGroup
	var mu sync.Mutex
	callCount := 0
	numListeners := 2
	numEmitters := 10

	// Register two listeners for a concurrent event.
	for i := 0; i < numListeners; i++ {
		emitter.RegisterListener("concurrent", func(e *types.Event) {
			mu.Lock()
			callCount++
			mu.Unlock()
			wg.Done()
		})
	}

	totalCalls := numListeners * numEmitters
	wg.Add(totalCalls)

	// Emit events concurrently.
	for i := 0; i < numEmitters; i++ {
		go emitter.Emit(types.NewEvent("concurrent", "concurrent event"))
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		assert.Equal(t, totalCalls, callCount)
	case <-time.After(1 * time.Second):
		t.Error("timeout waiting for concurrent emits")
	}
}
