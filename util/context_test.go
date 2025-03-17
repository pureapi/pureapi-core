package util

import (
	"context"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// ContextTestSuite groups tests for the context-related functions.
type ContextTestSuite struct {
	suite.Suite
}

// TestContextTestSuite runs the ContextTestSuite.
func TestContextTestSuite(t *testing.T) {
	suite.Run(t, new(ContextTestSuite))
}

// TestNewContext_CreatesCustomData verifies that NewContext injects a custom
// data map.
func (suite *ContextTestSuite) TestNewContext_CreatesCustomData() {
	baseCtx := context.Background()
	ctx := NewContext(baseCtx)

	// Try retrieving a non-existent key, should return default.
	val := GetContextValue(ctx, "nonexistent", "default")
	assert.Equal(suite.T(), "default", val)
}

// TestSetAndGetContextValue verifies that a value set in the context can be
// retrieved.
func (suite *ContextTestSuite) TestSetAndGetContextValue() {
	baseCtx := context.Background()
	ctx := NewContext(baseCtx)

	// Set a value and verify retrieval.
	key := "test-key"
	value := "test-value"
	ctx, err := SetContextValue(ctx, key, value)
	if err != nil {
		assert.Fail(suite.T(), err.Error())
	}
	ret := GetContextValue(ctx, key, "default")
	assert.Equal(suite.T(), value, ret)

	// For a missing key, the default is returned.
	ret2 := GetContextValue(ctx, "missing", 42)
	assert.Equal(suite.T(), 42, ret2)
}

// TestGetContextValue_TypeMismatch verifies that if the stored value is of a
// different type,
// the default value is returned.
func (suite *ContextTestSuite) TestGetContextValue_TypeMismatch() {
	baseCtx := context.Background()
	ctx := NewContext(baseCtx)

	key := "number"
	// Store a string.
	ctx, err := SetContextValue(ctx, key, "not a number")
	if err != nil {
		assert.Fail(suite.T(), err.Error())
	}

	// Try to get an int value; should return the default.
	val := GetContextValue(ctx, key, 100)
	assert.Equal(suite.T(), 100, val)
}

// TestSetContextValue_NilKey verifies that setting a nil key returns an error.
func (suite *ContextTestSuite) TestSetContextValue_NilKey() {
	baseCtx := context.Background()
	ctx := NewContext(baseCtx)
	_, err := SetContextValue(ctx, nil, "value")
	assert.Error(suite.T(), err)
}

// TestSetContextValue_NoCustomContext verifies that setting a value on a
// context without custom data returns an error.
func (suite *ContextTestSuite) TestSetContextValue_NoCustomContext() {
	// Use a plain background context.
	baseCtx := context.Background()
	_, err := SetContextValue(baseCtx, "key", "value")
	require.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "no custom context set in request")
}

// TestNewDataKey_Uniqueness verifies that successive calls to NewDataKey yield
// unique, increasing keys.
func (suite *ContextTestSuite) TestNewDataKey_Uniqueness() {
	key1 := NewDataKey()
	key2 := NewDataKey()
	assert.True(suite.T(), key2 > key1, "NewDataKey should produce increasing values")
}

// TestNewDataKey_Concurrent verifies that concurrent calls to NewDataKey yield
// unique keys.
func (suite *ContextTestSuite) TestNewDataKey_Concurrent() {
	const numGoroutines = 100
	keys := make(chan DataKey, numGoroutines)
	var wg sync.WaitGroup
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			keys <- NewDataKey()
		}()
	}
	wg.Wait()
	close(keys)

	collected := make(map[string]struct{})
	for key := range keys {
		kStr := strconv.Itoa(int(key))
		_, exists := collected[kStr]
		assert.False(suite.T(), exists, "Duplicate key found: %v", key)
		collected[kStr] = struct{}{}
	}
}
