package api

import (
	"context"
	"net/http"

	databasetypes "github.com/pureapi/pureapi-core/database/types"
	"github.com/pureapi/pureapi-core/input"
	repositorytypes "github.com/pureapi/pureapi-core/repository/types"
)

// CreateInvokeFn is the function invokes the create endpoint.
type CreateInvokeFn[Entity databasetypes.Mutator] func(
	ctx context.Context, entity Entity,
) (Entity, error)

// CreateEntityFactoryFn is the function that creates a new entity.
type CreateEntityFactoryFn[Input any, Entity databasetypes.Mutator] func(
	ctx context.Context, input *Input,
) (Entity, error)

// ToCreateOutputFn is the function that converts the entity to the endpoint
// output.
type ToCreateOutputFn[Entity any] func(entity Entity) (any, error)

// BeforeCreateCallback is the function that runs before the create operation.
// It can be used to modify the entity before it is created.
type BeforeCreateCallback[Input any, Entity databasetypes.Mutator] func(
	w http.ResponseWriter, r *http.Request, entity *Entity, input *Input,
) (Entity, error)

// createHandler is the handler implementation for the create endpoint.
type createHandler[Entity databasetypes.Mutator, Input any] struct {
	entityFactoryFn CreateEntityFactoryFn[Input, Entity]
	createInvokeFn  CreateInvokeFn[Entity]
	toOutputFn      ToCreateOutputFn[Entity]
	beforeCallback  BeforeCreateCallback[Input, Entity]
}

// NewCreateHandler creates a new create handler.
//
// Parameters:
//   - entityFactoryFn: The function that creates a new entity.
//   - createInvokeFn: The function that invokes the create endpoint.
//   - toOutputFn: The function that converts the entity to the endpoint output.
//   - beforeCallback: The function that runs before the create operation.
//
// Returns:
//   - *createHandler: The new create handler.
func NewCreateHandler[Entity databasetypes.Mutator, Input any](
	entityFactoryFn CreateEntityFactoryFn[Input, Entity],
	createInvokeFn CreateInvokeFn[Entity],
	toOutputFn ToCreateOutputFn[Entity],
	beforeCallback BeforeCreateCallback[Input, Entity],
) *createHandler[Entity, Input] {
	return &createHandler[Entity, Input]{
		entityFactoryFn: entityFactoryFn,
		createInvokeFn:  createInvokeFn,
		toOutputFn:      toOutputFn,
		beforeCallback:  beforeCallback,
	}
}

// Handle processes the create endpoint.
//
// Parameters:
//   - w: The response writer.
//   - r: The request.
//   - i: The input.
//
// Returns:
//   - any: The endpoint output.
//   - error: An error if the request fails.
func (h *createHandler[Mutator, Input]) Handle(
	w http.ResponseWriter, r *http.Request, i *Input,
) (any, error) {
	entity, err := h.entityFactoryFn(r.Context(), i)
	if err != nil {
		return nil, err
	}
	if h.beforeCallback != nil {
		entity, err = h.beforeCallback(w, r, &entity, i)
		if err != nil {
			return nil, err
		}
	}
	createdEntity, err := h.createInvokeFn(r.Context(), entity)
	if err != nil {
		return nil, err
	}
	return h.toOutputFn(createdEntity)
}

// GetInvokeFn is the function that invokes the get endpoint.
type GetInvokeFn[Entity databasetypes.Getter] func(
	ctx context.Context,
	parsedInput *input.ParsedGetEndpointInput,
	entityFactoryFn repositorytypes.GetterFactoryFn[Entity],
) ([]Entity, int, error)

// ToGetOutputFn is the function that converts the entities to the endpoint
// output.
type ToGetOutputFn[Entity any, Output any] func(
	entities []Entity, count int,
) (*Output, error)

// BeforeGetCallback is the function that runs before the get operation.
// It can be used to modify the parsed input before it is used.
type BeforeGetCallback[Input any, Entity databasetypes.Getter] func(
	w http.ResponseWriter,
	r *http.Request,
	parsedInput *input.ParsedGetEndpointInput,
	input *Input,
) (*input.ParsedGetEndpointInput, error)

// getHandler is the handler for the get endpoint.
type getHandler[Entity databasetypes.Getter, Input any, Output any] struct {
	parseInputFn    func(input *Input) (*input.ParsedGetEndpointInput, error)
	getInvokeFn     GetInvokeFn[Entity]
	toOutputFn      ToGetOutputFn[Entity, Output]
	entityFactoryFn repositorytypes.GetterFactoryFn[Entity]
	beforeCallback  BeforeGetCallback[Input, Entity]
}

