package quickwitgosdk

import "fmt"

// QuickwitError represents an error returned by the Quickwit API.
type QuickwitError struct {
	StatusCode int
	Message    string
}

func (e *QuickwitError) Error() string {
	return fmt.Sprintf("quickwit api error (status %d): %s", e.StatusCode, e.Message)
}
