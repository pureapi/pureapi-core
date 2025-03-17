package endpoint

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/pureapi/pureapi-core/apierror"
	"github.com/pureapi/pureapi-core/endpoint/types"
	"github.com/pureapi/pureapi-core/util"
	utiltypes "github.com/pureapi/pureapi-core/util/types"
)

// Constants for event apitypes.
const (
	// EventError event is emitted when an error occurs during request
	// processing.
	EventError utiltypes.EventType = "event_error"

	// EventOutputError event is emitted when an output error occurs.
	EventOutputError utiltypes.EventType = "event_output_error"
)

// InputFactoryFn returns a function that creates a new instance of the input.
type InputFactoryFn[Input any] func() Input

// HandlerLogicFn is a function for handling endpoint logic.
type HandlerLogicFn[Input any] func(
	w http.ResponseWriter, r *http.Request, i *Input,
) (any, error)

// endpointHandler represents an endpoint with input, business logic, and
// output.
type endpointHandler[Input any] struct {
	systemID       *string
	inputHandler   types.InputHandler[Input]
	inputFactoryFn InputFactoryFn[Input]
	handlerLogicFn HandlerLogicFn[Input]
	errorHandler   types.ErrorHandler
	outputHandler  types.OutputHandler
	emitterLogger  utiltypes.EmitterLogger
}

// NewEndpointHandler creates a new endpointHandler. During requst handling it
// executes common endpoints logic. It calls the input handler, handler
// logic, and output handler. Before calling the error handler it adds the
// system ID to any APIError instances passing through this handler. This can be
// useful for filtering errors based on the system ID in the error handler.
// If an error occurs during output handling, it will write a 500 error.
//
// Parameters:
//   - systemID: The optional system ID. It is used to add the system ID to any
//     APIError instances passing through this handler.
//   - inputHandler: The input handler.
//   - inputFactoryFn: The input factory function.
//   - handlerLogicFn: The handler logic function.
//   - errorHandler: The error handler.
//   - outputHandler: The output handler.
//   - emitterLogger: The emitter logger.
//
// Returns:
//   - *endpointHandler: The created endpointHandler.
func NewEndpointHandler[Input any](
	systemID *string,
	inputHandler types.InputHandler[Input],
	inputFactoryFn InputFactoryFn[Input],
	handlerLogicFn HandlerLogicFn[Input],
	errorHandler types.ErrorHandler,
	outputHandler types.OutputHandler,
	emitterLogger utiltypes.EmitterLogger,
) *endpointHandler[Input] {
	var useEmitterLogger utiltypes.EmitterLogger
	if emitterLogger == nil {
		useEmitterLogger = util.NewNoopEmitterLogger()
	} else {
		useEmitterLogger = emitterLogger
	}
	return &endpointHandler[Input]{
		systemID:       systemID,
		inputHandler:   inputHandler,
		inputFactoryFn: inputFactoryFn,
		handlerLogicFn: handlerLogicFn,
		errorHandler:   errorHandler,
		outputHandler:  outputHandler,
		emitterLogger:  useEmitterLogger,
	}
}

// Handle executes common endpoints logic. It calls the input handler, handler
// logic, and output handler.
//
// Parameters:
//   - w: The HTTP response writer.
//   - r: The HTTP request.
func (h *endpointHandler[Input]) Handle(
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
func (h *endpointHandler[Input]) handleError(
	w http.ResponseWriter, r *http.Request, err error,
) {
	// Add system ID to error if available.
	if h.systemID != nil {
		var apiError *apierror.APIError
		if ok := errors.As(err, &apiError); ok {
			err = apiError.WithOrigin(*h.systemID)
		}
	}
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
			map[string]any{"status": statusCode, "err": err, "out": outError},
		),
		r.Context(),
	)
	// Handle and write output.
	h.handleOutput(w, r, nil, outError, statusCode)
}

// handleOutput processes and writes the endpoint response.
func (h *endpointHandler[Input]) handleOutput(
	w http.ResponseWriter,
	r *http.Request,
	out any,
	outputError error,
	statusCode int,
) {
	if err := h.outputHandler.Handle(
		w, r, out, outputError, statusCode,
	); err != nil {
		h.emitterLogger.Trace(
			utiltypes.NewEvent(
				EventOutputError,
				fmt.Sprintf("Error handling output: %+v", err),
				map[string]any{"err": err},
			),
			r.Context(),
		)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