// NewGetHandler creates a new get handler.
//
// Parameters:
//   - parseInputFn: The function that parses the input.
//   - getInvokeFn: The function that invokes the get endpoint.
//   - toOutputFn: The function that converts the entities to the endpoint
//     output.
//   - entityFactoryFn: The function that creates a new entity.
//   - beforeCallback: The function that runs before the get operation.
//
// Returns:
//   - *getHandler: The new get handler.
func NewGetHandler[Entity databasetypes.Getter, Input any, Output any](
	parseInputFn func(input *Input) (*input.ParsedGetEndpointInput, error),
	getInvokeFn GetInvokeFn[Entity],
	toOutputFn ToGetOutputFn[Entity, Output],
	entityFactoryFn repositorytypes.GetterFactoryFn[Entity],
	beforeCallback BeforeGetCallback[Input, Entity],
) *getHandler[Entity, Input, Output] {
	return &getHandler[Entity, Input, Output]{
		parseInputFn:    parseInputFn,
		getInvokeFn:     getInvokeFn,
		toOutputFn:      toOutputFn,
		entityFactoryFn: entityFactoryFn,
		beforeCallback:  beforeCallback,
	}
}

// Handle processes the get endpoint.
//
// Parameters:
//   - w: The response writer.
//   - r: The request.
//   - i: The input.
//
// Returns:
//   - any: The endpoint output.
//   - error: An error if the request fails.
func (h *getHandler[Entity, Input, Output]) Handle(
	w http.ResponseWriter, r *http.Request, i *Input,
) (any, error) {
	parsedInput, err := h.parseInputFn(i)
	if err != nil {
		return nil, err
	}
	if h.beforeCallback != nil {
		parsedInput, err = h.beforeCallback(w, r, parsedInput, i)
		if err != nil {
			return nil, err
		}
	}
	entities, count, err := h.getInvokeFn(
		r.Context(), parsedInput, h.entityFactoryFn,
	)
	if err != nil {
		return nil, err
	}
	return h.toOutputFn(entities, count)
}

// UpdateInvokeFn is the function that invokes the update endpoint.
type ToUpdateOutputFn func(count int64) (any, error)

// UpdateEntityFactoryFn is the function that creates a new entity.
type UpdateEntityFactoryFn func() databasetypes.Mutator

// UpdateInvokeFn is the function that invokes the update endpoint.
type UpdateInvokeFn func(
	ctx context.Context,
	parsedInput *input.ParsedUpdateEndpointInput,
	updater databasetypes.Mutator,
) (int64, error)

// BeforeUpdateCallback is the function that runs before the update operation.
// It can be used to modify the parsed input and entity before they are used.
type BeforeUpdateCallback[Input any] func(
	w http.ResponseWriter,
	r *http.Request,
	parsedInput *input.ParsedUpdateEndpointInput,
	entity databasetypes.Mutator,
	input *Input,
) (*input.ParsedUpdateEndpointInput, databasetypes.Mutator, error)

// updateHandler is the handler implementation for the update endpoint.
type updateHandler[Input any] struct {
	parseInputFn    func(input *Input) (*input.ParsedUpdateEndpointInput, error)
	updateInvokeFn  UpdateInvokeFn
	toOutputFn      ToUpdateOutputFn
	entityFactoryFn UpdateEntityFactoryFn
	beforeCallback  BeforeUpdateCallback[Input]
}

// NewUpdateHandler creates a new update handler.
//
// Parameters:
//   - parseInputFn: The function that parses the input.
//   - updateInvokeFn: The function that invokes the update endpoint.
//   - toOutputFn: The function that converts the entities to the endpoint
//     output.
//   - entityFactoryFn: The function that creates a new entity.
//   - beforeCallback: The function that runs before the update operation.
//
// Returns:
//   - *updateHandler: The new update handler.
func NewUpdateHandler[Input any](
	parseInputFn func(input *Input) (*input.ParsedUpdateEndpointInput, error),
	updateInvokeFn UpdateInvokeFn,
	toOutputFn ToUpdateOutputFn,
	entityFactoryFn UpdateEntityFactoryFn,
	beforeCallback BeforeUpdateCallback[Input],
) *updateHandler[Input] {
	return &updateHandler[Input]{
		parseInputFn:    parseInputFn,
		updateInvokeFn:  updateInvokeFn,
		toOutputFn:      toOutputFn,
		entityFactoryFn: entityFactoryFn,
		beforeCallback:  beforeCallback,
	}
}

