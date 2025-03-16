package repository

import (
	"context"
	"fmt"

	databasetypes "github.com/pureapi/pureapi-core/database/types"
	repositorytypes "github.com/pureapi/pureapi-core/repository/types"
)

// Exec prepares and executes a query with parameters, returning the Result.
//
// Parameters:
//   - ctx: Context to use.
//   - preparer: The preparer to use for the query.
//   - query: The SQL query to execute.
//   - parameters: The query parameters.
//   - errorChecker: An optional ErrorChecker to check for errors.
//
// Returns:
//   - Result: The Result of the query.
//   - error: An error if the query fails.
func Exec(
	ctx context.Context,
	preparer databasetypes.Preparer,
	query string,
	parameters []any,
	errorChecker repositorytypes.ErrorChecker,
) (databasetypes.Result, error) {
	if preparer == nil {
		return nil, fmt.Errorf("Exec: preparer is nil")
	}
	result, err := doExec(ctx, preparer, query, parameters)
	if err != nil {
		if errorChecker == nil {
			return nil, err
		}
		return nil, errorChecker.Check(err)
	}
	return result, nil
}

// Query prepares and executes a query that returns rows. The caller is
// responsible for closing both the returned rows and statement.
//
// Parameters:
//   - ctx: Context to use.
//   - preparer: The preparer to use for the query.
//   - query: The SQL query to execute.
//   - parameters: The query parameters.
//   - errorChecker: An optional ErrorChecker to check for errors.
//
// Returns:
//   - Rows: The rows of the query.
//   - Stmt: The statement of the query.
//   - error: An error if the query fails.
func Query(
	ctx context.Context,
	preparer databasetypes.Preparer,
	query string,
	parameters []any,
	errorChecker repositorytypes.ErrorChecker,
) (databasetypes.Rows, databasetypes.Stmt, error) {
	if preparer == nil {
		return nil, nil, fmt.Errorf("Query: preparer is nil")
	}
	rows, stmt, err := doQuery(ctx, preparer, query, parameters)
	if err != nil {
		if errorChecker == nil {
			return nil, nil, err
		}
		return nil, nil, errorChecker.Check(err)
	}
	return rows, stmt, nil
}

// ExecRaw executes a query directly on the DB without explicit preparation.
//
// Parameters:
//   - ctx: Context to use.
//   - db: The DB to execute the query on.
//   - query: The SQL query to execute.
//   - parameters: The query parameters.
//   - errorChecker: An optional ErrorChecker to check for errors.
//
// Returns:
//   - Result: The Result of the query.
//   - error: An error if the query fails.
func ExecRaw(
	ctx context.Context,
	db databasetypes.DB,
	query string,
	parameters []any,
	errorChecker repositorytypes.ErrorChecker,
) (databasetypes.Result, error) {
	if db == nil {
		return nil, fmt.Errorf("ExecRaw: db is nil")
	}
	result, err := doExecRaw(ctx, db, query, parameters)
	if err != nil {
		if errorChecker == nil {
			return nil, err
		}
		return nil, errorChecker.Check(err)
	}
	return result, nil
}

// QueryRaw executes a query directly on the DB without preparation.
// The caller must close the returned rows.
//
// Parameters:
//   - ctx: Context to use.
//   - db: The DB to execute the query on.
//   - query: The SQL query to execute.
//   - parameters: The query parameters.
//   - errorChecker: An optional ErrorChecker to check for errors.
//
// Returns:
//   - Rows: The rows of the query.
//   - error: An error if the query fails.
func QueryRaw(
	ctx context.Context,
	db databasetypes.DB,
	query string,
	parameters []any,
	errorChecker repositorytypes.ErrorChecker,
) (databasetypes.Rows, error) {
	if db == nil {
		return nil, fmt.Errorf("QueryRaw: db is nil")
	}
	rows, err := doQueryRaw(ctx, db, query, parameters)
	if err != nil {
		if errorChecker == nil {
			return nil, err
		}
		return nil, errorChecker.Check(err)
	}
	return rows, nil
}

