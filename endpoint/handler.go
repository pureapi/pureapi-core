package endpoint

import (
	"fmt"
	"net/http"

	endpointtypes "github.com/pureapi/pureapi-core/endpoint/types"
	"github.com/pureapi/pureapi-core/util"
	utiltypes "github.com/pureapi/pureapi-core/util/types"
)

// Constants for event types.
const (
	// EventError event is emitted when an error occurs during request
	// processing.
	EventError utiltypes.EventType = "event_error"

	// EventOutputError event is emitted when an output error occurs.
	EventOutputError utiltypes.EventType = "event_output_error"
)

// HandlerLogicFn is a function for handling endpoint logic.
type HandlerLogicFn[Input any] func(
	w http.ResponseWriter, r *http.Request, i *Input,
) (any, error)

// defaultHandler represents an endpoint with input, business logic, and
// output.
type defaultHandler[Input any] struct {
	inputHandler   endpointtypes.InputHandler[Input]
	handlerLogicFn HandlerLogicFn[Input]
	errorHandler   endpointtypes.ErrorHandler
	outputHandler  endpointtypes.OutputHandler
	emitterLogger  utiltypes.EmitterLogger
}

// NewHandler creates a new handler. During requst handling it
// executes common endpoints logic. It calls the input handler, handler
// logic, and output handler. If an error occurs during output handling, it
// will write a 500 error.
//
// Parameters:
//   - inputHandler: The input handler.
//   - handlerLogicFn: The handler logic function.
//   - errorHandler: The error handler.
//   - outputHandler: The output handler.
//
// Returns:
//   - *defaultHandler: A new defaultHandler instance.
func NewHandler[Input any](
	inputHandler endpointtypes.InputHandler[Input],
	handlerLogicFn HandlerLogicFn[Input],
	errorHandler endpointtypes.ErrorHandler,
	outputHandler endpointtypes.OutputHandler,
) *defaultHandler[Input] {
	return &defaultHandler[Input]{
		inputHandler:   inputHandler,
		handlerLogicFn: handlerLogicFn,
		errorHandler:   errorHandler,
		outputHandler:  outputHandler,
		emitterLogger:  defaultEmitterLogger(),
	}
}

// WithEmitterLogger adds an emitter logger to the handler.
//
// Parameters:
//   - emitterLogger: The emitter logger.
//
// Returns:
//   - *handler: A new handler instance.
func (h *defaultHandler[Input]) WithEmitterLogger(
	emitterLogger utiltypes.EmitterLogger,
) *defaultHandler[Input] {
	new := *h
	if new.emitterLogger == nil {
		new.emitterLogger = defaultEmitterLogger()
	} else {
		new.emitterLogger = emitterLogger
	}
	return &new
}

// Handle executes common endpoints logic. It calls the input handler, handler
// logic, and output handler.
//
// Parameters:
//   - w: The HTTP response writer.
//   - r: The HTTP request.
func (h *defaultHandler[Input]) Handle(
	w http.ResponseWriter, r *http.Request,
) {
	// Handle input.
	input, err := h.inputHandler.Handle(w, r)
	if err != nil {
		h.handleError(w, r, err)
		return
	}
	// Call handler logic.
	out, err := h.handlerLogicFn(w, r, input)
	if err != nil {
		h.handleError(w, r, err)
		return
	}
	// Write output.
	h.handleOutput(w, r, out, nil, http.StatusOK)
}

// handleError maps apierror and writes the error response.
func (h *defaultHandler[Input]) handleError(
	w http.ResponseWriter, r *http.Request, err error,
) {
	// Handle error.
	statusCode, outError := h.errorHandler.Handle(err)
	h.emitterLogger.Trace(
		utiltypes.NewEvent(
			EventError,
			fmt.Sprintf(
				"Error, status: %d, err: %s, out: %s",
				statusCode,
				err,
				outError,
			),
		).WithData(
			map[string]any{"status": statusCode, "err": err, "out": outError},
		),
		r.Context(),
	)
	// Handle and write output.
	h.handleOutput(w, r, nil, outError, statusCode)
}

// handleOutput processes and writes the endpoint response.
func (h *defaultHandler[Input]) handleOutput(
	w http.ResponseWriter, r *http.Request, out any, outError error, status int,
) {
	if err := h.outputHandler.Handle(
		w, r, out, outError, status,
	); err != nil {
		// If error handling output, write 500 error.
		h.emitterLogger.Trace(
			utiltypes.NewEvent(
				EventOutputError,
				fmt.Sprintf("Error handling output: %+v", err),
			).WithData(map[string]any{"err": err}),
			r.Context(),
		)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// defaultEmitterLogger returns a noop emitter logger.
func defaultEmitterLogger() utiltypes.EmitterLogger {
	return util.NewNoopEmitterLogger()
}
