package server

import (
	"bufio"
	"context"
	"io"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/pureapi/pureapi-core/endpoint"
	"github.com/stretchr/testify/require"
)

// FuzzHTTPRequest sends fuzzed HTTP request data to the test server.
// It will help uncover vulnerabilities, crashes, or unexpected behavior.
// Run with: go test -fuzz=FuzzHTTPRequest

func FuzzHTTPRequest(f *testing.F) {
	// startTestServer accepts testing.TB so *testing.F works.
	addr, shutdown := startFuzzTestServer(f)
	defer shutdown()

	// Provide some seed inputs.
	seeds := []string{
		"GET / HTTP/1.1\r\nHost: localhost\r\n\r\n",
		"GET /test HTTP/1.1\r\nHost: localhost\r\n\r\n",
		"FOO / HTTP/1.1\r\nHost: localhost\r\n\r\n",
		"GET / HTTP/1.1\r\nHost: localhost\r\nX-Fuzz: " +
			strings.Repeat("A", 10) + "\r\n\r\n",
	}
	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, input string) {
		// Open a TCP connection to the test server.
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			t.Skipf("Skipping: could not connect to server: %v", err)
		}
		defer conn.Close()

		// Set a write deadline to avoid hanging on slow writes.
		_ = conn.SetWriteDeadline(time.Now().Add(2 * time.Second))
		_, err = conn.Write([]byte(input))
		if err != nil {
			// Some fuzzed inputs may trigger a connection reset.
			t.Skipf("Skipping write error: %v", err)
		}

		// Set a read deadline.
		_ = conn.SetReadDeadline(time.Now().Add(2 * time.Second))
		reader := bufio.NewReader(conn)
		// Try reading a line (the status line, for example).
		_, err = reader.ReadString('\n')
		// It's fine if we get an EOF or read error, as long as the server
		// didn't crash or panic.
		if err != nil && err != io.EOF {
			t.Logf("Read error: %v", err)
		}
	})
}

// startFuzzTestServer is a helper that starts a test server using
// DefaultHTTPServer. It registers a simple GET "/" endpoint returning "OK".
func startFuzzTestServer(tb testing.TB) (addr string, shutdown func()) {
	ep := endpoint.NewEndpoint("/", "GET", nil).WithHandler(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("OK"))
		},
	)
	handler := NewHandler(nil)
	// Use port 0 to let the OS select an available port.
	server := DefaultHTTPServer(handler, 0, []endpoint.Endpoint{*ep})
	ln, err := net.Listen("tcp", server.Addr)
	require.NoError(tb, err)
	go server.Serve(ln)
	return ln.Addr().String(), func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = server.Shutdown(ctx)
		_ = ln.Close()
	}
}
