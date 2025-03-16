package types

import (
	"fmt"
	"log"
	"strings"

	_ "github.com/mattn/go-sqlite3"

	"github.com/pureapi/pureapi-core/apierror"
	"github.com/pureapi/pureapi-core/database/types"
	"github.com/pureapi/pureapi-core/dbquery"
	repotypes "github.com/pureapi/pureapi-core/repository/types"
)

// User represents a simple entity for demonstration.
// It implements both Getter and Mutator interfaces.
type User struct {
	ID   int64
	Name string
}

// TableName returns the DB table name for User.
//
// Returns:
//   - string: The table name.
func (u *User) TableName() string {
	return "users"
}

// ScanRow populates the User from a DB row.
//
// Parameters:
//   - row: The DB row to scan.
//
// Returns:
//   - error: An error if the scan fails.
func (u *User) ScanRow(row types.Row) error {
	return row.Scan(&u.ID, &u.Name)
}

// InsertedValues returns the column names and values for insertion.
//
// Returns:
//   - []string: The column names.
//   - []any: The values.
func (u *User) InsertedValues() ([]string, []any) {
	// The "id" column is auto-generated.
	return []string{"name"}, []any{u.Name}
}

// ErrorUser is a type that deliberately returns a non-existent table name
// to force an error during retrieval.
type ErrorUser User

// TableName returns a table name that does not exist.
//
// Returns:
//   - string: The table name.
func (eu *ErrorUser) TableName() string {
	return "nonexistent_users"
}

// ScanRow populates the ErrorUser from a DB row.
//
// Parameters:
//   - row: The DB row to scan.
//
// Returns:
//   - error: An error if the scan fails.
func (eu *ErrorUser) ScanRow(row types.Row) error {
	// Reuse User's ScanRow implementation.
	return (*User)(eu).ScanRow(row)
}

// InsertedValues returns the column names and values for insertion.
//
// Returns:
//   - []string: The column names.
//   - []any: The values.
func (eu *ErrorUser) InsertedValues() ([]string, []any) {
	return (*User)(eu).InsertedValues()
}

// SimpleQueryBuilder implements the DataMutatorQuery and DataReaderQuery
// interfaces.
type SimpleQueryBuilder struct{}

// Get builds a simple SELECT query for the given table.
// The query is not parameterized and should only be used for testing.
//
// Parameters:
//   - table: The table to query.
//   - opts: Optional GetOptions.
//
// Returns:
//   - string: The query.
//   - []any: The values.
func (qb *SimpleQueryBuilder) Get(
	table string, opts *repotypes.GetOptions,
) (string, []any) {
	query := fmt.Sprintf("SELECT id, name FROM %s", table)
	log.Printf("Get query: %s", query)
	return query, nil
}

func (qb *SimpleQueryBuilder) Count(
	table string, opts *repotypes.CountOptions,
) (string, []any) {
	panic("not implemented")
}

// Insert builds an INSERT query based on the provided inserted values.
//
// Parameters:
//   - table: The table to insert into.
//   - insertedValuesFn: A function that returns the column names and values.
//
// Returns:
//   - string: The query.
//   - []any: The values.
func (qb *SimpleQueryBuilder) Insert(
	table string, insertedValuesFn repotypes.InsertedValuesFn,
) (string, []any) {
	cols, vals := insertedValuesFn()
	placeholders := make([]string, len(vals))
	for i := range vals {
		placeholders[i] = "?"
	}
	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s)",
		table,
		strings.Join(cols, ", "),
		strings.Join(placeholders, ", "),
	)
	log.Printf("Insert query: %s, values: %v", query, vals)
	return query, vals
}

func (qb *SimpleQueryBuilder) InsertMany(
	table string, valuesFuncs []repotypes.InsertedValuesFn,
) (query string, params []any) {
	panic("not implemented")
}
func (qb *SimpleQueryBuilder) UpsertMany(
	table string, valuesFuncs []repotypes.InsertedValuesFn,
	updateProjections []dbquery.Projection,
) (string, []any) {
	panic("not implemented")
}
func (qb *SimpleQueryBuilder) UpdateQuery(
	table string, updates []dbquery.Update, selectors []dbquery.Selector,
) (string, []any) {
	panic("not implemented")
}
func (qb *SimpleQueryBuilder) Delete(
	table string, selectors []dbquery.Selector, opts *repotypes.DeleteOptions,
) (string, []any) {
	panic("not implemented")
}

// SimpleErrorChecker is a trivial custom error checker that returns a custom
// translated error from the original error.
type SimpleErrorChecker struct{}

// Check returns the error without modification.
//
// Parameters:
//   - err: The error to check.
//
// Returns:
//   - error: The translated error.
func (ec *SimpleErrorChecker) Check(err error) error {
	return apierror.NewAPIError("MY_API_ERROR").
		WithData(err.Error()).WithOrigin("my_api")
}

// SimpleSchemaManager implements the SchemaManager interface for demonstration.
type SimpleSchemaManager struct{}

// CreateTableQuery builds a CREATE TABLE statement.
//
// Parameters:
//   - tableName: The name of the table to create.
//   - ifNotExists: Whether to create the table if it doesn't exist.
//   - columns: The column definitions for the table.
//   - constraints: Additional constraints for the table.
//   - opts: Optional table options.
//
// Returns:
//   - string: The query.
//   - []any: The values.
func (sm *SimpleSchemaManager) CreateTableQuery(
	tableName string,
	ifNotExists bool,
	columns []repotypes.ColumnDefinition,
	constraints []string,
	opts repotypes.TableOptions,
) (string, []any, error) {
	inClause := ""
	if ifNotExists {
		inClause = "IF NOT EXISTS"
	}

	var colDefs []string
	for _, col := range columns {
		def := fmt.Sprintf("%s %s", col.Name, col.Type)

		// Ensure AUTOINCREMENT is used correctly
		if col.AutoIncrement {
			if col.Type == "INTEGER" && col.PrimaryKey {
				def = fmt.Sprintf(
					"%s INTEGER PRIMARY KEY AUTOINCREMENT", col.Name,
				)
			} else {
				log.Printf(
					"Warning: AUTOINCREMENT can only be used with INTEGER PRIMARY KEY."+
						"Ignoring for column: %s",
					col.Name,
				)
			}
		} else {
			// Normal column definitions
			if col.PrimaryKey {
				def += " PRIMARY KEY"
			}
			if col.NotNull {
				def += " NOT NULL"
			}
			if col.Unique {
				def += " UNIQUE"
			}
			if col.Default != nil {
				def += fmt.Sprintf(" DEFAULT %s", *col.Default)
			}
			if col.Extra != "" {
				def += " " + col.Extra
			}
		}

		colDefs = append(colDefs, def)
	}

	// Add any constraints.
	colDefs = append(colDefs, constraints...)

	query := fmt.Sprintf(
		"CREATE TABLE %s %s (%s);",
		inClause,
		tableName,
		strings.Join(colDefs, ", "),
	)

	log.Printf("Create table query: %s\n", query)
	return query, nil, nil
}

func (sm *SimpleSchemaManager) CreateDatabaseQuery(
	dbName string, ifNotExists bool, charset string, collate string,
) (string, []any, error) {
	panic("not implemented")
}
func (sm *SimpleSchemaManager) UseDatabaseQuery(
	dbName string,
) (string, []any, error) {
	panic("not implemented")
}
func (sm *SimpleSchemaManager) SetVariableQuery(
	variable string, value string,
) (string, []any, error) {
	panic("not implemented")
}
