package tag

import (
	"fmt"
	"io"

	"github.com/pinmonl/pinmonl/handler/api/request"
	"github.com/pinmonl/pinmonl/model"
	"github.com/pinmonl/pinmonl/validate"
)

// Input defines the accepted data from client.
type Input struct {
	model.Tag
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
	case !model.TagNamePattern.MatchString(n):
		es = append(es, fmt.Errorf("name contains invalid characters"))
	}

	return es.Result()
}

// Fill copies input to tag.
func (in *Input) Fill(m *model.Tag) error {
	m.Name = in.Name
	m.ParentID = in.ParentID
	m.Sort = in.Sort
	return nil
}
