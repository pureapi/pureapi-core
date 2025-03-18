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

// Order represents an order in our orders table.
type Order struct {
	ID       int64
	Item     string
	Quantity int
}

// ScanRow scans a database row into the Order.
func (o *Order) ScanRow(row types.Row) error {
	return row.Scan(&o.ID, &o.Item, &o.Quantity)
}

// TableName returns the table name for the Order.
func (o *Order) TableName() string {
	return "orders"
}

// CustomErrorChecker wraps errors with a custom message.
type CustomErrorChecker struct{}

// Check wraps the provided error with additional context.
func (cec *CustomErrorChecker) Check(err error) error {
	return fmt.Errorf("custom error occurred: %w", err)
}

// This example demonstrates how to manage transactions manually and pass
// them to the repository database functions.
func main() {
	// Connect to the database.
	db, err := examples.Connect(examples.Cfg(), examples.DummyConnectionOpen)
	if err != nil {
		log.Fatalf("DB connection error: %v", err)
	}
	defer db.Close()

	// Create the "orders" table.
	createTableSQL := `
		CREATE TABLE IF NOT EXISTS orders (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			item TEXT NOT NULL,
			quantity INTEGER NOT NULL
		);`
	if _, err := db.Exec(createTableSQL); err != nil {
		log.Fatalf("Error creating orders table: %v", err)
	}
	log.Println("Table 'orders' created.")
	ctx := context.Background()

	// Run a successful transaction.
	RunSuccessfulTransaction(ctx, db)

	// Run a rolled back transaction.
	RunRolledBackTransaction(ctx, db)

	// Query the order count and order by ID.
	GetOrderCount(ctx, db)
	GetOrderByID(ctx, db, 1)

	// Should log a custom error since the transaction was rolled back.
	GetOrderByID(ctx, db, 2)
}

// Run a successful transaction.
//
// Parameters:
//   - ctx: The context for the transaction.
//   - db: The database connection.
func RunSuccessfulTransaction(ctx context.Context, db types.DB) {
	tx, err := BeginTx(ctx, db)
	if err != nil {
		log.Fatalf("Error beginning transaction: %v", err)
	}

	// Insert an order and update it.
	InsertOrder(ctx, tx)
	UpdateOrder(ctx, tx, 1)

	// Query the order count.
	GetOrderCount(ctx, tx)
	GetOrderByID(ctx, tx, 1)

	// Commit the transaction.
	if err := tx.Commit(); err != nil {
		log.Fatalf("Error committing transaction: %v", err)
	}

	log.Println("Transaction committed successfully")
}

// Run a rolled back transaction.
//
// Parameters:
//   - ctx: The context for the transaction.
//   - db: The database connection.
func RunRolledBackTransaction(ctx context.Context, db types.DB) {
	tx, err := BeginTx(ctx, db)
	if err != nil {
		log.Fatalf("Error beginning transaction: %v", err)
	}

	// Insert an order and update it.
	InsertOrder(ctx, tx)
	UpdateOrder(ctx, tx, 1)

	// Query the order count.
	GetOrderCount(ctx, tx)
	GetOrderByID(ctx, tx, 1)

	// Rollback the transaction.
	if err := tx.Rollback(); err != nil {
		log.Fatalf("Error rolling back transaction: %v", err)
	}

	log.Println("Transaction rolled back")
}

// Begin a transaction.
//
// Parameters:
//   - db: The database connection.
//
// Returns:
//   - types.Tx: The transaction object.
//   - error: An error if the transaction fails.
func BeginTx(ctx context.Context, db types.DB) (types.Tx, error) {
	// Begin a transaction.
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	return tx, nil
}

// InsertOrder inserts a new order. It demonstrates how to
// the Exec function can be used with transactions.
//
// Parameters:
//   - db: The database connection.
//
// Returns:
//   - int64: The ID of the inserted order.
func InsertOrder(ctx context.Context, tx types.Tx) int64 {
	insertSQL := "INSERT INTO orders (item, quantity) VALUES (?, ?);"
	res, err := database.Exec(
		ctx, tx, insertSQL, []any{"Book", 5}, &CustomErrorChecker{},
	)
	if err != nil {
		err := tx.Rollback()
		if err != nil {
			log.Fatalf("rollback error: %v", err)
		}
		log.Fatalf("insert order error: %v", err)
	}
	orderID, err := res.LastInsertId()
	if err != nil {
		err := tx.Rollback()
		if err != nil {
			log.Fatalf("rollback error: %v", err)
		}
		log.Fatalf("last insert id error: %v", err)
	}
	return orderID
}

// UpdateOrder updates the quantity of an order. It demonstrates how to
// the Exec function can be used with transactions.
//
// Parameters:
//   - db: The database connection.
//   - orderID: The ID of the order to update.
//
// Returns:
//   - string: A success message.
//   - error: An error if the update fails.
func UpdateOrder(
	ctx context.Context, tx types.Tx, orderID int64,
) {
	// Otherwise, update the order.
	updateSQL := "UPDATE orders SET quantity = ? WHERE id = ?;"
	_, err := database.Exec(
		ctx, tx, updateSQL, []any{10, orderID}, &CustomErrorChecker{},
	)
	if err != nil {
		err := tx.Rollback()
		if err != nil {
			log.Fatalf("rollback error: %v", err)
		}
		log.Fatalf("update order error: %v", err)
	}
	log.Printf("Order %d updated successfully", orderID)
}

// GetOrderCount uses QuerySingleValue to get the count of orders.
//
// Parameters:
//   - db: The database connection.
func GetOrderCount(ctx context.Context, preparer types.Preparer) {
	count, err := database.QuerySingleValue(
		ctx,
		preparer,
		"SELECT COUNT(*) FROM orders;",
		nil,
		&CustomErrorChecker{},
		func() *int { return new(int) },
	)
	if err != nil {
		log.Fatalf("query order count error: %v", err)
	}
	log.Printf("Order count: %d", *count)
}

// GetOrderByID uses QuerySingleEntity to get an order by ID.
//
// Parameters:
//   - db: The database connection.
func GetOrderByID(
	ctx context.Context, preparer types.Preparer, orderID int64,
) {
	order, err := database.QuerySingleEntity(
		ctx,
		preparer,
		"SELECT id, item, quantity FROM orders WHERE id = ?;",
		[]any{orderID},
		&CustomErrorChecker{},
		func() *Order { return &Order{} },
	)
	if err != nil {
		log.Fatalf("query order error: %v", err)
	}
	log.Printf("Order: %+v", order)
}
