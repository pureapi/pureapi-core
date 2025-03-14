package types

import (
	"net/http"
)

type ReqWrap interface {
	GetRequest() *http.Request
	GetBody() []byte
}
