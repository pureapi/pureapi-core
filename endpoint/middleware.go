package endpoint

import (
	"net/http"

	"github.com/pureapi/pureapi-core/endpoint/types"
)

// defaultMiddlewares is a slice of Middleware functions.
type defaultMiddlewares struct {
	middlewares []types.Middleware
}

// defaultMiddlewares implements the types.Middlewares interface.
var _ types.Middlewares = (*defaultMiddlewares)(nil)

// NewMiddlewares creates a new defaultMiddlewares instance with the provided
// middlewares.
//
// Parameters:
//   - middlewares: The middlewares to add to the list.
//
// Returns:
//   - *defaultMiddlewares: A new defaultMiddlewares instance.
func NewMiddlewares(middlewares ...types.Middleware) *defaultMiddlewares {
	return &defaultMiddlewares{
		middlewares: middlewares,
	}
}

// Middlewares returns the middlewares in the list.
//
// Returns:
//   - []types.Middleware: The middlewares in the list.
func (m defaultMiddlewares) Middlewares() []types.Middleware {
	return m.middlewares
}

// Chain applies a sequence of middlewares to an http.Handler. During a request
// the middlewaress are applied in the order they are provided.
// The middlewares are applied so that the first middleware in the list becomes
// the outermost wrapper.
//
// Example with middlewares m1, m2
//
//	Chain(finalHandler) yields m1(m2(finalHandler)).
//
// Parameters:
//   - h: The http.Handler to wrap.
//
// Returns:
//   - http.Handler: The wrapped http.Handler.
func (m defaultMiddlewares) Chain(h http.Handler) http.Handler {
	wrapped := h
	for i := len(m.middlewares) - 1; i >= 0; i-- {
		wrapped = m.middlewares[i](wrapped)
	}
	return wrapped
}

// Add adds one or more middlewares to the list and returns a new
// defaultMiddlewares instance.
//
// Parameters:
//   - middlewares: The middlewares to add to the list.
func (m defaultMiddlewares) Add(
	middlewares ...types.Middleware,
) *defaultMiddlewares {
	allMiddlewares := append([]types.Middleware{}, m.middlewares...)
	allMiddlewares = append(allMiddlewares, middlewares...)
	return NewMiddlewares(allMiddlewares...)
}
