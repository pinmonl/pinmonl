package store

import (
	"fmt"
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

func bindQueryIDs(prefix string, ids []string) ([]string, map[string]string) {
	ks := make([]string, len(ids))
	out := make(map[string]string)
	for i, id := range ids {
		k := fmt.Sprintf("%s%d", prefix, i)
		out[k] = id
		ks[i] = ":" + k
	}
	return ks, out
}
