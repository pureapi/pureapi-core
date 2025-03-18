package examples

import (
	"fmt"

	"github.com/pureapi/pureapi-core/util/types"
)

// ExampleLogger is an example implementation of the ILogger interface.
type ExampleLogger struct{}

// NewExampleLogger creates a new ExampleLogger.
//
// Returns:
//   - *ExampleLogger: A new ExampleLogger.
func NewExampleLogger() *ExampleLogger {
	return &ExampleLogger{}
}

// LoggerFactoryFn is a function that returns a logger.
//
// Returns:
//   - func() types.ILogger: A function that returns a logger.
func LoggerFactoryFn() func(params ...any) types.ILogger {
	return func(params ...any) types.ILogger {
		return NewExampleLogger()
	}
}

// Debug logs at the Debug level.
//
// Parameters:
//   - messages The messages to log.
func (l *ExampleLogger) Debug(messages ...any) {
	fmt.Println("Debug:", fmt.Sprint(messages...))
}

// Debugf logs at the Debug level.
//
// Parameters:
//   - message The message to log.
//   - params The parameters to use in the message.
func (l *ExampleLogger) Debugf(message string, params ...any) {
	fmt.Printf("Debug: %s\n", fmt.Sprintf(message, params...))
}

// Trace logs at the Trace level.
//
// Parameters:
//   - messages The messages to log.
func (l *ExampleLogger) Trace(messages ...any) {
	fmt.Println("Trace:", fmt.Sprint(messages...))
}

// Tracef logs at the Trace level.
//
// Parameters:
//   - message The message to log.
//   - params The parameters to use in the message.
func (l *ExampleLogger) Tracef(message string, params ...any) {
	fmt.Printf("Trace: %s\n", fmt.Sprintf(message, params...))
}

// Info logs at the Info level.
//
// Parameters:
//   - messages The messages to log.
func (l *ExampleLogger) Info(messages ...any) {
	fmt.Println("Info:", fmt.Sprint(messages...))
}

// Infof logs at the Info level.
//
// Parameters:
//   - message The message to log.
//   - params The parameters to use in the message.
func (l *ExampleLogger) Infof(message string, params ...any) {
	fmt.Printf("Info: %s\n", fmt.Sprintf(message, params...))
}

// Warn logs at the Warn level.
//
// Parameters:
//   - messages The messages to log.
func (l *ExampleLogger) Warn(messages ...any) {
	fmt.Println("Warn:", fmt.Sprint(messages...))
}

// Warnf logs at the Warn level.
//
// Parameters:
//   - message The message to log.
//   - params The parameters to use in the message.
func (l *ExampleLogger) Warnf(message string, params ...any) {
	fmt.Printf("Warn: %s\n", fmt.Sprintf(message, params...))
}

// Error logs at the Error level.
//
// Parameters:
//   - messages The messages to log.
func (l *ExampleLogger) Error(messages ...any) {
	fmt.Println("Error:", fmt.Sprint(messages...))
}

// Errorf logs at the Error level.
//
// Parameters:
//   - message The message to log.
//   - params The parameters to use in the message.
func (l *ExampleLogger) Errorf(message string, params ...any) {
	fmt.Printf("Error: %s\n", fmt.Sprintf(message, params...))
}

// Fatal logs at the Fatal level.
//
// Parameters:
//   - messages The messages to log.
func (l *ExampleLogger) Fatal(messages ...any) {
	fmt.Println("Fatal:", fmt.Sprint(messages...))
}

// Fatalf logs at the Fatal level.
//
// Parameters:
//   - message The message to log.
//   - params The parameters to use in the message.
func (l *ExampleLogger) Fatalf(message string, params ...any) {
	fmt.Printf("Fatal: %s\n", fmt.Sprintf(message, params...))
}
