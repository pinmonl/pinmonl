package store

import "errors"

var (
	// ErrMissingID indicates the ID is not provided.
	ErrMissingID = errors.New("ID is missing")
)
