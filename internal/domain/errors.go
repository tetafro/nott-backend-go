package domain

import "github.com/pkg/errors"

// ErrNotFound is returned when object is not found.
var ErrNotFound = errors.New("not found")
