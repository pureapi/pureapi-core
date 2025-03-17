package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pureapi/pureapi-core/middleware/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// dummyMiddleware is a simple middleware that adds a header.
func dummyMiddleware(header, value string) types.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add(header, value)
			next.ServeHTTP(w, r)
		})
	}
}

func TestNewStack(t *testing.T) {
	// Create initial wrappers.
	w1 := NewWrapper("mw1", dummyMiddleware("X-One", "one"))
	w2 := NewWrapper("mw2", dummyMiddleware("X-Two", "two"))
	stack := NewStack(w1, w2)
	require.NotNil(t, stack)
	assert.Len(t, stack.wrappers, 2)
}

func TestBuild(t *testing.T) {
	w1 := NewWrapper("mw1", dummyMiddleware("X-Test", "value"))
	stack := NewStack(w1)
	mws := stack.Middlewares()
	require.Len(t, mws, 1)
	// Test middleware behavior: create a handler that writes "ok", wrap it,
	// and verify that the header "X-Test" is added.
	final := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})
	wrapped := Chain(final, mws)
	req := httptest.NewRequest("GET", "/test", nil)
	rec := httptest.NewRecorder()
	wrapped.ServeHTTP(rec, req)
	assert.Equal(t, "value", rec.Header().Get("X-Test"))
	assert.Equal(t, "ok", rec.Body.String())
}

func TestStackClone(t *testing.T) {
	// Create a stack with one wrapper.
	orig := NewStack(
		NewWrapper("mw1", dummyMiddleware("X-Test", "value")),
	)
	clone := orig.Clone()
	// Ensure the clone is a deep copy.
	assert.False(t, orig == clone, "stack pointers should differ")
	require.Equal(t, len(orig.wrappers), len(clone.Wrappers()))
	for i := range orig.wrappers {
		// Compare the ID and Data fields.
		assert.Equal(t, orig.wrappers[i].ID(), clone.Wrappers()[i].ID())
		assert.Equal(t, orig.wrappers[i].Data(), clone.Wrappers()[i].Data())
		// We cannot compare the Middleware functions directly.
		assert.NotNil(t, orig.wrappers[i].Middleware())
		assert.NotNil(t, clone.Wrappers()[i].Middleware())
	}

	// Modify the original stack and verify the clone remains unchanged.
	cloneLength := len(clone.Wrappers())
	orig.AddWrapper(NewWrapper("mw2", dummyMiddleware("X-Extra", "extra")))
	assert.Equal(t, cloneLength+1, len(orig.wrappers))
	// The clone remains with the same number of wrappers.
	assert.Equal(t, cloneLength, len(clone.Wrappers()))
}

func TestAdd(t *testing.T) {
	stack := NewStack()
	assert.Len(t, stack.wrappers, 0)
	stack.AddWrapper(
		NewWrapper("mw-add", dummyMiddleware("X-Add", "added")),
	)
	assert.Len(t, stack.wrappers, 1)
	assert.Equal(t, "mw-add", stack.wrappers[0].ID())
}

func TestInsertBefore(t *testing.T) {
	// Start with a stack with one wrapper.
	stack := NewStack(NewWrapper("mw1", dummyMiddleware("X-1", "one")))
	// Insert a new wrapper before the one with ID "mw1".
	newWrapper := NewWrapper(
		"mw-before", dummyMiddleware("X-Before", "before"),
	)
	updated, found := stack.InsertBefore("mw1", newWrapper)
	assert.True(t, found)
	require.Len(t, updated.Wrappers(), 2)
	// The new wrapper should be at index 0.
	assert.Equal(t, "mw-before", updated.Wrappers()[0].ID())

	// Try inserting before a non-existent ID; should append.
	newWrapper2 := NewWrapper("mw-new", dummyMiddleware("X-New", "new"))
	updated, found = stack.InsertBefore("non-existent", newWrapper2)
	assert.False(t, found)
	assert.Equal(t, 3, len(updated.Wrappers()))
	assert.Equal(t, "mw-new", updated.Wrappers()[2].ID())
}

func TestInsertAfter(t *testing.T) {
	// Start with a stack with one wrapper.
	stack := NewStack(NewWrapper("mw1", dummyMiddleware("X-1", "one")))
	// Insert a new wrapper after the one with ID "mw1".
	newWrapper := NewWrapper(
		"mw-after", dummyMiddleware("X-After", "after"),
	)
	updated, found := stack.InsertAfter("mw1", newWrapper)
	assert.True(t, found)
	require.Len(t, updated.Wrappers(), 2)
	// The new wrapper should be at index 1.
	assert.Equal(t, "mw-after", updated.Wrappers()[1].ID())

	// Try inserting after a non-existent ID; should append.
	newWrapper2 := NewWrapper("mw-new", dummyMiddleware("X-New", "new"))
	updated, found = stack.InsertAfter("non-existent", newWrapper2)
	assert.False(t, found)
	assert.Equal(t, 3, len(updated.Wrappers()))
	assert.Equal(t, "mw-new", updated.Wrappers()[2].ID())
}

func TestRemove(t *testing.T) {
	// Create a stack with two wrappers.
	stack := NewStack(
		NewWrapper("mw1", dummyMiddleware("X-1", "one")),
		NewWrapper("mw2", dummyMiddleware("X-2", "two")),
	)
	// Remove the wrapper with ID "mw1".
	updated, found := stack.Remove("mw1")
	assert.True(t, found)
	assert.Len(t, updated.Wrappers(), 1)
	assert.Equal(t, "mw2", updated.Wrappers()[0].ID())

	// Try to remove a non-existent wrapper.
	updated, found = stack.Remove("non-existent")
	assert.False(t, found)
	assert.Len(t, updated.Wrappers(), 1)
}
