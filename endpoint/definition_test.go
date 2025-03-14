package endpoint

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pureapi/pureapi-core/middleware"
	"github.com/pureapi/pureapi-core/stack"
	types "github.com/pureapi/pureapi-core/stack/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// dummyMiddleware is a simple middleware that adds a header.
func dummyMiddleware(header, value string) middleware.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add(header, value)
			next.ServeHTTP(w, r)
		})
	}
}

// dummyHandler is a simple HTTP handler that writes "ok".
func dummyHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok"))
}

// executeHandler is a helper that executes an http.HandlerFunc and returns the
// response body as string.
func executeHandler(h http.HandlerFunc) string {
	req := httptest.NewRequest("GET", "/test", nil)
	rec := httptest.NewRecorder()
	h(rec, req)
	return rec.Body.String()
}

func TestNewDefinition(t *testing.T) {
	stack := stack.NewStack()
	def := NewDefinition("/test", "GET", stack, dummyHandler)
	assert.Equal(t, "/test", def.url)
	assert.Equal(t, "GET", def.method)
	assert.Equal(t, stack, def.stack)
	require.NotNil(t, def.handler)
	// Verify the handler behavior.
	output := executeHandler(def.handler)
	assert.Equal(t, "ok", output)
}

func TestDefinitionClone(t *testing.T) {
	// Create a definition with a non-nil stack and handler.
	st := stack.NewStack(
		stack.NewWrapper("mw1", dummyMiddleware("X-Test", "value")),
	)
	def := NewDefinition("/test", "GET", st, dummyHandler)

	// Clone without options.
	clone := def.Clone()
	assert.Equal(t, def.url, clone.url)
	assert.Equal(t, def.method, clone.method)
	// Verify handler behavior is preserved.
	origOutput := executeHandler(def.handler)
	cloneOutput := executeHandler(clone.handler)
	assert.Equal(t, origOutput, cloneOutput)

	// Ensure deep copy: the stack pointers should differ but the content equal.
	require.NotNil(t, def.stack)
	require.NotNil(t, clone.stack)
	assert.False(t, def.stack == clone.stack, "stack pointers should differ")
	require.Equal(t, len(def.stack.Wrappers()), len(clone.stack.Wrappers()))
	for i, w := range def.stack.Wrappers() {
		cloneWrapper := clone.stack.Wrappers()[i]
		assert.Equal(
			t, w.ID(), cloneWrapper.ID(),
			"Wrapper ID mismatch at index %d", i,
		)
		assert.Equal(
			t, w.Data(), cloneWrapper.Data(),
			"Wrapper Data mismatch at index %d", i,
		)
		assert.NotNil(
			t, w.Middleware(),
			"Original Middleware at index %d is nil", i,
		)
		assert.NotNil(
			t, cloneWrapper.Middleware(),
			"Cloned Middleware at index %d is nil", i,
		)
	}

	// Capture the length of the clone's stack before modifying the original.
	cloneLength := len(clone.stack.Wrappers())
	// Now modify the original stack.
	def.stack.AddWrapper(
		stack.NewWrapper("mw2", dummyMiddleware("X-Extra", "orig")),
	)
	// The original definition's stack now has more wrappers.
	assert.Equal(t, cloneLength+1, len(def.stack.Wrappers()))
	// The previously created clone should remain unchanged.
	assert.Equal(t, cloneLength, len(clone.stack.Wrappers()))

	// Clone with an option that changes the URL.
	clone2 := def.Clone(WithURL("/new"))
	assert.Equal(t, "/new", clone2.url)
}

func TestWithURLOption(t *testing.T) {
	def := NewDefinition("/old", "GET", nil, nil)
	opt := WithURL("")
	opt(def)
	// If empty, URL should default to "/".
	assert.Equal(t, "/", def.url)

	opt2 := WithURL("/new")
	opt2(def)
	assert.Equal(t, "/new", def.url)
}

func TestWithMethodOption(t *testing.T) {
	def := NewDefinition("/old", "GET", nil, nil)
	WithMethod("POST")(def)
	assert.Equal(t, "POST", def.method)
}

func TestWithMiddlewareStackOption(t *testing.T) {
	stack := stack.NewStack(stack.NewWrapper(
		"mw-stack", dummyMiddleware("X-Stack", "stack"),
	))
	def := NewDefinition("/test", "GET", nil, nil)
	WithMiddlewareStack(stack)(def)
	assert.Equal(t, stack, def.stack)
}

func TestWithMiddlewareWrappersFuncOption(t *testing.T) {
	def := NewDefinition("/test", "GET", nil, nil)
	f := func(d *Definition) types.Stack {
		// Return a new stack with one wrapper.
		return stack.NewStack(stack.NewWrapper(
			"mw-func", dummyMiddleware("X-Func", "func"),
		))
	}
	WithMiddlewareWrappersFunc(f)(def)
	require.NotNil(t, def.stack)
	assert.Len(t, def.stack.Wrappers(), 1)
	assert.Equal(t, "mw-func", def.stack.Wrappers()[0].ID())
}

func TestDefinitions_WithAndToEndpoints(t *testing.T) {
	// Create two definitions with stacks containing dummy middleware.
	stack1 := stack.NewStack(
		stack.NewWrapper("mw1", dummyMiddleware("X-One", "one")),
	)
	def1 := NewDefinition("/one", "GET", stack1, dummyHandler)
	stack2 := stack.NewStack(
		stack.NewWrapper("mw2", dummyMiddleware("X-Two", "two")),
	)
	def2 := NewDefinition("/two", "POST", stack2, dummyHandler)

	defs := Definitions{}.With(*def1)
	defs = defs.With(*def2)
	assert.Len(t, defs, 2)

	// Convert Definitions to core endpoints.
	endpoints := defs.ToEndpoints()
	assert.Len(t, endpoints, 2)
	// Verify that URL, method, and handler are correctly transferred.
	assert.Equal(t, "/one", endpoints[0].URL)
	assert.Equal(t, "GET", endpoints[0].Method)
	require.NotNil(t, endpoints[0].Handler)

	// To test middleware application, wrap a final handler.
	finalHandler := http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("final"))
		},
	)
	// Apply middleware from the first endpoint's stack.
	wrapped := endpoints[0].Middlewares.Chain(finalHandler)
	req := httptest.NewRequest("GET", "/one", nil)
	rec := httptest.NewRecorder()
	wrapped.ServeHTTP(rec, req)
	// Verify header added by the middleware.
	assert.Equal(t, "one", rec.Header().Get("X-One"))
	// Verify final response.
	assert.Equal(t, "final", rec.Body.String())
}
