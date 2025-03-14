package database

import (
	"context"
	"errors"
	"testing"

	"github.com/pureapi/pureapi-core/database/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// FakeTx is a fake transaction that implements the Tx interface.
type FakeTx struct {
	commitCalled   bool
	rollbackCalled bool
	commitErr      error
	rollbackErr    error
}

func (f *FakeTx) Commit() error {
	f.commitCalled = true
	return f.commitErr
}

func (f *FakeTx) Rollback() error {
	f.rollbackCalled = true
	return f.rollbackErr
}

// For our Transaction tests we don't need to implement Prepare and Exec.
func (f *FakeTx) Prepare(query string) (types.Stmt, error) {
	return nil, errors.New("not implemented")
}

func (f *FakeTx) Exec(query string, args ...any) (types.Result, error) {
	return nil, errors.New("not implemented")
}

// TransactionTestSuite is a test suite for transaction-related tests.
type TransactionTestSuite struct {
	suite.Suite
}

// TestTransactionTestSuite registers the test suite.
func TestTransactionTestSuite(t *testing.T) {
	suite.Run(t, new(TransactionTestSuite))
}

// TestTransaction_Success verifies that when txFn returns successfully,
// Transaction commits the transaction and returns the result.
func (s *TransactionTestSuite) TestTransaction_Success() {
	fakeTx := &FakeTx{}
	resultValue := 42
	txFn := func(ctx context.Context, tx types.Tx) (int, error) {
		// Successful transactional work.
		return resultValue, nil
	}
	res, err := Transaction(context.Background(), fakeTx, txFn)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), resultValue, res)
	assert.True(
		s.T(), fakeTx.commitCalled,
		"Commit should be called on success",
	)
	assert.False(
		s.T(), fakeTx.rollbackCalled,
		"Rollback should not be called on success",
	)
}

// TestTransaction_TxFnError verifies that if txFn returns an error,
// Transaction rolls back the transaction and returns the error.
func (s *TransactionTestSuite) TestTransaction_TxFnError() {
	fakeTx := &FakeTx{}
	txFnErr := errors.New("txFn error")
	txFn := func(ctx context.Context, tx types.Tx) (int, error) {
		return 0, txFnErr
	}
	res, err := Transaction(context.Background(), fakeTx, txFn)
	require.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), "txFn error")
	assert.True(
		s.T(), fakeTx.rollbackCalled,
		"Rollback should be called when txFn returns an error",
	)
	assert.False(
		s.T(), fakeTx.commitCalled,
		"Commit should not be called when txFn returns an error",
	)
	assert.Equal(s.T(), 0, res)
}

// TestTransaction_Panic verifies that if txFn panics,
// Transaction recovers, calls Rollback, and then re-panics.
func (s *TransactionTestSuite) TestTransaction_Panic() {
	fakeTx := &FakeTx{}
	panicValue := "panic occurred"
	txFn := func(ctx context.Context, tx types.Tx) (int, error) {
		panic(panicValue)
	}
	assert.PanicsWithValue(
		s.T(),
		panicValue,
		func() {
			_, _ = Transaction(context.Background(), fakeTx, txFn)
		},
	)
	assert.True(
		s.T(), fakeTx.rollbackCalled, "Rollback should be called on panic",
	)
	assert.False(
		s.T(), fakeTx.commitCalled, "Commit should not be called on panic",
	)
}

// TestTransaction_CommitError verifies that if txFn returns no error but Commit
// fails, Transaction returns a commit error.
func (s *TransactionTestSuite) TestTransaction_CommitError() {
	commitErr := errors.New("commit failed")
	fakeTx := &FakeTx{commitErr: commitErr}
	txFn := func(ctx context.Context, tx types.Tx) (int, error) {
		return 1, nil
	}
	res, err := Transaction(context.Background(), fakeTx, txFn)
	require.Error(s.T(), err)
	assert.Contains(
		s.T(), err.Error(), "commit error")
	assert.True(
		s.T(), fakeTx.commitCalled,
		"Commit should be attempted",
	)
	assert.False(
		s.T(), fakeTx.rollbackCalled,
		"Rollback should not be called when txFn returns nil",
	)
	assert.Equal(
		s.T(), 0, res,
		"Result should be zero value on commit error",
	)
}

// TestTransaction_RollbackError verifies that if txFn returns an error and
// Rollback fails, Transaction returns the rollback error.
func (s *TransactionTestSuite) TestTransaction_RollbackError() {
	rollbackErr := errors.New("rollback failed")
	fakeTx := &FakeTx{rollbackErr: rollbackErr}
	txFn := func(ctx context.Context, tx types.Tx) (int, error) {
		return 0, errors.New("txFn error")
	}
	res, err := Transaction(context.Background(), fakeTx, txFn)
	require.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), "rollback error")
	assert.True(
		s.T(), fakeTx.rollbackCalled,
		"Rollback should be called when txFn returns an error",
	)
	assert.False(
		s.T(), fakeTx.commitCalled,
		"Commit should not be called when txFn returns an error",
	)
	assert.Equal(s.T(), 0, res)
}
