package types

import "net/http"

// OutputHandler processes and writes the endpoint response.
type OutputHandler interface {
	Handle(
		w http.ResponseWriter,
		r *http.Request,
		out any,
		outputError error,
		statusCode int,
	) error
}

// InputHandler defines how to process the request input.
type Handler[Input any] interface {
	Handle(w http.ResponseWriter, r *http.Request) (*Input, error)
}
