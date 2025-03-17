package util

import (
	"fmt"
	"testing"

	"github.com/pureapi/pureapi-core/util/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// FakeEventEmitter is a dummy event emitter that records emitted events.
type FakeEventEmitter struct {
	EmittedEvents []*types.Event
}

func (f *FakeEventEmitter) RegisterListener(
	eventType types.EventType, callback types.EventCallback,
) types.EventEmitter {
	return f
}

func (f *FakeEventEmitter) RemoveListener(
	eventType types.EventType, id string,
) {
}

func (f *FakeEventEmitter) Emit(event *types.Event) {
	f.EmittedEvents = append(f.EmittedEvents, event)
}

// FakeLogger is a dummy logger that records which method was called and with
// what message.
type FakeLogger struct {
	LastCalledMethod string
	LastMessage      string
}

func (f *FakeLogger) Debug(messages ...any) {
	f.LastCalledMethod = "Debug"
	f.LastMessage = fmt.Sprint(messages...)
}
func (f *FakeLogger) Debugf(message string, params ...any) {
	f.LastCalledMethod = "Debug"
	f.LastMessage = fmt.Sprintf(message, params...)
}
func (f *FakeLogger) Trace(messages ...any) {
	f.LastCalledMethod = "Trace"
	f.LastMessage = fmt.Sprint(messages...)
}
func (f *FakeLogger) Tracef(message string, params ...any) {
	f.LastCalledMethod = "Trace"
	f.LastMessage = fmt.Sprintf(message, params...)
}
func (f *FakeLogger) Info(messages ...any) {
	f.LastCalledMethod = "Info"
	f.LastMessage = fmt.Sprint(messages...)
}
func (f *FakeLogger) Infof(message string, params ...any) {
	f.LastCalledMethod = "Info"
	f.LastMessage = fmt.Sprintf(message, params...)
}
func (f *FakeLogger) Warn(messages ...any) {
	f.LastCalledMethod = "Warn"
	f.LastMessage = fmt.Sprint(messages...)
}
func (f *FakeLogger) Warnf(message string, params ...any) {
	f.LastCalledMethod = "Warn"
	f.LastMessage = fmt.Sprintf(message, params...)
}
func (f *FakeLogger) Error(messages ...any) {
	f.LastCalledMethod = "Error"
	f.LastMessage = fmt.Sprint(messages...)
}
func (f *FakeLogger) Errorf(message string, params ...any) {
	f.LastCalledMethod = "Error"
	f.LastMessage = fmt.Sprintf(message, params...)
}
func (f *FakeLogger) Fatal(messages ...any) {
	f.LastCalledMethod = "Fatal"
	f.LastMessage = fmt.Sprint(messages...)
}
func (f *FakeLogger) Fatalf(message string, params ...any) {
	f.LastCalledMethod = "Fatal"
	f.LastMessage = fmt.Sprintf(message, params...)
}

// EmitterLoggerTestSuite is a test suite for emitterLogger.
type EmitterLoggerTestSuite struct {
	suite.Suite
	fakeEmitter *FakeEventEmitter
	fakeLogger  *FakeLogger
}

// TestEmitterLoggerTestSuite runs the test suite.
func TestEmitterLoggerTestSuite(t *testing.T) {
	suite.Run(t, new(EmitterLoggerTestSuite))
}

// SetupTest initializes fake dependencies.
func (suite *EmitterLoggerTestSuite) SetupTest() {
	suite.fakeEmitter = &FakeEventEmitter{}
	suite.fakeLogger = &FakeLogger{}
}

// fakeLoggerFactory returns the fake logger stored in the suite.
func (suite *EmitterLoggerTestSuite) fakeLoggerFactory(
	params ...any,
) types.ILogger {
	return suite.fakeLogger
}

// testLogging is a helper to test a specific logging level.
func (suite *EmitterLoggerTestSuite) testLogging(
	level string,
	logFn func(el *emitterLogger, event *types.Event),
) {
	// Create a sample event.
	event := types.NewEvent(types.EventType(level), level+" message")
	// Create a new emitterLogger with both fake event emitter and fake logger.
	el := NewEmitterLogger(suite.fakeEmitter, suite.fakeLoggerFactory)
	// Call the logging method.
	logFn(el, event)
	// Verify the event was emitted.
	assert.Len(
		suite.T(), suite.fakeEmitter.EmittedEvents, 1,
		"Expected one event emitted",
	)
	assert.Equal(
		suite.T(), event, suite.fakeEmitter.EmittedEvents[0],
		"Emitted event should match",
	)
	// Verify that the logger was called with the proper method and message.
	assert.Equal(
		suite.T(), level, suite.fakeLogger.LastCalledMethod,
		"Logger method should match",
	)
	assert.Equal(
		suite.T(), event.Message, suite.fakeLogger.LastMessage,
		"Logger message should match",
	)
}

// TestDebug verifies the Debug method.
func (suite *EmitterLoggerTestSuite) TestDebug() {
	suite.testLogging("Debug", func(el *emitterLogger, event *types.Event) {
		el.Debug(event, "param1", "param2")
	})
}

// TestTrace verifies the Trace method.
func (suite *EmitterLoggerTestSuite) TestTrace() {
	suite.testLogging("Trace", func(el *emitterLogger, event *types.Event) {
		el.Trace(event)
	})
}

// TestInfo verifies the Info method.
func (suite *EmitterLoggerTestSuite) TestInfo() {
	suite.testLogging("Info", func(el *emitterLogger, event *types.Event) {
		el.Info(event)
	})
}

// TestWarn verifies the Warn method.
func (suite *EmitterLoggerTestSuite) TestWarn() {
	suite.testLogging("Warn", func(el *emitterLogger, event *types.Event) {
		el.Warn(event)
	})
}

// TestError verifies the Error method.
func (suite *EmitterLoggerTestSuite) TestError() {
	suite.testLogging("Error", func(el *emitterLogger, event *types.Event) {
		el.Error(event)
	})
}

// TestFatal verifies the Fatal method.
func (suite *EmitterLoggerTestSuite) TestFatal() {
	suite.testLogging("Fatal", func(el *emitterLogger, event *types.Event) {
		el.Fatal(event)
	})
}

// TestNilLoggerFactory verifies that if the loggerFactoryFn is nil, only event
// emission happens.
func (suite *EmitterLoggerTestSuite) TestNilLoggerFactory() {
	event := types.NewEvent("Info", "info message")
	el := NewEmitterLogger(suite.fakeEmitter, nil)
	el.Info(event)
	// Verify event is emitted.
	assert.Len(
		suite.T(), suite.fakeEmitter.EmittedEvents, 1,
		"Expected event to be emitted",
	)
	// No logger is called; fakeLogger remains unchanged.
	assert.Equal(
		suite.T(), "", suite.fakeLogger.LastCalledMethod,
		"Logger should not be called",
	)
}

// TestNilEventEmitter verifies that if the event emitter is nil, logging still
// occurs.
func (suite *EmitterLoggerTestSuite) TestNilEventEmitter() {
	event := types.NewEvent("Warn", "warn message")
	el := NewEmitterLogger(nil, suite.fakeLoggerFactory)
	el.Warn(event)
	// No event emitted.
	assert.Len(
		suite.T(), suite.fakeEmitter.EmittedEvents, 0,
		"No event should be emitted",
	)
	// Logger should still be called.
	assert.Equal(
		suite.T(), "Warn", suite.fakeLogger.LastCalledMethod,
		"Logger method should be Warn",
	)
	assert.Equal(
		suite.T(), event.Message, suite.fakeLogger.LastMessage,
		"Logger message should match",
	)
}

// TestNilEventEmitterAndLoggerFactory verifies that if both the event emitter
// and logger factory are nil, nothing happens.
func (suite *EmitterLoggerTestSuite) TestNilEventEmitterAndLoggerFactory() {
	el := NewEmitterLogger(nil, nil)
	event := types.NewEvent("Info", "info message")
	el.Info(event)
	assert.Len(
		suite.T(), suite.fakeEmitter.EmittedEvents, 0,
		"No event should be emitted",
	)
	assert.Equal(
		suite.T(), "", suite.fakeLogger.LastCalledMethod,
		"Logger should not be called",
	)
}

// TestNoopEmitterLogger verifies that NewNoopEmitterLogger does nothing.
func (suite *EmitterLoggerTestSuite) TestNoopEmitterLogger() {
	el := NewNoopEmitterLogger()
	event := types.NewEvent("Info", "noop message")
	// Calling methods should not panic or emit events.
	el.Info(event)
	el.Debug(event)
	el.Trace(event)
	el.Warn(event)
	el.Error(event)
	el.Fatal(event)
	// Nothing to assert since it's a no-op.
}
