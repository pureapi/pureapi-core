package main

import (
	"context"
	"errors"
	"log"

	_ "github.com/mattn/go-sqlite3"

	"github.com/pureapi/pureapi-core/apierror"
	"github.com/pureapi/pureapi-core/database/examples"
	"github.com/pureapi/pureapi-core/database/types"
	"github.com/pureapi/pureapi-core/repository"
	exampletypes "github.com/pureapi/pureapi-core/repository/examples/repository_query_builder/types"
	repotypes "github.com/pureapi/pureapi-core/repository/types"
)

// This example demonstrates the usage of the QueryBuilder and ErrorChecker
// interfaces. It creates a user table and inserts and retrieves a user
// using the custom QueryBuilder and ErrorChecker implementations.
func main() {
	db, err := examples.Connect(examples.Cfg(), examples.DummyConnectionOpen)
	if err != nil {
		log.Printf("Connection failed: %v", err)
	}
	defer db.Close()

	// Create the "users" table. Then create and retrieve a user.
	CreateTable(db)
	CreateUser(db)
	GetUser(db)

	// Try to retrieve a user from a non-existent table.
	// This should log an error processed by the custom ErrorChecker.
	GetErrorUser(db)
}

// CreateTable creates the "users" table using the custom SchemaManager
// implementation. This demonstrates that you can use a custom query builder
// implementation and run custom SQL queries without the need for a repository.
//
// Parameters:
//   - db: The database handle.
func CreateTable(db types.DB) {
	schemaManager := &exampletypes.SimpleSchemaManager{}
	columns := []repotypes.ColumnDefinition{
		{
			Name:          "id",
			Type:          "INTEGER",
			NotNull:       true,
			AutoIncrement: true,
			PrimaryKey:    true,
		},
		{
			Name:    "name",
			Type:    "TEXT",
			NotNull: true,
		},
	}
	createTableQuery, _, err := schemaManager.CreateTableQuery(
		"users", true, columns, nil, repotypes.TableOptions{},
	)
	if err != nil {
		log.Printf("Create table query error: %v", err)
		return
	}
	if _, err = db.Exec(createTableQuery); err != nil {
		log.Printf("Create table execution error: %v", err)
		return
	}
	log.Println("Table 'users' created.")
}

// CreateUser inserts a new user into the database using the custom QueryBuilder
// and ErrorChecker implementations.
//
// Parameters:
//   - db: The database handle.
func CreateUser(db types.DB) {
	// Create a mutatorRepo for inserting users.
	// We use *User because it implements the Mutator interface.
	mutatorRepo := repository.NewMutatorRepo[*exampletypes.User](
		&exampletypes.SimpleQueryBuilder{}, &exampletypes.SimpleErrorChecker{},
	)

	// Create a new user and insert it.
	newUser := &exampletypes.User{Name: "Alice"}
	insertedUser, err := mutatorRepo.Insert(
		context.Background(), db, newUser,
	)
	if err != nil {
		log.Printf("Insert error: %v", err)
		return
	}
	log.Printf("Inserted user: %+v\n", insertedUser)
}

// GetUser retrieves a user from the database using the custom QueryBuilder
// and ErrorChecker implementations.
//
// Parameters:
//   - db: The database handle.
func GetUser(db types.DB) {
	// Create a readerRepo for retrieving users.
	// We use *User because it implements the Getter interface.
	readerRepo := repository.NewReaderRepo[*exampletypes.User](
		&exampletypes.SimpleQueryBuilder{}, &exampletypes.SimpleErrorChecker{},
	)

	// Retrieve a single user record.
	retrievedUser, err := readerRepo.GetOne(
		context.Background(),
		db,
		func() *exampletypes.User { return &exampletypes.User{} },
		nil,
	)
	if err != nil {
		log.Printf("GetOne error: %v", err)
	}
	log.Printf("Retrieved user: %+v\n", retrievedUser)
}

// GetErrorUser attempts to retrieve a record using ErrorUser which maps to a
// non-existent table. This demonstrates how the custom ErrorChecker can be used
// to handle database errors and translate them into application errors.
//
// Parameters:
//   - db: The database handle.
func GetErrorUser(db types.DB) {
	// Using ErrorUser will trigger an error because its TableName returns
	// "nonexistent_users".
	readerRepo := repository.NewReaderRepo[*exampletypes.ErrorUser](
		&exampletypes.SimpleQueryBuilder{}, &exampletypes.SimpleErrorChecker{},
	)
	_, err := readerRepo.GetOne(
		context.Background(),
		db,
		func() *exampletypes.ErrorUser { return &exampletypes.ErrorUser{} },
		nil,
	)
	if err != nil {
		// The error should be wrapped with our custom message.
		log.Printf("GetErrorUser error: %v", err)

		// Check if the error is an APIError and log it.
		var apiError *apierror.APIError
		if ok := errors.As(err, &apiError); ok {
			log.Printf(
				"APIError, ID: %v, data: %s, origin: %v\n",
				apiError.ID,
				apiError.Data,
				apiError.Origin,
			)
		}
	}
}
