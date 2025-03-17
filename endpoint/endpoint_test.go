package endpoint

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pureapi/pureapi-core/middleware"
	"github.com/pureapi/pureapi-core/middleware/types"
	"github.com/stretchr/testify/assert"
)

func TestNewEndpoint(t *testing.T) {
	// A dummy middleware that just passes control to the next handler.
	dummyMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
		})
	}

	// Table-driven test cases for NewEndpoint.
	testCases := []struct {
		name            string
		url             string
		method          string
		middlewares     types.Middlewares
		middlewareCheck func(t *testing.T, ep *defaultEndpoint)
	}{
		{
			name:        "Empty values",
			url:         "",
			method:      "",
			middlewares: nil,
			middlewareCheck: func(t *testing.T, ep *defaultEndpoint) {
				assert.Equal(t, ep.Middlewares(), middleware.NewMiddlewares(),
					"Middlewares should be empty default middlewares",
				)
			},
		},
		{
			name:        "Single middleware",
			url:         "/test",
			method:      "GET",
			middlewares: middleware.NewMiddlewares(dummyMiddleware),
		},
		{
			name:   "Multiple middlewares",
			url:    "/api",
			method: "POST",
			middlewares: middleware.NewMiddlewares(
				dummyMiddleware,
				dummyMiddleware,
			),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ep := NewEndpoint(tc.url, tc.method).WithMiddlewares(tc.middlewares)
			assert.Equal(t, tc.url, ep.URL(), "URL should match")
			assert.Equal(t, tc.method, ep.Method(), "HTTP method should match")
			if tc.middlewareCheck != nil {
				tc.middlewareCheck(t, ep)
			} else {
				assert.Equal(
					t, tc.middlewares, ep.Middlewares(),
					"Middlewares should match",
				)
			}
			// The initial Handler should be nil.
			assert.Nil(t, ep.Handler(), "Handler should be nil by default")
		})
	}
}

func TestEndpointWithHandler(t *testing.T) {
	// Create an endpoint without a handler.
	ep := NewEndpoint("/handler-test", "GET")

	// Create a dummy handler that writes a fixed response.
	handlerCalled := false
	dummyHandler := func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("hello"))
	}

	// Use WithHandler to attach the dummy handler.
	newEp := ep.WithHandler(dummyHandler)

	// Check that the new endpoint has the handler attached.
	assert.NotNil(
		t, newEp.Handler(), "New endpoint should have a non-nil handler",
	)
	// Ensure that the original endpoint is unchanged.
	assert.Nil(
		t, ep.Handler(),
		"Original endpoint should remain unchanged (nil Handler)",
	)

	// Test that the new endpoint's handler behaves as expected.
	req := httptest.NewRequest("GET", "/handler-test", nil)
	w := httptest.NewRecorder()
	newEp.handler(w, req)
	res := w.Result()
	defer res.Body.Close()

	assert.True(t, handlerCalled, "Dummy handler should be called")
	assert.Equal(
		t, http.StatusOK, res.StatusCode,
		"Response status code should be 200 OK",
	)
}
