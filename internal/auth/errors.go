package auth

import "errors"

// ErrNotFound returned when record not found in database.
var ErrNotFound = errors.New("record not found")
