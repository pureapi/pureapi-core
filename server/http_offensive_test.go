package server

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/pureapi/pureapi-core/endpoint"
	"github.com/pureapi/pureapi-core/endpoint/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// startOffensiveTestServer starts an HTTP server using our DefaultHTTPServer
// implementation with a simple GET / endpoint that returns "OK". It listens on
// an ephemeral port.
func startOffensiveTestServer(t *testing.T) (addr string, shutdown func()) {
	ep := endpoint.NewEndpoint("/", "GET").WithHandler(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("OK"))
		},
	)
	handler := NewHandler(nil)
	// Using port 0 to let the OS pick an available port.
	server := DefaultHTTPServer(handler, 0, []types.Endpoint{ep})
	server.ReadTimeout = 10 * time.Second
	server.WriteTimeout = 10 * time.Second
	server.IdleTimeout = 60 * time.Second
	server.MaxHeaderBytes = 1 << 16
	ln, err := net.Listen("tcp", server.Addr)
	require.NoError(t, err)
	go server.Serve(ln)
	return ln.Addr().String(), func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = server.Shutdown(ctx)
		_ = ln.Close()
	}
}

// TestExcessiveHeaderSize attempts to send a request with an extremely long
// header. Expectation: the server should respond with a 400 Bad Request
// (or close the connection) without crashing.
func TestExcessiveHeaderSize(t *testing.T) {
	addr, shutdown := startOffensiveTestServer(t)
	defer shutdown()

	conn, err := net.Dial("tcp", addr)
	require.NoError(t, err)
	defer conn.Close()

	// Create a header value larger than 64KB.
	hugeValue := strings.Repeat("A", 70*1024) // 70 KB
	request := fmt.Sprintf(
		"GET / HTTP/1.1\r\nHost: localhost\r\nX-Test: %s\r\n\r\n", hugeValue,
	)
	_, err = conn.Write([]byte(request))
	require.NoError(t, err)

	// Read response (the server should eventually close the connection or
	// return an error response)
	reader := bufio.NewReader(conn)
	line, err := reader.ReadString('\n')
	// In some cases the connection may be closed immediately.
	if err != nil && errors.Is(err, io.EOF) {
		t.Log("Connection closed by server as expected")
		return
	}
	require.NoError(t, err)
	// Check that the status line indicates a 400 Bad Request (or similar)
	// RFC 6585 defines 431 as the appropriate response on excessive headers.
	if !strings.Contains(line, "400") && !strings.Contains(line, "431") {
		t.Errorf(
			"expected status code '400' or '431' for huge header, got %s", line,
		)
	}

}

// TestMalformedRequest sends an incomplete/malformed HTTP request and ensures
// the server does not crash.
func TestMalformedRequest(t *testing.T) {
	addr, shutdown := startOffensiveTestServer(t)
	defer shutdown()

	conn, err := net.Dial("tcp", addr)
	require.NoError(t, err)
	defer conn.Close()

	// Write an incomplete HTTP request.
	_, err = conn.Write([]byte("BAD REQUEST\r\n\r\n"))
	require.NoError(t, err)

	reader := bufio.NewReader(conn)
	response, err := reader.ReadString('\n')
	// Even if the server closes the connection, that's acceptable.
	if err != nil && errors.Is(err, io.EOF) {
		t.Log("Connection closed after malformed request, as expected")
		return
	}
	require.NoError(t, err)
	// The response status line should indicate a client error (usually 400).
	assert.Contains(
		t, response, "400", "Expected 400 status for malformed request",
	)
}

// TestConcurrentRequests simulates many concurrent valid and invalid requests
// to ensure the server remains responsive and does not crash.
func TestConcurrentRequests(t *testing.T) {
	addr, shutdown := startOffensiveTestServer(t)
	defer shutdown()

	var wg sync.WaitGroup
	client := &http.Client{Timeout: 3 * time.Second}
	numRequests := 50

	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			// Alternate between valid GET and an invalid method.
			var req *http.Request
			var err error
			if i%2 == 0 {
				req, err = http.NewRequest("GET", "http://"+addr+"/", nil)
			} else {
				req, err = http.NewRequest("FOO", "http://"+addr+"/", nil)
			}
			require.NoError(t, err)
			resp, err := client.Do(req)
			if err != nil {
				t.Logf("Request %d error: %v", i, err)
				return
			}
			defer resp.Body.Close()
			if req.Method == "GET" {
				assert.Equal(t, http.StatusOK, resp.StatusCode)
			} else {
				// Our mux should respond with Method Not Allowed (405) for
				// unknown methods.
				assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
			}
		}(i)
	}
	wg.Wait()
}

// TestInvalidHTTPMethod sends a request with an invalid HTTP method to a
// defined endpoint and expects a 405 Method Not Allowed.
func TestInvalidHTTPMethod(t *testing.T) {
	addr, shutdown := startOffensiveTestServer(t)
	defer shutdown()

	req, err := http.NewRequest("INVALID", "http://"+addr+"/", nil)
	require.NoError(t, err)
	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
}

// TestSlowClient simulates a slow client that sends the request piecewise.
// The server's ReadTimeout should eventually trigger closing the connection.
func TestSlowClient(t *testing.T) {
	addr, shutdown := startOffensiveTestServer(t)
	defer shutdown()

	conn, err := net.Dial("tcp", addr)
	require.NoError(t, err)
	defer conn.Close()

	// Send request in parts with delays.
	parts := []string{
		"GET / HTTP/1.1\r\n",
		"Host: localhost\r\n",
		"User-Agent: slow-client\r\n",
		"\r\n",
	}

	for _, part := range parts {
		_, err := conn.Write([]byte(part))
		require.NoError(t, err)
		// Sleep for a duration shorter than the ReadTimeout but long enough to
		// simulate slowness.
		time.Sleep(500 * time.Millisecond)
	}

	// Read response; expect a valid response if completed in time,
	// or a timeout/connection reset if the delay exceeds the server's limits.
	reader := bufio.NewReader(conn)
	statusLine, err := reader.ReadString('\n')
	if err != nil {
		t.Logf(
			"Slow client test: connection error (expected if timed out): %v",
			err,
		)
		return
	}
	// If a response is received, it should be 200 OK.
	assert.Contains(
		t, statusLine, "200", "Expected 200 OK if request completes",
	)
}

// TestRandomGarbageRequest sends random binary data as a request and then
// checks that subsequent valid requests are handled properly.
func TestRandomGarbageRequest(t *testing.T) {
	addr, shutdown := startOffensiveTestServer(t)
	defer shutdown()

	// Send random data.
	conn, err := net.Dial("tcp", addr)
	require.NoError(t, err)
	randomData := make([]byte, 1024)
	_, err = rand.Read(randomData)
	require.NoError(t, err)
	_, _ = conn.Write(randomData)
	// Close the connection (server should handle this gracefully).
	_ = conn.Close()

	// Now send a valid request via http.Client to ensure the server is still
	// functional.
	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Get("http://" + addr + "/")
	require.NoError(t, err)
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.Equal(t, "OK", string(body))
}
