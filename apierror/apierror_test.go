package apierror

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

// APIErrorTestSuite defines a test suite for APIError-related tests.
type APIErrorTestSuite struct {
	suite.Suite
}

// TestAPIErrorTestSuite runs the test suite.
func TestAPIErrorTestSuite(t *testing.T) {
	suite.Run(t, new(APIErrorTestSuite))
}

// Test_NewAPIError verifies that NewAPIError returns an APIError with the
// correct initial values.
func (s *APIErrorTestSuite) Test_NewAPIError() {
	id := "ERROR_001"
	errObj := NewAPIError(id)
	s.Require().NotNil(errObj)
	s.Equal(id, errObj.id)
	s.Nil(errObj.data)
	s.Empty(errObj.message)
	s.Equal("-", errObj.origin)
}

// Test_WithData verifies that WithData returns a new APIError with the data
// field set and other fields unchanged.
func (s *APIErrorTestSuite) Test_WithData() {
	base := NewAPIError("E001")
	testCases := []struct {
		name     string
		data     any
		expected any
	}{
		{"nil", nil, nil},
		{"string", "sample", "sample"},
		{"int", 123, 123},
		{"struct", struct{ A int }{A: 10}, struct{ A int }{A: 10}},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			newErr := base.WithData(tc.data)
			s.NotSame(base, newErr, "WithData should return a new instance")
			s.Equal(base.id, newErr.id)
			s.Equal(base.message, newErr.message)
			s.Equal(base.origin, newErr.origin)
			s.Equal(tc.expected, newErr.data)
		})
	}
}

// Test_WithMessage verifies that WithMessage returns a new APIError with the
// message pointer set properly.
func (s *APIErrorTestSuite) Test_WithMessage() {
	base := NewAPIError("E002")
	testCases := []struct {
		name     string
		message  string
		expected string
	}{
		{"empty", "", ""},
		{"nonempty", "An error occurred", "An error occurred"},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			newErr := base.WithMessage(tc.message)
			s.NotSame(base, newErr, "WithMessage should return a new instance")
			s.Equal(base.id, newErr.id)
			s.Equal(base.data, newErr.data)
			s.Equal(base.origin, newErr.origin)
			s.NotNil(newErr.message)
			s.Equal(tc.expected, newErr.message)
		})
	}
}

// Test_WithOrigin verifies that WithOrigin returns a new APIError with the
// origin field updated.
func (s *APIErrorTestSuite) Test_WithOrigin() {
	base := NewAPIError("E003")
	testCases := []struct {
		name     string
		origin   string
		expected string
	}{
		{"default", "-", "-"},
		{"custom", "serviceA", "serviceA"},
		{"empty", "", ""},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			newErr := base.WithOrigin(tc.origin)
			s.NotSame(base, newErr, "WithOrigin should return a new instance")
			s.Equal(base.id, newErr.ID())
			s.Equal(base.data, newErr.Data())
			s.Equal(tc.expected, newErr.Origin())
		})
	}
}

// Test_Error checks that the Error() method returns just the ID when no
// message is set, and "ID: message" when a message is provided.
func (s *APIErrorTestSuite) Test_Error() {
	base := NewAPIError("E004")
	s.Equal("E004", base.Error())

	msg := "Something went wrong"
	errWithMsg := base.WithMessage(msg)
	s.Equal("E004: "+msg, errWithMsg.Error())
}
