package types

// APIError represents a custom error type.
type APIError interface {
	Error() string
	ID() string
	Data() any
	Message() string
	Origin() *string
}
