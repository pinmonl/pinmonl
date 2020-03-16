package user

import (
	"fmt"
	"io"

	"github.com/pinmonl/pinmonl/handler/api/request"
	"github.com/pinmonl/pinmonl/model"
	"github.com/pinmonl/pinmonl/pkg/password"
	"github.com/pinmonl/pinmonl/validate"
)

// Input defines the accepted data from client.
type Input struct {
	model.User
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

// Validate checks fields in data.
func (in *Input) Validate() error {
	var err validate.Errors
	err = append(err, in.validateLogin()...)
	err = append(err, in.validatePassword()...)
	err = append(err, in.validateName()...)
	return err.Result()
}

// ValidateLogin checks fields for login process.
func (in *Input) ValidateLogin() error {
	var err validate.Errors
	if in.Login == "" {
		err = append(err, fmt.Errorf("login cannot be empty"))
	}
	if in.Password == "" {
		err = append(err, fmt.Errorf("password cannot be empty"))
	}
	return err.Result()
}

// ValidateUpdate checks fields for updating.
func (in *Input) ValidateUpdate() error {
	var err validate.Errors
	if in.Login != "" {
		err = append(err, in.validateLogin()...)
	}
	if in.Password != "" {
		err = append(err, in.validatePassword()...)
	}
	if in.Name != "" {
		err = append(err, in.validateName()...)
	}
	return err.Result()
}

// Fill copies input to user.
func (in *Input) Fill(m *model.User) error {
	pw, err := password.Hash(in.Password)
	if err != nil {
		return err
	}

	m.Login = in.Login
	m.Password = pw
	m.Name = in.Name
	return nil
}

// FillDirty copies non-empty input only.
func (in *Input) FillDirty(m *model.User) error {
	if in.Password != "" {
		pw, err := password.Hash(in.Password)
		if err != nil {
			return err
		}
		m.Password = pw
	}
	if in.Name != "" {
		m.Name = in.Name
	}
	if in.Login != "" {
		m.Login = in.Login
	}
	return nil
}

func (in *Input) validateLogin() []error {
	var err []error
	switch login := in.Login; {
	case login == "":
		err = append(err, fmt.Errorf("login cannot be empty"))
	case !validate.IsAlphanumeric(login):
		err = append(err, fmt.Errorf("invalid characters in login"))
	}
	return err
}

func (in *Input) validatePassword() []error {
	var err []error
	switch pw := in.Password; {
	case len(pw) < 8:
		err = append(err, fmt.Errorf("password must be at least 8 characters long"))
	}
	return err
}

func (in *Input) validateName() []error {
	var err []error
	switch name := in.Name; {
	case name == "":
		err = append(err, fmt.Errorf("name cannot be empty"))
	}
	return err
}
