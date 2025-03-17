package apierror

import (
	"encoding/json"
	"fmt"

	"github.com/pureapi/pureapi-core/apierror/types"
)

// DefaultAPIError represents a JSON marshalable custom error type.
type DefaultAPIError struct {
	ErrID      string `json:"id"`
	ErrData    any    `json:"data,omitempty"`
	ErrMessage string `json:"message,omitempty"`
	ErrOrigin  string `json:"origin,omitempty"`
}

var _ types.APIError = (*DefaultAPIError)(nil)

// NewAPIError returns a new error with the given ID. The origin is set to "-"
// to control data leakage through empty origins.
//
// Parameters:
//   - id: The ID of the error.
//
// Returns:
//   - *defaultAPIError: A new defaultAPIError instance.
func NewAPIError(id string) *DefaultAPIError {
	return &DefaultAPIError{
		ErrID:      id,
		ErrData:    nil,
		ErrMessage: "",
		ErrOrigin:  "-", // Set to prevent empty origin.
	}
}

// From converts an APIError to a DefaultAPIError.
//
// Parameters:
//   - err: The APIError to convert.
//
// Returns:
//   - *defaultAPIError: A new defaultAPIError instance.
func From(err types.APIError) *DefaultAPIError {
	return &DefaultAPIError{
		ErrID:      err.ID(),
		ErrData:    err.Data(),
		ErrMessage: err.Message(),
		ErrOrigin:  err.Origin(),
	}
}

// Alias type to avoid infinite recursion in JSON marshaling.
type defaultAPIErrorAlias DefaultAPIError

// MarshalJSON implements custom JSON marshaling.
//
// Returns:
//   - []byte: The JSON representation of the error.
//   - error: An error if the marshaling fails.
func (e *DefaultAPIError) MarshalJSON() ([]byte, error) {
	alias := defaultAPIErrorAlias(*e)
	return json.Marshal(alias)
}

// UnmarshalJSON implements custom JSON unmarshaling.
//
// Parameters:
//   - data: The JSON data to unmarshal.
//
// Returns:
//   - error: An error if the unmarshaling fails.
func (e *DefaultAPIError) UnmarshalJSON(data []byte) error {
	var alias defaultAPIErrorAlias
	if err := json.Unmarshal(data, &alias); err != nil {
		return err
	}
	*e = DefaultAPIError(alias)
	return nil
}

// WithID returns a new error with the given ID.
//
// Parameters:
//   - id: The ID to include in the error.
//
// Returns:
//   - *defaultAPIError: A new defaultAPIError.
func (e *DefaultAPIError) WithID(id string) *DefaultAPIError {
	new := *e
	new.ErrID = id
	return &new
}

// WithData returns a new error with the given data.
//
// Parameters:
//   - data: The data to include in the error.
//
// Returns:
//   - *defaultAPIError: A new defaultAPIError.
func (e *DefaultAPIError) WithData(data any) *DefaultAPIError {
	new := *e
	new.ErrData = data
	return &new
}

// WithMessage returns a new error with the given message.
//
// Parameters:
//   - message: The message to include in the error.
//
// Returns:
//   - *defaultAPIError: A new defaultAPIError.
func (e *DefaultAPIError) WithMessage(message string) *DefaultAPIError {
	new := *e
	new.ErrMessage = message
	return &new
}

// WithMessage returns a new error with the given origin.
//
// Parameters:
//   - message: The message to include in the error.
//
// Returns:
//   - *defaultAPIError: A new defaultAPIError.
func (e *DefaultAPIError) WithOrigin(origin string) *DefaultAPIError {
	new := *e
	new.ErrOrigin = origin
	return &new
}

// Error returns the full error message as a string. If the error has a message,
// it returns the ID followed by the message. Otherwise, it returns just the ID.
//
// Returns:
//   - string: The full error message as a string.
func (e *DefaultAPIError) Error() string {
	if e.ErrMessage != "" {
		return fmt.Sprintf("%s: %s", e.ErrID, e.ErrMessage)
	}
	return e.ErrID
}

// ID returns the ID of the error.
//
// Returns:
//   - string: The ID of the error.
func (e *DefaultAPIError) ID() string {
	return e.ErrID
}

// Data returns the data associated with the error.
//
// Returns:
//   - any: The data associated with the error.
func (e *DefaultAPIError) Data() any {
	return e.ErrData
}

// Message returns the message associated with the error.
//
// Returns:
//   - string: The message associated with the error.
func (e *DefaultAPIError) Message() string {
	return e.ErrMessage
}

// Origin returns the origin associated with the error.
//
// Returns:
//   - string: The origin associated with the error.
func (e *DefaultAPIError) Origin() string {
	return e.ErrOrigin
}
