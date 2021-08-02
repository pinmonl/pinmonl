package monl

import "errors"

var (
	ErrNoMatch    = errors.New("no matched monler")
	ErrCannotOpen = errors.New("cannot open")
)
