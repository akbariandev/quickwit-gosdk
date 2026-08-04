// Package types defines the shared data types for the Quickwit Go SDK.
package types

import "fmt"

// Error represents an error returned by the Quickwit API.
type Error struct {
	StatusCode int
	Message    string
}

func (e *Error) Error() string {
	return fmt.Sprintf("quickwit api error (status %d): %s", e.StatusCode, e.Message)
}
