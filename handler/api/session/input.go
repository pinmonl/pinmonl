package session

import (
	"fmt"
	"io"

	"github.com/pinmonl/pinmonl/handler/api/request"
	"github.com/pinmonl/pinmonl/validate"
)

// Input defines the accepted data from client.
type Input struct {
	Login    string `json:"login"`
	Password string `json:"password"`
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

// Validate checks fields for login process.
func (in *Input) Validate() error {
	var err validate.Errors

	if in.Login == "" {
		err = append(err, fmt.Errorf("login cannot be empty"))
	}
	if in.Password == "" {
		err = append(err, fmt.Errorf("password cannot be empty"))
	}

	return err.Result()
}
