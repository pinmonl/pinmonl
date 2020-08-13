package provider

import "errors"

type Error error

var (
	ErrNotSupport  = Error(errors.New("provider not support"))
	ErrNotFound    = Error(errors.New("uri not found"))
	ErrNoPing      = Error(errors.New("provider does not ping"))
	ErrEndOfReport = Error(errors.New("end of report"))
)
