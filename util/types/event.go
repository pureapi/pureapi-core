package types

// EventType represents the type of event.
type EventType string

// Event represents an emitted event.
type Event struct {
	Type    EventType
	Message string
	Data    any
}

// NewEvent creates a new event.
//
// Parameters:
//   - eventType: The type of the event.
//   - message: The message of the event.
//   - data: The optional data of the event.
//
// Returns:
//   - *Event: A new Event instance.
func NewEvent(eventType EventType, message string, data ...any) *Event {
	return &Event{
		Type:    eventType,
		Message: message,
		Data:    data,
	}
}

// EventCallback is a function that handles an event.
type EventCallback func(event *Event)

// EventEmitter is responsible for emitting events.
type EventEmitter interface {
	RegisterListener(eventType EventType, callback EventCallback) EventEmitter
	RemoveListener(eventType EventType, id string)
	Emit(event *Event)
}
