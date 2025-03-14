package types

// TableNamer provides the table name for an entity.
type TableNamer interface {
	TableName() string
}

// Mutator provides values for insert and update operations.
type Mutator interface {
	TableNamer
	// InsertedValues returns the column names and values for insertion.
	InsertedValues() ([]string, []any)
}

// Getter can scan a database row into itself.
type Getter interface {
	TableNamer
	// ScanRow should populate the entity from the given Row.
	ScanRow(row Row) error
}

// CRUDEntity is a helper constraint for entities that can be both queried and
// altered.
type CRUDEntity interface {
	Getter
	Mutator
}
