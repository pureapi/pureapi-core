package dbquery

import "github.com/pureapi/pureapi-core/apierror"

// Commmon database errors.
var (
	ErrDuplicateEntry    = apierror.NewAPIError("DUPLICATE_ENTRY")
	ErrForeignConstraint = apierror.NewAPIError("FOREIGN_CONSTRAINT_ERROR")
	ErrNoRows            = apierror.NewAPIError("NO_ROWS")
)
