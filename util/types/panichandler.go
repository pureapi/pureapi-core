package types

import (
	"net/http"
)

type PanicHandler interface {
	HandlePanic(w http.ResponseWriter, r *http.Request, err any)
}
