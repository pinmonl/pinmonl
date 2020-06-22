package storeutils

import (
	"context"

	"github.com/pinmonl/pinmonl/model"
	"github.com/pinmonl/pinmonl/store"
)

func CleanupShare(ctx context.Context, sharetags *store.Sharetags, sharepins *store.Sharepins, pinls *store.Pinls, taggables *store.Taggables, shareID string, deletePinl bool) error {
	spList, err := sharepins.List(ctx, &store.SharepinOpts{
		ShareIDs: []string{shareID},
	})
	if err != nil {
		return err
	}

	// Clean up share's pins.
	for _, sp := range spList {
		_, err = taggables.DeleteByTarget(ctx, model.Pinl{ID: sp.PinlID})
		if err != nil {
			return err
		}
		_, err = sharepins.Delete(ctx, sp.ID)
		if err != nil {
			return err
		}

		if deletePinl {
			_, err = pinls.Delete(ctx, sp.PinlID)
			if err != nil {
				return err
			}
		}
	}

	stList, err := sharetags.List(ctx, &store.SharetagOpts{
		ShareIDs: []string{shareID},
	})
	if err != nil {
		return err
	}

	// Clean up share's tags.
	for _, st := range stList {
		_, err = sharetags.Delete(ctx, st.ID)
		if err != nil {
			return err
		}
	}

	return nil
}
