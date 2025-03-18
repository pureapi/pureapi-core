package main

import (
	"context"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"

	"github.com/pureapi/pureapi-core/database"
	"github.com/pureapi/pureapi-core/database/types"
	examples "github.com/pureapi/pureapi-core/doc/examples/database"
)

// Product represents a product entity.
// It implements the Getter interface (TableName and ScanRow).
type Product struct {
	ID    int64
	Name  string
	Price float64
}

// TableName returns the table name for the Product.
func (p *Product) TableName() string {
	return "products"
}

// ScanRow scans a database row into the Product.
func (p *Product) ScanRow(row types.Row) error {
	return row.Scan(&p.ID, &p.Name, &p.Price)
}

// CustomErrorChecker wraps errors with a custom message.
type CustomErrorChecker struct{}

// Check wraps the provided error with additional context.
func (cec *CustomErrorChecker) Check(err error) error {
	return fmt.Errorf("custom error occurred: %w", err)
}

// This example demonstrates basic database operations using common database
// functions from the database package.
func main() {
	// Connect to the database.
	db, err := examples.Connect(examples.Cfg(), examples.DummyConnectionOpen)
	if err != nil {
		log.Fatalf("DB connection error: %v", err)
	}
	defer db.Close()

	// Create the "products" table.
	log.Println("===== Creating table...")
	CreateTable(db)

	// Insert two products.
	log.Println("===== Inserting products...")
	InsertProducts(db)

	// Update a product.
	log.Println("===== Updating product...")
	UpdateProduct(db)

	// Update a product using ExecRaw.
	log.Println("===== Updating product using ExecRaw...")
	UpdateProductRaw(db)

	// Query the count of products and print it.
	log.Println("===== Querying product count...")
	QueryProductCount(db)

	// Query a single product and print it.
	log.Println("===== Querying product...")

	// Query a single product using QueryRaw.
	log.Println("===== Querying product using QueryRaw...")
	GetProductRaw(db)

	// Query all products and print them.
	log.Println("===== Querying all products...")
	GetAllProducts(db)

	// Lower level conversion functions.

	// Demo using RowToEntity.
	log.Println("===== Demo: RowToEntity =====")
	DemoRowToEntity(db)

	// Demo using RowToAny.
	log.Println("===== Demo: RowToAny =====")
	DemoRowToAny(db)

	// Demo using RowsToAny.
	log.Println("===== Demo: RowsToAny =====")
	DemoRowsToAny(db)

	// Demo using RowsToEntities.
	log.Println("===== Demo: RowsToEntities =====")
	DemoRowsToEntities(db)
}

// CreateTable creates the "products" table. It demonstrates how to use the
// Exec function to run a query that creates a table.
//
// Parameters:
//   - db: The database handle.
func CreateTable(db types.DB) {
	// Create the "products" table.
	createTableSQL := `
		CREATE TABLE products (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			price REAL NOT NULL
		);`
	if _, err := db.Exec(createTableSQL); err != nil {
		log.Fatalf("Error creating table: %v", err)
	}
	log.Println("Table 'products' created.")

}

// InsertProducts inserts two products using Exec. It demonstrates how to use
// the Exec function to run a query that inserts rows.
func InsertProducts(preparer types.Preparer) {
	// Insert two products using Exec.
	insertSQL := "INSERT INTO products (name, price) VALUES (?, ?);"
	if _, err := database.Exec(
		context.Background(),
		preparer,
		insertSQL,
		[]any{"Widget", 99.99}, &CustomErrorChecker{},
	); err != nil {
		log.Fatalf("Error inserting product 1: %v", err)
	}
	if _, err := database.Exec(
		context.Background(),
		preparer,
		insertSQL,
		[]any{"Gadget", 199.99}, &CustomErrorChecker{},
	); err != nil {
		log.Fatalf("Error inserting product 2: %v", err)
	}
	log.Println("Inserted 2 products.")
}

