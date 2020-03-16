package user

import (
	"github.com/pinmonl/pinmonl/model"
)

// Body defines the response body of User.
type Body struct {
	model.User
}

// NewBody creates response body.
func NewBody(u model.User) Body {
	return Body{User: u}
}
