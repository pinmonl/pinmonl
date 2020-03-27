package tui

import (
	"github.com/atotto/clipboard"
)

func WriteClipboard(text string) error {
	return clipboard.WriteAll(text)
}
