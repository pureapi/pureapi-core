package main

import (
	"context"
	"errors"
	"log"

	_ "github.com/mattn/go-sqlite3"

	"github.com/pureapi/pureapi-core/apierror"
	dbexamples "github.com/pureapi/pureapi-core/database/examples"
	"github.com/pureapi/pureapi-core/database/types"
	"github.com/pureapi/pureapi-core/repository"
	"github.com/pureapi/pureapi-core/repository/examples"
)

// This example demonstrates how the custom ErrorChecker can be used to handle
// database errors and translate them into customized application errors.
func main() {
	// Connect to the database.
	db, err := dbexamples.Connect(dbexamples.Cfg(), dbexamples.DummyConnectionOpen)
	if err != nil {
		log.Fatalf("Connection failed: %v", err)
	}
	defer db.Close()

	// Attempt to retrieve a user from a non-existent table.
	// Should see an error processed by the custom ErrorChecker.
	GetErrorUser(db)
}

// GetErrorUser attempts to retrieve a record using ErrorUser which maps to a
// non-existent table. This demonstrates how the custom ErrorChecker can be used
// to handle database errors and translate them into customized application
// errors.
//
// Parameters:
//   - db: The database handle.
func GetErrorUser(db types.DB) {
	readerRepo := repository.NewReaderRepo[*examples.ErrorUser](
		&examples.SimpleQueryBuilder{}, &examples.SimpleErrorChecker{},
	)
	_, err := readerRepo.GetOne(
		context.Background(),
		db,
		func() *examples.ErrorUser { return &examples.ErrorUser{} },
		nil,
	)
	if err != nil {
		log.Printf("GetErrorUser error: %v", err)
		// Check if the error is an APIError and log its details.
		var apiErr *apierror.APIError
		if errors.As(err, &apiErr) {
			log.Printf(
				"APIError, ID: %v, data: %s, origin: %v",
				apiErr.ID,
				apiErr.Data,
				apiErr.Origin,
			)
		}
	}
}
