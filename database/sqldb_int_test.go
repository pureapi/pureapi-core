package database_test

import (
	"context"
	"testing"
	"time"

	"github.com/pureapi/pureapi-core/database"
	"github.com/pureapi/pureapi-core/database/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	_ "github.com/mattn/go-sqlite3"
)

// openInMemoryDB is a helper that opens an in-memory sqlite3 database.
func openInMemoryDB(t *testing.T) types.DB {
	db, err := database.NewSQLDBAdapter("sqlite3", ":memory:")
	require.NoError(t, err)
	return db
}

// SQLDBIntTestSuite is a suite of tests for SQLDB-related integration tests.
// It uses an in-memory SQLite database.
type SQLDBIntTestSuite struct {
	suite.Suite
	db types.DB
}

// TestSQLDBIntTestSuite registers the test suite.
func TestSQLDBIntTestSuite(t *testing.T) {
	suite.Run(t, new(SQLDBIntTestSuite))
}

// SetupTest initializes common test configurations.
func (s *SQLDBIntTestSuite) SetupTest() {
	s.db = openInMemoryDB(s.T())
}

// TearDownTest tears down common test configurations.
func (s *SQLDBIntTestSuite) TearDownTest() {
	s.db.Close()
}

// TestSQLDB_Ping verifies that Ping works.
func TestSQLDB_Ping(t *testing.T) {
	db := openInMemoryDB(t)
	defer db.Close()

	err := db.Ping()
	require.NoError(t, err)
}

// TestSQLDB_Exec_Query tests Exec and Query methods by creating a table,
// inserting a record, and then querying the inserted data.
func (s *SQLDBIntTestSuite) Test_Exec_Query() {

	createTable := `
		CREATE TABLE IF NOT EXISTS test (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL
		);
	`
	_, err := s.db.Exec(createTable)
	require.NoError(s.T(), err)

	insertQuery := `INSERT INTO test (name) VALUES (?)`
	res, err := s.db.Exec(insertQuery, "Alice")
	require.NoError(s.T(), err)

	rowsAffected, err := res.RowsAffected()
	require.NoError(s.T(), err)
	assert.Equal(s.T(), int64(1), rowsAffected)

	lastID, err := res.LastInsertId()
	require.NoError(s.T(), err)
	assert.True(s.T(), lastID > 0)

	query := `SELECT id, name FROM test WHERE id = ?`
	rows, err := s.db.Query(query, lastID)
	require.NoError(s.T(), err)
	defer rows.Close()

	var id int
	var name string
	if rows.Next() {
		err = rows.Scan(&id, &name)
		require.NoError(s.T(), err)
		assert.Equal(s.T(), lastID, int64(id))
		assert.Equal(s.T(), "Alice", name)
	} else {
		s.T().Fatal("expected one row")
	}
	require.NoError(s.T(), rows.Err())
}

// TestSQLDB_Prepare verifies that prepared statements work as expected.
func (s *SQLDBIntTestSuite) Test_Prepare() {

	createTable := `
		CREATE TABLE IF NOT EXISTS user (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT NOT NULL
		);
	`
	_, err := s.db.Exec(createTable)
	require.NoError(s.T(), err)

	// Test prepared insert statement.
	stmt, err := s.db.Prepare(`INSERT INTO user (username) VALUES (?)`)
	require.NoError(s.T(), err)

	res, err := stmt.Exec("bob")
	require.NoError(s.T(), err)
	rowsAffected, err := res.RowsAffected()
	require.NoError(s.T(), err)
	assert.Equal(s.T(), int64(1), rowsAffected)

	lastID, err := res.LastInsertId()
	require.NoError(s.T(), err)
	assert.True(s.T(), lastID > 0)

	err = stmt.Close()
	require.NoError(s.T(), err)

	// Test prepared query statement.
	stmt, err = s.db.Prepare(`SELECT id, username FROM user WHERE id = ?`)
	require.NoError(s.T(), err)

	row := stmt.QueryRow(lastID)
	var id int
	var username string
	err = row.Scan(&id, &username)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), lastID, int64(id))
	assert.Equal(s.T(), "bob", username)

	err = stmt.Close()
	require.NoError(s.T(), err)
}

// TestSQLDB_BeginTx_Commit tests starting a transaction, inserting data,
// committing the transaction, and then verifying that the data persists.
func (s *SQLDBIntTestSuite) Test_BeginTx_Commit() {

	createTable := `
		CREATE TABLE IF NOT EXISTS account (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			balance INTEGER NOT NULL
		);
	`
	_, err := s.db.Exec(createTable)
	require.NoError(s.T(), err)

	tx, err := s.db.BeginTx(context.Background(), nil)
	require.NoError(s.T(), err)

	insertQuery := `INSERT INTO account (balance) VALUES (?)`
	res, err := tx.Exec(insertQuery, 100)
	require.NoError(s.T(), err)

	lastID, err := res.LastInsertId()
	require.NoError(s.T(), err)

	err = tx.Commit()
	require.NoError(s.T(), err)

	query := `SELECT balance FROM account WHERE id = ?`
	rows, err := s.db.Query(query, lastID)
	require.NoError(s.T(), err)
	defer rows.Close()

	var balance int
	if rows.Next() {
		err = rows.Scan(&balance)
		require.NoError(s.T(), err)
		assert.Equal(s.T(), 100, balance)
	} else {
		s.T().Fatal("expected one row")
	}
	require.NoError(s.T(), rows.Err())
}

// TestSQLDB_BeginTx_Rollback tests that a transaction rollback cancels
// the operations performed within the transaction.
func (s *SQLDBIntTestSuite) Test_BeginTx_Rollback() {

	defer s.db.Close()

	createTable := `
		CREATE TABLE IF NOT EXISTS product (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL
		);
	`
	_, err := s.db.Exec(createTable)
	require.NoError(s.T(), err)

	tx, err := s.db.BeginTx(context.Background(), nil)
	require.NoError(s.T(), err)

	insertQuery := `INSERT INTO product (name) VALUES (?)`
	_, err = tx.Exec(insertQuery, "Widget")
	require.NoError(s.T(), err)

	err = tx.Rollback()
	require.NoError(s.T(), err)

	query := `SELECT id FROM product WHERE name = ?`
	rows, err := s.db.Query(query, "Widget")
	require.NoError(s.T(), err)
	defer rows.Close()

	if rows.Next() {
		s.T().Fatal("expected no rows after rollback")
	}
	require.NoError(s.T(), rows.Err())
}

// TestSQLDB_ErrorHandling verifies that errors are properly wrapped when
// executing an invalid query.
func (s *SQLDBIntTestSuite) Test_ErrorHandling() {

	_, err := s.db.Query(`SELECT * FROM non_existing_table`)
	require.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), "DB.Query error")
}

// TestSQLDB_ConnectionConfig verifies that connection configuration
// functions execute without errors.
func (s *SQLDBIntTestSuite) Test_ConnectionConfig() {

	s.db.SetConnMaxLifetime(5 * time.Minute)
	s.db.SetConnMaxIdleTime(2 * time.Minute)
	s.db.SetMaxOpenConns(10)
	s.db.SetMaxIdleConns(5)

	err := s.db.Ping()
	require.NoError(s.T(), err)
}
