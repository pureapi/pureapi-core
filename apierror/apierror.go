package apierror

import (
	"encoding/json"
	"fmt"

	"github.com/pureapi/pureapi-core/apierror/types"
)

// DefaultAPIError represents a JSON marshalable custom error type.
type DefaultAPIError struct {
	id      string
	data    any
	message string
	origin  string
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
		id:      id,
		data:    nil,
		message: "",
		origin:  "-", // Set to prevent empty origin.
	}
}

// MarshalJSON implements custom JSON marshaling.
//
// Returns:
//   - []byte: The JSON representation of the error.
//   - error: An error if the marshaling fails.
func (e *DefaultAPIError) MarshalJSON() ([]byte, error) {
	// Create an anonymous struct with JSON tags.
	return json.Marshal(struct {
		ID      string `json:"id"`
		Data    any    `json:"data,omitempty"`
		Message string `json:"message,omitempty"`
		Origin  string `json:"origin,omitempty"`
	}{
		ID:      e.id,
		Data:    e.data,
		Message: e.message,
		Origin:  e.origin,
	})
}

// UnmarshalJSON implements custom JSON unmarshaling.
//
// Parameters:
//   - data: The JSON data to unmarshal.
//
// Returns:
//   - error: An error if the unmarshaling fails.
func (e *DefaultAPIError) UnmarshalJSON(data []byte) error {
	aux := struct {
		ID      string `json:"id"`
		Data    any    `json:"data,omitempty"`
		Message string `json:"message,omitempty"`
		Origin  string `json:"origin,omitempty"`
	}{}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	e.id = aux.ID
	e.data = aux.Data
	e.message = aux.Message
	e.origin = aux.Origin
	return nil
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
	new.data = data
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
	new.message = message
	return &new
}

// WithOrigin returns a new error with the given origin.
//
// Parameters:
//   - origin: The origin to include in the error.
//
// Returns:
//   - types.APIError: A new APIError.
func (e *DefaultAPIError) WithOrigin(origin string) types.APIError {
	new := *e
	new.origin = origin
	return &new
}

// Error returns the full error message as a string. If the error has a message,
// it returns the ID followed by the message. Otherwise, it returns just the ID.
//
// Returns:
//   - string: The full error message as a string.
func (e *DefaultAPIError) Error() string {
	if e.message != "" {
		return fmt.Sprintf("%s: %s", e.id, e.message)
	}
	return e.id
}

// ID returns the ID of the error.
//
// Returns:
//   - string: The ID of the error.
func (e *DefaultAPIError) ID() string {
	return e.id
}

// Data returns the data associated with the error.
//
// Returns:
//   - any: The data associated with the error.
func (e *DefaultAPIError) Data() any {
	return e.data
}

// Message returns the message associated with the error.
//
// Returns:
//   - string: The message associated with the error.
func (e *DefaultAPIError) Message() string {
	return e.message
}

// Origin returns the origin associated with the error.
//
// Returns:
//   - string: The origin associated with the error.
func (e *DefaultAPIError) Origin() string {
	return e.origin
}
