package errors

import "fmt"

// ErrNotFound is returned when object is not found.
var ErrNotFound = fmt.Errorf("not found")
