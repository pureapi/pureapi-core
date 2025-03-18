package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/pureapi/pureapi-core/endpoint"
	"github.com/pureapi/pureapi-core/endpoint/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// DummyHTTPServer implements HTTPServer for testing purposes.
type DummyHTTPServer struct {
	ListenCalled      bool
	ShutdownCalled    bool
	ListenAndServeErr error
	ShutdownErr       error
	listenCh          chan struct{}
}

func NewDummyHTTPServer() *DummyHTTPServer {
	return &DummyHTTPServer{
		listenCh: make(chan struct{}),
	}
}

func (d *DummyHTTPServer) ListenAndServe() error {
	d.ListenCalled = true
	// If an error is expected, return it immediately.
	if d.ListenAndServeErr != nil {
		return d.ListenAndServeErr
	}
	<-d.listenCh // otherwise, block until shutdown triggers channel closing.
	return nil
}

func (d *DummyHTTPServer) Shutdown(ctx context.Context) error {
	d.ShutdownCalled = true
	// Only close the channel if it's not already closed.
	select {
	case <-d.listenCh:
		// already closed, do nothing.
	default:
		close(d.listenCh)
	}
	return d.ShutdownErr
}

func TestDefaultHTTPServer(t *testing.T) {
	// Create a basic server handler.
	handler := NewHandler(nil)
	endpoints := []types.Endpoint{
		endpoint.NewEndpoint("/test", "GET"),
	}
	port := 8080
	server := DefaultHTTPServer(handler, port, endpoints)
	assert.Equal(t, fmt.Sprintf(":%d", port), server.Addr)
	assert.Equal(t, 10*time.Second, server.ReadTimeout)
	assert.Equal(t, 10*time.Second, server.WriteTimeout)
	assert.Equal(t, 60*time.Second, server.IdleTimeout)
	assert.Equal(t, 1<<16, server.MaxHeaderBytes)
	assert.NotNil(t, server.Handler)
}

func TestStartServer_Normal(t *testing.T) {
	// Simulate normal shutdown:
	// Dummy server blocks in ListenAndServe until shutdown is triggered.
	dummyServer := NewDummyHTTPServer()
	dummyServer.ListenAndServeErr = nil // normal operation
	shutdownTimeout := 100 * time.Millisecond

	// Create a stopChan so we can control shutdown.
	stopChan := make(chan os.Signal, 1)
	handler := NewHandler(nil)

	// Start the server in a separate goroutine.
	errCh := make(chan error, 1)
	go func() {
		err := handler.startServer(stopChan, dummyServer, shutdownTimeout)
		errCh <- err
	}()

	// Give the server a moment to start.
	time.Sleep(50 * time.Millisecond)
	// Trigger shutdown by sending a signal.
	stopChan <- os.Interrupt

	// Wait for startServer to return.
	err := <-errCh
	assert.NoError(t, err)
	assert.True(t, dummyServer.ShutdownCalled)
}

func TestStartServer_ListenError(t *testing.T) {
	// Test the case when ListenAndServe returns an error.
	expectedErr := errors.New("listen error")
	dummyServer := NewDummyHTTPServer()
	dummyServer.ListenAndServeErr = expectedErr

	// Create a stopChan for the server handler.
	stopChan := make(chan os.Signal, 1)
	handler := NewHandler(nil)

	// In this branch, listenAndServe returns a non-nil error,
	// which triggers sending a signal on stopChan.
	err := handler.startServer(stopChan, dummyServer, 100*time.Millisecond)
	// startServer should return the expected listen error.
	assert.Equal(t, expectedErr, err)
	assert.True(t, dummyServer.ShutdownCalled)
}

func TestStartServer_ShutdownError(t *testing.T) {
	// Test scenario where Shutdown returns an error.
	shutdownErr := errors.New("shutdown failure")
	dummyServer := NewDummyHTTPServer()
	dummyServer.ListenAndServeErr = nil
	dummyServer.ShutdownErr = shutdownErr

	// Create a stopChan.
	stopChan := make(chan os.Signal, 1)
	handler := NewHandler(nil)

	// Send shutdown signal after a short delay.
	go func() {
		time.Sleep(50 * time.Millisecond)
		stopChan <- os.Interrupt
	}()

	err := handler.startServer(stopChan, dummyServer, 100*time.Millisecond)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "shutdown error")
}

func TestSetupMuxAndHandlers(t *testing.T) {
	// Create two endpoints:
	// - /test with GET returns "OK".
	// - /panic with GET panics to test panic recovery.
	endpointOK := endpoint.NewEndpoint("/test", "GET").WithHandler(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		},
	)
	endpointPanic := endpoint.NewEndpoint("/panic", "GET").WithHandler(
		func(w http.ResponseWriter, r *http.Request) {
			panic("handler panic")
		},
	)
	endpoints := []types.Endpoint{endpointOK, endpointPanic}

	// Use a dummy logger to capture events.
	handler := NewHandler(nil)
	mux := handler.setupMux(endpoints)
	require.NotNil(t, mux)

	// Test /test with allowed GET method.
	req := httptest.NewRequest("GET", "/test", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "OK", rr.Body.String())

	// Test /test with disallowed POST method returns 405.
	req = httptest.NewRequest("POST", "/test", nil)
	rr = httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)

	// Test /panic with GET; the panic should be recovered and result in a 500.
	req = httptest.NewRequest("GET", "/panic", nil)
	rr = httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	// Test an unregistered URL returns 404.
	req = httptest.NewRequest("GET", "/nonexistent", nil)
	rr = httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestMultiplexEndpoints(t *testing.T) {
	// Create endpoints with different URLs and methods.
	end1 := endpoint.NewEndpoint("/a", "GET").WithHandler(
		func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("A GET"))
		},
	)
	end2 := endpoint.NewEndpoint("/a", "POST").WithHandler(
		func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("A POST"))
		},
	)
	end3 := endpoint.NewEndpoint("/b", "GET").WithHandler(
		func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("B GET"))
		},
	)
	endpoints := []types.Endpoint{end1, end2, end3}

	handler := NewHandler(nil)
	muxed := handler.multiplexEndpoints(endpoints)

	// Check that the keys "/a" and "/b" exist.
	assert.Contains(t, muxed, "/a")
	assert.Contains(t, muxed, "/b")

	// Check that "/a" contains both GET and POST.
	methodsA := muxed["/a"]
	assert.Contains(t, methodsA, "GET")
	assert.Contains(t, methodsA, "POST")

	// Check that "/b" contains GET.
	methodsB := muxed["/b"]
	assert.Contains(t, methodsB, "GET")
}

func TestServerPanicHandler(t *testing.T) {
	// Create a handler that panics.
	panicHandler := http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			panic("test panic")
		},
	)
	handler := NewHandler(nil)
	wrapped := handler.serverPanicHandler(panicHandler)

	req := httptest.NewRequest("GET", "/panic", nil)
	rr := httptest.NewRecorder()
	wrapped.ServeHTTP(rr, req)

	// The panic should be recovered and return an internal server error.
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}
