package input

import (
	"github.com/pureapi/pureapi-core/apierror"
	"github.com/pureapi/pureapi-core/dbquery"
	repositorytypes "github.com/pureapi/pureapi-core/repository/types"
)

// APIToDBFields maps API fields to database fields.
type APIToDBFields map[string]DBField

// Selector parse errors.
var (
	ErrNeedAtLeastOneSelector = apierror.NewAPIError("NEED_AT_LEAST_ONE_SELECTOR")
	ErrNeedAtLeastOneUpdate   = apierror.NewAPIError("NEED_AT_LEAST_ONE_UPDATE")
)

// ParsedGetEndpointInput represents a parsed get endpoint input.
type ParsedGetEndpointInput struct {
	Selectors dbquery.Selectors
	Orders    []dbquery.Order
	Page      *dbquery.Page
	Count     bool
}

// ParsedUpdateEndpointInput represents a parsed update endpoint input.
type ParsedUpdateEndpointInput struct {
	Selectors dbquery.Selectors
	Updates   []dbquery.Update
	Upsert    bool
}

// ParsedDeleteEndpointInput represents a parsed delete endpoint input.
type ParsedDeleteEndpointInput struct {
	Selectors  dbquery.Selectors
	DeleteOpts *repositorytypes.DeleteOptions
}

// ParseGetInput translates API parameters to DB parameters.
//
// Parameters:
//   - apiToDBFields: A map translating API field names to their corresponding
//     database field definitions.
//   - selectors: A slice of API-level selectors.
//   - orders: A slice of API-level orders.
//   - inputPage: A pointer to the input page.
//   - maxPage: The maximum page size.
//   - count: A boolean indicating whether to return the count.
//
// Returns:
//   - *ParsedGetEndpointInput: A pointer to the parsed get endpoint input.
//   - error: An error if the input is invalid.
func ParseGetInput(
	apiToDBFields APIToDBFields,
	selectors Selectors,
	orders Orders,
	inputPage *Page,
	maxPage int,
	count bool,
) (*ParsedGetEndpointInput, error) {
	dbOrders, err := orders.TranslateToDBOrders(apiToDBFields)
	if err != nil {
		return nil, err
	}
	if inputPage == nil {
		inputPage = &Page{Offset: 0, Limit: maxPage}
	}
	dbSelectors, err := selectors.ToDBSelectors(apiToDBFields)
	if err != nil {
		return nil, err
	}
	return &ParsedGetEndpointInput{
		Orders:    dbOrders,
		Selectors: dbSelectors,
		Page:      inputPage.ToDBPage(),
		Count:     count,
	}, nil
}

// ParseUpdateInput translates API update input into DB update input.
//
// Parameters:
//   - apiToDBFields: A map translating API field names to their corresponding
//     database field definitions.
//   - selectors: A slice of API-level selectors.
//   - updates: A map of API-level updates.
//   - upsert: A boolean indicating whether to upsert.
//
// Returns:
//   - *ParsedUpdateEndpointInput: A pointer to the parsed update endpoint
//     input.
//   - error: An error if the input is invalid.
func ParseUpdateInput(
	apiToDBFields APIToDBFields,
	selectors Selectors,
	updates Updates,
	upsert bool,
) (*ParsedUpdateEndpointInput, error) {
	dbSelectors, err := selectors.ToDBSelectors(apiToDBFields)
	if err != nil {
		return nil, err
	}
	if len(dbSelectors) == 0 {
		return nil, ErrNeedAtLeastOneSelector
	}
	dbUpdates, err := updates.ToDBUpdates(apiToDBFields)
	if err != nil {
		return nil, err
	}
	if len(dbUpdates) == 0 {
		return nil, ErrNeedAtLeastOneUpdate
	}
	return &ParsedUpdateEndpointInput{
		Selectors: dbSelectors,
		Updates:   dbUpdates,
		Upsert:    upsert,
	}, nil
}

// ParseDeleteInput translates API delete input into DB delete input.
//
// Parameters:
//   - apiToDBFields: A map translating API field names to their corresponding
//     database field definitions.
//   - selectors: A slice of API-level selectors.
//   - orders: A slice of API-level orders.
//   - limit: The maximum number of entities to delete.
//
// Returns:
//   - *ParsedDeleteEndpointInput: A pointer to the parsed delete endpoint
//     input.
//   - error: An error if the input is invalid.
func ParseDeleteInput(
	apiToDBFields APIToDBFields,
	selectors Selectors,
	orders Orders,
	limit int,
) (*ParsedDeleteEndpointInput, error) {
	dbSelectors, err := selectors.ToDBSelectors(apiToDBFields)
	if err != nil {
		return nil, err
	}
	if len(dbSelectors) == 0 {
		return nil, ErrNeedAtLeastOneSelector
	}
	dbOrders, err := orders.TranslateToDBOrders(apiToDBFields)
	if err != nil {
		return nil, err
	}
	return &ParsedDeleteEndpointInput{
		Selectors: dbSelectors,
		DeleteOpts: &repositorytypes.DeleteOptions{
			Limit:  limit,
			Orders: dbOrders,
		},
	}, nil
}
