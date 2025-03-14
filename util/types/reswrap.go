package types

import (
	"bufio"
	"net"
	"net/http"
)

type ResWrap interface {
	http.ResponseWriter
	StatusCode() int
	Header() http.Header
	WriteHeader(statusCode int)
	Body() []byte
	Write(data []byte) (int, error)
	Flush()
	Hijack() (net.Conn, *bufio.ReadWriter, error)
}
