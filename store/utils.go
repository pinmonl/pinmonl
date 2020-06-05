package store

import (
	"time"

	"github.com/pinmonl/pinmonl/model/field"
	"github.com/rs/xid"
)

func newID() string {
	return xid.New().String()
}

func timestamp() field.Time {
	t := time.Now().Round(time.Second).UTC()
	return (field.Time)(t)
}

func prefixStrings(arr []string, prefix string) []string {
	arr2 := make([]string, len(arr))
	for i, s := range arr {
		arr2[i] = prefix + s
	}
	return arr2
}
