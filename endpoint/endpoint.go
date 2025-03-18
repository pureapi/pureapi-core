package endpoint

import (
	"net/http"

	"github.com/pureapi/pureapi-core/endpoint/types"
)

// defaultEndpoint represents an API endpoint with middlewares.
type defaultEndpoint struct {
	url         string
	method      string
	middlewares types.Middlewares
	handler     http.HandlerFunc // Optional handler for the endpoint.
}

// defaultEndpoint implements the Endpoint interface.
var _ types.Endpoint = (*defaultEndpoint)(nil)

// NewEndpoint creates a new defaultEndpoint with the given details.
//
// Parameters:
//   - url: The URL of the endpoint.
//   - method: The HTTP method of the endpoint.
//   - middlewares: The middlewares to apply to the endpoint.
//
// Returns:
//   - *defaultEndpoint: A new defaultEndpoint instance.
func NewEndpoint(url string, method string) *defaultEndpoint {
	return &defaultEndpoint{
		url:         url,
		method:      method,
		middlewares: nil,
		handler:     nil,
	}
}

// URL returns the URL of the endpoint.
//
// Returns:
//   - string: The URL of the endpoint.
func (e *defaultEndpoint) URL() string {
	return e.url
}

// Method returns the HTTP method of the endpoint.
//
// Returns:
//   - string: The HTTP method of the endpoint.
func (e *defaultEndpoint) Method() string {
	return e.method
}

// Middlewares returns the middlewares of the endpoint. If no middlewares are
// set, it returns an empty Middlewares instance.
//
// Returns:
//   - Middlewares: The middlewares of the endpoint.
func (e *defaultEndpoint) Middlewares() types.Middlewares {
	if e.middlewares == nil {
		return NewMiddlewares()
	}
	return e.middlewares
}

// Handler returns the handler of the endpoint.
//
// Returns:
//   - http.HandlerFunc: The handler of the endpoint.
func (e *defaultEndpoint) Handler() http.HandlerFunc {
	return e.handler
}

// WithMiddlewares sets the middlewares for the endpoint. It returns a new
// endpoint.
//
// Parameters:
//   - middlewares: The middlewares to apply to the endpoint.
//
// Returns:
//   - Endpoint: A new endpoint.
func (e *defaultEndpoint) WithMiddlewares(
	middlewares types.Middlewares,
) *defaultEndpoint {
	new := *e
	new.middlewares = middlewares
	return &new
}

// WithHandler sets the handler for the endpoint. It returns a new endpoint.
//
// Parameters:
//   - handler: The handler for the endpoint.
//
// Returns:
//   - Endpoint: A new endpoint.
func (e *defaultEndpoint) WithHandler(
	handler http.HandlerFunc,
) *defaultEndpoint {
	new := *e
	new.handler = handler
	return &new
}
