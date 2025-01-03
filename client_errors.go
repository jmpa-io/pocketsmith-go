package pocketsmith

import "fmt"

// ErrClientEmptyToken is returned when no token is provided to the client.
type ErrClientEmptyToken struct {
}

func (e ErrClientEmptyToken) Error() string {
	return "the provided token is empty"
}

// ErrClientFailedToSetOption is returned when an option encounters an error
// when trying to be set with the client.
type ErrClientFailedToSetOption struct {
	err error
}

func (e ErrClientFailedToSetOption) Error() string {
	return fmt.Sprintf("failed to set option in client: %v", e.err)
}
