package database

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/pureapi/pureapi-core/database/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// fakePreparer implements types.Preparer.
type fakePreparer struct {
	prepareFunc func(query string) (types.Stmt, error)
}

func (fp *fakePreparer) Prepare(query string) (types.Stmt, error) {
	if fp.prepareFunc != nil {
		return fp.prepareFunc(query)
	}
	return nil, nil
}

// fakeStmt implements types.Stmt.
type fakeStmt struct {
	execFunc     func(args ...any) (types.Result, error)
	queryFunc    func(args ...any) (types.Rows, error)
	queryRowFunc func(args ...any) types.Row
	closeFunc    func() error
}

func (fs *fakeStmt) Exec(args ...any) (types.Result, error) {
	if fs.execFunc != nil {
		return fs.execFunc(args...)
	}
	return nil, nil
}

func (fs *fakeStmt) Query(args ...any) (types.Rows, error) {
	if fs.queryFunc != nil {
		return fs.queryFunc(args...)
	}
	return nil, nil
}

func (fs *fakeStmt) QueryRow(args ...any) types.Row {
	if fs.queryRowFunc != nil {
		return fs.queryRowFunc(args...)
	}
	return nil
}

func (fs *fakeStmt) Close() error {
	if fs.closeFunc != nil {
		return fs.closeFunc()
	}
	return nil
}

// fakeResult implements types.Result.
type fakeResult struct {
	lastInsertID int64
	rowsAffected int64
}

func (fr *fakeResult) LastInsertId() (int64, error) {
	return fr.lastInsertID, nil
}

func (fr *fakeResult) RowsAffected() (int64, error) {
	return fr.rowsAffected, nil
}

// fakeRows implements types.Rows.
type fakeRows struct {
	current   int
	total     int
	scanFunc  func(dest ...any) error
	returnErr error
}

func (fr *fakeRows) Next() bool {
	if fr.current < fr.total {
		fr.current++
		return true
	}
	return false
}

func (fr *fakeRows) Scan(dest ...any) error {
	if fr.scanFunc != nil {
		return fr.scanFunc(dest...)
	}
	return nil
}

func (fr *fakeRows) Close() error {
	return nil
}

func (fr *fakeRows) Err() error {
	return fr.returnErr
}

// fakeRow implements types.Row.
type fakeRow struct {
	scanFunc func(dest ...any) error
	err      error
}

func (fr *fakeRow) Scan(dest ...any) error {
	if fr.scanFunc != nil {
		return fr.scanFunc(dest...)
	}
	return nil
}

func (fr *fakeRow) Err() error {
	return fr.err
}

// fakeDB implements types.DB for ExecRaw and QueryRaw.
type fakeDB struct {
	execFunc  func(query string, args ...any) (types.Result, error)
	queryFunc func(query string, args ...any) (types.Rows, error)
	closeFunc func() error
}

func (fdb *fakeDB) Exec(query string, args ...any) (types.Result, error) {
	if fdb.execFunc != nil {
		return fdb.execFunc(query, args...)
	}
	return nil, nil
}

func (fdb *fakeDB) Query(query string, args ...any) (types.Rows, error) {
	if fdb.queryFunc != nil {
		return fdb.queryFunc(query, args...)
	}
	return nil, nil
}

func (fdb *fakeDB) QueryRow(query string, args ...any) types.Row {
	// Not implemented
	return nil
}

func (fdb *fakeDB) Close() error {
	if fdb.closeFunc != nil {
		return fdb.closeFunc()
	}
	return nil
}

func (fdb *fakeDB) Ping() error                        { return nil }
func (fdb *fakeDB) SetConnMaxLifetime(d time.Duration) {}
func (fdb *fakeDB) SetConnMaxIdleTime(d time.Duration) {}
func (fdb *fakeDB) SetMaxOpenConns(n int)              {}
func (fdb *fakeDB) SetMaxIdleConns(n int)              {}
func (fdb *fakeDB) Prepare(query string) (types.Stmt, error) {
	return nil, nil
}
func (fdb *fakeDB) BeginTx(ctx context.Context,
	opts *sql.TxOptions) (types.Tx, error) {
	return nil, errors.New("not implemented")
}

