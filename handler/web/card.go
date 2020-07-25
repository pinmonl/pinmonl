package web

import (
	"encoding/base64"
	"net/http"
	"net/url"

	"github.com/go-chi/chi"
	"github.com/pinmonl/pinmonl/pkgs/card"
	"github.com/pinmonl/pinmonl/pkgs/response"
)

func (s *Server) fetchCardHandler(w http.ResponseWriter, r *http.Request) {
	dest, err := url.Parse(chi.URLParam(r, "*"))
	if err != nil {
		response.JSON(w, err, http.StatusBadRequest)
		return
	}

	c, err := card.NewCard(dest.String())
	if err != nil {
		response.JSON(w, err, http.StatusBadRequest)
		return
	}

	var body struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		ImageData   string `json:"imageData"`
	}

	body.Title = c.Title()
	body.Description = c.Description()

	if img, err := c.Image(); err == nil && img != nil {
		body.ImageData = base64.StdEncoding.EncodeToString(img)
	}

	response.JSON(w, body, http.StatusOK)
}
