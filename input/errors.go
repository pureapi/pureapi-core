package input

import "github.com/pureapi/pureapi-core/apierror"

// Commmon input errors.
var (
	ErrValidation    = apierror.NewAPIError("VALIDATION_ERROR")
	ErrInputDecoding = apierror.NewAPIError("ERROR_DECODING_INPUT")
	ErrInvalidInput  = apierror.NewAPIError("INVALID_INPUT")
)
