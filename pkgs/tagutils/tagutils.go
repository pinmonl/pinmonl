package tagutils

import "strings"

func ToNamePattern(tagName string) string {
	// Remove all '%'.
	pattern := strings.ReplaceAll(tagName, "%", "")

	// Remove leading '/'.
	if strings.HasPrefix(pattern, "/") {
		pattern = strings.TrimPrefix(pattern, "/")
		// Prepend '%' if does not have leading '/'.
	} else {
		pattern = "%" + pattern
	}

	// Remove trailing '/'.
	if strings.HasSuffix(pattern, "/") {
		pattern = strings.TrimSuffix(pattern, "/")
		// Replace as wildcard search for trailing '/*'.
	} else if strings.HasSuffix(pattern, "/*") {
		pattern = strings.TrimSuffix(pattern, "*") + "%"
		// Append '%' if does not have trailing '/'
	} else {
		pattern += "%"
	}

	return pattern
}
