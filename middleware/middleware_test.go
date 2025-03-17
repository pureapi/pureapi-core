package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pureapi/pureapi-core/middleware/types"
	"github.com/stretchr/testify/assert"
)

// makeMiddleware returns a Middleware that appends pre and post markers
// to the events slice.
func makeMiddleware(label string, events *[]string) types.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			*events = append(*events, label+"-pre")
			next.ServeHTTP(w, r)
			*events = append(*events, label+"-post")
		})
	}
}

// TestChain tests the Chain function.
func TestChain(t *testing.T) {
	type testCase struct {
		name     string
		labels   []string
		expected []string
	}

	tests := []testCase{
		{
			name:     "No middlewares",
			labels:   []string{},
			expected: []string{"final"},
		},
		{
			name:     "Single middleware",
			labels:   []string{"m1"},
			expected: []string{"m1-pre", "final", "m1-post"},
		},
		{
			name:   "Two middlewares",
			labels: []string{"m1", "m2"},
			// Applied as m1(m2(final))
			expected: []string{
				"m1-pre", "m2-pre", "final", "m2-post", "m1-post",
			},
		},
		{
			name:   "Three middlewares",
			labels: []string{"m1", "m2", "m3"},
			// Applied as m1(m2(m3(final)))
			expected: []string{
				"m1-pre", "m2-pre", "m3-pre", "final",
				"m3-post", "m2-post", "m1-post",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var events []string
			var mws defaultMiddlewares

			// Create middleware for each label.
			allMiddlewares := []types.Middleware{}
			for _, label := range tc.labels {
				allMiddlewares = append(
					allMiddlewares, (makeMiddleware(label, &events)),
				)
			}
			mws = *NewMiddlewares(allMiddlewares...)

			// Final handler that appends "final" to events.
			final := http.HandlerFunc(
				func(w http.ResponseWriter, r *http.Request) {
					events = append(events, "final")
				},
			)

			// Wrap final handler with middlewares.
			wrapped := mws.Chain(final)
			req := httptest.NewRequest("GET", "/", nil)
			rr := httptest.NewRecorder()
			wrapped.ServeHTTP(rr, req)

			assert.Equal(
				t, len(events), len(tc.expected),
				"expected %d events, got %d", len(tc.expected), len(events),
			)
			for i, exp := range tc.expected {
				assert.Equal(
					t, events[i], exp,
					"at index %d, expected %q, got %q", i, exp, events[i],
				)
			}
		})
	}
}

// TestAdd tests that the Add function creates a new instance
// combining the original and added middlewares, while leaving the
// original instance unchanged.
func TestAddx(t *testing.T) {
	var events []string

	// Create an original middleware instance with one middleware.
	mwOriginal := NewMiddlewares(makeMiddleware("m1", &events))
	// Create a new middleware to add.
	additionalMiddleware := makeMiddleware("m2", &events)

	// Create a new instance by adding the additional middleware.
	mwNew := mwOriginal.Add(additionalMiddleware)

	// Final handler that appends "final" to the shared events.
	final := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		events = append(events, "final")
	})

	// Wrap final handler with the new middleware chain.
	wrapped := mwNew.Chain(final)
	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()
	wrapped.ServeHTTP(rr, req)

	expected := []string{"m1-pre", "m2-pre", "final", "m2-post", "m1-post"}
	assert.Equal(t, expected, events,
		"Expected chain to be %v, but got %v", expected, events)
}
