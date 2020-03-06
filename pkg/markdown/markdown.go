package markdown

import (
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
)

// SafeHTML produces sanitized HTML from Markdown.
func SafeHTML(input string) string {
	unsafe := Render(input)
	html := bluemonday.UGCPolicy().Sanitize(unsafe)
	return html
}

// Render produces HTML from Markdown.
func Render(input string) string {
	output := blackfriday.Run([]byte(input))
	return string(output)
}
