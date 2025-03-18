package endpoint

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pureapi/pureapi-core/endpoint/types"
	"github.com/stretchr/testify/suite"
)

// noopMiddleware is a middleware that does nothing.
func noopMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

// makeTestMiddleware returns a middleware that appends markers to events.
func makeTestMiddleware(label string, events *[]string) types.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			*events = append(*events, label+"-pre")
			next.ServeHTTP(w, r)
			*events = append(*events, label+"-post")
		})
	}
}

// StackTestSuite is the test suite for the defaultStack.
type StackTestSuite struct {
	suite.Suite
}

// TestStackTestSuite runs the test suite.
func TestStackTestSuite(t *testing.T) {
	suite.Run(t, new(StackTestSuite))
}

// TestNewStack verifies that NewStack initializes a stack with the provided
// wrappers.
func (s *StackTestSuite) TestNewStack() {
	// Create two wrappers.
	w1 := NewWrapper("w1", noopMiddleware)
	w2 := NewWrapper("w2", noopMiddleware)
	stack := NewStack(w1, w2)

	wrappers := stack.Wrappers()
	s.Require().Len(wrappers, 2)
	s.Equal("w1", wrappers[0].ID())
	s.Equal("w2", wrappers[1].ID())
}

// TestMiddlewares verifies that Middlewares returns the correct slice of
// middleware functions.
func (s *StackTestSuite) TestMiddlewares() {
	var events []string
	// Create wrappers with test middleware that logs events.
	w1 := NewWrapper("w1", makeTestMiddleware("w1", &events))
	w2 := NewWrapper("w2", makeTestMiddleware("w2", &events))
	stack := NewStack(w1, w2)

	mws := stack.Middlewares()
	// mws should be a Middlewares that wraps w1 and w2.
	// Chain them over a final handler that appends "final" to events.
	final := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		events = append(events, "final")
	})
	// Note: Chain applies middlewares in reverse order: first middleware
	// becomes outermost.
	wrapped := mws.Chain(final)
	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()
	wrapped.ServeHTTP(rr, req)
	// Expected execution:
	// Outer: w1-pre, inner: w2-pre, then final, then w2-post, then w1-post.
	expected := []string{"w1-pre", "w2-pre", "final", "w2-post", "w1-post"}
	s.Equal(expected, events)
}

// TestClone verifies that Clone returns a deep copy of the stack.
func (s *StackTestSuite) TestClone() {
	// Create a stack with two wrappers.
	w1 := NewWrapper("w1", noopMiddleware)
	w2 := NewWrapper("w2", noopMiddleware)
	orig := NewStack(w1, w2)
	clone := orig.Clone()

	// Check that the wrappers are equal in content.
	origWrappers := orig.Wrappers()
	cloneWrappers := clone.Wrappers()
	s.Require().Equal(len(origWrappers), len(cloneWrappers))
	for i := range origWrappers {
		s.Equal(origWrappers[i].ID(), cloneWrappers[i].ID())
	}
	// Verify that modifying the original's slice does not affect the clone.
	orig.AddWrapper(NewWrapper("w3", noopMiddleware))
	s.Len(orig.Wrappers(), len(cloneWrappers)+1)
	s.Len(clone.Wrappers(), len(cloneWrappers))
}

// TestAddWrapper verifies that AddWrapper appends a wrapper.
func (s *StackTestSuite) TestAddWrapper() {
	stack := NewStack()
	s.Len(stack.Wrappers(), 0)
	w := NewWrapper("add", noopMiddleware)
	stack = stack.AddWrapper(w).(*defaultStack)
	s.Len(stack.Wrappers(), 1)
	s.Equal("add", stack.Wrappers()[0].ID())
}

// TestInsertBefore verifies that InsertBefore inserts before a matching
// wrapper, and appends if not found.
func (s *StackTestSuite) TestInsertBefore() {
	// Create a stack with two wrappers.
	w1 := NewWrapper("w1", noopMiddleware)
	w2 := NewWrapper("w2", noopMiddleware)
	stack := NewStack(w1, w2)

	// Insert a new wrapper before "w1".
	newWrapper := NewWrapper("w0", noopMiddleware)
	updated, found := stack.InsertBefore("w1", newWrapper)
	s.True(found)
	s.Require().Len(updated.Wrappers(), 3)
	s.Equal("w0", updated.Wrappers()[0].ID())
	s.Equal("w1", updated.Wrappers()[1].ID())
	s.Equal("w2", updated.Wrappers()[2].ID())

	// Try inserting before a non-existent id; should append.
	newWrapper2 := NewWrapper("w3", noopMiddleware)
	updated, found = stack.InsertBefore("non-existent", newWrapper2)
	s.False(found)
	s.Equal("w3", updated.Wrappers()[len(updated.Wrappers())-1].ID())
}

// TestInsertAfter verifies that InsertAfter inserts after a matching wrapper,
// and appends if not found.
func (s *StackTestSuite) TestInsertAfter() {
	// Create a stack with two wrappers.
	w1 := NewWrapper("w1", noopMiddleware)
	w2 := NewWrapper("w2", noopMiddleware)
	stack := NewStack(w1, w2)

	// Insert a new wrapper after "w1".
	newWrapper := NewWrapper("w1.5", noopMiddleware)
	updated, found := stack.InsertAfter("w1", newWrapper)
	s.True(found)
	s.Require().Len(updated.Wrappers(), 3)
	s.Equal("w1", updated.Wrappers()[0].ID())
	s.Equal("w1.5", updated.Wrappers()[1].ID())
	s.Equal("w2", updated.Wrappers()[2].ID())

	// Try inserting after a non-existent id; should append.
	newWrapper2 := NewWrapper("w3", noopMiddleware)
	updated, found = stack.InsertAfter("non-existent", newWrapper2)
	s.False(found)
	s.Equal("w3", updated.Wrappers()[len(updated.Wrappers())-1].ID())
}

// TestRemove verifies that Remove deletes a wrapper with the given ID.
func (s *StackTestSuite) TestRemove() {
	// Create a stack with three wrappers.
	w1 := NewWrapper("w1", noopMiddleware)
	w2 := NewWrapper("w2", noopMiddleware)
	w3 := NewWrapper("w3", noopMiddleware)
	stack := NewStack(w1, w2, w3)

	// Remove "w2".
	updated, found := stack.Remove("w2")
	s.True(found)
	s.Require().Len(updated.Wrappers(), 2)
	s.Equal("w1", updated.Wrappers()[0].ID())
	s.Equal("w3", updated.Wrappers()[1].ID())

	// Attempt to remove a non-existent wrapper.
	updated, found = stack.Remove("non-existent")
	s.False(found)
	// The stack remains unchanged.
	s.Require().Len(updated.Wrappers(), 2)
}
