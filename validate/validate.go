package validate

import (
	"regexp"

	"github.com/asaskevich/govalidator"
)

// IsURL checks if the string is an URL.
func IsURL(str string) bool {
	return govalidator.IsURL(str)
}

// IsEmail checks if the string is an email.
func IsEmail(str string) bool {
	return govalidator.IsEmail(str)
}

// IsAlpha checks if the string contains only alpha characters.
func IsAlpha(str string) bool {
	return govalidator.IsAlpha(str)
}

// IsAlphanumeric checks if the string contains alpha and numeric characters.
func IsAlphanumeric(str string) bool {
	return govalidator.IsAlphanumeric(str)
}

// IsPattern checks if the string matches with the pattern using RegExp.
func IsPattern(str, pattern string) bool {
	r := regexp.MustCompile(pattern)
	return r.MatchString(str)
}
