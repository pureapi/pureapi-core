package util

import "github.com/pureapi/pureapi-core/util/types"

// emitterLogger is a struct that can emit events and log messages.
type emitterLogger struct {
	eventEmitter    types.EventEmitter
	loggerFactoryFn types.LoggerFactoryFn
}

// NewEmitterLogger creates a new EmitterLogger.
//
// Parameters:
//   - eventEmitter: An EventEmitter.
//   - loggerFactoryFn: A LoggerFactoryFn.
//
// Returns:
//   - *emitterLogger: A new emitterLogger.
func NewEmitterLogger(
	eventEmitter types.EventEmitter, loggerFactoryFn types.LoggerFactoryFn,
) *emitterLogger {
	return &emitterLogger{
		eventEmitter:    eventEmitter,
		loggerFactoryFn: loggerFactoryFn,
	}
}

// NewNoopEmitterLogger creates a new EmitterLogger that does nothing.
//
// Returns:
//   - *emitterLogger: A new emitterLogger.
func NewNoopEmitterLogger() *emitterLogger {
	return &emitterLogger{}
}

// WithEventEmitter sets the event emitter for the EmitterLogger. It returns
// a new EmitterLogger.
//
// Parameters:
//   - eventEmitter: An EventEmitter.
//
// Returns:
//   - *emitterLogger: A new emitterLogger.
func (e *emitterLogger) WithEventEmitter(
	eventEmitter types.EventEmitter,
) *emitterLogger {
	return NewEmitterLogger(eventEmitter, e.loggerFactoryFn)
}

// Debug emits an event and logs at the Debug level.
//
// Parameters:
//   - event The event to emit and log.
//   - factoryParams: The parameters to pass to the logger factory function.
func (e *emitterLogger) Debug(event *types.Event, factoryParams ...any) {
	e.emitIfCan(event)
	if e.loggerFactoryFn != nil {
		e.loggerFactoryFn(factoryParams...).Debug(event.Message)
	}
}

// Trace emits an event and logs at the Trace level.
//
// Parameters:
//   - event The event to emit and log.
//   - factoryParams: The parameters to pass to the logger factory function.
func (e *emitterLogger) Trace(event *types.Event, factoryParams ...any) {
	e.emitIfCan(event)
	if e.loggerFactoryFn != nil {
		e.loggerFactoryFn(factoryParams...).Trace(event.Message)
	}
}

// Info emits an event and logs at the Info level.
//
// Parameters:
//   - event The event to emit and log.
//   - factoryParams: The parameters to pass to the logger factory function.
func (e *emitterLogger) Info(event *types.Event, factoryParams ...any) {
	e.emitIfCan(event)
	if e.loggerFactoryFn != nil {
		e.loggerFactoryFn(factoryParams...).Info(event.Message)
	}
}

// Warn emits an event and logs at the Warn level.
//
// Parameters:
//   - event The event to emit and log.
//   - factoryParams: The parameters to pass to the logger factory function.
func (e *emitterLogger) Warn(event *types.Event, factoryParams ...any) {
	e.emitIfCan(event)
	if e.loggerFactoryFn != nil {
		e.loggerFactoryFn(factoryParams...).Warn(event.Message)
	}
}

// Error emits an event and logs at the Error level.
//
// Parameters:
//   - event The event to emit and log.
//   - factoryParams: The parameters to pass to the logger factory function.
func (e *emitterLogger) Error(event *types.Event, factoryParams ...any) {
	e.emitIfCan(event)
	if e.loggerFactoryFn != nil {
		e.loggerFactoryFn(factoryParams...).Error(event.Message)
	}
}

// Fatal emits an event and logs at the Fatal level.
//
// Parameters:
//   - event The event to emit and log.
//   - factoryParams: The parameters to pass to the logger factory function.
func (e *emitterLogger) Fatal(event *types.Event, factoryParams ...any) {
	e.emitIfCan(event)
	if e.loggerFactoryFn != nil {
		e.loggerFactoryFn(factoryParams...).Fatal(event.Message)
	}
}

// emitIfCan emits the event if the event emitter is not nil.
func (e *emitterLogger) emitIfCan(event *types.Event) {
	if e.eventEmitter != nil {
		e.eventEmitter.Emit(event)
	}
}
