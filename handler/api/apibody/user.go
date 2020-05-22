package apibody

import "github.com/pinmonl/pinmonl/model"

// User defines the response body of User.
type User struct {
	model.User
}

// NewUser creates response body.
func NewUser(u model.User) User {
	return User{User: u}
}
