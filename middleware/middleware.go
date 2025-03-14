package middleware

import "net/http"

// Middleware represents a function that wraps an http.Handler with additional
// behavior. A Middleware typically performs actions before and/or after calling
// the next handler.
//
// Example:
//
//	finalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//	    // final processing
//	})
//	wrappedHandler := Apply(finalHandler, middleware1, middleware2)
type Middleware func(http.Handler) http.Handler

// Middlewares is a slice of Middleware functions.
type Middlewares []Middleware

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
//
// Returns:
//   - http.Handler: The wrapped http.Handler.
func (m Middlewares) Chain(h http.Handler) http.Handler {
	wrapped := h
	for i := len(m) - 1; i >= 0; i-- {
		wrapped = m[i](wrapped)
	}
	return wrapped
}
