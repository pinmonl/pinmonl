package validate

import "strings"

// Errors is slice of error.
type Errors []error

var _ error = Errors{}

// Error concatenates all errors into multi-line message.
func (es Errors) Error() string {
	return strings.Join(es.Errors(), "\n")
}

// Errors returns slice of error messages.
func (es Errors) Errors() []string {
	ess := make([]string, len(es))
	for i, err := range es {
		ess[i] = err.Error()
	}
	return ess
}

// Result returns nil if length is zero or else, returns self.
func (es Errors) Result() error {
	if len(es) == 0 {
		return nil
	}
	return es
}
