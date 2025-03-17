package util

import (
	"bufio"
	"bytes"
	"errors"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// fakeResponseWriter is a custom ResponseWriter that supports Flush and Hijack.
type fakeResponseWriter struct {
	header          http.Header
	body            bytes.Buffer
	status          int
	flushed         bool
	hijackSupported bool
}

func newFakeResponseWriter(hijackSupported bool) *fakeResponseWriter {
	return &fakeResponseWriter{
		header:          make(http.Header),
		hijackSupported: hijackSupported,
	}
}

func (f *fakeResponseWriter) Header() http.Header {
	return f.header
}

func (f *fakeResponseWriter) Write(data []byte) (int, error) {
	return f.body.Write(data)
}

func (f *fakeResponseWriter) WriteHeader(statusCode int) {
	f.status = statusCode
}

func (f *fakeResponseWriter) Flush() {
	f.flushed = true
}

func (f *fakeResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if !f.hijackSupported {
		return nil, nil, errors.New("hijack not supported")
	}
	// For testing, return dummy values.
	return nil, bufio.NewReadWriter(nil, nil), nil
}

// ResWrapTestSuite is the test suite for the defaultResWrap.
type ResWrapTestSuite struct {
	suite.Suite
}

// TestResWrapTestSuite runs the test suite.
func TestResWrapTestSuite(t *testing.T) {
	suite.Run(t, new(ResWrapTestSuite))
}

// TestNewResWrapInitialState verifies that NewResWrap sets the default values.
func (suite *ResWrapTestSuite) TestNewResWrapInitialState() {
	underlying := httptest.NewRecorder()
	rw := NewResWrap(underlying)
	// Default status code should be 200.
	assert.Equal(suite.T(), http.StatusOK, rw.StatusCode())
	// The captured headers should be an empty map.
	assert.NotNil(suite.T(), rw.Header())
	assert.Equal(suite.T(), 0, len(rw.Header()))
	// No body has been written yet.
	assert.Equal(suite.T(), 0, len(rw.Body()))
}

// TestWriteHeader verifies that WriteHeader sets the status code and writes
// headers.
func (suite *ResWrapTestSuite) TestWriteHeader() {
	underlying := httptest.NewRecorder()
	rw := NewResWrap(underlying)
	// Set some headers on the reswrap.
	rw.Header().Set("X-Test", "value1")
	// Call WriteHeader.
	rw.WriteHeader(http.StatusCreated)
	// Check that the status code is updated.
	assert.Equal(suite.T(), http.StatusCreated, rw.StatusCode())
	// Underlying writer should have received the header.
	assert.Equal(suite.T(), "value1", underlying.Header().Get("X-Test"))
	// Calling WriteHeader again should not modify the headers/status.
	rw.Header().Set("X-Test", "value2")
	rw.WriteHeader(http.StatusBadRequest)
	assert.Equal(suite.T(), http.StatusCreated, rw.StatusCode())
	assert.Equal(suite.T(), "value1", underlying.Header().Get("X-Test"))
}

// TestWrite verifies that Write writes the data, captures it in the body, and
// auto-writes headers.
func (suite *ResWrapTestSuite) TestWrite() {
	underlying := httptest.NewRecorder()
	rw := NewResWrap(underlying)
	// Write some data without calling WriteHeader explicitly.
	data := []byte("Hello, World!")
	n, err := rw.Write(data)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), len(data), n)
	// The internal body should have captured the data.
	assert.Equal(suite.T(), data, rw.Body())
	// Underlying writer's body should also have the data.
	assert.Equal(suite.T(), data, underlying.Body.Bytes())
	// Headers should have been auto-written.
	assert.True(
		suite.T(), rw.headerWritten, "Headers should be written automatically",
	)
}

// TestMultipleWrites verifies that multiple writes append to the body.
func (suite *ResWrapTestSuite) TestMultipleWrites() {
	underlying := httptest.NewRecorder()
	rw := NewResWrap(underlying)
	data1 := []byte("Hello")
	data2 := []byte(", World!")
	_, err1 := rw.Write(data1)
	_, err2 := rw.Write(data2)
	assert.NoError(suite.T(), err1)
	assert.NoError(suite.T(), err2)
	expected := append(data1, data2...)
	assert.Equal(suite.T(), expected, rw.Body())
	assert.Equal(suite.T(), expected, underlying.Body.Bytes())
}

// TestFlush verifies that Flush() forwards the flush call when supported.
func (suite *ResWrapTestSuite) TestFlush() {
	// Create a fake underlying writer that supports Flush.
	fake := newFakeResponseWriter(true)
	rw := NewResWrap(fake)
	// Call Flush on our reswrap.
	rw.Flush()
	assert.True(
		suite.T(), fake.flushed, "Underlying writer should have been flushed",
	)
}

// TestHijackSuccess verifies that Hijack() is forwarded when supported.
func (suite *ResWrapTestSuite) TestHijackSuccess() {
	// Create a fake underlying writer that supports Hijack.
	fake := newFakeResponseWriter(true)
	rw := NewResWrap(fake)
	conn, rwBuf, err := rw.Hijack()
	// Our fake returns nil conn and a non-nil ReadWriter.
	assert.NoError(suite.T(), err)
	assert.Nil(suite.T(), conn)
	assert.NotNil(suite.T(), rwBuf)
}

// TestHijackFailure verifies that Hijack() returns an error when not supported.
func (suite *ResWrapTestSuite) TestHijackFailure() {
	// Use an underlying writer that does NOT support Hijack.
	// httptest.NewRecorder does not implement Hijacker.
	underlying := httptest.NewRecorder()
	rw := NewResWrap(underlying)
	conn, rwBuf, err := rw.Hijack()
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), conn)
	assert.Nil(suite.T(), rwBuf)
	assert.Contains(
		suite.T(), err.Error(),
		"underlying ResponseWriter does not support hijacking",
	)
}

// TestHeaderModification verifies that modifications to the reswrap header
// are isolated from the underlying writer until WriteHeader is called.
func (suite *ResWrapTestSuite) TestHeaderModification() {
	underlying := httptest.NewRecorder()
	rw := NewResWrap(underlying)
	// Initially, underlying header is empty.
	assert.Equal(suite.T(), 0, len(underlying.Header()))
	// Set a header in the reswrap.
	rw.Header().Set("X-Custom", "foo")
	// Underlying header should still be empty.
	assert.Equal(suite.T(), 0, len(underlying.Header()))
	// Now call WriteHeader.
	rw.WriteHeader(http.StatusAccepted)
	// Underlying header should now contain X-Custom.
	assert.Equal(suite.T(), "foo", underlying.Header().Get("X-Custom"))
}
