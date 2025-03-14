package input

import (
	"fmt"
	"strings"

	"github.com/pureapi/pureapi-core/apierror"
	"github.com/pureapi/pureapi-core/dbquery"
)

// Predicate is a string representation of a filtering predicate.
type Predicate string

// String returns the string representation of the predicate.
//
// Returns:
//   - The string representation of the predicate.
func (p Predicate) String() string {
	return string(p)
}

// Predicates is a slice of Predicate values.
type Predicates []Predicate

// String returns a string representation of the predicates.
//
// Returns:
//   - A comma-separated string of the predicates.
func (p Predicates) String() string {
	str := make([]string, len(p))
	for i, predicate := range p {
		str[i] = predicate.String()
	}
	return strings.Join(str, ",")
}

// StrSlice returns a slice of strings representing the predicates.
//
// Returns:
//   - A slice of strings representing the predicates.
func (p Predicates) StrSlice() []string {
	str := make([]string, len(p))
	for i, predicate := range p {
		str[i] = predicate.String()
	}
	return str
}

// Predicates for filtering data.
const (
	Greater        Predicate = ">"
	Gt             Predicate = "gt"
	GreaterOrEqual Predicate = ">="
	Ge             Predicate = "ge"
	Equal          Predicate = "="
	Eq             Predicate = "eq"
	NotEqual       Predicate = "!="
	Ne             Predicate = "ne"
	Less           Predicate = "<"
	Lt             Predicate = "LT"
	LessOrEqual    Predicate = "<="
	Le             Predicate = "le"
	In             Predicate = "in"
	NotIn          Predicate = "not_in"
)

// ToDBPredicates maps API-level predicates to database predicates.
var ToDBPredicates = map[Predicate]dbquery.Predicate{
	Greater:        dbquery.Greater,
	Gt:             dbquery.Greater,
	GreaterOrEqual: dbquery.GreaterOrEqual,
	Ge:             dbquery.GreaterOrEqual,
	Equal:          dbquery.Equal,
	Eq:             dbquery.Equal,
	NotEqual:       dbquery.NotEqual,
	Ne:             dbquery.NotEqual,
	Less:           dbquery.Less,
	Lt:             dbquery.Less,
	LessOrEqual:    dbquery.LessOrEqual,
	Le:             dbquery.LessOrEqual,
	In:             dbquery.In,
	NotIn:          dbquery.NotIn,
}

// AllPredicates is a slice of all available predicates.
var AllPredicates = []Predicate{
	Greater, Gt,
	GreaterOrEqual, Ge,
	Equal, Eq,
	NotEqual, Ne,
	Less, Lt,
	LessOrEqual, Le,
	In,
	NotIn,
}

// OnlyEqualPredicates is a slice of predicates that only allow equality.
var OnlyEqualPredicates = []Predicate{Equal, Eq}

// EqualAndNotEqualPredicates is a slice of predicates that allow both equality
// and inequality.
var EqualAndNotEqualPredicates = []Predicate{Equal, Eq, NotEqual, Ne}

// OnlyGreaterPredicates is a slice of predicates that only allow greater
// values.
var OnlyGreaterPredicates = []Predicate{GreaterOrEqual, Ge, Greater, Gt}

// OnlyLessPredicates is a slice of predicates that only allow less values.
var OnlyLessPredicates = []Predicate{LessOrEqual, Le, Less, Lt}

// OnlyInAndNotInPredicates is a slice of predicates that only allow
// IN and NOT_IN.
var OnlyInAndNotInPredicates = []Predicate{In, NotIn}

// DBField is used to translate between API field and database field.
type DBField struct {
	Table  string
	Column string
}

// ErrMaxPageLimitExceededData is the data for the ErrMaxPageLimitExceeded
// error.
type ErrMaxPageLimitExceededData struct {
	MaxLimit int `json:"max_limit"`
}

// ErrMaxPageLimitExceeded is returned when a page limit is exceeded.
var ErrMaxPageLimitExceeded = apierror.NewAPIError("MAX_PAGE_LIMIT_EXCEEDED")

// Page represents a pagination input.
type Page struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

// ToDBPage converts a Page to database Page.
//
// Returns:
//   - The database Page.
func (p *Page) ToDBPage() *dbquery.Page {
	return &dbquery.Page{
		Offset: p.Offset,
		Limit:  p.Limit,
	}
}

// ErrInvalidOrderFieldData is the data for the ErrInvalidOrderField error.
type ErrInvalidOrderFieldData struct {
	Field string `json:"field"`
}

// ErrInvalidOrderField is returned when a field is not allowed.
var ErrInvalidOrderField = apierror.NewAPIError("INVALID_ORDER_FIELD")

