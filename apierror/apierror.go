package apierror

import (
	"fmt"
)

// APIError represents a JSON marshalable custom error type with an ID and
// other data.
type APIError struct {
	ID      string  `json:"id"`
	Data    any     `json:"data,omitempty"`
	Message *string `json:"message,omitempty"`
	Origin  string  `json:"origin,omitempty"` // Origin of the error.
}

// NewAPIError returns a new error with the given ID. The origin is set to "-"
// to prevent empty origins and data leakage.
//
// Parameters:
//   - id: The ID of the error.
//
// Returns:
//   - *APIError: A new APIError instance.
func NewAPIError(id string) *APIError {
	return &APIError{
		ID:      id,
		Data:    nil,
		Message: nil,
		Origin:  "-", // Set to prevent empty origin.
	}
}

// WithData returns a new error with the given data.
//
// Parameters:
//   - data: The data to include in the error.
//
// Returns:
//   - *APIError: A new APIError.
func (e *APIError) WithData(data any) *APIError {
	new := *e
	new.Data = data
	return &new
}

// WithMessage returns a new error with the given message.
//
// Parameters:
//   - message: The message to include in the error.
//
// Returns:
//   - *APIError: A new APIError.
func (e *APIError) WithMessage(message string) *APIError {
	new := *e
	new.Message = &message
	return &new
}

// WithOrigin returns a new error with the given origin.
//
// Parameters:
//   - origin: The origin to include in the error.
//
// Returns:
//   - *APIError: A new APIError.
func (e *APIError) WithOrigin(origin string) *APIError {
	new := *e
	new.Origin = origin
	return &new
}

// Error returns the full error message as a string. If the error has a message,
// it returns the ID followed by the message. Otherwise, it returns just the ID.
//
// Returns:
//   - string: The full error message as a string.
func (e *APIError) Error() string {
	if e.Message != nil {
		return fmt.Sprintf("%s: %s", e.ID, *e.Message)
	}
	return e.ID
}
