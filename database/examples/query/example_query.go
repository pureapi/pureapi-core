package main

import (
	"log"

	"github.com/pureapi/pureapi-core/database/examples"

	// Using the SQLite3 driver as an example.
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := examples.Connect(examples.Cfg(), examples.DummyConnectionOpen)
	if err != nil {
		log.Fatalf("Connection failed: %v", err)
	}
	defer db.Close()

	// Create a sample "users" table.
	createTable := `
CREATE TABLE users (
  id   INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL,
  age  INTEGER NOT NULL
);`
	_, err = db.Exec(createTable)
	if err != nil {
		log.Fatalf("Create table error: %v", err)
	}
	log.Println("Table 'users' created.")

	// Insert a new user.
	insertUser := `INSERT INTO users(name, age) VALUES(?, ?);`
	res, err := db.Exec(insertUser, "Alice", 30)
	if err != nil {
		log.Fatalf("Insert error: %v", err)
	}
	lastID, err := res.LastInsertId()
	if err != nil {
		log.Fatalf("LastInsertId error: %v", err)
	}
	log.Printf("Inserted user with ID %d\n", lastID)

	// Query the record.
	query := `SELECT id, name, age FROM users;`
	rows, err := db.Query(query)
	if err != nil {
		log.Fatalf("Query error: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int64
		var name string
		var age int
		if err := rows.Scan(&id, &name, &age); err != nil {
			log.Fatalf("Scan error: %v", err)
		}
		log.Printf("Query User: ID=%d, Name=%s, Age=%d\n", id, name, age)
	}
	if err := rows.Err(); err != nil {
		log.Fatalf("Rows error: %v", err)
	}
}
