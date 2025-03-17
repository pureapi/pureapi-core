package endpoint

import (
	"net/http"
	"net/http/httptest"
	"testing"

	middlewaretypes "github.com/pureapi/pureapi-core/middleware/types"
	"github.com/stretchr/testify/suite"
)

// dummyStack is a minimal implementation of middlewaretypes.Stack
// used for testing.
type dummyStack struct {
	id string
}

func (ds *dummyStack) Wrappers() []middlewaretypes.Wrapper { return nil }
func (ds *dummyStack) Middlewares() middlewaretypes.Middlewares {
	return nil
}
func (ds *dummyStack) Clone() middlewaretypes.Stack {
	return &dummyStack{id: ds.id + "_clone"}
}
func (ds *dummyStack) AddWrapper(w middlewaretypes.Wrapper) middlewaretypes.Stack {
	return nil
}
func (ds *dummyStack) InsertBefore(id string, w middlewaretypes.Wrapper) (middlewaretypes.Stack, bool) {
	return nil, false
}
func (ds *dummyStack) InsertAfter(id string, w middlewaretypes.Wrapper) (middlewaretypes.Stack, bool) {
	return nil, false
}
func (ds *dummyStack) Remove(id string) (middlewaretypes.Stack, bool) { return nil, false }

// --------------------
// DefinitionTestSuite
// --------------------

// DefinitionTestSuite is a test suite for the Definition type.
type DefinitionTestSuite struct {
	suite.Suite
}

// TestDefinitionsTestSuite runs the test suite.
func TestDefinitionsTestSuite(t *testing.T) {
	suite.Run(t, new(DefinitionsTestSuite))
}

// Test_NewDefinition tests that NewDefinition returns a valid definition.
func (s *DefinitionTestSuite) Test_NewDefinition() {
	ds := &dummyStack{id: "test"}
	dummyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})
	url := "/api/test"
	method := http.MethodGet
	def := NewDefinition(url, method, ds, dummyHandler)
	s.Require().NotNil(def)
	s.Equal(url, def.URL())
	s.Equal(method, def.Method())
	s.Equal(ds, def.Stack())

	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", url, nil)
	def.Handler()(rr, req)
	s.Equal("ok", rr.Body.String())
}

// Test_Clone tests that Clone returns a deep copy of the definition.
func (s *DefinitionTestSuite) Test_Clone() {
	ds := &dummyStack{id: "cloneTest"}
	dummyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("original"))
	})
	def := NewDefinition("/test", "POST", ds, dummyHandler)
	clone := def.Clone()

	s.NotSame(def, clone)
	s.Equal(def.URL(), clone.URL())
	s.Equal(def.Method(), clone.Method())

	// Instead of comparing handler functions directly (which is invalid),
	// invoke both handlers and compare their output.
	rr1 := httptest.NewRecorder()
	rr2 := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	def.Handler()(rr1, req)
	clone.Handler()(rr2, req)
	s.Equal(rr1.Body.String(), rr2.Body.String(), "cloned handler should produce same output")

	// Ensure that the middleware stack was cloned.
	s.NotEqual(def.Stack(), clone.Stack())
	clonedStack, ok := clone.Stack().(*dummyStack)
	s.True(ok)
	s.Equal(ds.id+"_clone", clonedStack.id)
}

// Test_WithURL tests that WithURL returns a new definition with the provided
// URL.
func (s *DefinitionTestSuite) Test_WithURL() {
	// When empty URL is provided, defaultURL should return "/".
	def := NewDefinition("", "GET", nil, nil)
	s.Equal("/", def.URL())

	newDef := def.WithURL("/new")
	s.NotSame(def, newDef)
	s.Equal("/new", newDef.URL())
	// Original remains unchanged.
	s.Equal("/", def.URL())

	// Test with an empty string.
	newDef2 := def.WithURL("")
	s.Equal("/", newDef2.URL())
}

// Test_WithMethod tests that WithMethod returns a new definition with the
// provided method.
func (s *DefinitionTestSuite) Test_WithMethod() {
	def := NewDefinition("/path", "GET", nil, nil)
	newDef := def.WithMethod("POST")
	s.NotSame(def, newDef)
	s.Equal("POST", newDef.Method())
	s.Equal("GET", def.Method())
}

// Test_WithHandler tests that WithHandler returns a new definition with the
// provided handler.
func (s *DefinitionTestSuite) Test_WithHandler() {
	handler1 := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("first"))
	})
	handler2 := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("second"))
	})
	def := NewDefinition("/path", "GET", nil, handler1)
	newDef := def.WithHandler(handler2)
	s.NotSame(def, newDef)

	rr1 := httptest.NewRecorder()
	def.Handler()(rr1, httptest.NewRequest("GET", "/path", nil))
	s.Equal("first", rr1.Body.String())

	rr2 := httptest.NewRecorder()
	newDef.Handler()(rr2, httptest.NewRequest("GET", "/path", nil))
	s.Equal("second", rr2.Body.String())
}

// Test_WithMiddlewareStack tests that WithMiddlewareStack returns a new
// definition with the provided middleware stack.
func (s *DefinitionTestSuite) Test_WithMiddlewareStack() {
	ds1 := &dummyStack{id: "stack1"}
	ds2 := &dummyStack{id: "stack2"}
	def := NewDefinition("/path", "GET", ds1, nil)
	newDef := def.WithMiddlewareStack(ds2)
	s.NotSame(def, newDef)
	s.Equal(ds2, newDef.Stack())
	s.Equal(ds1, def.Stack())
}
