package share

import (
	"fmt"
	"io"

	"github.com/pinmonl/pinmonl/handler/api/request"
	"github.com/pinmonl/pinmonl/model"
	"github.com/pinmonl/pinmonl/validate"
)

// Input defines the accepted data from client.
type Input struct {
	model.Share
	MustTags []string `json:"mustTags"`
	AnyTags  []string `json:"anyTags"`
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

	switch n := in.Name; {
	case n == "":
		es = append(es, fmt.Errorf("name cannot be empty"))
	case !model.ShareNamePattern.MatchString(n):
		es = append(es, fmt.Errorf("name contains invalid characters"))
	}

	if len(in.MustTags) == 0 {
		es = append(es, fmt.Errorf("must-tags cannot be empty"))
	}

	vt := func(ts []string) {
		for _, t := range ts {
			if !model.TagNamePattern.MatchString(t) {
				es = append(es, fmt.Errorf("tag contains invalid characters"))
				break
			}
		}
	}

	vt(in.MustTags)
	vt(in.AnyTags)

	return es.Result()
}

// Fill copies data to tag.
func (in *Input) Fill(m *model.Share) error {
	m.Name = in.Name
	m.Description = in.Description
	m.Readme = in.Readme
	return nil
}
