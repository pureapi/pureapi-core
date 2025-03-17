package middleware

import (
	"github.com/pureapi/pureapi-core/middleware/types"
)

// defaultWrapper encapsulates a middleware with an identifier and optional
// metadata. ID can be used to identify the middleware type (e.g. for reordering
// or documentation). Data can carry any type of additional information.
type defaultWrapper struct {
	id         string
	middleware types.Middleware
	data       any
}

// defaultWrapper implements the Wrapper interface.
var _ types.Wrapper = (*defaultWrapper)(nil)

// NewWrapper creates a new middleware defaultWrapper.
//
// Parameters:
//   - m: The middleware to wrap.
//   - id: The ID of the wrapper.
//
// Returns:
//   - *defaultWrapper: A new defaultWrapper instance.
func NewWrapper(
	id string, middleware types.Middleware,
) *defaultWrapper {
	defaultWrapper := &defaultWrapper{
		id:         id,
		middleware: middleware,
		data:       nil,
	}
	return defaultWrapper
}

// WithData returns a new defaultWrapper with the given data and returns a new
// defaultWrapper.
//
// Parameters:
//   - data: The data to attach to the wrapper.
//
// Returns:
//   - *defaultWrapper: A new defaultWrapper instance.
func (m *defaultWrapper) WithData(data any) *defaultWrapper {
	new := *m
	new.data = data
	return &new
}

// ID returns the ID of the wrapper.
//
// Returns:
//   - string: The ID of the wrapper.
func (m *defaultWrapper) Middleware() types.Middleware {
	return m.middleware
}

// ID returns the ID of the wrapper.
//
// Returns:
//   - string: The ID of the wrapper.
func (m *defaultWrapper) ID() string {
	return m.id
}

// Data returns the data attached to the wrapper.
//
// Returns:
//   - any: The data attached to the wrapper.
func (m *defaultWrapper) Data() any {
	return m.data
}
