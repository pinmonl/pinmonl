package webui

import "github.com/gobuffalo/packr/v2"

var box = packr.New("webui", "../webui/dist")

// PackrBox exports webui in packr.Box.
func PackrBox() *packr.Box {
	return box
}
