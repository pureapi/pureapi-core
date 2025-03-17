package middleware

import (
	"github.com/pureapi/pureapi-core/middleware/types"
)

// defaultWrapper encapsulates a middleware with an identifier and optional
// metadata. ID can be used to identify the middleware type (e.g. for reordering
// or documentation). Data can carry any type of additional information.
type defaultWrapper struct {
	middleware types.Middleware
	id         string
	data       any
}

// defaultWrapper implements the Wrapper interface.
var _ types.Wrapper = (*defaultWrapper)(nil)

// NewWrapper creates a new middleware defaultWrapper.
//
// Parameters:
//   - m: The middleware to wrap.
//   - id: The ID of the wrapper.
//   - options: Optional configuration functions.
//
// Returns:
//   - *defaultWrapper: A new defaultWrapper instance.
func NewWrapper(
	id string,
	middleware types.Middleware,
	opts ...Option,
) *defaultWrapper {
	w := &defaultWrapper{
		middleware: middleware,
		id:         id,
	}
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt(w)
	}
	return w
}

type Option func(*defaultWrapper)

// DataOption returns a function that sets the data for the wrapper.
//
// Parameters:
//   - data: The data to set.
//
// Returns:
//   - func(*defaultWrapper): A function that sets the data for the wrapper.
func DataOption(data any) Option {
	return func(w *defaultWrapper) {
		w.data = data
	}
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
