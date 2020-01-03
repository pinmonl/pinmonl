package pinl

import (
	"fmt"
	"io"

	"github.com/pinmonl/pinmonl/handler/api/request"
	"github.com/pinmonl/pinmonl/model"
	"github.com/pinmonl/pinmonl/validate"
)

// Input defines the accepted data from client.
type Input struct {
	URL         string   `json:"url"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Readme      string   `json:"readme"`
	Tags        []string `json:"tags"`
}

// ReadInput parses string as JSON.
func ReadInput(r io.Reader) (*Input, error) {
	var in Input
	err := request.JSON(r, &in)
	if err != nil {
		return nil, err
	}
	return &in, nil
}

// Validate checks fields in data.
func (in *Input) Validate() error {
	var es validate.Errors

	switch u := in.URL; {
	case u == "":
		es = append(es, fmt.Errorf("url cannot be empty"))
	case !validate.IsURL(u):
		es = append(es, fmt.Errorf("url is not valid"))
	}

	for _, t := range in.Tags {
		if !model.TagNamePattern.MatchString(t) {
			es = append(es, fmt.Errorf("tag contains invalid characters"))
			break
		}
	}

	return es.Result()
}

// Fill copies data to tag.
func (in *Input) Fill(m *model.Pinl) error {
	m.URL = in.URL
	m.Title = in.Title
	m.Description = in.Description
	m.Readme = in.Readme
	return nil
}
