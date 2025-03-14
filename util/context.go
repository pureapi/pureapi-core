package util

import (
	"context"
	"sync"
)

// DataKey is a unique key for storing custom data in the context.
type DataKey int

var (
	// base is the base value for generating unique data keys.
	base DataKey = 1

	// lock is used to synchronize access to base.
	lock sync.Mutex

	// mainDataKey is the main key for storing custom data in the context.
	mainDataKey = NewDataKey()
)

// contextData is a map of custom data stored in the context.
type contextData struct {
	data sync.Map
}

// NewDataKey safely increments and returns the next value of base.
// It is used to create a unique key for storing custom context data.
//
// Returns:
//   - The next data key value.
func NewDataKey() DataKey {
	lock.Lock()
	defer lock.Unlock()
	base++
	return base
}

// NewContext initializes a new context with an empty contextData map.
//
// Parameters:
//   - fromCtx: The context from which the new context is derived.
//
// Returns:
//   - A new context with an initialized custom data map.
func NewContext(fromCtx context.Context) context.Context {
	return context.WithValue(fromCtx, mainDataKey, &contextData{})
}

// GetContextValue tries to retrieve a value from the custom data of the context
// for a given key.
// If the key exists and the value matches the expected type, it returns the
// value. Otherwise, it returns the provided default value.
//
// Parameters:
//   - ctx: The context from which to retrieve the value.
//   - key: The key for which to retrieve the value.
//   - returnOnNull: The default value to return if the key does not exist or
//     the type does not match.
//
// Returns:
//   - The value from the context if it exists and matches the expected type,
//     otherwise the default value.
func GetContextValue[T any](ctx context.Context, key any, returnOnNull T) T {
	cd, ok := getContextData(ctx)
	if !ok {
		return returnOnNull
	}
	value, exists := cd.data.Load(key)
	if !exists {
		return returnOnNull
	}
	typedValue, isType := value.(T)
	if !isType {
		return returnOnNull
	}
	return typedValue
}

// SetContextValue sets a value in the custom data of the context for the
// provided key. It panics if the key is nil or if the custom context is not
// set.
//
// Parameters:
//   - ctx: The context in which to set the value.
//   - key: The key for which to set the value.
//   - data: The value to set in the context.
//
// Returns:
//   - The updated context.
func SetContextValue(ctx context.Context, key any, data any) context.Context {
	if key == nil {
		panic("set context value: key cannot be nil")
	}
	cd, ok := getContextData(ctx)
	if !ok {
		panic("set context value: no custom context set in request")
	}
	cd.data.Store(key, data)
	return ctx
}

// getContextData tries to retrieve the custom data from the context.
func getContextData(ctx context.Context) (*contextData, bool) {
	cd, ok := ctx.Value(mainDataKey).(*contextData)
	return cd, ok && cd != nil
}
