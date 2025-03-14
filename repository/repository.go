package repository

import (
	"context"
	"fmt"

	"github.com/pureapi/pureapi-core/database"
	databasetypes "github.com/pureapi/pureapi-core/database/types"
	"github.com/pureapi/pureapi-core/dbquery"
	repositorytypes "github.com/pureapi/pureapi-core/repository/types"
)

// defaultReaderRepo implements read operations.
type defaultReaderRepo[Entity databasetypes.Getter] struct {
	queryBuilder repositorytypes.QueryBuilder
	errorChecker repositorytypes.ErrorChecker
}

// DefaultReaderRepo implements ReaderRepo.
var _ repositorytypes.ReaderRepo[databasetypes.Getter] = (*defaultReaderRepo[databasetypes.Getter])(nil)

// NewDefaultReaderRepo creates a new defaultReaderRepo.
//
// Parameters:
//   - ctx: Context to use.
//   - queryBuilder: The QueryBuilder to use for building queries.
//   - errorChecker: The ErrorChecker to use for checking errors.
//
// Returns:
//   - *defaultReaderRepo: A new defaultReaderRepo.
func NewDefaultReaderRepo[Entity databasetypes.Getter](
	queryBuilder repositorytypes.QueryBuilder, errorChecker repositorytypes.ErrorChecker,
) *defaultReaderRepo[Entity] {
	return &defaultReaderRepo[Entity]{
		queryBuilder: queryBuilder,
		errorChecker: errorChecker,
	}
}

// GetOne retrieves a single record from the DB by delegating to dbOps.
//
// Parameters:
//   - ctx: Context to use.
//   - preparer: The preparer to use for the query.
//   - factoryFn: A function that returns a new instance of T.
//   - getOpts: The GetOptions to use for the query.
//
// Returns:
//   - T: The entity scanned from the query.
//   - error: An error if the query fails.
func (r *defaultReaderRepo[Entity]) GetOne(
	ctx context.Context,
	preparer databasetypes.Preparer,
	factoryFn repositorytypes.GetterFactoryFn[Entity],
	getOpts *repositorytypes.GetOptions,
) (Entity, error) {
	tableName := factoryFn().TableName()
	query, params := r.queryBuilder.Get(tableName, getOpts)
	return QuerySingleEntity(
		ctx, preparer, query, params, r.errorChecker, factoryFn,
	)
}

// GetMany retrieves multiple records from the DB.
//
// Parameters:
//   - ctx: Context to use.
//   - preparer: The preparer to use for the query.
//   - factoryFn: A function that returns a new instance of T.
//   - getOpts: The GetOptions to use for the query.
//
// Returns:
//   - []T: A slice of entities scanned from the query.
//   - error: An error if the query fails.
func (r *defaultReaderRepo[Entity]) GetMany(
	ctx context.Context,
	preparer databasetypes.Preparer,
	factoryFn repositorytypes.GetterFactoryFn[Entity],
	getOpts *repositorytypes.GetOptions,
) ([]Entity, error) {
	tableName := factoryFn().TableName()
	query, params := r.queryBuilder.Get(tableName, getOpts)
	return QueryEntities(
		ctx, preparer, query, params, r.errorChecker, factoryFn,
	)
}

// Count returns the count of matching records.
//
// Parameters:
//   - ctx: Context to use.
//   - preparer: The preparer to use for the query.
//   - selectors: The Selectors to use for the query.
//   - page: The Page to use for the query.
//   - factoryFn: A function that returns a new instance of T.
//
// Returns:
//   - int: The count of matching records.
//   - error: An error if the query fails.
func (r *defaultReaderRepo[Entity]) Count(
	ctx context.Context,
	preparer databasetypes.Preparer,
	selectors dbquery.Selectors,
	page *dbquery.Page,
	factoryFn repositorytypes.GetterFactoryFn[Entity],
) (int, error) {
	tableName := factoryFn().TableName()
	countOpts := &repositorytypes.CountOptions{
		Selectors: selectors,
		Page:      page,
	}
	query, params := r.queryBuilder.Count(tableName, countOpts)
	result, err := QuerySingleValue(
		ctx,
		preparer,
		query,
		params,
		r.errorChecker, func() *int {
			return new(int)
		},
	)
	if err != nil {
		return 0, err
	}
	return *result, nil
}

// Query executes a custom SQL query that is already built.
//
// Parameters:
//   - ctx: Context to use.
//   - preparer: The preparer to use for the query.
//   - query: The SQL query to execute.
//   - parameters: The query parameters.
//   - factoryFn: A function that returns a new instance of T.
//
// Returns:
//   - []T: A slice of entities scanned from the query.
//   - error: An error if the query fails.
func (r *defaultReaderRepo[Entity]) Query(
	ctx context.Context,
	preparer databasetypes.Preparer,
	query string,
	parameters []any,
	factoryFn repositorytypes.GetterFactoryFn[Entity],
) ([]Entity, error) {
	return QueryEntities(
		ctx, preparer, query, parameters, r.errorChecker, factoryFn,
	)
}

