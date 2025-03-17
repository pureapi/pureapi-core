package types

import (
	"net/http"

	"github.com/pureapi/pureapi-core/apierror"
)

// ErrorHandler handles apierror and maps them to appropriate HTTP responses.
type ErrorHandler interface {
	Handle(err error) (int, *apierror.APIError)
}

// EndpointHandler is an interface for handling endpoints.
type EndpointHandler[Input any] interface {
	Handle(w http.ResponseWriter, r *http.Request)
}
