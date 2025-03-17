package util

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/pureapi/pureapi-core/util/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// DummyLogger is a dummy logger that records which method was called and with
// what message.
type DummyLogger struct {
	LastCalledMethod string
	LastMessage      string
}

func (f *DummyLogger) Debug(messages ...any) {
	f.LastCalledMethod = "Debug"
	f.LastMessage = fmt.Sprint(messages...)
}
func (f *DummyLogger) Debugf(message string, params ...any) {
	f.LastCalledMethod = "Debug"
	f.LastMessage = fmt.Sprintf(message, params...)
}
func (f *DummyLogger) Trace(messages ...any) {
	f.LastCalledMethod = "Trace"
	f.LastMessage = fmt.Sprint(messages...)
}
func (f *DummyLogger) Tracef(message string, params ...any) {
	f.LastCalledMethod = "Trace"
	f.LastMessage = fmt.Sprintf(message, params...)
}
func (f *DummyLogger) Info(messages ...any) {
	f.LastCalledMethod = "Info"
	f.LastMessage = fmt.Sprint(messages...)
}
func (f *DummyLogger) Infof(message string, params ...any) {
	f.LastCalledMethod = "Info"
	f.LastMessage = fmt.Sprintf(message, params...)
}
func (f *DummyLogger) Warn(messages ...any) {
	f.LastCalledMethod = "Warn"
	f.LastMessage = fmt.Sprint(messages...)
}
func (f *DummyLogger) Warnf(message string, params ...any) {
	f.LastCalledMethod = "Warn"
	f.LastMessage = fmt.Sprintf(message, params...)
}
func (f *DummyLogger) Error(messages ...any) {
	f.LastCalledMethod = "Error"
	f.LastMessage = fmt.Sprint(messages...)
}
func (f *DummyLogger) Errorf(message string, params ...any) {
	f.LastCalledMethod = "Error"
	f.LastMessage = fmt.Sprintf(message, params...)
}
func (f *DummyLogger) Fatal(messages ...any) {
	f.LastCalledMethod = "Fatal"
	f.LastMessage = fmt.Sprint(messages...)
}
func (f *DummyLogger) Fatalf(message string, params ...any) {
	f.LastCalledMethod = "Fatal"
	f.LastMessage = fmt.Sprintf(message, params...)
}

// PanicHandlerTestSuite is a test suite for the PanicHandler.
type PanicHandlerTestSuite struct {
	suite.Suite
	dummyLogger *DummyLogger
}

// TestPanicHandlerTestSuite runs the test suite.
func TestPanicHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(PanicHandlerTestSuite))
}

// fakeCtxLoggerFactory is a fake context logger factory that creates a new
// FakeLogger
// and assigns it to the suite for later assertions.
func (suite *PanicHandlerTestSuite) fakeCtxLoggerFactory(
	ctx context.Context,
) types.ILogger {
	suite.dummyLogger = &DummyLogger{}
	return suite.dummyLogger
}

// fakeResponseWrapper returns a ResWrap wrapping an httptest.ResponseRecorder.
func (suite *PanicHandlerTestSuite) fakeResponseWrapper(
	r *http.Request,
) ResWrap {
	rr := httptest.NewRecorder()
	return NewResWrap(rr)
}

// customFakeResponseWrapper returns a ResWrap with a long header value.
func (suite *PanicHandlerTestSuite) customFakeResponseWrapper(
	r *http.Request,
) ResWrap {
	rr := httptest.NewRecorder()
	rw := NewResWrap(rr)
	// Set a header value longer than our max size.
	longRespHeader := strings.Repeat("y", 20)
	rw.Header().Set("X-Long-Resp", longRespHeader)
	// Write a dummy body so that it is non-empty.
	_, err := rw.Write([]byte("dummy"))
	if err != nil {
		panic(err)
	}
	return rw
}

// Test_HandlePanic uses table tests to check various request configurations.
func (suite *PanicHandlerTestSuite) Test_HandlePanic() {
	testCases := []struct {
		name               string
		requestURL         string
		requestBody        string
		requestHeader      http.Header
		useResponseWrapper bool
	}{
		{
			name:        "With response wrapper, valid body and headers",
			requestURL:  "http://example.com/test",
			requestBody: "test body",
			requestHeader: http.Header{
				"Content-Type": {"application/json"},
				"X-Test":       {"value1"},
			},
			useResponseWrapper: true,
		},
		{
			name:        "Without response wrapper",
			requestURL:  "http://example.com/nil",
			requestBody: "",
			requestHeader: http.Header{
				"Accept": {"*/*"},
			},
			useResponseWrapper: false,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			req := httptest.NewRequest(
				"GET", tc.requestURL, strings.NewReader(tc.requestBody),
			)
			req.Header = tc.requestHeader
			rr := httptest.NewRecorder()

			var respWrapFn func(*http.Request) ResWrap
			if tc.useResponseWrapper {
				respWrapFn = suite.fakeResponseWrapper
			} else {
				respWrapFn = func(r *http.Request) ResWrap { return nil }
			}

			// Create the panic handler with a max dump size of 1024 bytes.
			panicHandler := NewPanicHandler(
				suite.fakeCtxLoggerFactory, 1024, respWrapFn,
			)
			testErr := errors.New("test panic")
			panicHandler.HandlePanic(rr, req, testErr)

			// Verify the HTTP response.
			result := rr.Result()
			bodyBytes, _ := io.ReadAll(result.Body)
			assert.Equal(
				suite.T(), http.StatusInternalServerError, result.StatusCode,
				"HTTP status should be 500",
			)
			assert.Equal(
				suite.T(),
				http.StatusText(http.StatusInternalServerError),
				strings.TrimSpace(string(bodyBytes)),
				"Response body should match error text",
			)

			// Verify that the fake logger was created and called.
			assert.NotNil(
				suite.T(), suite.dummyLogger,
				"Logger should have been created",
			)
			assert.Equal(
				suite.T(), "Error", suite.dummyLogger.LastCalledMethod,
				"Logger method should be Error",
			)
			// Verify that the logged message contains the error and request URL.
			assert.Contains(
				suite.T(), suite.dummyLogger.LastMessage, "test panic",
				"Logged message should contain error message",
			)
			assert.Contains(
				suite.T(), suite.dummyLogger.LastMessage, tc.requestURL,
				"Logged message should contain request URL",
			)
		})
	}
}

