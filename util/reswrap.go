package util

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
)

type ResWrap interface {
	http.ResponseWriter
	StatusCode() int
	Header() http.Header
	WriteHeader(statusCode int)
	Body() []byte
	Write(data []byte) (int, error)
	Flush()
	Hijack() (net.Conn, *bufio.ReadWriter, error)
}

// defaultResWrap wraps an http.ResponseWriter to capture response data for
// inspection.
type defaultResWrap struct {
	http.ResponseWriter             // Embedded ResponseWriter.
	headers             http.Header // Captured headers.
	statusCode          int         // Captured status code.
	body                []byte      // Captured response body.
	headerWritten       bool        // Indicates if headers have been written.
}

// defaultResWrap implements the ResWrap interface.
var _ ResWrap = (*defaultResWrap)(nil)

// NewResWrap creates a new &defaultResWrap instance wrapping the given
// http.ResponseWriter.
//
// Parameters:
//   - w: The original http.ResponseWriter.
//
// Returns:
//   - *defaultResWrap: A new defaultResWrap instance.
func NewResWrap(w http.ResponseWriter) *defaultResWrap {
	return &defaultResWrap{
		ResponseWriter: w,
		headers:        make(http.Header),
		statusCode:     http.StatusOK,
		headerWritten:  false,
	}
}

// StatusCode returns the captured status code.
//
// Returns:
//   - int: The captured status code.
func (rw *defaultResWrap) StatusCode() int {
	return rw.statusCode
}

// Header overrides the Header method of the http.ResponseWriter interface.
// It returns the captured headers without modifying the underlying
// ResponseWriter's headers.
//
// Returns:
//   - The captured http.Header that can be modified before writing.
func (rw *defaultResWrap) Header() http.Header {
	return rw.headers
}

// WriteHeader captures the status code to be written, delaying its execution.
// It ensures that headers are only written once and applies the captured
// headers to the underlying ResponseWriter.
//
// Parameters:
//   - statusCode: The HTTP status code to write.
func (rw *defaultResWrap) WriteHeader(statusCode int) {
	if !rw.headerWritten { // Only write headers once
		rw.statusCode = statusCode
		// Apply the captured headers to the underlying ResponseWriter.
		for key, values := range rw.headers {
			for _, value := range values {
				rw.ResponseWriter.Header().Add(key, value)
			}
		}
		rw.ResponseWriter.WriteHeader(statusCode)
		rw.headerWritten = true
	}
}

// Body returns the captured response body.
//
// Returns:
//   - []byte: The captured response body.
func (rw *defaultResWrap) Body() []byte {
	return rw.body
}

// Write writes the response body and ensures headers and status code are
// written.
// It captures the response body and writes the data to the underlying
// ResponseWriter after ensuring that the headers and status code are written.
//
// Parameters:
//   - data: The response body data to write.
//
// Returns:
//   - int: The number of bytes written.
//   - error: An error if writing to the underlying ResponseWriter fails.
func (rw *defaultResWrap) Write(data []byte) (int, error) {
	// Ensure headers and status code are written before writing the body.
	if !rw.headerWritten {
		rw.WriteHeader(rw.statusCode)
	}

	// Append the data to the response body buffer.
	rw.body = append(rw.body, data...)

	// Write the data to the underlying ResponseWriter.
	return rw.ResponseWriter.Write(data)
}

// Flush forwards the flush call to the underlying ResponseWriter if supported.
func (rw *defaultResWrap) Flush() {
	if flusher, ok := rw.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}

// Hijack forwards the hijack call to the underlying ResponseWriter if
// supported.
//
// Returns:
//   - net.Conn: The hijacked connection.
//   - error: If the underlying ResponseWriter does not support hijacking.
func (rw *defaultResWrap) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if hijacker, ok := rw.ResponseWriter.(http.Hijacker); ok {
		return hijacker.Hijack()
	}
	return nil, nil, fmt.Errorf(
		"Hijack: underlying ResponseWriter does not support hijacking",
	)
}
