package apierror

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAPIError(t *testing.T) {
	errID := "ERR001"
	apiErr := NewAPIError(errID)
	assert.Equal(
		t, apiErr.ID, errID,
		"expected ID %s, got %s", errID, apiErr.ID,
	)
	assert.Equal(
		t, apiErr.Origin, "-",
		"expected Origin '-' got %s", apiErr.Origin,
	)
	assert.Equal(
		t, apiErr.Data, nil,
		"expected Data nil, got %v", apiErr.Data,
	)
	assert.Nil(
		t, apiErr.Message,
		"expected nil Message, got %v", apiErr.Message,
	)
}

func TestWithData(t *testing.T) {
	apiErr := NewAPIError("ERR002")
	data := map[string]int{"count": 42}
	newErr := apiErr.WithData(data)
	assert.EqualValues(
		t, data, newErr.Data,
		"expected Data %v, got %v", data, newErr.Data,
	)
	// Ensure immutability: original error should not have Data set.
	assert.Nil(
		t, apiErr.Data,
		"expected original Data nil, got %v", apiErr.Data,
	)
}

func TestWithMessage(t *testing.T) {
	apiErr := NewAPIError("ERR003")
	msg := "something went wrong"
	newErr := apiErr.WithMessage(msg)
	if newErr.Message == nil || *newErr.Message != msg {
		t.Errorf("expected Message '%s', got %v", msg, newErr.Message)
	}
	// Check error string format.
	expectedStr := "ERR003: " + msg
	if newErr.Error() != expectedStr {
		assert.EqualError(
			t, newErr, expectedStr,
			"expected error string '%s', got '%s'",
		)
	}
	// Original error should not have Message.
	assert.Nil(
		t, apiErr.Message,
		"expected original Message nil, got %v", apiErr.Message,
	)
}

func TestWithOrigin(t *testing.T) {
	apiErr := NewAPIError("ERR004")
	newOrigin := "database"
	newErr := apiErr.WithOrigin(newOrigin)
	assert.Equal(
		t, newErr.Origin, newOrigin,
		"expected Origin '%s', got '%s'", newOrigin, newErr.Origin,
	)
	// Ensure original error is unchanged.
	assert.Equal(
		t, apiErr.Origin, "-",
		"expected original Origin '-' got '%s'", apiErr.Origin,
	)
}

func TestErrorMethodWithoutMessage(t *testing.T) {
	apiErr := NewAPIError("ERR005")
	assert.EqualError(
		t, apiErr, "ERR005",
		"expected error string 'ERR005', got '%s'",
	)
}

func TestChainMethods(t *testing.T) {
	// Test chaining multiple methods.
	apiErr := NewAPIError("ERR006")
	msg := "chained error"
	data := []int{1, 2, 3}
	newOrigin := "service"
	chainedErr := apiErr.
		WithMessage(msg).
		WithData(data).
		WithOrigin(newOrigin)
	// Verify all fields.
	assert.Equal(
		t, chainedErr.ID, "ERR006",
		"expected ID 'ERR006', got '%s'", chainedErr.ID,
	)
	if chainedErr.Message == nil || *chainedErr.Message != msg {
		t.Errorf("expected Message '%s', got %v", msg, chainedErr.Message)
	}
	if !reflect.DeepEqual(chainedErr.Data, data) {
		t.Errorf("expected Data %v, got %v", data, chainedErr.Data)
	}
	assert.Equal(
		t, chainedErr.Origin, newOrigin,
		"expected Origin '%s', got '%s'", newOrigin, chainedErr.Origin,
	)
	// Original APIError should remain unchanged.
	if apiErr.Message != nil || apiErr.Data != nil || apiErr.Origin != "-" {
		t.Error("expected original APIError to be unchanged")
	}
}

func TestJSONMarshalling(t *testing.T) {
	// Test JSON marshalling for APIError.
	msg := "json error"
	data := map[string]string{"key": "value"}
	apiErr := NewAPIError("ERR007").
		WithMessage(msg).
		WithData(data).
		WithOrigin("handler")
	// Marshal the APIError to JSON.
	b, err := json.Marshal(apiErr)
	assert.NoError(t, err, "json.Marshal returned error: %v", err)
	// Unmarshal to a map to check fields.
	var result map[string]any
	assert.NoError(
		t, json.Unmarshal(b, &result),
		"json.Unmarshal returned error: %v", err,
	)
	// Check required fields.
	assert.Equal(
		t, result["id"],
		"ERR007", "expected id 'ERR007', got %v", result["id"],
	)
	// Message should be present.
	assert.Equal(
		t, result["message"], msg,
		"expected message '%s', got %v", msg, result["message"],
	)
	// Data should be present.
	dataMap, ok := result["data"].(map[string]any)
	assert.True(t, ok, "expected data to be a map, got %T", result["data"])
	assert.Equal(
		t, dataMap["key"], "value",
		"expected data key 'value', got %v", dataMap["key"],
	)
	// Origin should be present.
	assert.Equal(
		t, result["origin"], "handler",
		"expected origin 'handler', got %v", result["origin"],
	)
}