// QuerySingleValue executes a query that is expected to return a single scalar
// value. It prepares the query, executes it using QueryRow, scans the result
// using the provided factory function, and checks for errors.
//
// Parameters:
//   - ctx: Context to use.
//   - preparer: The preparer to use for the query.
//   - query: The SQL query to execute.
//   - parameters: The query parameters.
//   - errorChecker: An optional ErrorChecker to check for errors.
//   - factoryFn: A function that returns a new instance of T
//     (typically a pointer).
//
// Returns:
//   - T: The scanned scalar value.
//   - error: An error if the query or scan fails.
func QuerySingleValue[T any](
	ctx context.Context,
	preparer databasetypes.Preparer,
	query string,
	parameters []any,
	errorChecker repositorytypes.ErrorChecker,
	factoryFn func() T,
) (T, error) {
	var zero T
	if preparer == nil {
		return zero, fmt.Errorf("QuerySingleValue: preparer is nil")
	}
	stmt, err := preparer.Prepare(query)
	if err != nil {
		return zero, err
	}
	defer stmt.Close()

	row := stmt.QueryRow(parameters...)
	result, err := RowToAny(ctx, row, factoryFn)
	if err != nil {
		if errorChecker != nil {
			return result, errorChecker.Check(err)
		}
		return result, err
	}
	return result, nil
}

// QuerySingleEntity executes a query and scans a single entity of type T,
// handling statement and row closures internally.
//
// Parameters:
//   - ctx: Context to use.
//   - preparer: The preparer to use for the query.
//   - query: The SQL query to execute.
//   - parameters: The query parameters.
//   - errorChecker: An optional ErrorChecker to check for errors.
//   - factoryFn: A function that returns a new instance of T.
//
// Returns:
//   - T: The entity scanned from the query.
//   - error: An error if the query fails.
func QuerySingleEntity[Entity databasetypes.Getter](
	ctx context.Context,
	preparer databasetypes.Preparer,
	query string,
	parameters []any,
	errorChecker repositorytypes.ErrorChecker,
	factoryFn func() Entity,
) (Entity, error) {
	var zero Entity
	if preparer == nil {
		return zero, fmt.Errorf("QuerySingleEntity: preparer is nil")
	}
	stmt, err := preparer.Prepare(query)
	if err != nil {
		return zero, err
	}
	defer stmt.Close()
	entity, err := RowToEntity(ctx, stmt.QueryRow(parameters...), factoryFn)
	if err != nil {
		if errorChecker == nil {
			return zero, err
		}
		return zero, errorChecker.Check(err)
	}
	return entity, nil
}

// QueryEntities executes a query and scans all entities of type T,
// handling statement and row closures internally.
//
// Parameters:
//   - ctx: Context to use.
//   - preparer: The preparer to use for the query.
//   - query: The SQL query to execute.
//   - parameters: The query parameters.
//   - errorChecker: An optional ErrorChecker to check for errors.
//   - factoryFn: A function that returns a new instance of T.
//
// Returns:
//   - []T: A slice of entities scanned from the query.
//   - error: An error if the query fails.
func QueryEntities[Entity databasetypes.Getter](
	ctx context.Context,
	preparer databasetypes.Preparer,
	query string,
	parameters []any,
	errorChecker repositorytypes.ErrorChecker,
	factoryFn func() Entity,
) ([]Entity, error) {
	rows, stmt, err := Query(ctx, preparer, query, parameters, errorChecker)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	defer rows.Close()
	return RowsToEntities(ctx, rows, factoryFn)
}

// RowToEntity scans a single row into a new entity.
//
// Parameters:
//   - ctx: Context to use.
//   - row: The row to scan.
//   - factoryFn: A function that returns a new instance of T.
//
// Returns:
//   - T: The entity scanned from the row.
//   - error: An error if the scan fails.
func RowToEntity[T databasetypes.Getter](
	_ context.Context,
	row databasetypes.Row,
	factoryFn func() T,
) (T, error) {
	var zero T
	entity := factoryFn()
	if err := entity.ScanRow(row); err != nil {
		return zero, err
	}
	if err := row.Err(); err != nil {
		return zero, err
	}
	return entity, nil
}