// defaultMutatorRepo implements mutation operations.
type defaultMutatorRepo[Entity databasetypes.Mutator] struct {
	queryBuilder repositorytypes.QueryBuilder
	errorChecker repositorytypes.ErrorChecker
}

// DefaultMutatorRepo implements MutatorRepo.
var _ repositorytypes.MutatorRepo[databasetypes.Mutator] = (*defaultMutatorRepo[databasetypes.Mutator])(nil)

// NewDefaultMutatorRepo creates a new defaultMutatorRepo.
//
// Parameters:
//   - ctx: Context to use.
//   - queryBuilder: The query builder to use for the repository.
//   - errorChecker: The error checker to use for the repository.
//
// Returns:
//   - *defaultMutatorRepo: A new defaultMutatorRepo.
func NewDefaultMutatorRepo[Entity databasetypes.Mutator](
	queryBuilder repositorytypes.QueryBuilder,
	errorChecker repositorytypes.ErrorChecker,
) *defaultMutatorRepo[Entity] {
	return &defaultMutatorRepo[Entity]{
		queryBuilder: queryBuilder,
		errorChecker: errorChecker,
	}
}

// Insert builds an insert query and executes it.
//
// Parameters:
//   - ctx: Context to use.
//   - preparer: The preparer to use for the query.
//   - mutator: The entity to insert.
//
// Returns:
//   - T: The inserted entity.
//   - error: An error if the query fails.
func (r *defaultMutatorRepo[Entity]) Insert(
	ctx context.Context, preparer databasetypes.Preparer, mutator Entity,
) (Entity, error) {
	query, params := r.queryBuilder.Insert(
		mutator.TableName(), mutator.InsertedValues,
	)
	result, err := Exec(ctx, preparer, query, params, r.errorChecker)
	if err != nil {
		return mutator, err
	}
	_, err = result.LastInsertId()
	if err != nil && r.errorChecker != nil {
		return mutator, r.errorChecker.Check(err)
	}
	return mutator, err
}

// Update builds an update query and executes it.
//
// Parameters:
//   - ctx: Context to use.
//   - preparer: The preparer to use for the query.
//   - updater: The entity to update.
//   - selectors: The selectors to use for the update.
//   - updates: The updates to apply to the entity.
//
// Returns:
//   - int64: The number of rows affected by the update.
//   - error: An error if the query fails.
func (r *defaultMutatorRepo[Entity]) Update(
	ctx context.Context,
	preparer databasetypes.Preparer,
	updater Entity,
	selectors dbquery.Selectors,
	updates dbquery.Updates,
) (int64, error) {
	query, params := r.queryBuilder.UpdateQuery(
		updater.TableName(), updates, selectors,
	)
	result, err := Exec(ctx, preparer, query, params, r.errorChecker)
	if err != nil {
		return 0, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		if r.errorChecker != nil {
			return 0, r.errorChecker.Check(err)
		}
		return 0, err
	}
	return rowsAffected, nil
}

// Delete builds a delete query and executes it.
//
// Parameters:
//   - ctx: Context to use.
//   - preparer: The preparer to use for the query.
//   - deleter: The entity to delete.
//   - selectors: The selectors to use for the delete.
//   - deleteOpts: The delete options.
//
// Returns:
//   - int64: The number of rows affected by the delete.
//   - error: An error if the query fails.
func (r *defaultMutatorRepo[Entity]) Delete(
	ctx context.Context,
	preparer databasetypes.Preparer,
	deleter Entity,
	selectors dbquery.Selectors,
	deleteOpts *repositorytypes.DeleteOptions,
) (int64, error) {
	query, params := r.queryBuilder.Delete(
		deleter.TableName(), selectors, deleteOpts,
	)
	result, err := Exec(ctx, preparer, query, params, r.errorChecker)
	if err != nil {
		return 0, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		if r.errorChecker != nil {
			return 0, r.errorChecker.Check(err)
		}
		return 0, err
	}
	return rowsAffected, nil
}

// defaultCustomRepo implements the CustomRepo interface.
type defaultCustomRepo[T any] struct {
	errorChecker repositorytypes.ErrorChecker
}

// defaultCustomRepo implements the CustomRepo interface.
var _ repositorytypes.CustomRepo[any] = (*defaultCustomRepo[any])(nil)

// NewDefaultCustomRepo creates a new defaultCustomRepo.
// It requires an optional ErrorChecker to translate database-specific errors.
func NewDefaultCustomRepo[T any](errorChecker repositorytypes.ErrorChecker) repositorytypes.CustomRepo[T] {
	return &defaultCustomRepo[T]{errorChecker: errorChecker}
}

