package types

import "context"

// HTTPServer represents an HTTP server.
type HTTPServer interface {
	ListenAndServe() error              // Start the server.
	Shutdown(ctx context.Context) error // Shut down the server.
}
