package util

import (
	"bytes"
	"io"
	"net/http"
)

type ReqWrap interface {
	GetRequest() *http.Request
	GetBody() []byte
}

// defaultReqWrap wraps an http.Request, capturing its body for multiple reads
// and inspection.
type defaultReqWrap struct {
	*http.Request        // Embedded request.
	bodyContent   []byte // Captured request body.
}

// defaultReqWrap implements the ReqWrap interface.
var _ ReqWrap = (*defaultReqWrap)(nil)

// NewReqWrap creates a new defaultResWrap instance and captures the request
// body, enforcing a maximum size to prevent excessive memory usage. If the
// request body is larger than the maximum size, an error is returned.
//
// Parameters:
//   - r: The original http.Request.
//
// Returns:
//   - *defaultReqWrap: A new defaultReqWrap instance.
//   - error: Any error encountered during body reading.
func NewReqWrap(
	r *http.Request, maxRequestBodySize int64,
) (*defaultReqWrap, error) {
	bodyBytes, err := io.ReadAll(io.LimitReader(r.Body, maxRequestBodySize))
	if err != nil {
		return nil, err
	}
	// Replace the body so it can be read again.
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	return &defaultReqWrap{
		Request:     r,
		bodyContent: bodyBytes,
	}, nil
}

// GetRequest returns the wrapped http.Request.
//
// Returns:
//   - *http.Request: The wrapped request.
func (r *defaultReqWrap) GetRequest() *http.Request {
	return r.Request
}

// GetBody returns the captured request body.
//
// Returns:
//   - []byte: The captured request body.
func (r *defaultReqWrap) GetBody() []byte {
	return r.bodyContent
}
