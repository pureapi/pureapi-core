package types

import (
	"net/http"

	"github.com/pureapi/pureapi-core/util/types"
)

// Handler is an interface for handling endpoints.
type Handler[Input any] interface {
	Handle(w http.ResponseWriter, r *http.Request)
}

// ErrorHandler handles apierror and maps them to appropriate HTTP responses.
type ErrorHandler interface {
	Handle(err error) (int, types.APIError)
}