// fakeErrorChecker implements types.ErrorChecker.
type fakeErrorChecker struct {
	prefix string
}

func (fec *fakeErrorChecker) Check(err error) error {
	if err == nil {
		return nil
	}
	return errors.New(fec.prefix + err.Error())
}

// fakeEntity is used for testing QuerySingleEntity.
type fakeEntity struct {
	Value int
}

func (fe *fakeEntity) TableName() string {
	return "fake"
}

func (fe *fakeEntity) ScanRow(row types.Row) error {
	return row.Scan(&fe.Value)
}

// ConnectionTestSuite is a suite of tests for dbops-related tests.
type DBOpsTestSuite struct {
	suite.Suite
	ctx          context.Context
	errorChecker types.ErrorChecker
}

// TestDBOpsTestSuite runs the test suite.
func TestDBOpsTestSuite(t *testing.T) {
	suite.Run(t, new(DBOpsTestSuite))
}

// SetupTest sets up the test suite.
func (s *DBOpsTestSuite) SetupTest() {
	s.ctx = context.Background()
	s.errorChecker = &fakeErrorChecker{prefix: "checked: "}
}

// TestExec_NilPreparer tests that Exec returns an error if the preparer is nil.
func (s *DBOpsTestSuite) TestExec_NilPreparer() {
	result, err := Exec(s.ctx, nil, "SELECT 1", nil, nil)
	require.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), "Exec: preparer is nil")
	assert.Nil(s.T(), result)
}

// TestExec_Success tests that Exec returns the result if the query is
// successful.
func (s *DBOpsTestSuite) TestExec_Success() {
	fakeStmt := &fakeStmt{
		execFunc: func(args ...any) (types.Result, error) {
			return &fakeResult{lastInsertID: 100, rowsAffected: 1}, nil
		},
		closeFunc: func() error { return nil },
	}
	fakePrep := &fakePreparer{
		prepareFunc: func(query string) (types.Stmt, error) {
			return fakeStmt, nil
		},
	}
	result, err := Exec(
		s.ctx, fakePrep, "UPDATE table SET col=?", []any{"value"}, nil,
	)
	require.NoError(s.T(), err)
	res, ok := result.(*fakeResult)
	require.True(s.T(), ok)
	assert.Equal(s.T(), int64(100), res.lastInsertID)
}

// TestExec_ErrorChecker tests that Exec returns an error if the error checker
// returns an error.
func (s *DBOpsTestSuite) TestExec_ErrorChecker() {
	fakeErr := errors.New("exec error")
	fakeStmt := &fakeStmt{
		execFunc: func(args ...any) (types.Result, error) {
			return nil, fakeErr
		},
		closeFunc: func() error { return nil },
	}
	fakePrep := &fakePreparer{
		prepareFunc: func(query string) (types.Stmt, error) {
			return fakeStmt, nil
		},
	}
	result, err := Exec(
		s.ctx, fakePrep, "DELETE FROM table", nil, s.errorChecker,
	)
	require.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), "checked: exec error")
	assert.Nil(s.T(), result)
}

// TestQuery_NilPreparer tests that Query returns an error if the preparer is
// nil.
func (s *DBOpsTestSuite) TestQuery_NilPreparer() {
	rows, stmt, err := Query(
		s.ctx, nil, "SELECT * FROM table", nil, nil,
	)
	require.Error(s.T(), err)
	assert.Nil(s.T(), rows)
	assert.Nil(s.T(), stmt)
	assert.Contains(s.T(), err.Error(), "Query: preparer is nil")
}

