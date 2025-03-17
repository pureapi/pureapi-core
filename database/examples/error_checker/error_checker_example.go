package main

import (
	"context"
	"errors"
	"log"

	_ "github.com/mattn/go-sqlite3"

	"github.com/pureapi/pureapi-core/apierror"
	apierrortypes "github.com/pureapi/pureapi-core/apierror/types"
	"github.com/pureapi/pureapi-core/database"
	dbexamples "github.com/pureapi/pureapi-core/database/examples"
	databasetypes "github.com/pureapi/pureapi-core/database/types"
)

// SimpleErrorChecker is a trivial custom error checker that returns a custom
// translated error from the original error.
type SimpleErrorChecker struct{}

// Check translates all database errors into a custom APIError.
//
// Parameters:
//   - err: The error to check.
//
// Returns:
//   - error: The translated error.
func (ec *SimpleErrorChecker) Check(err error) error {
	return apierror.NewAPIError("MY_API_ERROR").
		WithData(err.Error()).WithOrigin("my_api")
}

// This example demonstrates how the ErrorChecker can be used to handle database
// errors and translate them into customized application errors.
func main() {
	// Connect to the database.
	db, err := dbexamples.Connect(
		dbexamples.Cfg(), dbexamples.DummyConnectionOpen,
	)
	if err != nil {
		log.Fatalf("Connection failed: %v", err)
	}
	defer db.Close()

	// Attempt to run an invalid query.
	// Should see an error processed by the custom ErrorChecker.
	InvalidQuery(db)
}

// InvalidQuery attempts to run an invalid query. This demonstrates how the
// ErrorChecker can be used to handle database errors and translate them into
// customized application errors.
//
// Parameters:
//   - db: The database handle.
func InvalidQuery(db databasetypes.DB) {
	_, err := database.Exec(
		context.Background(),
		db,
		"INVALID_QUERY",
		nil,
		&SimpleErrorChecker{},
	)
	if err != nil {
		log.Printf("InvalidQuery error: %v", err)

		// Check if the error is an APIError and log its details.
		var apiErr apierrortypes.APIError
		if errors.As(err, &apiErr) {
			log.Printf(
				"APIError, ID: %v, data: %s, origin: %v",
				apiErr.ID(),
				apiErr.Data(),
				apiErr.Origin(),
			)
		}
	}
}
