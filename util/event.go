package util

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/pureapi/pureapi-core/util/types"
)

// eventListener wraps a listener callback with an ID.
type eventListener struct {
	id       string
	callback func(*types.Event)
}

// EventEmitter is responsible for emitting events.
type defaultEventEmitter struct {
	listeners map[types.EventType][]eventListener
	mu        sync.RWMutex   // Mutex for thread safety when emitting events.
	counter   int            // Used to generate unique IDs for listeners.
	timeout   *time.Duration // Optional timeout for each callback.
}

// defaultEventEmitter implements the EventEmitter interface.
var _ types.EventEmitter = (*defaultEventEmitter)(nil)

// NewEventEmitter creates a new defaultEventEmitter.
//
// Parameters:
//   - opts: Options to configure the defaultEventEmitter.
//
// Returns:
//   - *defaultEventEmitter: A new defaultEventEmitter.
func NewEventEmitter() *defaultEventEmitter {
	eventEmitter := &defaultEventEmitter{
		listeners: make(map[types.EventType][]eventListener),
		mu:        sync.RWMutex{},
		counter:   0,
		timeout:   nil,
	}
	return eventEmitter
}

// WithTimeout sets the timeout for each callback. If the timeout is exceeded,
// an error message will be printed to stderr. It will return a new
// eventEmitterOption.
//
// Parameters:
//   - timeout: The timeout duration.
//
// Returns:
//   - *defaultEventEmitter: A new defaultEventEmitter.
func (e *defaultEventEmitter) WithTimeout(
	timeout *time.Duration,
) *defaultEventEmitter {
	new := NewEventEmitter()
	new.timeout = timeout
	return new
}

// RegisterListener registers a listener for a specific event type.
//
// Parameters:
//   - eventType: The type of the event.
//   - callback: The function to call when the event is emitted.
//
// Returns:
//   - *eventEmitter: The eventEmitter.
func (e *defaultEventEmitter) RegisterListener(
	eventType types.EventType, callback types.EventCallback,
) types.EventEmitter {
	// Generate a unique ID for the listener.
	e.mu.Lock()
	defer e.mu.Unlock()
	e.counter++
	id := fmt.Sprintf("%s-%d", eventType, e.counter)

	// Add the listener to the list.
	e.listeners[eventType] = append(e.listeners[eventType], eventListener{
		id:       id,
		callback: callback,
	})
	return e
}

// RemoveListener removes a listener for a specific event type.
//
// Parameters:
//   - eventType: The type of the event.
//   - listener: The listener function.
//
// Returns:
//   - *eventEmitter: The eventEmitter.
func (e *defaultEventEmitter) RemoveListener(
	eventType types.EventType, id string,
) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if list, found := e.listeners[eventType]; found {
		for i, l := range list {
			if l.id == id {
				// Remove the listener with the matching ID.
				e.listeners[eventType] = append(list[:i], list[i+1:]...)
				break
			}
		}
	}
}

// Emit emits an event to all registered listeners. It runs each callback in a
// separate goroutine. If timeout is set for the eventEmitter, the callbacks
// will be run with the specified timeout. If the timeout is exceeded, an error
// message will be printed to stderr.
//
// Parameters:
//   - event: The event to emit.
//
// Returns:
//   - *eventEmitter: The eventEmitter.
func (e *defaultEventEmitter) Emit(event *types.Event) {
	e.mu.RLock()
	listeners := e.listeners[event.Type]
	e.mu.RUnlock()
	// Determine the timeout for each callback.
	var timeout *time.Duration
	if e.timeout != nil {
		timeout = new(time.Duration)
		*timeout = *e.timeout
	}
	// Run each callback in a separate goroutine.
	for _, l := range listeners {
		go func(cb types.EventCallback, timeout *time.Duration) {
			runCallback(event, cb, timeout)
		}(l.callback, timeout)
	}
}

// runCallback runs a callback with an optional timeout.
func runCallback(
	event *types.Event, cb types.EventCallback, timeout *time.Duration,
) {
	done := make(chan struct{})
	go func() {
		cb(event)
		close(done)
	}()
	if timeout != nil {
		select {
		case <-done:
			// Callback completed within the timeout.
		case <-time.After(*timeout):
			// Timeout reached; the callback might still be in the background.
			fmt.Fprintf(
				os.Stderr,
				"Callback for event %v timed out after %v, event type: %v\n",
				event.Type,
				*timeout,
				event.Type,
			)
		}
	}
}
