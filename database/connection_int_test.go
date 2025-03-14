// connection_int_test.go
package database_test

import (
	"testing"
	"time"

	"github.com/pureapi/pureapi-core/database"
	"github.com/pureapi/pureapi-core/database/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	_ "github.com/mattn/go-sqlite3"
)

// ConnectionIntTestSuite is a suite of integration tests for database connection.
type ConnectionIntTestSuite struct {
	suite.Suite
	db  types.DB
	cfg database.ConnectConfig
}

// SetupTest initializes the test by creating an in-memory SQLite3 database.
func (s *ConnectionIntTestSuite) SetupTest() {
	s.cfg = database.ConnectConfig{
		Driver:          "sqlite3",
		ConnMaxLifetime: 5 * time.Minute,
		ConnMaxIdleTime: 2 * time.Minute,
		MaxOpenConns:    10,
		MaxIdleConns:    5,
	}
	var err error
	s.db, err = database.Connect(s.cfg, database.NewSQLDBAdapter, ":memory:")
	require.NoError(s.T(), err)
	require.NotNil(s.T(), s.db)
}

// TearDownTest closes the database connection.
func (s *ConnectionIntTestSuite) TearDownTest() {
	if s.db != nil {
		s.db.Close()
	}
}

// Test_Ping verifies that the connection can be pinged.
func (s *ConnectionIntTestSuite) Test_Ping() {
	err := s.db.Ping()
	require.NoError(s.T(), err)
}

// Test_Exec_Query verifies that SQL execution works by creating a table,
// inserting a row, and querying that row.
func (s *ConnectionIntTestSuite) Test_Exec_Query() {
	createStmt := `
		CREATE TABLE IF NOT EXISTS test_conn (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL
		);
	`
	_, err := s.db.Exec(createStmt)
	require.NoError(s.T(), err)

	res, err := s.db.Exec(`INSERT INTO test_conn (name) VALUES (?)`, "Bob")
	require.NoError(s.T(), err)
	lastID, err := res.LastInsertId()
	require.NoError(s.T(), err)
	require.True(s.T(), lastID > 0)

	rows, err := s.db.Query(`SELECT name FROM test_conn WHERE id = ?`, lastID)
	require.NoError(s.T(), err)
	defer rows.Close()

	var name string
	if rows.Next() {
		err = rows.Scan(&name)
		require.NoError(s.T(), err)
		require.Equal(s.T(), "Bob", name)
	} else {
		s.T().Fatal("expected a row")
	}
}

// Test_ConnectionConfig verifies that updating connection settings does not break the connection.
func (s *ConnectionIntTestSuite) Test_ConnectionConfig() {
	s.db.SetConnMaxLifetime(10 * time.Minute)
	s.db.SetConnMaxIdleTime(5 * time.Minute)
	s.db.SetMaxOpenConns(20)
	s.db.SetMaxIdleConns(10)

	err := s.db.Ping()
	require.NoError(s.T(), err)
}

func TestConnectionIntTestSuite(t *testing.T) {
	suite.Run(t, new(ConnectionIntTestSuite))
}
