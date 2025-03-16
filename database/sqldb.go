package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/pureapi/pureapi-core/database/types"
)

// sqlDB wraps *sql.DB for database operations.
type sqlDB struct {
	*sql.DB
}

// sqlDB implements DB interface.
var _ types.DB = (*sqlDB)(nil)

// NewSQLDB creates a new instance of SQLDB and connects to the database using
// the provided driver and connection string.
//
// Parameters:
//   - driver: The database driver name.
//   - dsn: The database connection string.
//
// Returns:
//   - *sqlDB: A new instance of sqlDB.
//   - error: An error if the connection fails.
func NewSQLDB(driver string, dsn string) (*sqlDB, error) {
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, fmt.Errorf("sql.Open error: %w", err)
	}
	return &sqlDB{DB: db}, nil
}

// NewSQLDBAdapter creates a new instance of SQLDB and connects to the database
// using the provided driver and connection string. It adapts the NewSQLDB
// function to the DB interface.
//
// Parameters:
//   - driver: The database driver name.
//   - dsn: The database connection string.
//
// Returns:
//   - DB: A new instance of DB.
//   - error: An error if the connection fails.
func NewSQLDBAdapter(driver string, dsn string) (types.DB, error) {
	db, err := NewSQLDB(driver, dsn)
	if err != nil {
		return nil, fmt.Errorf("NewSQLDBAdapter error: %w", err)
	}
	return db, err
}

// Ping sends a ping to the database to check if it is alive.
//
// Returns:
//   - error: An error if the ping fails.
func (db *sqlDB) Ping() error {
	return db.DB.Ping()
}

// SetConnMaxLifetime sets the maximum time a connection may be reused.
//
// Parameters:
//   - d: The maximum lifetime of a connection.
func (db *sqlDB) SetConnMaxLifetime(d time.Duration) {
	db.DB.SetConnMaxLifetime(d)
}

// SetConnMaxIdleTime sets the maximum time an idle connection may be reused.
//
// Parameters:
//   - d: The maximum idle time of a connection.
func (db *sqlDB) SetConnMaxIdleTime(d time.Duration) {
	db.DB.SetConnMaxIdleTime(d)
}

// SetMaxOpenConns sets the maximum number of open connections to the database.
//
// Parameters:
//   - n: The maximum number of open connections.
func (db *sqlDB) SetMaxOpenConns(n int) {
	db.DB.SetMaxOpenConns(n)
}

// SetMaxIdleConns sets the maximum number of idle connections to the database.
//
// Parameters:
//   - n: The maximum number of idle connections.
func (db *sqlDB) SetMaxIdleConns(n int) {
	db.DB.SetMaxIdleConns(n)
}

// Prepare creates a prepared statement for later queries or executions.
//
// Parameters:
//   - query: The SQL query string to prepare.
//
// Returns:
//   - Stmt: The prepared statement.
//   - error: An error if the statement cannot be prepared.
func (db *sqlDB) Prepare(query string) (types.Stmt, error) {
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("DB.Prepare error: %w", err)
	}
	return &RealStmt{Stmt: stmt}, nil
}

// BeginTx creates a transaction and returns it.
//
// Parameters:
//   - ctx: The context for the transaction.
//   - opts: The transaction options.
//
// Returns:
//   - Tx: The transaction.
//   - error: An error if the transaction cannot be created.
func (db *sqlDB) BeginTx(
	ctx context.Context, opts *sql.TxOptions,
) (types.Tx, error) {
	tx, err := db.DB.BeginTx(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("DB.BeginTx error: %w", err)
	}
	return &RealTx{Tx: tx}, nil
}

// Exec executes a query without returning rows.
//
// Parameters:
//   - query: The SQL query string to execute.
//   - args: The query parameters.
//
// Returns:
//   - Result: The result of the query.
//   - error: An error if the query fails.
func (db *sqlDB) Exec(query string, args ...any) (types.Result, error) {
	res, err := db.DB.Exec(query, args...)
	if err != nil {
		return nil, fmt.Errorf("DB.Exec error: %w", err)
	}
	return &RealResult{Result: res}, nil
}

// Query executes a query that returns rows.
//
// Parameters:
//   - query: The SQL query string to execute.
//   - args: The query parameters.
//
// Returns:
//   - Rows: The rows of the query.
//   - error: An error if the query fails.
func (db *sqlDB) Query(query string, args ...any) (types.Rows, error) {
	rows, err := db.DB.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("DB.Query error: %w", err)
	}
	return &RealRows{Rows: rows}, nil
}

// QueryRow executes a query that returns a single row.
//
// Parameters:
//   - query: The SQL query string to execute.
//   - args: The query parameters.
//
// Returns:
//   - Row: The row of the query.
func (db *sqlDB) QueryRow(query string, args ...any) types.Row {
	return db.DB.QueryRow(query, args...)
}

// RealStmt wraps *sql.Stmt to implement the Stmt interface.
type RealStmt struct {
	*sql.Stmt
}

// Close closes the statement.
//
// Returns:
//   - error: An error if the statement cannot be closed.
func (s *RealStmt) Close() error {
	err := s.Stmt.Close()
	if err != nil {
		return fmt.Errorf("Stmt.Close error: %w", err)
	}
	return nil
}

