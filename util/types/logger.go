package types

import (
	"context"
)

// ILogger represents a logger with different logging levels.
// DebugLogger defines debug and trace level logging.
type DebugLogger interface {
	Debug(messages ...any)
	Debugf(message string, params ...any)
	Trace(messages ...any)
	Tracef(message string, params ...any)
}

// InfoLogger defines info and warning level logging.
type InfoLogger interface {
	Info(messages ...any)
	Infof(message string, params ...any)
	Warn(messages ...any)
	Warnf(message string, params ...any)
}

// ErrorLogger defines error and fatal level logging.
type ErrorLogger interface {
	Error(messages ...any)
	Errorf(message string, params ...any)
	Fatal(messages ...any)
	Fatalf(message string, params ...any)
}

// ILogger combines all logging levels.
type ILogger interface {
	DebugLogger
	InfoLogger
	ErrorLogger
}

// LoggerFactoryFn is a function that returns a logger.
type LoggerFactoryFn func(params ...any) ILogger

// CtxLoggerFactoryFn is a function that returns a logger with context.
type CtxLoggerFactoryFn func(ctx context.Context) ILogger
