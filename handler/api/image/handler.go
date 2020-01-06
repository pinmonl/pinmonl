package image

import (
	"context"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"

	"github.com/pinmonl/pinmonl/handler/api/request"
	"github.com/pinmonl/pinmonl/model"
	"github.com/pinmonl/pinmonl/store"
)

// HandleFind returns the content of image.
func HandleFind(images store.ImageStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m, _ := request.ImageFrom(r.Context())
		w.WriteHeader(200)
		w.Write(m.Content)
	}
}

// Upload creates image from multipart form.
func Upload(
	ctx context.Context,
	images store.ImageStore,
	fileheader *multipart.FileHeader,
) (*model.Image, error) {
	file, err := fileheader.Open()
	if err != nil {
		return nil, err
	}
	return UploadFromReader(ctx, images, file)
}

// UploadFromReader creates image from io.Reader.
func UploadFromReader(
	ctx context.Context,
	images store.ImageStore,
	r io.Reader,
) (*model.Image, error) {
	c, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	m := model.Image{Content: c}
	err = images.Create(ctx, &m)
	if err != nil {
		return nil, err
	}
	return &m, nil
}
