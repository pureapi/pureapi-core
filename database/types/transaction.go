package types

import (
	"context"
)

// TxFn is a function that takes in a transaction, and returns a result and an
// error.
type TxFn[Result any] func(ctx context.Context, tx Tx) (Result, error)
