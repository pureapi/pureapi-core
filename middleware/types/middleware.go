package types

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

// Wrapper is an interface for a middleware wrapper. It encapsulates a
// Middleware with an identifier and optional metadata.
type Wrapper interface {
	ID() string
	Middleware() Middleware
	Data() any
}

// Stack is an interface for managing a list of middleware wrappers.
type Stack interface {
	Wrappers() []Wrapper
	Middlewares() Middlewares
	Clone() Stack
	AddWrapper(w Wrapper) Stack
	InsertBefore(id string, w Wrapper) (Stack, bool)
	InsertAfter(id string, w Wrapper) (Stack, bool)
	Remove(id string) (Stack, bool)
}
