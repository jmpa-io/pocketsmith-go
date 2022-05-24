package pocketsmith

import (
	"fmt"
)

// ErrMissingToken is returned when no token is provided to the client.
type ErrMissingToken struct {
}

func (e ErrMissingToken) Error() string {
	return "missing token"
}

// ErrFailedOptionSet is returned when an option encounters an error when trying to be set.
type ErrFailedOptionSet struct {
	err error
}

func (e ErrFailedOptionSet) Error() string {
	return fmt.Sprintf("failed to set option: %v", e.err)
}

// ErrFailedLoggerSet is returned when the client fails to setup the logger.
type ErrFailedLoggerSetup struct {
	err error
}

func (e ErrFailedLoggerSetup) Error() string {
	return fmt.Sprintf("failed to setup client logger: %v", e.err)
}

// ErrFailedFindingLogLevel is returned when this package is unable to find a log level from a given string.
type ErrFailedFindingLogLevel struct {
	given string
}

func (e ErrFailedFindingLogLevel) Error() string {
	return fmt.Sprintf("failed to find log level for given string: %s is not a valid logging level for this package", e.given)
}

// ErrFailedMarshal is returned whenever this package has an error returned from json.Marshal.
type ErrFailedMarshal struct {
	err error
}

func (e ErrFailedMarshal) Error() string {
	return fmt.Sprintf("failed to marshal data: %v", e.err)
}

// ErrFailedUnmarshal is returned whenever this package has an error returned from json.Unmarshal.
type ErrFailedUnmarshal struct {
	err error
}

func (e ErrFailedUnmarshal) Error() string {
	return fmt.Sprintf("failed to unmarshal data: %v", e.err)
}

// ErrSenderFailedSetupRequest is returned whenever the sender fails to setup a request.
type ErrSenderFailedSetupRequest struct {
	err error
}

func (e ErrSenderFailedSetupRequest) Error() string {
	return fmt.Sprintf("failed to setup request: %v", e.err)
}

// ErrSenderFailedSendRequest is returned whenever the sender fails to send a request.
type ErrSenderFailedSendRequest struct {
	err error
}

func (e ErrSenderFailedSendRequest) Error() string {
	return fmt.Sprintf("failed to send request: %v", e.err)
}

// ErrSenderFailedParseResponse is returned when the sender fails to parse a response
// from the API, for a given request.
type ErrSenderFailedParseResponse struct {
	err error
}

func (e ErrSenderFailedParseResponse) Error() string {
	return fmt.Sprintf("failed to parse response: %v", e.err)
}

// ErrSenderInvalidResponse is returned when the sender receives an uncaught or
// unexpected error from the API.
type ErrSenderInvalidResponse struct {
	err        string
	statusCode int
}

func (e ErrSenderInvalidResponse) Error() string {
	return fmt.Sprintf("error returned from API; status_code=%v, error=%q", e.statusCode, e.err)
}
