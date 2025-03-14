package types

import "github.com/pureapi/pureapi-core/middleware"

// Wrapper is an interface for a middleware wrapper. It encapsulates a
// Middleware with an identifier and optional metadata.
type Wrapper interface {
	ID() string
	Middleware() middleware.Middleware
	Data() any
}

// Stack is an interface for managing a list of middleware wrappers.
type Stack interface {
	Wrappers() []Wrapper
	Middlewares() middleware.Middlewares
	Clone() Stack
	AddWrapper(w Wrapper) Stack
	InsertBefore(id string, w Wrapper) (Stack, bool)
	InsertAfter(id string, w Wrapper) (Stack, bool)
	Remove(id string) (Stack, bool)
}
