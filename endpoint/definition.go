package endpoint

import (
	"net/http"

	"github.com/pureapi/pureapi-core/middleware"
	types "github.com/pureapi/pureapi-core/stack/types"
)

// Definition represents an endpoint definition.
type Definition struct {
	url     string
	method  string
	stack   types.Stack
	handler http.HandlerFunc // Optional handler for the endpoint.
}

// NewDefinition creates a new endpoint definition.
//
// Parameters:
//   - url: The URL of the endpoint.
//   - method: The HTTP method of the endpoint.
//   - stack: The middleware stack for the endpoint.
//   - handler: The optional handler for the endpoint.
//
// Returns:
//   - *defaultDefinition: A new defaultDefinition instance.
func NewDefinition(
	url string,
	method string,
	stack types.Stack,
	handler http.HandlerFunc,
) *Definition {
	return &Definition{
		url:     url,
		method:  method,
		stack:   stack,
		handler: handler,
	}
}

// URL returns the URL of the endpoint.
//
// Returns:
//   - string: The URL of the endpoint.
func (d *Definition) URL() string {
	return d.url
}

// Method returns the HTTP method of the endpoint.
//
// Returns:
//   - string: The HTTP method of the endpoint.
func (d *Definition) Method() string {
	return d.method
}

// Stack returns the middleware stack of the endpoint.
//
// Returns:
//   - middleware.Stack: The middleware stack of the endpoint.
func (d *Definition) Stack() types.Stack {
	return d.stack
}

// Handler returns the handler of the endpoint.
//
// Returns:
//   - http.HandlerFunc: The handler of the endpoint.
func (d *Definition) Handler() http.HandlerFunc {
	return d.handler
}

// Option is a function that modifies a definition.
type Option func(*Definition)

// Clone creates a deep copy of an endpoint definition with options.
//
// Parameters:
//   - opts: Options to apply to the cloned definition.
//
// Returns:
//   - *Definition: the cloned definition.
func (d *Definition) Clone(opts ...Option) *Definition {
	cloned := *d
	if d.stack != nil {
		cloned.stack = d.stack.Clone()
	}
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt(&cloned)
	}
	return &cloned
}

// WithURL returns an option that sets the URL of the endpoint. If the URL is
// empty, it will be set to "/"
//
// Parameters:
//   - url: The URL of the endpoint.
//
// Returns:
//   - func(*Definition): a function that sets the URL of the endpoint.
func WithURL(url string) Option {
	return func(e *Definition) {
		if url == "" {
			e.url = "/"
		} else {
			e.url = url
		}
	}
}

// WithMethod returns an option that sets the method of the endpoint.
//
// Parameters:
//   - method: The method of the endpoint.
//
// Returns:
//   - func(*Definition): a function that sets the method of the endpoint.
func WithMethod(method string) Option {
	return func(e *Definition) {
		e.method = method
	}
}

// WithMiddlewareStack return an option that sets the middleware stack.
//
// Parameters:
//   - stack: The middleware stack.
//
// Returns:
//   - func(*Definition): a function that sets the middleware stack.
func WithMiddlewareStack(stack types.Stack) Option {
	return func(e *Definition) {
		e.stack = stack
	}
}

// WithMiddlewareWrappersFunc returns an option that sets the middleware stack.
//
// Parameters:
//   - middlewareWrappersFunc: A function that returns the middleware stack.
//
// Returns:
//   - func(*Definition): a function that sets the middleware stack.
func WithMiddlewareWrappersFunc(
	middlewareWrappersFunc func(
		definition *Definition,
	) types.Stack,
) Option {
	return func(e *Definition) {
		wrappers := middlewareWrappersFunc(e)
		e.stack = wrappers
	}
}

// Definitions is a new list of endpoint definitions.
type Definitions []Definition

// With returns an option that sets the endpoint definitions.
//
// Parameters:
//   - definitions: The new endpoint definitions.
//
// Returns:
//   - func(*Definition): a function that sets the endpoint definitions.
func (d Definitions) With(definitions ...Definition) Definitions {
	return append(d, definitions...)
}

// ToEndpoints converts a list of endpoint definitions to a list of API
// endpoints.
//
// Returns:
//   - []api.Endpoint: a list of API endpoints.
func (d Definitions) ToEndpoints() []Endpoint {
	endpoints := []Endpoint{}
	for _, definition := range d {
		middlewares := middleware.Middlewares{}
		if definition.Stack() != nil {
			for _, mw := range definition.Stack().Wrappers() {
				middlewares = append(middlewares, mw.Middleware())
			}
		}
		endpoints = append(
			endpoints,
			*NewEndpoint(definition.URL(), definition.Method(), middlewares).
				WithHandler(definition.Handler()),
		)
	}
	return endpoints
}
