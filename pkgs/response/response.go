package response

type Body map[string]interface{}

func Error(err error) Body {
	return Body{"error": err}
}
