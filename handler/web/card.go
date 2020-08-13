package web

import (
	"encoding/base64"
	"errors"
	"net/http"
	"net/url"

	"github.com/pinmonl/pinmonl/pkgs/card"
	"github.com/pinmonl/pinmonl/pkgs/response"
)

func (s *Server) fetchCardHandler(w http.ResponseWriter, r *http.Request) {
	rawurl := r.URL.Query().Get("url")
	if rawurl == "" {
		response.JSON(w, errors.New("url is required"), http.StatusBadRequest)
		return
	}
	u, err := url.Parse(rawurl)
	if err != nil {
		response.JSON(w, err, http.StatusBadRequest)
		return
	}

	// Use dest and pass down query string.
	c, err := card.NewCard(u.String())
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
