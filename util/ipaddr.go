package util

import (
	"net"
	"net/http"
	"strings"
)

const xForwardedFor = "X-Forwarded-For"

// RequestIPAddress returns the IP address of the request.
// It checks for the X-Forwarded-For header and falls back to the RemoteAddr.
//
// Parameters:
//   - request: The HTTP request.
//
// Returns:
//   - string: The IP address of the request.
func RequestIPAddress(request *http.Request) string {
	forwarded := request.Header.Get(xForwardedFor)
	if forwarded != "" {
		ips := strings.Split(forwarded, ",")
		return strings.TrimSpace(ips[0])
	}
	ip, _, err := net.SplitHostPort(request.RemoteAddr)
	if err != nil {
		return request.RemoteAddr
	}
	return ip
}
