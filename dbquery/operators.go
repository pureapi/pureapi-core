package dbquery

// Predicates for filtering data.
const (
	Greater        Predicate = ">"
	GreaterOrEqual Predicate = ">="
	Equal          Predicate = "="
	NotEqual       Predicate = "!="
	Less           Predicate = "<"
	LessOrEqual    Predicate = "<="
	In             Predicate = "IN"
	NotIn          Predicate = "NOT IN"
	Like           Predicate = "LIKE"
	NotLike        Predicate = "NOT LIKE"
)

// Order directions.
const (
	OrderAsc  OrderDirection = "ASC"
	OrderDesc OrderDirection = "DESC"
)

// Join types.
const (
	JoinTypeInner JoinType = "INNER"
	JoinTypeLeft  JoinType = "LEFT"
	JoinTypeRight JoinType = "RIGHT"
	JoinTypeFull  JoinType = "FULL"
)
