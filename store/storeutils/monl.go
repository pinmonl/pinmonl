package storeutils

import (
	"context"

	"github.com/pinmonl/pinmonl/model"
	"github.com/pinmonl/pinmonl/store"
)

func FindOrCreateMonl(ctx context.Context, monls *store.Monls, rawurl string) (*model.Monl, bool, error) {
	found, err := monls.List(ctx, &store.MonlOpts{
		URL: rawurl,
	})
	if err != nil {
		return nil, false, err
	}

	var (
		monl  *model.Monl
		isNew bool
	)

	if len(found) > 0 {
		monl = found[0]
	} else {
		monl = &model.Monl{URL: rawurl}
		err := monls.Create(ctx, monl)
		if err != nil {
			return nil, false, err
		}
		isNew = true
	}

	return monl, isNew, nil
}
