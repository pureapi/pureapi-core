package types

import "github.com/pureapi/pureapi-core/dbquery"

// GetOptions is used for get queries.
type GetOptions struct {
	Selectors   dbquery.Selectors
	Orders      dbquery.Orders
	Page        *dbquery.Page
	Joins       dbquery.Joins
	Projections dbquery.Projections
	Lock        bool
}

// CountOptions is used for count queries.
type CountOptions struct {
	Selectors dbquery.Selectors
	Page      *dbquery.Page
	Joins     dbquery.Joins
}

// DeleteOptions is used for delete queries.
type DeleteOptions struct {
	Limit  int
	Orders dbquery.Orders
}

// ColumnDefinition defines the properties for a table column in a table.
// creation query.
type ColumnDefinition struct {
	Name          string  // Column name
	Type          string  // Data type (with length/precision, e.g. "CHAR(36)")
	NotNull       bool    // Whether to add NOT NULL (if false, NULL is allowed)
	Default       *string // Optional default value (pass nil if not needed)
	AutoIncrement bool    // Whether to add AUTO_INCREMENT
	Extra         string  // Extra column options (e.g. "CHARACTER SET utf8mb4 COLLATE utf8mb4_bin")
	PrimaryKey    bool    // Marks this column as primary key (inline)
	Unique        bool    // Marks this column as UNIQUE (unless already primary key)
}

// TableOptions holds additional options for a table creation query.
type TableOptions struct {
	Engine  string // e.g. "InnoDB"
	Charset string // e.g. "utf8mb4"
	Collate string // e.g. "utf8mb4_bin"
}

// InsertedValuesFn defines a function that returns column names and values for
// an insert. This allows deferred evaluation of values and consistent ordering
// of parameters.
type InsertedValuesFn func() ([]string, []any)

// DataMutatorQuery provides methods for modifying data in a dbquery.
type DataMutatorQuery interface {
	// Insert builds an INSERT statement for a single row.
	Insert(
		table string, insertedValuesFunc InsertedValuesFn,
	) (query string, params []any)

	// InsertMany builds a batch INSERT for multiple rows.
	InsertMany(
		table string, valuesFuncs []InsertedValuesFn,
	) (query string, params []any)

	// UpsertMany builds an UPSERT (insert or update) statement.
	UpsertMany(
		table string, valuesFuncs []InsertedValuesFn,
		updateProjections []dbquery.Projection,
	) (query string, params []any)

	// UpdateQuery builds an UPDATE statement for given selectors.
	UpdateQuery(
		table string, updates []dbquery.Update, selectors []dbquery.Selector,
	) (query string, params []any)

	// Delete builds a DELETE statement for given selectors.
	Delete(
		table string, selectors []dbquery.Selector, opts *DeleteOptions,
	) (query string, params []any)
}

// DataReaderQuery provides methods for querying data.
type DataReaderQuery interface {
	// Get builds a SELECT statement with optional filters.
	Get(table string, opts *GetOptions) (query string, params []any)

	// Count builds a SELECT COUNT(*) statement with optional filters.
	Count(table string, opts *CountOptions) (query string, params []any)
}

// SchemaManager provides methods for managing the database schema.
type SchemaManager interface {
	// CreateDatabaseQuery builds a CREATE DATABASE statement.
	CreateDatabaseQuery(
		dbName string, ifNotExists bool, charset string, collate string,
	) (string, []any, error)

	// CreateTableQuery builds a CREATE TABLE statement.
	CreateTableQuery(
		tableName string, ifNotExists bool, columns []ColumnDefinition,
		constraints []string, opts TableOptions,
	) (string, []any, error)

	// UseDatabaseQuery builds a USE DATABASE statement.
	UseDatabaseQuery(dbName string) (string, []any, error)

	// SetVariableQuery builds a SET statement for a variable.
	SetVariableQuery(
		variable string, value string,
	) (string, []any, error)
}

// AdvisoryLocker provides methods for advisory locking.
type AdvisoryLocker interface {
	// AdvisoryLock builds an advisory lock statement.
	AdvisoryLock(lockName string, timeout int) (string, []any, error)

	// AdvisoryUnlock builds an advisory unlock statement.
	AdvisoryUnlock(lockName string) (string, []any, error)
}

// QueryBuilder combines all query-building interfaces.
type QueryBuilder interface {
	DataMutatorQuery
	DataReaderQuery
	SchemaManager
	AdvisoryLocker
}
