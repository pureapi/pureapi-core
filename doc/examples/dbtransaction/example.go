package main

import (
	"context"
	"fmt"
	"log"

	"github.com/pureapi/pureapi-core/database"
	"github.com/pureapi/pureapi-core/database/types"
	"github.com/pureapi/pureapi-core/doc/examples"

	// Using the SQLite3 driver as an example.
	_ "github.com/mattn/go-sqlite3"
)

// This example shows how to execute a transaction using the Transaction
// function. It demonstrates multiple operations within a transaction and
// handles errors gracefully.
func main() {
	db, err := examples.Connect(examples.Cfg(), examples.DummyConnectionOpen)
	if err != nil {
		log.Printf("Connection failed: %v", err)
	}
	defer db.Close()

	// Create a sample "products" table.
	_, err = db.Exec(`
CREATE TABLE products (
  id       INTEGER PRIMARY KEY AUTOINCREMENT,
  name     TEXT NOT NULL,
  quantity INTEGER NOT NULL
);`)
	if err != nil {
		log.Printf("Create table error: %v", err)
	}
	log.Println("Table 'products' created.")

	// Execute a successful transaction.
	SuccessfulTransaction(db, "Gadget", 100)
	QueryAllProducts(db)

	// Execute a failed transaction. Should see an "intentional error" message.
	FailedTransaction(db, "Widget", 50)

	// Should report only the row created by the successful transaction.
	QueryAllProducts(db)
}

// SuccessfulTransaction executes a successful transaction.
// It inserts a new product and then updates its quantity to demonstrate
// multiple operations within a transaction.
//
// Parameters:
//   - db: The database handle.
//   - product: The name of the product.
//   - quantity: The quantity of the product.
func SuccessfulTransaction(db types.DB, product string, quantity int) {
	// Define a transactional function that inserts and then updates a product.
	txFn := func(ctx context.Context, tx types.Tx) (int64, error) {
		// Insert a new product.
		productID, err := InsertProduct(tx, product, quantity)
		if err != nil {
			return 0, err
		}
		// Update the quantity.
		_, err = tx.Exec(
			`UPDATE products SET quantity = ? WHERE id = ?;`,
			quantity-1,
			productID,
		)
		if err != nil {
			return 0, err
		}
		return productID, nil
	}

	// Begin a transaction.
	tx, err := db.BeginTx(context.Background(), nil)
	if err != nil {
		log.Printf("Failed to begin transaction: %v", err)
		return
	}
	// Execute the transaction.
	productID, err := database.Transaction(context.Background(), tx, txFn)
	if err != nil {
		log.Printf("Transaction failed: %v", err)
		return
	}
	log.Printf("Transaction succeeded, product ID: %d\n", productID)
}

// FailedTransaction executes a failed transaction.
// It inserts a new product and then fails intentionally.
//
// Parameters:
//   - db: The database handle.
//   - product: The name of the product.
//   - quantity: The quantity of the product.
func FailedTransaction(db types.DB, product string, quantity int) {
	// Define a transactional function that inserts a product and then fails.
	txFn := func(ctx context.Context, tx types.Tx) (int64, error) {
		// Insert a new product.
		_, err := InsertProduct(tx, product, quantity)
		if err != nil {
			return 0, err
		}
		return 0, fmt.Errorf("intentional error")
	}

	// Begin a transaction.
	tx, err := db.BeginTx(context.Background(), nil)
	if err != nil {
		log.Printf("Failed to begin transaction: %v", err)
	}
	// Execute the transaction.
	_, err = database.Transaction(context.Background(), tx, txFn)
	if err != nil {
		log.Printf("Transaction failed: %v", err)
		return
	}
	// Rollback the transaction.
	if err := tx.Rollback(); err != nil {
		log.Printf("Failed to rollback transaction: %v", err)
		return
	}
}

// InsertProduct inserts a new product into the database.
//
// Parameters:
//   - tx: The transaction handle.
//   - product: The name of the product.
//   - quantity: The quantity of the product.
//
// Returns:
//   - int64: The ID of the inserted product.
//   - error: An error if the insertion fails.
func InsertProduct(tx types.Tx, product string, quantity int) (int64, error) {
	// Insert a new product.
	res, err := tx.Exec(
		`INSERT INTO products(name, quantity) VALUES(?, ?);`,
		product,
		quantity,
	)
	if err != nil {
		return 0, err
	}
	productID, err := res.LastInsertId()
	log.Printf("Product ID: %d", productID)
	if err != nil {
		return 0, err
	}
	return productID, nil
}

// QueryAllProducts queries all products. It prints each product's name and
// quantity found in the database.
//
// Parameters:
//   - db: The database handle.
func QueryAllProducts(db types.DB) {
	log.Println("Querying all products...")

	// Query all products.
	rows, err := db.Query(`SELECT name, quantity FROM products`)
	if err != nil {
		log.Printf("Query failed: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var name string
		var quantity int
		if err := rows.Scan(&name, &quantity); err != nil {
			log.Printf("Scan failed: %v", err)
		}
		log.Printf("Product: %s, Quantity: %d", name, quantity)
	}
}
