package types

// ErrorChecker translates database-specific errors into application errors.
type ErrorChecker interface {
	Check(err error) error
}