// OrderDirection is used to specify the direction of the order.
type OrderDirection string

// String returns the string representation of the order direction.
//
// Returns:
//   - The string representation of the order direction.
func (o OrderDirection) String() string {
	return string(o)
}

// Available order directions.
const (
	DirectionAsc        OrderDirection = "asc"
	DirectionAscending  OrderDirection = "ascending"
	DirectionDesc       OrderDirection = "desc"
	DirectionDescending OrderDirection = "descending"
)

// DirectionsToDB is a map of order directions to database order directions.
var DirectionsToDB = map[OrderDirection]dbquery.OrderDirection{
	DirectionAsc:        dbquery.OrderAsc,
	DirectionAscending:  dbquery.OrderAsc,
	DirectionDesc:       dbquery.OrderDesc,
	DirectionDescending: dbquery.OrderDesc,
}

// Orders is a map of field names to order directions.
type Orders map[string]OrderDirection

// TranslateToDBOrders translates the provided orders into database orders.
// It also returns an error if any of the orders are invalid.
//
//   - orders: The list of orders to translate.
//   - allowedOrderFields: The list of allowed order fields.
//   - apiToDBFieldMap: The mapping of API field names to database field names.
func (o Orders) TranslateToDBOrders(
	apiToDBFieldMap map[string]DBField,
) ([]dbquery.Order, error) {
	dbOrders, err := o.dedup().ToDBOrders(apiToDBFieldMap)
	if err != nil {
		return nil, fmt.Errorf("TranslateToDBOrders: %w", err)
	}
	return dbOrders, nil
}

// dedup deduplicates the provided orders.
func (o Orders) dedup() Orders {
	dedup := map[string]OrderDirection{}
	existing := make(map[string]bool)
	for field := range o {
		order := o[field]
		if !existing[field] {
			dedup[field] = order
			existing[field] = true
		}
	}
	return dedup
}

// ToDBOrders translates the provided orders into database orders.
// It returns an error if any of the orders are invalid.
//
//   - apiToDBFieldMap: The mapping of API field names to database field names.
func (o Orders) ToDBOrders(
	apiToDBFieldMap map[string]DBField,
) ([]dbquery.Order, error) {
	dbOrders := []dbquery.Order{}

	for field, direction := range o {
		translatedField := apiToDBFieldMap[field]

		// Translate field.
		dbColumn := translatedField.Column
		if dbColumn == "" {
			return nil, ErrInvalidOrderField.
				WithData(
					ErrInvalidOrderFieldData{Field: field},
				).
				WithMessage(fmt.Sprintf(
					"cannot translate field: %s", field,
				))
		}

		lowerDir := OrderDirection(strings.ToLower(string(direction)))
		dbOrders = append(dbOrders, dbquery.Order{
			Table:     translatedField.Table,
			Field:     dbColumn,
			Direction: DirectionsToDB[lowerDir],
		})
	}

	return dbOrders, nil
}

// ErrInvalidDatabaseSelectorTranslationData is the data for the
// ErrInvalidDatabaseSelectorTranslation error.
type ErrInvalidDatabaseSelectorTranslationData struct {
	Field string `json:"field"`
}

// ErrInvalidDatabaseSelectorTranslation is returned when a field is not
// allowed.
var ErrInvalidDatabaseSelectorTranslation = apierror.NewAPIError(
	"INVALID_DATABASE_SELECTOR_TRANSLATION",
)

// ErrInvalidPredicateData is the data for the ErrInvalidPredicate error.
type ErrInvalidPredicateData struct {
	Predicate Predicate `json:""`
}

// ErrInvalidPredicate is returned when a predicate is not allowed.
var ErrInvalidPredicate = apierror.NewAPIError("INVALID_PREDICATE")

// ErrInvalidSelectorFieldData is the data for the ErrInvalidSelectorField
// error.
type ErrInvalidSelectorFieldData struct {
	Field string `json:"field"`
}

// ErrInvalidSelectorField is returned when a field is not allowed.
var ErrInvalidSelectorField = apierror.NewAPIError("INVALID_SELECTOR_FIELD")

// ErrPredicateNotAllowedData is the data for the ErrPredicateNotAllowed
// error.
type ErrPredicateNotAllowedData struct {
	Predicate Predicate `json:"predicate"`
}

// ErrPredicateNotAllowed is returned when a predicate is not allowed.
var ErrPredicateNotAllowed = apierror.NewAPIError("PREDICATE_NOT_ALLOWED")

