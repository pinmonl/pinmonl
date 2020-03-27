package tui

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

type StatusLine struct {
	*View
	msg Message
}

func NewStatusLine(app *App) *StatusLine {
	v := NewView(app, "status_line", SetFrame(false))
	return &StatusLine{View: v}
}

func (s *StatusLine) Attach(g *gocui.Gui) error {
	v, err := s.Draw(g)
	if v == nil {
		return err
	}

	s.BindKeys(g)
	s.render(g)

	return nil
}

func (s *StatusLine) render(g *gocui.Gui) error {
	v, err := s.UIView(g)
	if err != nil {
		return err
	}
	v.Clear()
	switch {
	case s.msg != nil:
		fmt.Fprintf(v, s.renderMessage(v, s.msg))
	default:
		fmt.Fprintf(v, s.helpMessage(v))
	}
	return nil
}

func (s *StatusLine) renderMessage(v *gocui.View, msg Message) string {
	sx, _ := v.Size()
	maxX := -(sx - 1)
	return fmt.Sprintf("%s%*s\x1b[0m", msg.Styles(), maxX, msg.Text())
}

func (s *StatusLine) update(g *gocui.Gui) error {
	return s.render(g)
}

func (s *StatusLine) BindKeys(g *gocui.Gui) error {
	return nil
}

func (s *StatusLine) UnbindKeys(g *gocui.Gui) error {
	return nil
}

func (s *StatusLine) helpMessage(v *gocui.View) string {
	sx, _ := v.Size()
	maxX := -(sx - 1)
	return fmt.Sprintf("\x1b[48;5;240m%*s\x1b[0m", maxX, "Press '?' for help")
}
