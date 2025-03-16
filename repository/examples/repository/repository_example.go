package main

import (
	"context"
	"log"

	_ "github.com/mattn/go-sqlite3"

	dbexamples "github.com/pureapi/pureapi-core/database/examples"
	"github.com/pureapi/pureapi-core/database/types"
	"github.com/pureapi/pureapi-core/repository"
	"github.com/pureapi/pureapi-core/repository/examples"
	repotypes "github.com/pureapi/pureapi-core/repository/types"
)

// This example demonstrates the usage of the Repository interface. It creates
// a database table and inserts and retrieves a user using repositories and uses
// custom implementations of the QueryBuilder and ErrorChecker interfaces to
// handle the database operations.
func main() {
	// Connect to the database.
	db, err := dbexamples.Connect(
		dbexamples.Cfg(), dbexamples.DummyConnectionOpen,
	)
	if err != nil {
		log.Fatalf("Connection failed: %v", err)
	}
	defer db.Close()

	// Create the table.
	CreateTable(db)
	// Insert a new user.
	CreateUser(db)
	// Retrieve the inserted user.
	GetUser(db)
}

// CreateTable creates the "users" table.
func CreateTable(db types.DB) {
	schemaManager := &examples.SimpleSchemaManager{}
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
	mutatorRepo := repository.NewMutatorRepo[*examples.User](
		&examples.SimpleQueryBuilder{}, &examples.SimpleErrorChecker{},
	)
	newUser := &examples.User{Name: "Alice"}
	insertedUser, err := mutatorRepo.Insert(context.Background(), db, newUser)
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
	readerRepo := repository.NewReaderRepo[*examples.User](
		&examples.SimpleQueryBuilder{}, &examples.SimpleErrorChecker{},
	)
	retrievedUser, err := readerRepo.GetOne(
		context.Background(),
		db,
		func() *examples.User { return &examples.User{} },
		nil,
	)
	if err != nil {
		log.Printf("GetOne error: %v", err)
		return
	}
	log.Printf("Retrieved user: %+v\n", retrievedUser)
}
