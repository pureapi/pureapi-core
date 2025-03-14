package dbquery

// Predicate represents the predicate of a database selector.
type Predicate string

// OrderDirection is used to specify the order of the result set.
type OrderDirection string

// Order is used to specify the order of the result set.
type Order struct {
	Table     string
	Field     string
	Direction OrderDirection
}

// Orders is a list of orders.
type Orders []Order

// ColumnSelector represents a column selector.
type ColumnSelector struct {
	Table  string
	Column string
}

// Projection represents a projected column in a query.
type Projection struct {
	Table  string
	Column string
	Alias  string
}

// Projections is a list of projections.
type Projections []Projection

// Selector represents a database selector.
type Selector struct {
	Table     string
	Column    string
	Predicate Predicate
	Value     any
}

// NewSelector creates a new selector with the given parameters.
//
// Parameters:
//   - column: the column name
//   - predicate: the predicate
//   - value: the value
//
// Returns:
//   - *Selector: The new selector
func NewSelector(column string, predicate Predicate, value any) *Selector {
	return &Selector{
		Column:    column,
		Predicate: predicate,
		Value:     value,
	}
}

// WithTable returns a new selector with the provided table name.
//
// Parameters:
//   - table: the table name
//
// Returns:
//   - *Selector: The new selector
func (s *Selector) WithTable(table string) *Selector {
	newSelector := *s
	newSelector.Table = table
	return &newSelector
}

// Selectors represents a list of database selectors.
type Selectors []Selector

// NewSelectors returns a new list of selectors.
//
// Parameters:
//   - selectors: The selectors
//
// Returns:
//   - Selectors: The new list of selectors
func NewSelectors(selectors ...Selector) Selectors {
	return selectors
}

// Add adds a new selector to the list.
//
// Parameters:
//   - column: The column name
//   - predicate: The predicate
//   - value: The value
//
// Returns:
//   - Selectors: The new list of selectors
func (s Selectors) Add(
	column string, predicate Predicate, value any,
) Selectors {
	return append(s, *NewSelector(column, predicate, value))
}

// GetByField returns selector with the given field.
//
// Parameters:
//   - field: the field to search for
//
// Returns:
//   - *Selector: The selector
func (s Selectors) GetByField(field string) *Selector {
	for j := range s {
		if s[j].Column == field {
			return &s[j]
		}
	}
	return nil
}

// GetByFields returns selectors with the given fields.
//
// Parameters:
//   - fields: the fields to search for
//
// Returns:
//   - []Selector: A list of selectors
func (s Selectors) GetByFields(fields ...string) []Selector {
	var result []Selector
	for _, field := range fields {
		for i := range s {
			if s[i].Column == field {
				result = append(result, s[i])
			}
		}
	}
	return result
}

// Update is the options struct used for update queries.
type Update struct {
	Field string
	Value any
}

// NewUpdate creates a new update field.
//
// Parameters:
//   - field: The field
//   - value: The value
//
// Returns:
//   - Update: The new update field
func NewUpdate(field string, value any) Update {
	return Update{
		Field: field,
		Value: value,
	}
}

// Updates is a list of update fields
type Updates []Update

// NewUpdates creates a new list of updates
//
// Parameters:
//   - updates: The updates
//
// Returns:
//   - Updates: The new list of updates
func NewUpdates(updates ...Update) Updates {
	return updates
}

// Add adds a new update field to the list.
//
// Parameters:
//   - field: The field
//   - value: The value
//
// Returns:
//   - Updates: The new list of updates
func (u Updates) Add(field string, value any) Updates {
	return append(u, Update{Field: field, Value: value})
}

// Page is used to specify the page of the result set.
type Page struct {
	Offset int
	Limit  int
}

// JoinType represents the type of join
type JoinType string

// Join represents a database join clause.
type Join struct {
	JoinType JoinType
	Table    string
	OnLeft   ColumnSelector
	OnRight  ColumnSelector
}

// NewJoin creates a new join clause.
//
// Parameters:
//   - joinType: The type of join
//   - table: The table name
//   - onLeft: The left column selector
//   - onRight: The right column selector
//
// Returns:
//   - Join: The new join clause
func NewJoin(
	joinType JoinType, table string, onLeft, onRight ColumnSelector,
) Join {
	return Join{
		JoinType: joinType,
		Table:    table,
		OnLeft:   onLeft,
		OnRight:  onRight,
	}
}

// Joins is a list of joins
type Joins []Join