// TestQuery_Success tests that Query returns the rows and stmt if the query is
// successful.
func (s *DBOpsTestSuite) TestQuery_Success() {
	fakeRowsObj := &fakeRows{
		total:   1,
		current: 0,
		scanFunc: func(dest ...any) error {
			if len(dest) > 0 {
				if ptr, ok := dest[0].(*int); ok {
					*ptr = 42
				}
			}
			return nil
		},
	}
	fakeStmt := &fakeStmt{
		queryFunc: func(args ...any) (types.Rows, error) {
			return fakeRowsObj, nil
		},
		closeFunc: func() error { return nil },
	}
	fakePrep := &fakePreparer{
		prepareFunc: func(query string) (types.Stmt, error) {
			return fakeStmt, nil
		},
	}
	rows, stmt, err := Query(
		s.ctx, fakePrep, "SELECT col FROM table", nil, nil,
	)
	require.NoError(s.T(), err)
	assert.NotNil(s.T(), rows)
	assert.NotNil(s.T(), stmt)
	var val int
	if rows.Next() {
		err = rows.Scan(&val)
		require.NoError(s.T(), err)
		assert.Equal(s.T(), 42, val)
	}
	_ = stmt.Close()
	_ = rows.Close()
}

// TestExecRaw_NilDB tests that ExecRaw returns an error if the db is nil.
func (s *DBOpsTestSuite) TestExecRaw_NilDB() {
	result, err := ExecRaw(
		s.ctx, nil, "INSERT INTO table (col) VALUES(?)", []any{"value"}, nil,
	)
	require.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), "ExecRaw: db is nil")
	assert.Nil(s.T(), result)
}

// TestExecRaw_Success tests that ExecRaw returns the result if the query is
// successful.
func (s *DBOpsTestSuite) TestExecRaw_Success() {
	fakeDBObj := &fakeDB{
		execFunc: func(query string, args ...any) (types.Result, error) {
			return &fakeResult{lastInsertID: 200, rowsAffected: 1}, nil
		},
	}
	result, err := ExecRaw(
		s.ctx,
		fakeDBObj,
		"INSERT INTO table (col) VALUES(?)",
		[]any{"value"},
		nil,
	)
	require.NoError(s.T(), err)
	res, ok := result.(*fakeResult)
	require.True(s.T(), ok)
	assert.Equal(s.T(), int64(200), res.lastInsertID)
}

// TestQueryRaw_NilDB tests that QueryRaw returns an error if the db is nil.
func (s *DBOpsTestSuite) TestQueryRaw_NilDB() {
	rows, err := QueryRaw(s.ctx, nil, "SELECT * FROM table", nil, nil)
	require.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), "QueryRaw: db is nil")
	assert.Nil(s.T(), rows)
}

// TestQueryRaw_Success tests that QueryRaw returns the rows if the query is
// successful.
func (s *DBOpsTestSuite) TestQueryRaw_Success() {
	fakeRowsObj := &fakeRows{
		total:   1,
		current: 0,
		scanFunc: func(dest ...any) error {
			if len(dest) > 0 {
				if ptr, ok := dest[0].(*string); ok {
					*ptr = "hello"
				}
			}
			return nil
		},
	}
	fakeDBObj := &fakeDB{
		queryFunc: func(query string, args ...any) (types.Rows, error) {
			return fakeRowsObj, nil
		},
	}
	rows, err := QueryRaw(s.ctx, fakeDBObj, "SELECT col FROM table", nil, nil)
	require.NoError(s.T(), err)
	assert.NotNil(s.T(), rows)
	var sVal string
	if rows.Next() {
		err = rows.Scan(&sVal)
		require.NoError(s.T(), err)
		assert.Equal(s.T(), "hello", sVal)
	}
	_ = rows.Close()
}

// TestQuerySingleValue_NilPreparer tests that QuerySingleValue returns an error
// if the preparer is nil.
func (s *DBOpsTestSuite) TestQuerySingleValue_NilPreparer() {
	result, err := QuerySingleValue(
		s.ctx, nil, "SELECT 1", nil, nil, func() int { return 0 },
	)
	require.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), "QuerySingleValue: preparer is nil")
	assert.Equal(s.T(), 0, result)
}

