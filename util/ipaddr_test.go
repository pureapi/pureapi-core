package util

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestRequestIPAddress tests the RequestIPAddress function.
func TestRequestIPAddress(t *testing.T) {
	testCases := []struct {
		name           string
		remoteAddr     string
		xForwardedFor  string
		expectedOutput string
	}{
		{
			name:           "X-Forwarded-For present with single IP",
			remoteAddr:     "9.9.9.9:12345",
			xForwardedFor:  "1.2.3.4",
			expectedOutput: "1.2.3.4",
		},
		{
			name:           "X-Forwarded-For present with multiple IPs",
			remoteAddr:     "9.9.9.9:12345",
			xForwardedFor:  " 1.2.3.4 , 5.6.7.8",
			expectedOutput: "1.2.3.4",
		},
		{
			name:           "No X-Forwarded-For header, valid RemoteAddr",
			remoteAddr:     "9.10.11.12:34567",
			xForwardedFor:  "",
			expectedOutput: "9.10.11.12",
		},
		{
			name:           "No X-Forwarded-For header, invalid RemoteAddr",
			remoteAddr:     "invalid-addr",
			xForwardedFor:  "",
			expectedOutput: "invalid-addr",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "http://example.com", nil)
			req.RemoteAddr = tc.remoteAddr
			if tc.xForwardedFor != "" {
				req.Header.Set("X-Forwarded-For", tc.xForwardedFor)
			}

			actual := RequestIPAddress(req)
			assert.Equal(t, tc.expectedOutput, actual)
		})
	}
}
