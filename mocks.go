package pocketsmith

// This file contains common mocks used across both public tests.

import (
	"fmt"
	"net/http"
	"strings"
)

// mockRoundTripper is a mock implementation of the http.RoundTripper interface,
// used in tests in this package to simulate HTTP responses.
type mockRoundTripper struct {
	MockFunc func(req *http.Request) *http.Response
}

// RoundTrip executes the mock logic for handling HTTP requests in tests.
// It simulates the behavior of a real RoundTripper by either invoking the
// MockFunc or returning an error if the HTTP method contains the word "error".
func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if strings.Contains(req.Method, "error") {
		return nil, fmt.Errorf("an error occurred")
	}
	return m.MockFunc(req), nil
}