// TestQuerySingleValue_Success tests that QuerySingleValue returns the value if
// the query is successful.
func (s *DBOpsTestSuite) TestQuerySingleValue_Success() {
	fakeRowObj := &fakeRow{
		scanFunc: func(dest ...any) error {
			if len(dest) > 0 {
				if ptr, ok := dest[0].(*int); ok {
					*ptr = 55
				}
			}
			return nil
		},
	}
	fakeStmt := &fakeStmt{
		queryRowFunc: func(args ...any) types.Row {
			return fakeRowObj
		},
		closeFunc: func() error { return nil },
	}
	fakePrep := &fakePreparer{
		prepareFunc: func(query string) (types.Stmt, error) {
			return fakeStmt, nil
		},
	}
	result, err := QuerySingleValue(
		s.ctx,
		fakePrep,
		"SELECT val",
		nil,
		nil,
		func() *int { return new(int) },
	)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), 55, *result)
}

// TestQuerySingleEntity_NilPreparer tests that QuerySingleEntity returns an
// error if the preparer is nil.
func (s *DBOpsTestSuite) TestQuerySingleEntity_NilPreparer() {
	_, err := QuerySingleEntity(
		s.ctx,
		nil,
		"SELECT val",
		nil,
		nil,
		func() *fakeEntity { return new(fakeEntity) },
	)
	require.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), "QuerySingleEntity: preparer is nil")
}

// TestQuerySingleEntity_Success tests that QuerySingleEntity returns the entity
// if the query is successful.
func (s *DBOpsTestSuite) TestQuerySingleEntity_Success() {
	fakeRowObj := &fakeRow{
		scanFunc: func(dest ...any) error {
			if len(dest) > 0 {
				if ptr, ok := dest[0].(*int); ok {
					*ptr = 77
				}
			}
			return nil

		},
	}
	fakeStmt := &fakeStmt{
		queryRowFunc: func(args ...any) types.Row {
			return fakeRowObj
		},
		closeFunc: func() error { return nil },
	}
	fakePrep := &fakePreparer{
		prepareFunc: func(query string) (types.Stmt, error) {
			return fakeStmt, nil
		},
	}
	entity, err := QuerySingleEntity(
		s.ctx,
		fakePrep,
		"SELECT val",
		nil,
		nil,
		func() *fakeEntity { return new(fakeEntity) },
	)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), 77, entity.Value)
}

// TestQueryEntities_NilPreparer tests that QueryEntities returns an error if
// the preparer is nil.
func (s *DBOpsTestSuite) TestQueryEntities_Success() {
	// Simulate two rows returning values 10 and 20.
	count := 0
	fakeRowsObj := &fakeRows{
		total:   2,
		current: 0,
		scanFunc: func(dest ...any) error {
			count++
			if len(dest) > 0 {
				if ptr, ok := dest[0].(*int); ok {
					if count == 1 {
						*ptr = 10
					} else {
						*ptr = 20
					}
				}
			}
			return nil
		},
	}
	fakeStmt := &fakeStmt{
		queryFunc: func(args ...any) (types.Rows, error) {
			return fakeRowsObj, nil
		},
		closeFunc: func() error { return nil },
	}
	fakePrep := &fakePreparer{
		prepareFunc: func(query string) (types.Stmt, error) {
			return fakeStmt, nil
		},
	}
	entities, err := QueryEntities(
		s.ctx,
		fakePrep,
		"SELECT val FROM table",
		nil,
		nil,
		func() *fakeEntity { return new(fakeEntity) },
	)
	require.NoError(s.T(), err)
	require.Len(s.T(), entities, 2)
	assert.Equal(s.T(), &fakeEntity{Value: 10}, entities[0])
	assert.Equal(s.T(), &fakeEntity{Value: 20}, entities[1])
}
