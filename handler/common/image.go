package common

import (
	"context"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"github.com/pinmonl/pinmonl/model"
	"github.com/pinmonl/pinmonl/pkgs/request"
	"github.com/pinmonl/pinmonl/pkgs/response"
	"github.com/pinmonl/pinmonl/store"
	"github.com/pinmonl/pinmonl/store/storeutils"
)

func BindImage(images *store.Images, paramName string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			var (
				ctx     = r.Context()
				imageId = chi.URLParam(r, paramName)
			)

			image, err := images.Find(ctx, imageId)
			if err != nil {
				response.JSON(w, err, http.StatusInternalServerError)
				return
			}
			if image == nil {
				response.JSON(w, nil, http.StatusNotFound)
				return
			}

			ctx = request.WithImage(ctx, image)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}

func ImageHandler() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		var (
			ctx   = r.Context()
			image = request.ImageFrom(ctx)
		)

		w.Write(image.Content)
	}
	return http.HandlerFunc(fn)
}

func ImageUpload(ctx context.Context, r *http.Request, images *store.Images, target model.Morphable, maxMemory int64, replace bool) (image *model.Image, code int, outerr error) {
	r.ParseMultipartForm(maxMemory)
	file, _, err := r.FormFile("file")
	if err != nil {
		code = http.StatusBadRequest
		return
	}
	defer file.Close()

	content, err := ioutil.ReadAll(file)
	if err != nil {
		code = http.StatusBadRequest
		return
	}

	mime := http.DetectContentType(content)
	if !strings.HasPrefix(mime, "image/") {
		code = http.StatusBadRequest
		return
	}

	image2, err := storeutils.SaveImage(ctx, images, content, target, replace)
	if err != nil {
		outerr, code = err, http.StatusInternalServerError
		return
	}
	image = image2
	return
}