// UpdateProduct updates a product using Exec. It demonstrates how to use the
// Exec function to run a query that updates a row.
//
// Parameters:
//   - db: The database handle.
func UpdateProduct(preparer types.Preparer) {
	updateSQL := "UPDATE products SET price = ? WHERE id = ?;"
	res, err := database.Exec(
		context.Background(),
		preparer,
		updateSQL,
		[]any{1.99, 1},
		&CustomErrorChecker{},
	)
	if err != nil {
		log.Fatalf("Error updating product: %v", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Fatalf("Error getting rows affected: %v", err)
	}
	log.Printf("Updated product id=1, rows affected: %d", rowsAffected)
}

// UpdateProductRaw updates a product using ExecRaw. It demonstrates using
// ExecRaw to update a database row.
//
// Parameters:
//   - db: The database handle.
func UpdateProductRaw(db types.DB) {
	updateSQL := "UPDATE products SET price = ? WHERE id = ?;"
	res, err := database.ExecRaw(
		context.Background(),
		db,
		updateSQL,
		[]any{49.99, 2},
		&CustomErrorChecker{},
	)
	if err != nil {
		log.Fatalf("Error updating product with ExecRaw: %v", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Fatalf("Error getting rows affected (ExecRaw): %v", err)
	}
	log.Printf(
		"Updated product id=2 using ExecRaw, rows affected: %d", rowsAffected,
	)
}

// QueryProductCount queries the count of products and prints it. It
// demonstrates how to use the QuerySingleValue function to query a scalar
// value such as an integer.
//
// Parameters:
//   - db: The database handle.
func QueryProductCount(preparer types.Preparer) {
	count, err := database.QuerySingleValue(
		context.Background(),
		preparer,
		"SELECT COUNT(*) FROM products;",
		nil,
		&CustomErrorChecker{},
		func() *int { return new(int) },
	)
	if err != nil {
		log.Fatalf("Error querying count: %v", err)
	}
	log.Printf("Count of products: %d", *count)
}

// GetProduct queries a single product and prints it. It demonstrates how to
// use the QuerySingleEntity function.
//
// Parameters:
//   - db: The database handle.
func GetProduct(preparer types.Preparer) {
	product, err := database.QuerySingleEntity(
		context.Background(),
		preparer,
		"SELECT id, name, price FROM products WHERE id = ?;",
		[]any{1},
		&CustomErrorChecker{},
		func() *Product { return &Product{} },
	)
	if err != nil {
		log.Fatalf("Error querying product: %v", err)
	}
	log.Printf("Queried product: %+v", product)
}

// GetProductRaw gets a product using QueryRaw. It demonstrates using the
// QueryRaw function to manually iterate over rows.
//
// Parameters:
//   - db: The database handle.
func GetProductRaw(db types.DB) {
	rows, err := database.QueryRaw(
		context.Background(),
		db,
		"SELECT id, name, price FROM products WHERE price > ?;",
		[]any{10.0},
		&CustomErrorChecker{},
	)
	if err != nil {
		log.Fatalf("Error in QueryRaw: %v", err)
	}
	defer rows.Close()

	log.Println("QueryRaw results:")
	for rows.Next() {
		var id int64
		var name string
		var price float64
		if err := rows.Scan(&id, &name, &price); err != nil {
			log.Fatalf("Error scanning row: %v", err)
		}
		log.Printf("Raw row - id: %d, name: %s, price: %.2f", id, name, price)
	}
	if err := rows.Err(); err != nil {
		log.Fatalf("Rows iteration error: %v", err)
	}
}

// GetAllProducts queries all products and prints them. It demonstrates how to
// use the QueryEntities function.
//
// Parameters:
//   - db: The database handle.
func GetAllProducts(preparer types.Preparer) {
	// Query all products and map them into a slice of Product structs.
	products, err := database.QueryEntities(
		context.Background(),
		preparer,
		"SELECT id, name, price FROM products;",
		nil,
		&CustomErrorChecker{},
		func() *Product { return &Product{} },
	)
	if err != nil {
		log.Fatalf("Error querying products: %v", err)
	}
	log.Printf("Queried all products: %d", len(products))
	for _, product := range products {
		log.Printf("Product: %+v", product)
	}
}

// DemoRowToEntity demonstrates using RowToEntity to convert a single row into
// a struct.
//
// Parameters:
//   - db: The database handle.
func DemoRowToEntity(db types.DB) {
	stmt, err := db.Prepare(
		"SELECT id, name, price FROM products WHERE id = ?;",
	)
	if err != nil {
		log.Fatalf("Prepare error in DemoRowToEntity: %v", err)
	}
	defer stmt.Close()
	row := stmt.QueryRow(1)
	product, err := database.RowToEntity(
		context.Background(), row, func() *Product { return &Product{} },
	)
	if err != nil {
		log.Fatalf("RowToEntity error: %v", err)
	}
	log.Printf("DemoRowToEntity result: %+v", product)
}

// DemoRowToAny demonstrates using RowToAny to convert a single row into a
// generic map.
//
// Parameters:
//   - db: The database handle.
func DemoRowToAny(db types.DB) {
	stmt, err := db.Prepare("SELECT name FROM products WHERE id = ?;")
	if err != nil {
		log.Fatalf("Prepare error in DemoRowToAny: %v", err)
	}
	defer stmt.Close()
	row := stmt.QueryRow(2)
	namePtr, err := database.RowToAny(
		context.Background(), row, func() *string { return new(string) },
	)
	if err != nil {
		log.Fatalf("RowToAny error: %v", err)
	}
	log.Printf("DemoRowToAny result (product name): %s", *namePtr)
}

// DemoRowsToAny demonstrates using RowsToAny to convert multiple rows into a
// slice of scalar values.
//
// Parameters:
//   - db: The database handle.
func DemoRowsToAny(db types.DB) {
	rows, err := db.Query("SELECT name FROM products WHERE price > ?;", 0)
	if err != nil {
		log.Fatalf("Query error in DemoRowsToAny: %v", err)
	}
	defer rows.Close()
	names, err := database.RowsToAny(
		context.Background(), rows, func() *string { return new(string) },
	)
	if err != nil {
		log.Fatalf("RowsToAny error: %v", err)
	}
	log.Printf("DemoRowsToAny found %d product names:", len(names))
	for _, name := range names {
		log.Printf("Product name: %s", *name)
	}
}

// DemoRowsToEntities demonstrates using RowsToEntities to convert multiple rows
// into a slice of composite structs.
//
// Parameters:
//   - db: The database handle.
func DemoRowsToEntities(db types.DB) {
	rows, err := db.Query("SELECT id, name, price FROM products WHERE price > ?;", 0)
	if err != nil {
		log.Fatalf("Query error in DemoRowsToEntities: %v", err)
	}
	defer rows.Close()
	products, err := database.RowsToEntities(
		context.Background(), rows, func() *Product { return &Product{} },
	)
	if err != nil {
		log.Fatalf("RowsToEntities error: %v", err)
	}
	log.Printf("DemoRowsToEntities found %d products:", len(products))
	for _, p := range products {
		log.Printf("Product: %+v", p)
	}
}
