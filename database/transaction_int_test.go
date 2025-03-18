package database_test

import (
	"context"
	"errors"
	"testing"

	"github.com/pureapi/pureapi-core/database"
	"github.com/pureapi/pureapi-core/database/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	_ "github.com/mattn/go-sqlite3"
)

// TransactionIntTestSuite is a test suite for transaction-related integration
// tests.
type TransactionIntTestSuite struct {
	suite.Suite
	db types.DB
	tx types.Tx
}

// TestTransactionIntTestSuite runs the test suite.
func TestTransactionIntTestSuite(t *testing.T) {
	suite.Run(t, new(TransactionIntTestSuite))
}

// SetupTest initializes an in-memory SQLite database.
func (s *TransactionIntTestSuite) SetupTest() {
	db, err := database.NewSQLDBAdapter("sqlite3", "file::memory:?cache=shared")
	require.NoError(s.T(), err)
	s.db = db

	// Create the test table.
	createStmt := `
		CREATE TABLE IF NOT EXISTS Test_ (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			value TEXT NOT NULL
		);
	`
	_, err = s.db.Exec(createStmt)
	require.NoError(s.T(), err)

	// Begin a transaction.
	tx, err := s.db.BeginTx(context.Background(), nil)
	require.NoError(s.T(), err)
	s.tx = tx
}

// TearDownTest closes the database connection.
func (s *TransactionIntTestSuite) TearDownTest() {
	s.db.Close()
}

// Test_Success verifies that a successful transaction commits its changes and
// returns the expected result.
func (s *TransactionIntTestSuite) Test_Success() {
	// Execute Transaction with a txFn that inserts a row and returns its
	// last inserted id.
	lastID, err := database.Transaction(
		context.Background(),
		s.tx,
		func(ctx context.Context, tx types.Tx) (int64, error) {
			res, err := tx.Exec(
				`INSERT INTO Test_ (value) VALUES (?)`, "success",
			)
			if err != nil {
				return 0, err
			}
			return res.LastInsertId()
		},
	)
	require.NoError(s.T(), err)
	require.True(s.T(), lastID > 0, "expected a valid last insert id")

	// Verify that the row was committed.
	rows, err := s.db.Query(
		`SELECT value FROM Test_ WHERE id = ?`, lastID,
	)
	require.NoError(s.T(), err)
	defer rows.Close()

	var value string
	if rows.Next() {
		err = rows.Scan(&value)
		require.NoError(s.T(), err)
		assert.Equal(s.T(), "success", value)
	} else {
		s.T().Fatal("expected one row in the table")
	}
	require.NoError(s.T(), rows.Err())
}

// Test_TxFnError verifies that if the txFn returns an error, the transaction is
// rolled back and no changes persist.
func (s *TransactionIntTestSuite) Test_TxFnError() {
	// Execute Transaction that returns an error after inserting a row.
	_, err := database.Transaction(
		context.Background(),
		s.tx,
		func(ctx context.Context, tx types.Tx) (int, error) {
			_, err := tx.Exec(
				`INSERT INTO Test_ (value) VALUES (?)`, "fail",
			)
			if err != nil {
				return 0, err
			}
			return 0, errors.New("txFn error")
		},
	)
	require.Error(s.T(), err)

	// Verify that the row was not persisted.
	rows, err := s.db.Query(
		`SELECT COUNT(*) FROM Test_ WHERE value = ?`,
		"fail",
	)
	require.NoError(s.T(), err)
	defer rows.Close()

	var count int
	if rows.Next() {
		err = rows.Scan(&count)
		require.NoError(s.T(), err)
	}
	assert.Equal(
		s.T(), 0, count,
		"expected no rows since transaction should be rolled back",
	)
}

// Test_Panic verifies that if the txFn panics, the transaction is rolled back
// and the panic is re-propagated.
func (s *TransactionIntTestSuite) Test_Panic() {
	panicMsg := "panic in txFn"
	assert.PanicsWithValue(s.T(), panicMsg, func() {
		_, _ = database.Transaction(context.Background(), s.tx, func(
			ctx context.Context, tx types.Tx,
		) (int, error) {
			_, err := tx.Exec(
				`INSERT INTO Test_ (value) VALUES (?)`, "panic",
			)
			require.NoError(s.T(), err)
			panic(panicMsg)
		})
	})

	// Verify that the row was not persisted after the panic.
	rows, err := s.db.Query(
		`SELECT COUNT(*) FROM Test_ WHERE value = ?`, "panic",
	)
	require.NoError(s.T(), err)
	defer rows.Close()

	var count int
	if rows.Next() {
		err = rows.Scan(&count)
		require.NoError(s.T(), err)
	}
	assert.Equal(
		s.T(), 0, count,
		"expected no rows since transaction should be rolled back after panic",
	)
}
