package validate

import "github.com/asaskevich/govalidator"

// IsURL checks if the string is an URL.
func IsURL(str string) bool {
	return govalidator.IsURL(str)
}
