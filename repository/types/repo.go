package types

import (
	"context"

	"github.com/pureapi/pureapi-core/database/types"
	"github.com/pureapi/pureapi-core/dbquery"
)

// ConnFn returns a database connection.
type ConnFn func() (types.DB, error)

// GetterFactoryFn returns a Getter factory function.
type GetterFactoryFn[Entity types.Getter] func() Entity

// ReaderRepo defines retrieval-related operations.
type ReaderRepo[Entity types.Getter] interface {
	// GetOne retrieves a single record from the DB.
	GetOne(
		ctx context.Context,
		preparer types.Preparer,
		factoryFn GetterFactoryFn[Entity],
		getOptions *GetOptions,
	) (Entity, error)

	// GetMany retrieves multiple records from the DB.
	GetMany(
		ctx context.Context,
		preparer types.Preparer,
		factoryFn GetterFactoryFn[Entity],
		getOptions *GetOptions,
	) ([]Entity, error)

	// Count returns a record count.
	Count(
		ctx context.Context,
		preparer types.Preparer,
		selectors dbquery.Selectors,
		page *dbquery.Page,
		factoryFn GetterFactoryFn[Entity],
	) (int, error)
}

// MutatorRepo defines mutation-related operations.
type MutatorRepo[Entity types.Mutator] interface {
	// Insert builds an insert query and executes it.
	Insert(
		ctx context.Context, preparer types.Preparer, mutator Entity,
	) (Entity, error)

	// Update builds an update query and executes it.
	Update(
		ctx context.Context,
		preparer types.Preparer,
		updater Entity,
		selectors dbquery.Selectors,
		updates dbquery.Updates,
	) (int64, error)

	// Delete builds a delete query and executes it.
	Delete(
		ctx context.Context,
		preparer types.Preparer,
		deleter Entity,
		selectors dbquery.Selectors,
		deleteOpts *DeleteOptions,
	) (int64, error)
}

// CustomRepo defines methods for executing custom SQL queries and mapping the
// results into custom entities.
type CustomRepo[Entity any] interface {
	// QueryCustom executes a custom SQL query.
	// - ctx: context for the query.
	// - preparer: a Preparer (like a DB or Tx) to prepare the statement.
	// - query: the raw SQL query.
	// - parameters: any parameters for the query.
	// - factoryFn: a function that returns a new instance of T.
	//
	// It returns a slice of T or an error if the query or scan fails.
	QueryCustom(
		ctx context.Context,
		preparer types.Preparer,
		query string,
		parameters []any,
		factoryFn func() Entity,
	) ([]Entity, error)
}

// RawQueryer defines generic methods for executing raw queries and commands.
type RawQueryer interface {
	// Exec executes a query using a prepared statement that does not return
	// rows.
	Exec(
		ctx context.Context,
		preparer types.Preparer,
		query string,
		parameters []any,
	) (types.Result, error)

	// ExecRaw executes a query directly on the DB without explicit preparation.
	ExecRaw(
		ctx context.Context,
		db types.DB,
		query string,
		parameters []any,
	) (types.Result, error)

	// Query prepares and executes a query that returns rows. Returns rows.
	// The caller is responsible for closing the returned rows.
	Query(
		ctx context.Context,
		preparer types.Preparer,
		query string,
		parameters []any,
	) (types.Rows, error)

	// QueryRaw executes a query directly on the DB without preparation and
	// returns rows. The caller is responsible for closing the returned rows.
	QueryRaw(ctx context.Context, db types.DB, query string, parameters []any,
	) (types.Rows, error)
}

// TxManager is an interface for transaction management.
type TxManager[Entity any] interface {
	// WithTransaction wraps a function call in a DB transaction.
	WithTransaction(
		ctx context.Context, connFn ConnFn, callback types.TxFn[Entity],
	) (Entity, error)
}
