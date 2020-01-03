package store

import (
	"time"

	"github.com/pinmonl/pinmonl/model/field"
	"github.com/rs/xid"
)

func newUID() string {
	return xid.New().String()
}

func timestamp() field.Time {
	t := time.Now().Truncate(time.Nanosecond)
	return (field.Time)(t)
}
