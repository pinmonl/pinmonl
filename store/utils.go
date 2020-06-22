package store

import (
	"github.com/pinmonl/pinmonl/model/field"
	"github.com/rs/xid"
)

func newID() string {
	return xid.New().String()
}

func timestamp() field.Time {
	return field.Now()
}

func prefixStrings(arr []string, prefix string) []string {
	arr2 := make([]string, len(arr))
	for i, s := range arr {
		arr2[i] = prefix + s
	}
	return arr2
}
