package database

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// SQLDBTestSuite is a suite of tests for SQLDB-related tests.
type SQLDBTestSuite struct {
	suite.Suite
}

// TestSQLDBTestSuite runs the test suite.
func TestSQLDBTestSuite(t *testing.T) {
	suite.Run(t, new(SQLDBTestSuite))
}

// TestSQLDB_Ping verifies that Ping works.
func (s *SQLDBTestSuite) TestSQLDB_Ping() {
	// Create a new sqlmock database.
	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	require.NoError(s.T(), err)
	defer db.Close()

	// Wrap it with our SQLDB implementation.
	sqlDB := &sqlDB{DB: db}

	// Expect a Ping call.
	mock.ExpectPing().WillReturnError(nil)

	// Call Ping.
	err = sqlDB.Ping()
	require.NoError(s.T(), err)
	require.NoError(s.T(), mock.ExpectationsWereMet())
}

// Test_PrepareAndExec verifies that Prepare returns a RealStmt and that Exec
// works.
func (s *SQLDBTestSuite) Test_PrepareAndExec() {
	// This test verifies that Prepare returns a RealStmt and that Exec works.
	db, mock, err := sqlmock.New()
	require.NoError(s.T(), err)
	defer db.Close()
	sqlDB := &sqlDB{DB: db}

	query := "SELECT 1"
	// Expect a Prepare call.
	mock.ExpectPrepare(query)

	prep, err := sqlDB.Prepare(query)
	require.NoError(s.T(), err)
	require.NotNil(s.T(), prep)

	// For the prepared statement, expect an Exec call.
	mock.ExpectExec(query).WithArgs().WillReturnResult(sqlmock.NewResult(1, 1))
	res, err := prep.Exec()
	require.NoError(s.T(), err)
	id, err := res.LastInsertId()
	require.NoError(s.T(), err)
	assert.Equal(s.T(), int64(1), id)

	// Close the prepared statement.
	err = prep.Close()
	require.NoError(s.T(), err)
	require.NoError(s.T(), mock.ExpectationsWereMet())
}

// Test_Exec verifies that Exec works.
func (s *SQLDBTestSuite) Test_Exec() {
	// Test the Exec method of SQLDB.
	db, mock, err := sqlmock.New()
	require.NoError(s.T(), err)
	defer db.Close()
	sqlDB := &sqlDB{DB: db}

	query := "UPDATE test SET name = ? WHERE id = ?"
	// Use regexp.QuoteMeta to escape special regex characters in the query.
	expectedQuery := regexp.QuoteMeta(query)
	mock.ExpectExec(expectedQuery).WithArgs("new", 123).
		WillReturnResult(sqlmock.NewResult(0, 1))

	res, err := sqlDB.Exec(query, "new", 123)
	require.NoError(s.T(), err)
	affected, err := res.RowsAffected()
	require.NoError(s.T(), err)
	assert.Equal(s.T(), int64(1), affected)
	require.NoError(s.T(), mock.ExpectationsWereMet())
}

// Test_Query verifies that Query and QueryRow work.
func (s *SQLDBTestSuite) Test_QueryAndQueryRow() {
	// Test the Query method of SQLDB.
	db, mock, err := sqlmock.New()
	require.NoError(s.T(), err)
	defer db.Close()
	sqlDB := &sqlDB{DB: db}

	// Test Query.
	query := "SELECT id, name FROM test WHERE id = ?"
	rows := sqlmock.NewRows([]string{"id", "name"}).AddRow(123, "Alice")
	mock.ExpectQuery(query).WithArgs(123).WillReturnRows(rows)

	resultRows, err := sqlDB.Query(query, 123)
	require.NoError(s.T(), err)
	defer resultRows.Close()

	require.True(s.T(), resultRows.Next())
	var id int
	var name string
	err = resultRows.Scan(&id, &name)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), 123, id)
	assert.Equal(s.T(), "Alice", name)
	require.NoError(s.T(), mock.ExpectationsWereMet())

	// Test QueryRow.
	mock.ExpectQuery(query).WithArgs(123).WillReturnRows(
		sqlmock.NewRows([]string{"id", "name"}).AddRow(123, "Alice"),
	)

	row := sqlDB.QueryRow(query, 123)
	err = row.Scan(&id, &name)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), 123, id)
	assert.Equal(s.T(), "Alice", name)

	// Ensure all expectations were met
	require.NoError(s.T(), mock.ExpectationsWereMet())
}

// Test_BeginTxAndCommit tests starting a transaction and committing.
func (s *SQLDBTestSuite) Test_BeginTxAndCommit() {
	// Test starting a transaction and committing.
	db, mock, err := sqlmock.New()
	require.NoError(s.T(), err)
	defer db.Close()
	sqlDB := &sqlDB{DB: db}

	ctx := context.Background()
	opts := &sql.TxOptions{}
	mock.ExpectBegin()

	tx, err := sqlDB.BeginTx(ctx, opts)
	require.NoError(s.T(), err)
	require.NotNil(s.T(), tx)

	// Expect a commit on the transaction.
	mock.ExpectCommit()
	err = tx.Commit()
	require.NoError(s.T(), err)
	require.NoError(s.T(), mock.ExpectationsWereMet())
}

// Test_BeginTxAndRollback tests starting a transaction and rolling back.
func (s *SQLDBTestSuite) Test_BeginTxAndRollback() {
	// Test starting a transaction and rolling back.
	db, mock, err := sqlmock.New()
	require.NoError(s.T(), err)
	defer db.Close()
	sqlDB := &sqlDB{DB: db}

	ctx := context.Background()
	opts := &sql.TxOptions{}
	mock.ExpectBegin()

	tx, err := sqlDB.BeginTx(ctx, opts)
	require.NoError(s.T(), err)
	require.NotNil(s.T(), tx)

	// Expect a rollback on the transaction.
	mock.ExpectRollback()
	err = tx.Rollback()
	require.NoError(s.T(), err)
	require.NoError(s.T(), mock.ExpectationsWereMet())
}

// Test_Close verifies that Close works.
func (s *SQLDBTestSuite) Test_Close() {
	db, mock, err := sqlmock.New()
	require.NoError(s.T(), err)
	sqlDB := &sqlDB{DB: db}

	// Set expectation that Close will be called.
	mock.ExpectClose()

	err = sqlDB.Close()
	require.NoError(s.T(), err)
}

// Test_Prepare_Error verifies that Prepare wraps errors correctly.
func (s *SQLDBTestSuite) Test_Prepare_Error() {
	db, mock, err := sqlmock.New()
	require.NoError(s.T(), err)
	defer db.Close()
	sqlDB := &sqlDB{DB: db}

	query := "INVALID QUERY"
	mock.ExpectPrepare(query).WillReturnError(errors.New("prepare error"))

	_, err = sqlDB.Prepare(query)
	require.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), "DB.Prepare error")
	require.NoError(s.T(), mock.ExpectationsWereMet())
}
