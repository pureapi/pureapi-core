package types

import (
	"net/http"

	"github.com/pureapi/pureapi-core/database/types"
)

// CreateHandler is the handler interface for the create endpoint.
type CreateHandler[Entity types.Mutator, Input any] interface {
	Handle(w http.ResponseWriter, r *http.Request, i *Input) (any, error)
}

// GetHandler is the handler interface for the get endpoint.
type GetHandler[Entity types.Getter, Input any, Output any] interface {
	Handle(w http.ResponseWriter, r *http.Request, i *Input) (any, error)
}

// UpdateHandler is the handler interface for the update endpoint.
type UpdateHandler[Input any] interface {
	Handle(w http.ResponseWriter, r *http.Request, i *Input) (any, error)
}

// DeleteHandler is the handler interface for the delete endpoint.
type DeleteHandler[Input any] interface {
	Handle(w http.ResponseWriter, r *http.Request, i *Input) (any, error)
}