// QueryRow executes a prepared query statement with the given arguments.
//
// Parameters:
//   - args: The query parameters.
//
// Returns:
//   - Row: The row of the query.
func (s *RealStmt) QueryRow(args ...any) types.Row {
	return s.Stmt.QueryRow(args...)
}

// Exec executes a prepared statement with the given arguments.
//
// Parameters:
//   - args: The query parameters.
//
// Returns:
//   - Result: The result of the query.
func (s *RealStmt) Exec(args ...any) (types.Result, error) {
	res, err := s.Stmt.Exec(args...)
	if err != nil {
		return nil, fmt.Errorf("Stmt.Exec error: %w", err)
	}
	return &RealResult{Result: res}, nil
}

// Query executes a prepared query statement with the given arguments.
//
// Parameters:
//   - args: The query parameters.
//
// Returns:
//   - Rows: The rows of the query.
func (s *RealStmt) Query(args ...any) (types.Rows, error) {
	rows, err := s.Stmt.Query(args...)
	if err != nil {
		return nil, fmt.Errorf("Stmt.Query error: %w", err)
	}
	return &RealRows{Rows: rows}, nil
}

// RealRows wraps *sql.Rows to implement the Rows interface.
type RealRows struct {
	*sql.Rows
}

// Scan scans the rows into dest.
//
// Parameters:
//   - dest: The destination slice to scan into.
//
// Returns:
//   - error: An error if the rows cannot be scanned.
func (r *RealRows) Scan(dest ...any) error {
	err := r.Rows.Scan(dest...)
	if err != nil {
		return fmt.Errorf("Rows.Scan error: %w", err)
	}
	return nil
}

// Next advances the rows.
//
// Returns:
//   - bool: True if there are more rows, false otherwise.
func (r *RealRows) Next() bool {
	return r.Rows.Next()
}

// Close closes the rows.
//
// Returns:
//   - error: An error if the rows cannot be closed.
func (r *RealRows) Close() error {
	err := r.Rows.Close()
	if err != nil {
		return fmt.Errorf("Rows.Close error: %w", err)
	}
	return nil
}

// Err returns the error, if any, that was encountered during iteration.
//
// Returns:
//   - error: The error, if any, that was encountered during iteration.
func (r *RealRows) Err() error {
	err := r.Rows.Err()
	if err != nil {
		return fmt.Errorf("Rows.Err error: %w", err)
	}
	return nil
}

// RealRow wraps *sql.Row to implement the Row interface.
type RealRow struct {
	*sql.Row
}

// Scan scans the row into dest.
//
// Parameters:
//   - dest: The destination slice to scan into.
//
// Returns:
//   - error: An error if the row cannot be scanned.
func (r *RealRow) Scan(dest ...any) error {
	err := r.Row.Scan(dest...)
	if err != nil {
		return fmt.Errorf("Row.Scan error: %w", err)
	}
	return nil
}

// RealResult wraps sql.Result to implement the Result interface.
type RealResult struct {
	Result sql.Result
}

// LastInsertId returns the last inserted id.
//
// Returns:
//   - int64: The last inserted id.
//   - error: An error if the last inserted id cannot be retrieved.
func (r *RealResult) LastInsertId() (int64, error) {
	id, err := r.Result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("Result.LastInsertId error: %w", err)
	}
	return id, nil
}

// RowsAffected returns the number of rows affected.
//
// Returns:
//   - int64: The number of rows affected.
//   - error: An error if the number of rows affected cannot be retrieved.
func (r *RealResult) RowsAffected() (int64, error) {
	n, err := r.Result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("Result.RowsAffected error: %w", err)
	}
	return n, nil
}

// RealTx wraps *sql.Tx to implement the Tx interface.
type RealTx struct {
	*sql.Tx
}

// Prepare commits the transaction.
//
// Returns:
//   - error: An error if the transaction cannot be committed.
func (tx *RealTx) Commit() error {
	if err := tx.Tx.Commit(); err != nil {
		return fmt.Errorf("Tx.Commit error: %w", err)
	}
	return nil
}

// Rollback rollbacks the transaction.
//
// Returns:
//   - error: An error if the transaction cannot be rolled back.
func (tx *RealTx) Rollback() error {
	if err := tx.Tx.Rollback(); err != nil {
		return fmt.Errorf("Tx.Rollback error: %w", err)
	}
	return nil
}

// Prepare prepares the statement.
//
// Parameters:
//   - query: The SQL query string to prepare.
//
// Returns:
//   - Stmt: The prepared statement.
//   - error: An error if the statement cannot be prepared.
func (tx *RealTx) Prepare(query string) (types.Stmt, error) {
	stmt, err := tx.Tx.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("Tx.Prepare error: %w", err)
	}
	return &RealStmt{Stmt: stmt}, nil
}

// Exec executes a query without returning rows.
//
// Parameters:
//   - query: The SQL query string to execute.
//   - args: The query parameters.
//
// Returns:
//   - Result: The result of the query.
//   - error: An error if the query cannot be executed.
func (tx *RealTx) Exec(query string, args ...any) (types.Result, error) {
	res, err := tx.Tx.Exec(query, args...)
	if err != nil {
		return nil, fmt.Errorf("Tx.Exec error: %w", err)
	}
	return &RealResult{Result: res}, nil
}
