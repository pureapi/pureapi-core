package types

import (
	"net/http"
)

// Endpoint represents an API endpoint with middlewares.
type Endpoint interface {
	URL() string
	Method() string
	Middlewares() Middlewares
	Handler() http.HandlerFunc
}
