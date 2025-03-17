package types

import (
	"net/http"

	"github.com/pureapi/pureapi-core/middleware/types"
)

// Endpoint represents an API endpoint with middlewares.
type Definition interface {
	URL() string
	Method() string
	Stack() types.Stack
	Handler() http.HandlerFunc
}

// Definitions is a new list of endpoint definitions.
type Definitions interface {
	ToEndpoints() []Endpoint
}