// Selector represents a data selector that specifies criteria for filtering
// data based on fields, predicates, and values.
type Selector struct {
	Predicate Predicate `json:"predicate"` // The predicate to use.
	Value     any       `json:"value"`     // The value to filter on.
}

// NewSelector creates a new selector with the provided predicate and value.
//
// Parameters:
//   - predicate: The predicate to use.
//   - value: The value to filter on.
//
// Returns:
//   - A new selector.
func NewSelector(predicate Predicate, value any) *Selector {
	return &Selector{
		Predicate: predicate,
		Value:     value,
	}
}

// String returns a string representation of the selector.
// It is useful for debugging and logging purposes.
//
// Returns:
//   - A formatted string showing the field, predicate, and value.
func (s Selector) String() string {
	return fmt.Sprintf("%s %v", s.Predicate, s.Value)
}

// Selectors represents a collection of selectors used for filtering data.
// It is a map where the key is the field name and the value is the selector.
type Selectors map[string]Selector

// AddSelector adds a new selector to the collection of selectors.
//
// Parameters:
//   - field: The field name.
//   - predicate: The predicate to use.
//   - value: The value to filter on.
//
// Returns:
//   - A new collection of selectors with the new selector added.
func (s Selectors) AddSelector(
	field string, predicate Predicate, value any,
) Selectors {
	s[field] = Selector{
		Predicate: predicate,
		Value:     value,
	}
	return s
}

// ToDBSelectors converts a slice of API-level selectors to database selectors.
//
// Selectors
// Parameters:
//   - apiToDBFieldMap: A map translating API field names to their corresponding
//     database field definitions.
//
// Returns:
//   - A slice of types.Selector, which represents the translated database
//     selectors.
//   - An error if any validation fails, such as invalid predicates or unknown
//     fields.
func (s Selectors) ToDBSelectors(
	apiToDBFieldMap map[string]DBField,
) ([]dbquery.Selector, error) {
	databaseSelectors := []dbquery.Selector{}

	for field := range s {
		// Translate the predicate.
		selector := s[field]
		lowerPredicate := Predicate(strings.ToLower(string(selector.Predicate)))
		dbPredicate, ok := ToDBPredicates[lowerPredicate]
		if !ok {
			return nil, ErrInvalidPredicate.
				WithData(
					ErrInvalidPredicateData{Predicate: selector.Predicate},
				).
				WithMessage(fmt.Sprintf(
					"cannot translate predicate: %s", selector.Predicate,
				))
		}
		// Translate the field.
		dbField, ok := apiToDBFieldMap[field]
		if !ok {
			return nil, ErrInvalidDatabaseSelectorTranslation.
				WithData(
					ErrInvalidDatabaseSelectorTranslationData{Field: field},
				).
				WithMessage(fmt.Sprintf(
					"cannot translate field: %s", field,
				))
		}
		// Create the database selector.
		databaseSelectors = append(databaseSelectors, dbquery.Selector{
			Table:     dbField.Table,
			Column:    dbField.Column,
			Predicate: dbPredicate,
			Value:     selector.Value,
		})
	}

	return databaseSelectors, nil
}

// ErrInvalidDatabaseUpdateTranslationData is the data for the
// ErrInvalidDatabaseUpdateTranslation error.
type ErrInvalidDatabaseUpdateTranslationData struct {
	Field string `json:"field"`
}

// ErrInvalidDatabaseUpdateTranslation is used to indicate that the
// translation of a database update failed.
var ErrInvalidDatabaseUpdateTranslation = apierror.NewAPIError(
	"INVALID_DATABASE_UPDATE_TRANSLATION",
)

// Updates represents a list of updates to apply to a database entity.
type Updates map[string]any

// ToDBUpdates translates a list of updates to a database update list
// and returns an error if the translation fails.
//
// Parameters:
//   - updates: The list of updates to translate.
//   - apiToDBFieldMap: The mapping of API field names to database field names.
//
// Returns:
//   - A list of database entity updates.
//   - An error if any field translation fails.
func (updates Updates) ToDBUpdates(
	apiToDBFieldMap map[string]DBField,
) ([]dbquery.Update, error) {
	dbUpdates := []dbquery.Update{}

	for field := range updates {
		value := updates[field]

		// Translate the field.
		dbField, ok := apiToDBFieldMap[field]
		if !ok {
			return nil, ErrInvalidDatabaseUpdateTranslation.
				WithData(
					ErrInvalidDatabaseUpdateTranslationData{Field: field},
				).
				WithMessage(fmt.Sprintf(
					"cannot translate field: %s", field,
				))
		}

		dbUpdates = append(dbUpdates, dbquery.Update{
			Field: dbField.Column,
			Value: value,
		})
	}

	return dbUpdates, nil
}
