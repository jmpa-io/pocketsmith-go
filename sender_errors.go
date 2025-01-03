package pocketsmith

import (
	"fmt"
)

// ErrSenderFailedSetupRequest is returned whenever the sender fails to
// setup a new *http.Request.
type ErrSenderFailedSetupRequest struct {
	err error
}

func (e ErrSenderFailedSetupRequest) Error() string {
	return fmt.Sprintf("failed to setup http request: %v", e.err)
}

// ErrSenderFailedSendRequest is returned whenever the sender fails to send
// a new *http.Request to the API.
type ErrSenderFailedSendRequest struct {
	err error
}

func (e ErrSenderFailedSendRequest) Error() string {
	return fmt.Sprintf("failed to send http request: %v", e.err)
}

// ErrSenderFailedParseResponse is returned when the sender fails to parse a
// response from the API.
type ErrSenderFailedParseResponse struct {
	err error
}

func (e ErrSenderFailedParseResponse) Error() string {
	return fmt.Sprintf("failed to parse response: %v", e.err)
}

// ErrSenderInvalidResponse is returned when the sender receives an error
// response specifically from the API.
type ErrSenderInvalidResponse struct {
	resp       apiErrorResponse
	statusCode int
}

func (e ErrSenderInvalidResponse) Error() string {
	return fmt.Sprintf(
		"error response returned from API; status_code=%v, error=%s",
		e.statusCode,
		e.resp.Error,
	)
}
