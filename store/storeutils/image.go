package storeutils

import (
	"context"
	"net/http"

	"github.com/pinmonl/pinmonl/model"
	"github.com/pinmonl/pinmonl/store"
)

func SaveImage(ctx context.Context, images *store.Images, content []byte, target model.Morphable, replace bool) (*model.Image, error) {
	if replace {
		_, err := images.DeleteByTarget(ctx, target)
		if err != nil {
			return nil, err
		}
	}

	image := &model.Image{
		Content:     content,
		Size:        len(content),
		ContentType: http.DetectContentType(content),
		TargetID:    target.MorphKey(),
		TargetName:  target.MorphName(),
	}

	err := images.Create(ctx, image)
	if err != nil {
		return nil, err
	}
	return image, nil
}
