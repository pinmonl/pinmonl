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
	Login    string `json:"login"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Name     string `json:"name"`
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

	switch login := in.Login; {
	case login == "":
		err = append(err, fmt.Errorf("login cannot be empty"))
	case !validate.IsAlphanumeric(login):
		err = append(err, fmt.Errorf("invalid characters in login"))
	}

	switch pw := in.Password; {
	case len(pw) < 8:
		err = append(err, fmt.Errorf("password must be at least 8 characters long"))
	}

	switch email := in.Email; {
	case email == "":
		err = append(err, fmt.Errorf("email cannot be empty"))
	case !validate.IsEmail(email):
		err = append(err, fmt.Errorf("invalid email format"))
	}

	switch name := in.Name; {
	case name == "":
		err = append(err, fmt.Errorf("name cannot be empty"))
	}

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

// Fill copies input to user.
func (in *Input) Fill(m *model.User) error {
	pw, err := password.Hash(in.Password)
	if err != nil {
		return err
	}

	m.Login = in.Login
	m.Password = pw
	m.Email = in.Email
	m.Name = in.Name
	return nil
}
