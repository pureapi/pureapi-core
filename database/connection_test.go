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

// FakeDB is a fake implementation of the DB interface.
type FakeDB struct {
	DriverName      string
	DSN             string
	pingErr         error
	connMaxLifetime time.Duration
	connMaxIdleTime time.Duration
	maxOpenConns    int
	maxIdleConns    int
}

func NewFakeDB(driver, dsn string) *FakeDB {
	return &FakeDB{
		DriverName: driver,
		DSN:        dsn,
	}
}

// Ping returns the configured ping error (nil if no error).
func (f *FakeDB) Ping() error {
	return f.pingErr
}

func (f *FakeDB) SetConnMaxLifetime(d time.Duration) {
	f.connMaxLifetime = d
}

func (f *FakeDB) SetConnMaxIdleTime(d time.Duration) {
	f.connMaxIdleTime = d
}

func (f *FakeDB) SetMaxOpenConns(n int) {
	f.maxOpenConns = n
}

func (f *FakeDB) SetMaxIdleConns(n int) {
	f.maxIdleConns = n
}

func (f *FakeDB) BeginTx(
	ctx context.Context, options *sql.TxOptions,
) (types.Tx, error) {
	panic("not implemented")
}
func (f *FakeDB) Exec(query string, args ...any) (types.Result, error) {
	panic("not implemented")
}
func (f *FakeDB) Query(query string, args ...any) (types.Rows, error) {
	panic("not implemented")
}
func (f *FakeDB) Close() error {
	return nil
}
func (f *FakeDB) Prepare(query string) (types.Stmt, error) {
	panic("not implemented")
}

func fakeConnOpenFn(driver string, dsn string) (types.DB, error) {
	return NewFakeDB(driver, dsn), nil
}

func fakeConnOpenFnError(driver string, dsn string) (types.DB, error) {
	return nil, errors.New("fake connOpenFn error")
}

// ConnectionTestSuite is a suite of tests for connection-related tests.
type ConnectionTestSuite struct {
	suite.Suite
	cfg ConnectConfig
}

// TestConnectionSuite registers the test suite.
func TestConnectionSuite(t *testing.T) {
	suite.Run(t, new(ConnectionTestSuite))
}

// SetupTest initializes common test configurations.
func (s *ConnectionTestSuite) SetupTest() {
	s.cfg = ConnectConfig{
		User:            "user",
		Password:        "pass",
		Host:            "localhost",
		Port:            3306,
		Database:        "mydb",
		Parameters:      "charset=utf8",
		ConnMaxLifetime: 30 * time.Minute,
		ConnMaxIdleTime: 10 * time.Minute,
		MaxOpenConns:    100,
		MaxIdleConns:    10,
	}
}

// TestConnect_Success verifies that Connect creates a database connection.
func (s *ConnectionTestSuite) Test_Success() {
	// Use fakeConnOpenFn to create a FakeDB.
	dsn := "dummy"
	db, err := Connect(s.cfg, fakeConnOpenFn, dsn)
	require.NoError(s.T(), err)
	require.NotNil(s.T(), db)

	// Type assert to FakeDB to verify that configuration was applied.
	fake, ok := db.(*FakeDB)
	require.True(s.T(), ok, "expected FakeDB type")

	// Check that the DSN passed to fakeDB is what we generated.
	assert.Equal(s.T(), dsn, fake.DSN)
	// Check that connection configuration was applied.
	assert.Equal(s.T(), s.cfg.ConnMaxLifetime, fake.connMaxLifetime)
	assert.Equal(s.T(), s.cfg.ConnMaxIdleTime, fake.connMaxIdleTime)
	assert.Equal(s.T(), s.cfg.MaxOpenConns, fake.maxOpenConns)
	assert.Equal(s.T(), s.cfg.MaxIdleConns, fake.maxIdleConns)
}

// TestConnect_connOpenFnError verifies that Connect returns an error if the
// connOpenFn returns an error.
func (s *ConnectionTestSuite) Test_connOpenFnError() {
	cfg := ConnectConfig{}
	dsn := "dummy"
	// Use a factory that returns an error.
	db, err := Connect(cfg, fakeConnOpenFnError, dsn)
	require.Error(s.T(), err)
	assert.Nil(s.T(), db)
}

// TestConnect_PingError verifies that Connect returns an error if the
// database fails to ping.
func (s *ConnectionTestSuite) Test_PingError() {
	dsn := "dummy"
	// Use a factory that returns a FakeDB with a ping error.
	db, err := Connect(
		s.cfg,
		func(driver string, dsn string) (types.DB, error) {
			f := NewFakeDB(driver, dsn)
			f.pingErr = errors.New("ping failed")
			return f, nil
		},
		dsn,
	)
	require.Error(s.T(), err)
	assert.Nil(s.T(), db)
}
