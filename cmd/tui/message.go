package tui

type Message interface {
	Text() string
	Styles() string
}

func NewMessage(text string) Message {
	return newMessage(text, "\x1b[48;5;240m")
}

func NewSuccessMessage(text string) Message {
	return newMessage(text, "\x1b[48;5;2m\x1b[38;5;0m")
}

func NewErrorMessage(text string) Message {
	return newMessage(text, "\x1b[48;5;1m")
}

func newMessage(text string, styles string) Message {
	return &basicMessage{text: text, styles: styles}
}

type basicMessage struct {
	text   string
	styles string
}

func (b *basicMessage) Text() string { return b.text }

func (b *basicMessage) Styles() string { return b.styles }
