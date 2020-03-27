package pmapi

import "strings"

// ErrorBody containes general message from server.
type ErrorBody struct {
	Message     string   `json:"error"`
	Validations []string `json:"validations"`
}

// Error implements error interface.
func (eb ErrorBody) Error() string {
	switch {
	case eb.Message != "":
		return eb.Message
	case len(eb.Validations) > 0:
		return strings.Join(eb.Validations, ", ")
	default:
		return "no message"
	}
}
