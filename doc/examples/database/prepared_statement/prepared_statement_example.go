package main

import (
	"log"

	_ "github.com/mattn/go-sqlite3"
	examples "github.com/pureapi/pureapi-core/doc/examples/database"
)

// This example shows how to use a prepared statement using the Prepare
// function. It uses the sqlDB DB implementation.
func main() {
	db, err := examples.Connect(examples.Cfg(), examples.DummyConnectionOpen)
	if err != nil {
		log.Printf("Connection failed: %v", err)
	}
	defer db.Close()

	// Create a sample "orders" table.
	_, err = db.Exec(`
CREATE TABLE orders (
  id       INTEGER PRIMARY KEY AUTOINCREMENT,
  customer TEXT NOT NULL,
  amount   REAL NOT NULL
);`)
	if err != nil {
		log.Printf("Create table error: %v", err)
	}
	log.Println("Table 'orders' created.")

	// Prepare an insert statement.
	stmt, err := db.Prepare(
		`INSERT INTO orders(customer, amount) VALUES(?, ?);`,
	)
	if err != nil {
		log.Printf("Prepare statement error: %v", err)
	}
	defer stmt.Close()

	// Execute the prepared statement for multiple customers.
	customers := []string{"Bob", "Carol", "Dave"}
	for _, customer := range customers {
		res, err := stmt.Exec(customer, 99.99)
		if err != nil {
			log.Printf("Statement Exec error: %v", err)
		}
		id, err := res.LastInsertId()
		if err != nil {
			log.Printf("LastInsertId error: %v", err)
		}
		log.Printf("Inserted order for %s with ID %d\n", customer, id)
	}
}
