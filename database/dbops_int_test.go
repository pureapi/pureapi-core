package database

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/pureapi/pureapi-core/database/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	_ "github.com/mattn/go-sqlite3"
)

// TestEntity is a simple type used for scanning rows into an entity.
// It implements the types.Getter interface.
type TestEntity struct {
	ID   int
	Name string
}

func (te *TestEntity) TableName() string {
	return "test_entities"
}

func (te *TestEntity) ScanRow(row types.Row) error {
	return row.Scan(&te.ID, &te.Name)
}

// DBOpsIntTestSuite defines the suite for dbops-related integration tests.
type DBOpsIntTestSuite struct {
	suite.Suite
	ctx context.Context
	db  types.DB
}

// TestDBOpsIntTestSuite runs the test suite.
func TestDBOpsIntTestSuite(t *testing.T) {
	suite.Run(t, new(DBOpsIntTestSuite))
}

// SetupTest opens an in-memory SQLite3 connection.
func (s *DBOpsIntTestSuite) SetupTest() {
	s.ctx = context.Background()
	cfg := ConnectConfig{
		Driver:          "sqlite3",
		ConnMaxLifetime: 5 * time.Minute,
		ConnMaxIdleTime: 2 * time.Minute,
		MaxOpenConns:    10,
		MaxIdleConns:    5,
	}
	var err error
	s.db, err = Connect(
		cfg, NewSQLDBAdapter, "file::memory:?cache=shared",
	)
	require.NoError(s.T(), err)
	require.NotNil(s.T(), s.db)
}

// TearDownTest closes the database connection.
func (s *DBOpsIntTestSuite) TearDownTest() {
	if s.db != nil {
		s.db.Close()
	}
}

// Test_Exec uses Exec to create a table and insert a row.
func (s *DBOpsIntTestSuite) Test_Exec() {
	createStmt := `
		CREATE TABLE test_exec (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			value INTEGER NOT NULL
		);
	`
	_, err := s.db.Exec(createStmt)
	require.NoError(s.T(), err)

	insertQuery := "INSERT INTO test_exec (value) VALUES (?)"
	result, err := Exec(s.ctx, s.db, insertQuery, []any{123}, nil)
	require.NoError(s.T(), err)
	lastID, err := result.LastInsertId()
	require.NoError(s.T(), err)
	require.True(s.T(), lastID > 0)

	// Verify insertion
	query := "SELECT value FROM test_exec WHERE id = ?"
	rows, err := s.db.Query(query, lastID)
	require.NoError(s.T(), err)
	defer rows.Close()

	var value int
	if rows.Next() {
		err = rows.Scan(&value)
		require.NoError(s.T(), err)
		require.Equal(s.T(), 123, value)
	} else {
		s.T().Fatal("expected a row")
	}
}

// Test_Query uses Query to select data.
func (s *DBOpsIntTestSuite) Test_Query() {
	createStmt := `
		CREATE TABLE test_query (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			text_val TEXT NOT NULL
		);
	`
	_, err := s.db.Exec(createStmt)
	require.NoError(s.T(), err)

	insertQuery := "INSERT INTO test_query (text_val) VALUES (?)"
	res, err := s.db.Exec(insertQuery, "hello world")
	require.NoError(s.T(), err)
	lastID, err := res.LastInsertId()
	require.NoError(s.T(), err)

	rows, stmt, err := Query(
		s.ctx,
		s.db,
		"SELECT text_val FROM test_query WHERE id = ?",
		[]any{lastID},
		nil,
	)
	require.NoError(s.T(), err)
	defer stmt.Close()
	defer rows.Close()

	var text string
	if rows.Next() {
		err = rows.Scan(&text)
		require.NoError(s.T(), err)
		require.Equal(s.T(), "hello world", text)
	} else {
		s.T().Fatal("expected a row")
	}
}

// Test_ExecRaw uses ExecRaw to insert data.
func (s *DBOpsIntTestSuite) Test_ExecRaw() {
	createStmt := `
		CREATE TABLE test_execraw (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			num_val INTEGER NOT NULL
		);
	`
	_, err := s.db.Exec(createStmt)
	require.NoError(s.T(), err)

	insertQuery := "INSERT INTO test_execraw (num_val) VALUES (?)"
	result, err := ExecRaw(s.ctx, s.db, insertQuery, []any{456}, nil)
	require.NoError(s.T(), err)
	lastID, err := result.LastInsertId()
	require.NoError(s.T(), err)
	require.True(s.T(), lastID > 0)

	// Verify insertion
	query := "SELECT num_val FROM test_execraw WHERE id = ?"
	rows, err := s.db.Query(query, lastID)
	require.NoError(s.T(), err)
	defer rows.Close()

	var num int
	if rows.Next() {
		err = rows.Scan(&num)
		require.NoError(s.T(), err)
		require.Equal(s.T(), 456, num)
	} else {
		s.T().Fatal("expected a row")
	}
}

