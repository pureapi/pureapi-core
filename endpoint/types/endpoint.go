package types

import (
	"net/http"
)

// Endpoint represents an API endpoint with middlewares.
type Definition interface {
	URL() string
	Method() string
	Stack() Stack
	Handler() http.HandlerFunc
}

// Definitions is a new list of endpoint definitions.
type Definitions interface {
	ToEndpoints() []Endpoint
}
