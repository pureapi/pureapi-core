package middleware

import (
	"net/http"

	"github.com/pureapi/pureapi-core/middleware/types"
)

// Chain applies a sequence of middlewares to an http.Handler. During a request
// the middlewaress are applied in the order they are provided.
// The middlewares are applied so that the first middleware in the list becomes
// the outermost wrapper.
//
// Example:
//
//	Chain(finalHandler, m1, m2) yields m1(m2(finalHandler)).
//
// Parameters:
//   - h: The http.Handler to wrap.
//   - middlewares: The middlewares to apply.
//
// Returns:
//   - http.Handler: The wrapped http.Handler.
func Chain(h http.Handler, middlewares types.Middlewares) http.Handler {
	wrapped := h
	for i := len(middlewares) - 1; i >= 0; i-- {
		wrapped = middlewares[i](wrapped)
	}
	return wrapped
}