// Handle processes the update endpoint.
//
// Parameters:
//   - w: The response writer.
//   - r: The request.
//   - i: The input.
//
// Returns:
//   - any: The endpoint output.
//   - error: An error if the request fails.
func (h *updateHandler[Input]) Handle(
	w http.ResponseWriter, r *http.Request, i *Input,
) (any, error) {
	parsedInput, err := h.parseInputFn(i)
	if err != nil {
		return nil, err
	}
	entity := h.entityFactoryFn()
	if h.beforeCallback != nil {
		parsedInput, entity, err = h.beforeCallback(
			w, r, parsedInput, entity, i,
		)
		if err != nil {
			return nil, err
		}
	}
	count, err := h.updateInvokeFn(r.Context(), parsedInput, entity)
	if err != nil {
		return nil, err
	}
	return h.toOutputFn(count)
}

// DeleteInvokeFn is the function that invokes the delete endpoint.
type ToDeleteOutputFn func(count int64) (any, error)

// DeleteEntityFactoryFn is the function that creates a new entity.
type DeleteEntityFactoryFn func() databasetypes.Mutator

// DeleteInvokeFn is the function that invokes the delete endpoint.
type DeleteInvokeFn func(
	ctx context.Context,
	parsedInput *input.ParsedDeleteEndpointInput,
	entity databasetypes.Mutator,
) (int64, error)

// BeforeDeleteCallback is the function that runs before the delete operation.
// It can be used to modify the parsed input and entity before they are used.
type BeforeDeleteCallback[Input any] func(
	w http.ResponseWriter,
	r *http.Request,
	parsedInput *input.ParsedDeleteEndpointInput,
	entity databasetypes.Mutator,
	input *Input,
) (*input.ParsedDeleteEndpointInput, databasetypes.Mutator, error)

// deleteHandler is the handler implementation for the delete endpoint.
type deleteHandler[Input any] struct {
	parseInputFn    func(input *Input) (*input.ParsedDeleteEndpointInput, error)
	deleteInvokeFn  DeleteInvokeFn
	toOutputFn      ToDeleteOutputFn
	entityFactoryFn DeleteEntityFactoryFn
	beforeCallback  BeforeDeleteCallback[Input]
}

// NewDeleteHandler creates a new delete handler.
//
// Parameters:
//   - parseInputFn: The function that parses the input.
//   - deleteInvokeFn: The function that invokes the delete endpoint.
//   - toOutputFn: The function that converts the entities to the endpoint
//     output.
//   - entityFactoryFn: The function that creates a new entity.
//   - beforeCallback: The function that runs before the delete operation.
//
// Returns:
//   - *deleteHandler: The new delete handler.
func NewDeleteHandler[Input any](
	parseInputFn func(input *Input) (*input.ParsedDeleteEndpointInput, error),
	deleteInvokeFn DeleteInvokeFn,
	toOutputFn ToDeleteOutputFn,
	entityFactoryFn DeleteEntityFactoryFn,
	beforeCallback BeforeDeleteCallback[Input],
) *deleteHandler[Input] {
	return &deleteHandler[Input]{
		parseInputFn:    parseInputFn,
		deleteInvokeFn:  deleteInvokeFn,
		toOutputFn:      toOutputFn,
		entityFactoryFn: entityFactoryFn,
		beforeCallback:  beforeCallback,
	}
}

// Handle processes the delete endpoint.
//
// Parameters:
//   - w: The response writer.
//   - r: The request.
//   - i: The input.
//
// Returns:
//   - any: The endpoint output.
//   - error: An error if the request fails.
func (h *deleteHandler[Input]) Handle(
	w http.ResponseWriter, r *http.Request, i *Input,
) (any, error) {
	parsedInput, err := h.parseInputFn(i)
	if err != nil {
		return nil, err
	}
	entity := h.entityFactoryFn()
	if h.beforeCallback != nil {
		parsedInput, entity, err = h.beforeCallback(
			w, r, parsedInput, entity, i,
		)
		if err != nil {
			return nil, err
		}
	}
	count, err := h.deleteInvokeFn(r.Context(), parsedInput, entity)
	if err != nil {
		return nil, err
	}
	return h.toOutputFn(count)
}
