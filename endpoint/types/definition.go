package types

import (
	"net/http"

	"github.com/pureapi/pureapi-core/middleware/types"
)

// Endpoint represents an API endpoint with middlewares.
type Endpoint interface {
	URL() string
	Method() string
	Middlewares() types.Middlewares
	Handler() http.HandlerFunc
}