// Test_QueryRaw uses QueryRaw to select data.
func (s *DBOpsIntTestSuite) Test_QueryRaw() {
	createStmt := `
		CREATE TABLE test_queryraw (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			str_val TEXT NOT NULL
		);
	`
	_, err := s.db.Exec(createStmt)
	require.NoError(s.T(), err)

	insertQuery := "INSERT INTO test_queryraw (str_val) VALUES (?)"
	res, err := s.db.Exec(insertQuery, "raw query test")
	require.NoError(s.T(), err)
	lastID, err := res.LastInsertId()
	require.NoError(s.T(), err)

	rows, err := QueryRaw(
		s.ctx,
		s.db,
		"SELECT str_val FROM test_queryraw WHERE id = ?",
		[]any{lastID},
		nil,
	)
	require.NoError(s.T(), err)
	defer rows.Close()

	var sVal string
	if rows.Next() {
		err = rows.Scan(&sVal)
		require.NoError(s.T(), err)
		require.Equal(s.T(), "raw query test", sVal)
	} else {
		s.T().Fatal("expected a row")
	}
}

// Test_QuerySingleValue uses QuerySingleValue to fetch a scalar.
func (s *DBOpsIntTestSuite) Test_QuerySingleValue() {
	createStmt := `
		CREATE TABLE test_single_value (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			count_val INTEGER NOT NULL
		);
	`
	_, err := s.db.Exec(createStmt)
	require.NoError(s.T(), err)

	insertQuery := "INSERT INTO test_single_value (count_val) VALUES (?)"
	res, err := s.db.Exec(insertQuery, 789)
	require.NoError(s.T(), err)
	lastID, err := res.LastInsertId()
	require.NoError(s.T(), err)

	result, err := QuerySingleValue(
		s.ctx,
		s.db,
		"SELECT count_val FROM test_single_value WHERE id = ?",
		[]any{lastID},
		nil,
		func() *int { return new(int) },
	)
	require.NoError(s.T(), err)
	require.Equal(s.T(), 789, *result)
}

// User is a simple entity used for testing QuerySingleEntity and QueryEntities.
type User struct {
	ID   int
	Name string
}

func (u *User) TableName() string {
	return "users"
}

func (u *User) ScanRow(row types.Row) error {
	return row.Scan(&u.ID, &u.Name)
}

// Test_QuerySingleEntity uses QuerySingleEntity to fetch an entity.
func (s *DBOpsIntTestSuite) Test_QuerySingleEntity() {
	createStmt := `
		CREATE TABLE users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL
		);
	`
	_, err := s.db.Exec(createStmt)
	require.NoError(s.T(), err)

	insertQuery := "INSERT INTO users (name) VALUES (?)"
	res, err := s.db.Exec(insertQuery, "Alice")
	require.NoError(s.T(), err)
	lastID, err := res.LastInsertId()
	require.NoError(s.T(), err)

	user, err := QuerySingleEntity(
		s.ctx,
		s.db,
		"SELECT id, name FROM users WHERE id = ?",
		[]any{lastID},
		nil,
		func() *User { return new(User) },
	)
	require.NoError(s.T(), err)
	require.Equal(s.T(), "Alice", user.Name)
}

// Test_QueryEntities uses QueryEntities to fetch multiple entities.
func (s *DBOpsIntTestSuite) Test_QueryEntities() {
	createStmt := `
		CREATE TABLE users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL
		);
	`
	_, err := s.db.Exec(createStmt)
	require.NoError(s.T(), err)

	// Insert multiple users.
	names := []string{"Alice", "Bob", "Charlie"}
	for _, name := range names {
		_, err = s.db.Exec("INSERT INTO users (name) VALUES (?)", name)
		require.NoError(s.T(), err)
	}

	users, err := QueryEntities(
		s.ctx,
		s.db,
		"SELECT id, name FROM users",
		nil,
		nil,
		func() *User { return new(User) },
	)
	require.NoError(s.T(), err)
	require.Len(s.T(), users, len(names))

	// Check that each returned user has a valid ID and matching name.
	nameMap := make(map[string]bool)
	for _, u := range users {
		require.True(s.T(), u.ID > 0)
		nameMap[u.Name] = true
	}
	for _, name := range names {
		require.True(
			s.T(), nameMap[name], fmt.Sprintf("expected user %s", name),
		)
	}
}

