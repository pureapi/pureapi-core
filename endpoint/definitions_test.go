package endpoint

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/suite"
)

// DefinitionsTestSuite tests the collection of endpoint definitions.
type DefinitionsTestSuite struct {
	suite.Suite
}

// TestDefinitionTestSuite runs the test suite.
func TestDefinitionTestSuite(t *testing.T) {
	suite.Run(t, new(DefinitionTestSuite))
}

// Test_NewDefinitionsAndAdd tests the NewDefinitions and Add methods that
// the definitions are correctly created and added to the collection.
func (s *DefinitionsTestSuite) Test_NewDefinitionsAndAdd() {
	// Create two definitions with simple handlers.
	d1 := NewDefinition("/one", "GET", nil,
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("one"))
		}))
	d2 := NewDefinition("/two", "POST", nil,
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("two"))
		}))
	defs := NewDefinitions(d1)
	s.Equal(1, len(defs.definitions))

	newDefs := defs.Add(d2)
	s.Equal(2, len(newDefs.definitions))
	// Ensure original remains unchanged.
	s.Equal(1, len(defs.definitions))
}

// Test_ToEndpoints tests that ToEndpoints returns a slice of endpoints with
// the correct URL, method and handler.
func (s *DefinitionsTestSuite) Test_ToEndpoints() {
	// Create two definitions.
	d1 := NewDefinition("/one", "GET", nil,
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("one"))
		}))
	d2 := NewDefinition("/two", "POST", nil,
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("two"))
		}))
	defs := NewDefinitions(d1, d2)
	endpoints := defs.ToEndpoints()
	s.Equal(2, len(endpoints))

	// For each endpoint, verify URL, method and handler output.
	for i, ep := range endpoints {
		var expectedURL, expectedMethod, expectedBody string
		if i == 0 {
			expectedURL = "/one"
			expectedMethod = "GET"
			expectedBody = "one"
		} else {
			expectedURL = "/two"
			expectedMethod = "POST"
			expectedBody = "two"
		}
		// Assume that the endpoint implements the following methods.
		epAccessor, ok := ep.(interface {
			URL() string
			Method() string
			Handler() http.HandlerFunc
		})
		s.True(ok, "endpoint should implement URL, Method, and Handler")
		s.Equal(expectedURL, epAccessor.URL())
		s.Equal(expectedMethod, epAccessor.Method())

		rr := httptest.NewRecorder()
		epAccessor.Handler()(rr, httptest.NewRequest(expectedMethod,
			expectedURL, nil))
		s.Equal(expectedBody, rr.Body.String())
	}
}
