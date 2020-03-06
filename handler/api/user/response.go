package user

import "github.com/pinmonl/pinmonl/model"

// Resp returns data of user.
func Resp(m model.User) map[string]interface{} {
	resp := map[string]interface{}{
		"id":      m.ID,
		"login":   m.Login,
		"email":   m.Email,
		"name":    m.Name,
		"imageId": m.ImageID,
	}
	return resp
}
