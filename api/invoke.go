package api

import (
	"context"

	databasetypes "github.com/pureapi/pureapi-core/database/types"
	"github.com/pureapi/pureapi-core/input"
	repositorytypes "github.com/pureapi/pureapi-core/repository/types"
)

// CreateInvoke executes the create operation.
//
// Parameters:
//   - ctx: The context.
//   - connFn: The database connection function.
//   - entity: The entity to create.
//   - mutatorRepo: The mutator repository.
//   - txManager: The transaction manager.
//
// Returns:
//   - Entity: The created entity.
//   - error: Any error that occurred during the operation.
func CreateInvoke[Entity databasetypes.Mutator](
	ctx context.Context,
	connFn repositorytypes.ConnFn,
	entity Entity,
	mutatorRepo repositorytypes.MutatorRepo[Entity],
	txManager repositorytypes.TxManager[Entity],
) (Entity, error) {
	return txManager.WithTransaction(
		ctx,
		connFn,
		func(ctx context.Context, tx databasetypes.Tx) (Entity, error) {
			return mutatorRepo.Insert(ctx, tx, entity)
		},
	)
}

// GetInvoke executes the get operation.
//
// Parameters:
//   - ctx: The context.
//   - connFn: The database connection function.
//   - entityFactoryFn: The entity factory function.
//   - readerRepo: The reader repository.
//   - txManager: The transaction manager.
//
// Returns:
//   - []Entity: The entities.
//   - error: Any error that occurred during the operation.
func GetInvoke[Getter databasetypes.Getter](
	ctx context.Context,
	parsedInput *input.ParsedGetEndpointInput,
	connFn repositorytypes.ConnFn,
	entityFactoryFn repositorytypes.GetterFactoryFn[Getter],
	readerRepo repositorytypes.ReaderRepo[Getter],
	_ repositorytypes.TxManager[Getter],
) ([]Getter, int, error) {
	conn, err := connFn()
	if err != nil {
		return nil, 0, err
	}
	tx, err := conn.BeginTx(ctx, nil)
	if err != nil {
		return nil, 0, err
	}
	if parsedInput.Count {
		count, err := readerRepo.Count(
			ctx, tx, parsedInput.Selectors, parsedInput.Page, entityFactoryFn,
		)
		if err != nil {
			return nil, 0, err
		}
		return nil, count, nil
	}
	entities, err := readerRepo.GetMany(
		ctx,
		tx,
		entityFactoryFn,
		&repositorytypes.GetOptions{
			Selectors: parsedInput.Selectors,
			Orders:    parsedInput.Orders,
			Page:      parsedInput.Page,
		},
	)
	if err != nil {
		return nil, 0, err
	}
	return entities, len(entities), nil
}

// UpdateInvoke executes the update operation.
//
// Parameters:
//   - ctx: The context.
//   - connFn: The database connection function.
//   - entity: The entity to update.
//   - mutatorRepo: The mutator repository.
//   - txManager: The transaction manager.
//
// Returns:
//   - int64: The number of entities updated.
//   - error: Any error that occurred during the operation.
func UpdateInvoke(
	ctx context.Context,
	parsedInput *input.ParsedUpdateEndpointInput,
	connFn repositorytypes.ConnFn,
	entity databasetypes.Mutator,
	mutatorRepo repositorytypes.MutatorRepo[databasetypes.Mutator],
	txManager repositorytypes.TxManager[*int64],
) (int64, error) {
	count, err := txManager.WithTransaction(
		ctx,
		connFn,
		func(ctx context.Context, tx databasetypes.Tx) (*int64, error) {
			c, err := mutatorRepo.Update(
				ctx, tx, entity, parsedInput.Selectors, parsedInput.Updates,
			)
			return &c, err
		})
	if err != nil {
		return 0, err
	}
	return *count, nil
}

// DeleteInvoke executes the delete operation.
//
// Parameters:
//   - ctx: The context.
//   - connFn: The database connection function.
//   - entity: The entity to delete.
//   - mutatorRepo: The mutator repository.
//   - txManager: The transaction manager.
//
// Returns:
//   - int64: The number of entities deleted.
//   - error: Any error that occurred during the operation.
func DeleteInvoke[Entity databasetypes.Mutator](
	ctx context.Context,
	parsedInput *input.ParsedDeleteEndpointInput,
	connFn repositorytypes.ConnFn,
	entity Entity,
	mutatorRepo repositorytypes.MutatorRepo[databasetypes.Mutator],
	txManager repositorytypes.TxManager[*int64],
) (int64, error) {
	count, err := txManager.WithTransaction(
		ctx,
		connFn,
		func(ctx context.Context, tx databasetypes.Tx) (*int64, error) {
			c, err := mutatorRepo.Delete(
				ctx, tx, entity, parsedInput.Selectors, parsedInput.DeleteOpts,
			)
			return &c, err
		})
	if err != nil {
		return 0, err
	}
	return *count, nil
}