// QueryCustom executes a custom SQL query and maps the results into a slice of custom entities.
// It leverages the existing Query function and RowsToAny helper from the repository package.
func (r *defaultCustomRepo[T]) QueryCustom(
	ctx context.Context,
	preparer databasetypes.Preparer,
	query string,
	parameters []any,
	factoryFn func() T,
) ([]T, error) {
	if preparer == nil {
		return nil, fmt.Errorf("QueryCustom: preparer is nil")
	}
	rows, stmt, err := Query(ctx, preparer, query, parameters, r.errorChecker)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	defer rows.Close()
	return RowsToAny(ctx, rows, factoryFn)
}

// defaultRawQueryer provides direct query execution.
type defaultRawQueryer struct{}

// DefaultRawQueryer implements RawQueryer.
var _ repositorytypes.RawQueryer = (*defaultRawQueryer)(nil)

// NewDefaultRawQueryer creates a new defaultRawQueryer.
//
// Returns:
//   - *defaultRawQueryer: A new defaultRawQueryer.
func NewDefaultRawQueryer() *defaultRawQueryer {
	return &defaultRawQueryer{}
}

// defaultRawQueryer implements RawQueryer.
var _ repositorytypes.RawQueryer = (*defaultRawQueryer)(nil)

// Exec executes a query using a prepared statement.
//
// Parameters:
//   - ctx: Context to use.
//   - preparer: The preparer to use for the query.
//   - query: The SQL query to execute.
//   - parameters: The query parameters.
//
// Returns:
//   - Result: The Result of the query.
//   - error: An error if the query fails.
func (rq *defaultRawQueryer) Exec(
	ctx context.Context,
	preparer databasetypes.Preparer,
	query string,
	parameters []any,
) (databasetypes.Result, error) {
	return Exec(ctx, preparer, query, parameters, nil)
}

// ExecRaw executes a query directly on the DB.
//
// Parameters:
//   - ctx: Context to use.
//   - db: The DB to execute the query on.
//   - query: The SQL query to execute.
//   - parameters: The query parameters.
//
// Returns:
//   - Result: The Result of the query.
//   - error: An error if the query fails.
func (rq *defaultRawQueryer) ExecRaw(
	ctx context.Context, db databasetypes.DB, query string, parameters []any,
) (databasetypes.Result, error) {
	return ExecRaw(ctx, db, query, parameters, nil)
}

// Query executes a query that returns rows. The caller is responsible for
// closing the rows.
//
// Parameters:
//   - ctx: Context to use.
//   - preparer: The preparer to use for the query.
//   - query: The SQL query to execute.
//   - parameters: The query parameters.
//
// Returns:
//   - Rows: The rows of the query.
//   - error: An error if the query fails.
func (rq *defaultRawQueryer) Query(
	ctx context.Context,
	preparer databasetypes.Preparer,
	query string,
	parameters []any,
) (databasetypes.Rows, error) {
	rows, stmt, err := Query(ctx, preparer, query, parameters, nil)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	return rows, nil
}

// QueryRaw executes a query directly on the DB without preparation.
//
// Parameters:
//   - ctx: Context to use.
//   - db: The DB to execute the query on.
//   - query: The SQL query to execute.
//   - parameters: The query parameters.
//
// Returns:
//   - Rows: The rows of the query.
//   - error: An error if the query fails.
//
//nolint:ireturn
func (rq *defaultRawQueryer) QueryRaw(
	ctx context.Context, db databasetypes.DB, query string, parameters []any,
) (databasetypes.Rows, error) {
	return QueryRaw(ctx, db, query, parameters, nil)
}

// defaultTxManager is the default transaction manager.
type defaultTxManager[Entity any] struct{}

// DefaultTxManager implements TxManager.
var _ repositorytypes.TxManager[any] = (*defaultTxManager[any])(nil)

// NewDefaultTxManager returns a new defaultTxManager.
//
// Returns:
//   - *defaultTxManager[Entity]: The new defaultTxManager.
func NewDefaultTxManager[Entity any]() *defaultTxManager[Entity] {
	return &defaultTxManager[Entity]{}
}

// WithTransaction wraps a function call in a DB transaction.
//
// Parameters:
//   - ctx: Context to use.
//   - ctx: The context to use for the transaction.
//   - connFn: The function to get a DB connection.
//   - callback: The function to call in the transaction.
//
// Returns:
//   - Entity: The result of the callback.
//   - error: An error if the transaction fails.
func (t *defaultTxManager[Entity]) WithTransaction(
	ctx context.Context,
	connFn repositorytypes.ConnFn,
	callback databasetypes.TxFn[Entity],
) (Entity, error) {
	conn, err := connFn()
	if err != nil {
		var zero Entity
		return zero, err
	}
	tx, err := conn.BeginTx(ctx, nil)
	if err != nil {
		var zero Entity
		return zero, err
	}
	return database.Transaction(ctx, tx, callback)
}
