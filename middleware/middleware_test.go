package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// makeMiddleware returns a Middleware that appends pre and post markers
// to the events slice.
func makeMiddleware(label string, events *[]string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			*events = append(*events, label+"-pre")
			next.ServeHTTP(w, r)
			*events = append(*events, label+"-post")
		})
	}
}

func TestApply(t *testing.T) {
	// Define test cases with middleware labels and expected event order.
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
			var mws Middlewares

			// Create middleware for each label.
			for _, label := range tc.labels {
				mws = append(mws, makeMiddleware(label, &events))
			}

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
				"expected %d events, got %d",
				len(tc.expected), len(events),
			)
			for i, exp := range tc.expected {
				assert.Equal(
					t, events[i], exp,
					"at index %d, expected %q, got %q",
					i, exp, events[i],
				)
			}
		})
	}
}