// Test_RowToEntity verifies that RowToEntity correctly scans a single row into
// an entity.
func (s *DBOpsIntTestSuite) Test_RowToEntity() {
	// Create a table and insert a row.
	createStmt := `
		CREATE TABLE test_entity (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL
		);
	`
	_, err := s.db.Exec(createStmt)
	require.NoError(s.T(), err)

	insertQuery := "INSERT INTO test_entity (name) VALUES (?)"
	res, err := s.db.Exec(insertQuery, "TestName")
	require.NoError(s.T(), err)
	lastID, err := res.LastInsertId()
	require.NoError(s.T(), err)
	require.True(s.T(), lastID > 0)

	row := s.db.QueryRow(
		"SELECT id, name FROM test_entity WHERE id = ?", lastID,
	)
	entity, err := RowToEntity(
		s.ctx, row, func() *TestEntity { return new(TestEntity) },
	)
	require.NoError(s.T(), err)
	require.Equal(s.T(), "TestName", entity.Name)
	require.Equal(s.T(), int(lastID), entity.ID)
}

// Test_RowToAny verifies that RowToAny correctly scans a single scalar value.
func (s *DBOpsIntTestSuite) Test_RowToAny() {
	// Create a table for a scalar value.
	createStmt := `
		CREATE TABLE test_scalar (
			val INTEGER NOT NULL
		);
	`
	_, err := s.db.Exec(createStmt)
	require.NoError(s.T(), err)

	insertQuery := "INSERT INTO test_scalar (val) VALUES (?)"
	res, err := s.db.Exec(insertQuery, 777)
	require.NoError(s.T(), err)
	lastID, err := res.LastInsertId()
	require.NoError(s.T(), err)
	require.True(s.T(), lastID > 0)

	row := s.db.QueryRow("SELECT val FROM test_scalar WHERE rowid = ?", lastID)
	result, err := RowToAny(s.ctx, row, func() *int { return new(int) })
	require.NoError(s.T(), err)
	require.Equal(s.T(), 777, *result)
}

// Test_RowsToAny verifies that RowsToAny correctly scans multiple scalar rows.
func (s *DBOpsIntTestSuite) Test_RowsToAny() {
	// Create a table and insert several scalar values.
	createStmt := `
		CREATE TABLE test_scalars (
			val INTEGER NOT NULL
		);
	`
	_, err := s.db.Exec(createStmt)
	require.NoError(s.T(), err)

	expected := []int{1, 2, 3, 4, 5}
	for _, v := range expected {
		_, err = s.db.Exec("INSERT INTO test_scalars (val) VALUES (?)", v)
		require.NoError(s.T(), err)
	}

	rows, err := s.db.Query("SELECT val FROM test_scalars ORDER BY val ASC")
	require.NoError(s.T(), err)
	defer rows.Close()

	results, err := RowsToAny(s.ctx, rows, func() *int { return new(int) })
	require.NoError(s.T(), err)
	require.Len(s.T(), results, len(expected))
	for i, r := range results {
		require.Equal(s.T(), expected[i], *r)
	}
}

// Test_RowsToEntities verifies that RowsToEntities correctly scans multiple
// entities.
func (s *DBOpsIntTestSuite) Test_RowsToEntities() {
	// Create a table for entities.
	createStmt := `
		CREATE TABLE test_entities (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL
		);
	`
	_, err := s.db.Exec(createStmt)
	require.NoError(s.T(), err)

	names := []string{"Alice", "Bob", "Charlie"}
	for _, name := range names {
		_, err = s.db.Exec("INSERT INTO test_entities (name) VALUES (?)", name)
		require.NoError(s.T(), err)
	}

	rows, err := s.db.Query(
		"SELECT id, name FROM test_entities ORDER BY id ASC",
	)
	require.NoError(s.T(), err)
	defer rows.Close()

	entities, err := RowsToEntities(
		s.ctx, rows, func() *TestEntity { return new(TestEntity) },
	)
	require.NoError(s.T(), err)
	require.Len(s.T(), entities, len(names))
	for i, e := range entities {
		require.Equal(s.T(), names[i], e.Name)
		require.True(s.T(), e.ID > 0)
	}
}
