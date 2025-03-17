package util

import (
	"errors"
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// errorReader simulates an io.Reader that always returns an error.
type errorReader struct{}

func (e errorReader) Read(p []byte) (int, error) {
	return 0, errors.New("read error")
}

// ReqWrapTestSuite is a test suite for the defaultReqWrap functionality.
type ReqWrapTestSuite struct {
	suite.Suite
}

// TestReqWrapTestSuite runs the test suite.
func TestReqWrapTestSuite(t *testing.T) {
	suite.Run(t, new(ReqWrapTestSuite))
}

// TestGetBody verifies that NewReqWrap captures the body properly and that
// the request body can be re-read.
func (suite *ReqWrapTestSuite) TestGetBody() {
	body := "hello world"
	req := httptest.NewRequest(
		"POST",
		"http://example.com",
		strings.NewReader(body),
	)
	maxSize := int64(1024)
	rw, err := NewReqWrap(req, maxSize)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), body, string(rw.GetBody()))

	// Verify that the request body can be read again.
	data, err := io.ReadAll(req.Body)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), body, string(data))
}

// TestEmptyBody verifies that an empty request body is handled correctly.
func (suite *ReqWrapTestSuite) TestEmptyBody() {
	req := httptest.NewRequest("GET", "http://example.com", nil)
	maxSize := int64(1024)
	rw, err := NewReqWrap(req, maxSize)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "", string(rw.GetBody()))

	// The request body should also be empty.
	data, err := io.ReadAll(req.Body)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "", string(data))
}

// TestTruncatedBody verifies that NewReqWrap truncates the body when it
// exceeds maxRequestBodySize.
func (suite *ReqWrapTestSuite) TestTruncatedBody() {
	body := "abcdefghijklmnopqrstuvwxyz"
	maxSize := int64(10) // Force truncation.
	req := httptest.NewRequest(
		"POST", "http://example.com",
		strings.NewReader(body),
	)
	rw, err := NewReqWrap(req, maxSize)
	assert.NoError(suite.T(), err)
	expected := body[:10]
	assert.Equal(suite.T(), expected, string(rw.GetBody()))

	// Verify that the request body reflects the truncated content.
	data, err := io.ReadAll(req.Body)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expected, string(data))
}

// TestErrorDuringRead verifies that NewReqWrap returns an error if reading the
// request body fails.
func (suite *ReqWrapTestSuite) TestErrorDuringRead() {
	req := httptest.NewRequest(
		"POST",
		"http://example.com",
		io.NopCloser(errorReader{}),
	)
	maxSize := int64(1024)
	rw, err := NewReqWrap(req, maxSize)
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), rw)
}
