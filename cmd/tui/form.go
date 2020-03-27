package tui

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/jroimartin/gocui"
	"github.com/pinmonl/pinmonl/pkg/generate"
)

type Field struct {
	*View
	label      string
	value      string
	tainted    bool
	SingleLine bool
	OnEnter    CommandFunc
}

func NewField(app *App, opts ...ViewOpt) *Field {
	f := &Field{tainted: true}
	v := NewView(app, "field-"+generate.RandomString(10), opts...)
	v.Editable = true
	v.Editor = f
	f.View = v
	return f
}

func (f *Field) Attach(g *gocui.Gui) error {
	v, err := f.Draw(g)
	if v == nil {
		return err
	}
	f.BindKeys(g)
	f.render(g)
	return nil
}

func (f *Field) Detach(g *gocui.Gui) error {
	f.UnbindKeys(g)
	f.View.Detach(g)
	return nil
}

func (f *Field) BindKeys(g *gocui.Gui) error {
	if f.OnEnter != nil {
		g.SetKeybinding(f.Name(), gocui.KeyEnter, gocui.ModNone, f.OnEnter)
	}
	return nil
}

func (f *Field) UnbindKeys(g *gocui.Gui) error {
	if f.OnEnter != nil {
		g.DeleteKeybinding(f.Name(), gocui.KeyEnter, gocui.ModNone)
	}
	return nil
}

func (f *Field) render(g *gocui.Gui) error {
	v, err := f.UIView(g)
	if err != nil {
		return err
	}
	v.Clear()
	if f.label != "" {
		fmt.Fprint(v, f.label)
	}
	fmt.Fprint(v, f.value)
	if f.tainted {
		f.moveToEnd(v)
		f.tainted = false
	}
	return nil
}

func (f *Field) length() int {
	return len(f.value)
}

func (f *Field) viewLines(v *gocui.View) [][]rune {
	maxX, _ := v.Size()
	vRunes := [][]rune{}
	for _, bufLine := range v.BufferLines() {
		line := bytes.Runes([]byte(bufLine))
		for i := 0; i <= len(line); i += maxX {
			switch {
			case len(line[i:]) <= maxX:
				vRunes = append(vRunes, line[i:])
			default:
				vRunes = append(vRunes, line[i:i+maxX])
			}
		}
	}
	return vRunes
}

func (f *Field) moveToEnd(v *gocui.View) error {
	llen := len(f.label)
	if len(v.Buffer()) == 0 {
		v.SetOrigin(0, 0)
		return v.SetCursor(llen, 0)
	}
	vlines := f.viewLines(v)
	cy := len(vlines) - 1
	cx := len(vlines[cy])
	return v.SetCursor(cx, cy)
}

func (f *Field) moveToLineStart(v *gocui.View) error {
	_, cy := v.Cursor()
	if !f.hasLabel() || cy > 0 {
		return v.SetCursor(0, cy)
	}
	return v.SetCursor(f.labelLen(), 0)
}

func (f *Field) moveToLineEnd(v *gocui.View) error {
	_, cy := v.Cursor()
	vlines := f.viewLines(v)
	cx := len(vlines[cy])
	return v.SetCursor(cx, cy)
}

func (f *Field) Value(g *gocui.Gui) (string, error) {
	v, err := f.UIView(g)
	if err != nil {
		return "", err
	}
	bl := v.BufferLines()
	value := strings.Join(bl, "\n")
	return value[f.labelLen():], nil
}

func (f *Field) SetValue(g *gocui.Gui, value string) error {
	f.value = value
	f.tainted = true
	return f.render(g)
}

func (f *Field) Focus(g *gocui.Gui) (*gocui.View, error) {
	gv, err := f.View.Focus(g)
	return gv, err
}

func (f *Field) Edit(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	switch {
	case ch != 0 && mod == 0:
		v.EditWrite(ch)
	case key == gocui.KeySpace:
		v.EditWrite(' ')
	case key == gocui.KeyBackspace || key == gocui.KeyBackspace2:
		if !f.cursorOnLabel(v, -1, 0) {
			v.EditDelete(true)
		}
	case key == gocui.KeyDelete:
		v.EditDelete(false)
	case !f.SingleLine && key == gocui.KeyEnter:
		v.EditNewLine()
	case key == gocui.KeyHome:
		f.moveToLineStart(v)
	case key == gocui.KeyEnd:
		f.moveToLineEnd(v)
	case key == gocui.KeyArrowDown:
		v.MoveCursor(0, 1, false)
	case key == gocui.KeyArrowUp:
		if !f.cursorOnLabel(v, 0, -1) {
			v.MoveCursor(0, -1, false)
		} else {
			f.moveToLineStart(v)
		}
	case key == gocui.KeyArrowLeft:
		if !f.cursorOnLabel(v, -1, 0) {
			v.MoveCursor(-1, 0, false)
		}
	case key == gocui.KeyArrowRight:
		v.MoveCursor(1, 0, false)
	}
}

func (f *Field) labelLen() int {
	return len(f.label)
}

func (f *Field) hasLabel() bool {
	return f.labelLen() > 0
}

func (f *Field) SetLabel(g *gocui.Gui, label string) error {
	f.label = label
	f.tainted = true
	return f.render(g)
}

func (f *Field) cursorOnLabel(v *gocui.View, dx int, dy int) bool {
	if !f.hasLabel() {
		return false
	}
	cx, cy := v.Cursor()
	nx, ny := cx+dx, cy+dy
	switch {
	case ny < 0:
		return true
	case ny == 0:
		return 0 <= nx && nx < f.labelLen()
	default:
		return false
	}
}

type TagField struct {
	*Field
}

func NewTagField(app *App, opts ...ViewOpt) *TagField {
	f := NewField(app, opts...)
	f.Wrap = true
	f.SingleLine = true
	return &TagField{Field: f}
}

func (t *TagField) SetValue(g *gocui.Gui, value []string) error {
	t.Field.SetValue(g, strings.Join(value, ", "))
	return t.render(g)
}

func (t *TagField) Value(g *gocui.Gui) ([]string, error) {
	value, err := t.Field.Value(g)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(value) == "" {
		return []string{}, nil
	}
	tags := strings.Split(value, ",")
	for i, tag := range tags {
		tags[i] = strings.TrimSpace(tag)
	}
	return tags, nil
}
