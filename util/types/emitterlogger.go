package types

// EmitterLogger is an interface that can emit events and log messages.
type EmitterLogger interface {
	Debug(event *Event, factoryParams ...any)
	Info(event *Event, factoryParams ...any)
	Warn(event *Event, factoryParams ...any)
	Error(event *Event, factoryParams ...any)
	Fatal(event *Event, factoryParams ...any)
	Trace(event *Event, factoryParams ...any)
}