// Test_BodyTruncation verifies that a long request body is truncated.
func (suite *PanicHandlerTestSuite) Test_BodyTruncation() {
	smallMaxSize := int64(10)
	longBody := "This is a very long request body that should be truncated."
	req := httptest.NewRequest(
		"POST", "http://example.com/truncate", strings.NewReader(longBody),
	)
	req.Header = http.Header{"Content-Length": {strconv.Itoa(len(longBody))}}
	rr := httptest.NewRecorder()

	panicHandler := NewPanicHandler(
		suite.fakeCtxLoggerFactory, smallMaxSize, suite.fakeResponseWrapper,
	)
	panicHandler.HandlePanic(rr, req, "truncate error")

	result := rr.Result()
	bodyBytes, _ := io.ReadAll(result.Body)
	assert.Equal(
		suite.T(), http.StatusInternalServerError, result.StatusCode,
	)
	assert.Equal(
		suite.T(), http.StatusText(http.StatusInternalServerError),
		strings.TrimSpace(string(bodyBytes)),
	)
	assert.NotNil(
		suite.T(), suite.dummyLogger, "Logger should have been created",
	)
	// Verify that the logged message contains the truncation notice.
	assert.Contains(
		suite.T(), suite.dummyLogger.LastMessage, "... (truncated)",
		"Request body should be truncated",
	)
}

// Test_QueryParameterTruncation verifies that long query parameters are
// truncated.
func (suite *PanicHandlerTestSuite) Test_QueryParameterTruncation() {
	smallMaxSize := int64(10)
	query := "abcdefghijklmnop"
	req := httptest.NewRequest("GET", "http://example.com/test?"+query, nil)
	req.Header = http.Header{}
	rr := httptest.NewRecorder()

	panicHandler := NewPanicHandler(
		suite.fakeCtxLoggerFactory, smallMaxSize, suite.fakeResponseWrapper,
	)
	panicHandler.HandlePanic(rr, req, "query param error")
	assert.NotNil(
		suite.T(), suite.dummyLogger, "Logger should have been created",
	)
	expected := query[:10] + "... (truncated)"
	assert.Contains(
		suite.T(), suite.dummyLogger.LastMessage, expected,
		"Query parameters should be truncated",
	)
}

// Test_RequestHeaderTruncation verifies that long request header values are
// truncated.
func (suite *PanicHandlerTestSuite) Test_RequestHeaderTruncation() {
	smallMaxSize := int64(10)
	longHeaderValue := strings.Repeat("x", 20)
	req := httptest.NewRequest("GET", "http://example.com/test", nil)
	req.Header = http.Header{"X-Long": {longHeaderValue}}
	rr := httptest.NewRecorder()

	panicHandler := NewPanicHandler(
		suite.fakeCtxLoggerFactory, smallMaxSize, suite.fakeResponseWrapper,
	)
	panicHandler.HandlePanic(rr, req, "header error")
	assert.NotNil(
		suite.T(), suite.dummyLogger, "Logger should have been created",
	)
	expected := longHeaderValue[:10] + "... (truncated)"
	assert.Contains(
		suite.T(), suite.dummyLogger.LastMessage, expected,
		"Request header should be truncated",
	)
}

// Test_ResponseHeaderTruncation verifies that long response header values are
// truncated.
func (suite *PanicHandlerTestSuite) Test_ResponseHeaderTruncation() {
	smallMaxSize := int64(10)
	req := httptest.NewRequest("GET", "http://example.com/test", nil)
	req.Header = http.Header{}
	rr := httptest.NewRecorder()

	panicHandler := NewPanicHandler(
		suite.fakeCtxLoggerFactory, smallMaxSize,
		suite.customFakeResponseWrapper,
	)
	panicHandler.HandlePanic(rr, req, "response header error")
	assert.NotNil(
		suite.T(), suite.dummyLogger, "Logger should have been created",
	)
	expected := strings.Repeat("y", 10) + "... (truncated)"
	assert.Contains(
		suite.T(), suite.dummyLogger.LastMessage, expected,
		"Response header should be truncated",
	)
}
