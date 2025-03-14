package input

import (
	"strings"
	"testing"

	"github.com/pureapi/pureapi-core/dbquery"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPredicate_String(t *testing.T) {
	var p Predicate = "gt"
	assert.Equal(t, "gt", p.String())
}

func TestPredicates_StringAndStrSlice(t *testing.T) {
	preds := Predicates{Greater, Equal, Less}
	expectedString := strings.Join(
		[]string{Greater.String(), Equal.String(), Less.String()}, ",",
	)
	assert.Equal(t, expectedString, preds.String())

	expectedSlice := []string{Greater.String(), Equal.String(), Less.String()}
	assert.Equal(t, expectedSlice, preds.StrSlice())
}

func TestPage_ToDBPage(t *testing.T) {
	p := &Page{Offset: 10, Limit: 50}
	dbPage := p.ToDBPage()
	require.NotNil(t, dbPage)
	assert.Equal(t, 10, dbPage.Offset)
	assert.Equal(t, 50, dbPage.Limit)
}

func TestOrderDirection_StringAndMapping(t *testing.T) {
	var od OrderDirection = DirectionAsc
	assert.Equal(t, "asc", od.String())

	// Verify that the DirectionsToDB mapping works.
	assert.Equal(t, dbquery.OrderAsc, DirectionsToDB[DirectionAsc])
	assert.Equal(t, dbquery.OrderAsc, DirectionsToDB[DirectionAscending])
	assert.Equal(t, dbquery.OrderDesc, DirectionsToDB[DirectionDesc])
	assert.Equal(t, dbquery.OrderDesc, DirectionsToDB[DirectionDescending])
}

func TestOrders_ToDBOrders_Success(t *testing.T) {
	// Create an API-to-DB field map.
	apiToDBFieldMap := map[string]DBField{
		"field1": {Table: "tbl", Column: "col1"},
		"field2": {Table: "tbl", Column: "col2"},
	}
	orders := Orders{
		"field1": DirectionAsc,
		"field2": DirectionDesc,
	}
	dbOrders, err := orders.TranslateToDBOrders(apiToDBFieldMap)
	require.NoError(t, err)
	require.Len(t, dbOrders, 2)
	// Check that each order is translated correctly.
	for _, o := range dbOrders {
		switch o.Field {
		case "col1":
			assert.Equal(t, dbquery.OrderAsc, o.Direction)
			assert.Equal(t, "tbl", o.Table)
		case "col2":
			assert.Equal(t, dbquery.OrderDesc, o.Direction)
			assert.Equal(t, "tbl", o.Table)
		default:
			t.Errorf("unexpected column: %s", o.Field)
		}
	}
}

func TestOrders_ToDBOrders_InvalidField(t *testing.T) {
	// Missing mapping for "unknown"
	apiToDBFieldMap := map[string]DBField{
		"field1": {Table: "tbl", Column: "col1"},
	}
	orders := Orders{
		"unknown": DirectionAsc,
	}
	_, err := orders.TranslateToDBOrders(apiToDBFieldMap)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "cannot translate field")
}

func TestSelectors_AddSelectorAndToDBSelectors_Success(t *testing.T) {
	// Create an empty collection of selectors.
	sel := Selectors{}
	sel = sel.AddSelector("age", Equal, 18)
	sel = sel.AddSelector("name", Equal, "Alice")

	// Create an API-to-DB field map.
	apiToDBFieldMap := map[string]DBField{
		"age":  {Table: "users", Column: "age"},
		"name": {Table: "users", Column: "username"},
	}

	dbSelectors, err := sel.ToDBSelectors(apiToDBFieldMap)
	require.NoError(t, err)
	require.Len(t, dbSelectors, 2)

	// Verify translation.
	for _, ds := range dbSelectors {
		switch ds.Column {
		case "age":
			// "equal" should map to dbquery.Equal.
			assert.Equal(t, dbquery.Equal, ds.Predicate)
			assert.Equal(t, 18, ds.Value)
		case "username":
			assert.Equal(t, dbquery.Equal, ds.Predicate)
			assert.Equal(t, "Alice", ds.Value)
		default:
			t.Errorf("unexpected column: %s", ds.Column)
		}
	}
}

func TestSelectors_ToDBSelectors_InvalidPredicate(t *testing.T) {
	// Setup API-to-DB field map.
	apiToDBFieldMap := map[string]DBField{
		"status": {Table: "orders", Column: "status"},
	}
	// Create a selector with an invalid predicate.
	sel := Selectors{
		"status": {Predicate: "invalid", Value: "active"},
	}
	_, err := sel.ToDBSelectors(apiToDBFieldMap)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "cannot translate predicate")
}

func TestSelectors_ToDBSelectors_InvalidField(t *testing.T) {
	// Setup API-to-DB field map without the "name" field.
	apiToDBFieldMap := map[string]DBField{
		"age": {Table: "users", Column: "age"},
	}
	sel := Selectors{
		"name": {Predicate: "eq", Value: "Bob"},
	}
	_, err := sel.ToDBSelectors(apiToDBFieldMap)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "cannot translate field")
}

func TestUpdates_ToDBUpdates_Success(t *testing.T) {
	// Setup API-to-DB field map.
	apiToDBFieldMap := map[string]DBField{
		"name":  {Table: "users", Column: "username"},
		"score": {Table: "users", Column: "score"},
	}
	updates := Updates{
		"name":  "Charlie",
		"score": 99,
	}
	dbUpdates, err := updates.ToDBUpdates(apiToDBFieldMap)
	require.NoError(t, err)
	require.Len(t, dbUpdates, 2)
	for _, u := range dbUpdates {
		switch u.Field {
		case "username":
			assert.Equal(t, "Charlie", u.Value)
		case "score":
			assert.Equal(t, 99, u.Value)
		default:
			t.Errorf("unexpected field: %s", u.Field)
		}
	}
}

func TestUpdates_ToDBUpdates_InvalidField(t *testing.T) {
	apiToDBFieldMap := map[string]DBField{
		"name": {Table: "users", Column: "username"},
	}
	updates := Updates{
		"invalid": "value",
	}
	_, err := updates.ToDBUpdates(apiToDBFieldMap)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "cannot translate field")
}