// RowToAny scans a single row into a new entity of type T.
// It uses factoryFn to create an instance of T, scans the row into it,
// and then checks for any scanning errors.
//
// Parameters:
//   - ctx: Context to use.
//   - row: The row to scan.
//   - factoryFn: A function that returns a new instance of T
//     (typically a pointer).
//
// Returns:
//   - T: The entity scanned from the row.
//   - error: An error if scanning fails.
func RowToAny[T any](
	_ context.Context, row databasetypes.Row, factoryFn func() T,
) (T, error) {
	var zero T
	entity := factoryFn()
	// Capture any scanning error.
	if err := row.Scan(entity); err != nil {
		return zero, err
	}
	if err := row.Err(); err != nil {
		return zero, err
	}
	return entity, nil
}

// RowsToAny scans all rows into a slice of entities of type T.
// It repeatedly calls RowToAny for each row returned by rows.
//
// Parameters:
//   - ctx: Context to use.
//   - rows: The rows to scan.
//   - factoryFn: A function that returns a new instance of T
//     (typically a pointer).
//
// Returns:
//   - []T: A slice of entities scanned from the rows.
//   - error: An error if scanning any row fails.
func RowsToAny[T any](
	ctx context.Context, rows databasetypes.Rows, factoryFn func() T,
) ([]T, error) {
	var results []T
	for rows.Next() {
		entity, err := RowToAny(ctx, rows, factoryFn)
		if err != nil {
			return nil, err
		}
		results = append(results, entity)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return results, nil
}

// RowsToEntities scans all rows into a slice of entities.
//
// Parameters:
//   - ctx: Context to use.
//   - rows: The rows to scan.
//   - factoryFn: A function that returns a new instance of T.
//
// Returns:
//   - []T: A slice of entities scanned from the rows.
//   - error: An error if the scan fails.
func RowsToEntities[T databasetypes.Getter](
	_ context.Context, rows databasetypes.Rows, factoryFn func() T,
) ([]T, error) {
	results := []T{}
	for rows.Next() {
		entity := factoryFn()
		if err := entity.ScanRow(rows); err != nil {
			return nil, err
		}
		results = append(results, entity)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return results, nil
}

// doExec executes a query with parameters.
func doExec(
	_ context.Context,
	preparer databasetypes.Preparer,
	query string,
	parameters []any,
) (databasetypes.Result, error) {
	stmt, err := preparer.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	result, err := stmt.Exec(parameters...)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// doExecRaw executes a query directly on the DB without preparation.
func doExecRaw(
	_ context.Context, db databasetypes.DB, query string, parameters []any,
) (databasetypes.Result, error) {
	result, err := db.Exec(query, parameters...)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// doQuery executes a query with parameters.
func doQuery(
	_ context.Context,
	preparer databasetypes.Preparer,
	query string,
	parameters []any,
) (databasetypes.Rows, databasetypes.Stmt, error) {
	stmt, err := preparer.Prepare(query)
	if err != nil {
		return nil, nil, err
	}
	rows, err := stmt.Query(parameters...)
	if err != nil {
		if closeErr := stmt.Close(); closeErr != nil {
			return nil, nil, fmt.Errorf(
				"query error: %w; additionally, stmt.Close error: %w",
				err,
				closeErr,
			)
		}
		return nil, nil, err
	}
	return rows, stmt, nil
}

// doQueryRaw executes a query directly on the DB without preparation.
func doQueryRaw(
	_ context.Context, db databasetypes.DB, query string, parameters []any,
) (databasetypes.Rows, error) {
	rows, err := db.Query(query, parameters...)
	if err != nil {
		return nil, err
	}
	return rows, nil
}
