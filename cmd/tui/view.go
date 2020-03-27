package tui

import (
	"github.com/jroimartin/gocui"
)

type ViewOpt func(*View) *View

func SetTitle(title string) ViewOpt {
	return func(v *View) *View {
		v.Title = title
		return v
	}
}

func SetFrame(frame bool) ViewOpt {
	return func(v *View) *View {
		v.Frame = frame
		return v
	}
}

type View struct {
	App *App

	x0, y0     int
	x1, y1     int
	name       string
	editIndex  int
	editInsert bool

	BgColor    gocui.Attribute
	FgColor    gocui.Attribute
	SelBgColor gocui.Attribute
	SelFgColor gocui.Attribute

	Editable   bool
	Editor     gocui.Editor
	Frame      bool
	Wrap       bool
	Autoscroll bool
	Title      string
}

func NewView(app *App, name string, opts ...ViewOpt) *View {
	v := &View{
		App:       app,
		name:      name,
		editIndex: -1,
		Frame:     true,
	}
	return v.apply(opts...)
}

func (v *View) apply(opts ...ViewOpt) *View {
	for _, opt := range opts {
		v = opt(v)
	}
	return v
}

func (v *View) Name() string { return v.name }

func (v *View) SetPosition(x0, y0, x1, y1 int) {
	v.x0, v.y0, v.x1, v.y1 = x0, y0, x1, y1
}

func (v *View) X0() int { return v.x0 }

func (v *View) Y0() int { return v.y0 }

func (v *View) X1() int { return v.x1 }

func (v *View) Y1() int { return v.y1 }

func (v *View) Attach(g *gocui.Gui) error {
	_, err := v.Draw(g)
	return err
}

func (v *View) Draw(g *gocui.Gui) (*gocui.View, error) {
	gv, err := g.SetView(v.Name(), v.X0(), v.Y0(), v.X1(), v.Y1())
	if err == gocui.ErrUnknownView {
		gv.Editable = v.Editable
		gv.Frame = v.Frame
		gv.Wrap = v.Wrap
		gv.Autoscroll = v.Autoscroll
		gv.Editor = v.Editor
		gv.Title = v.Title
		gv.BgColor = v.BgColor
		gv.FgColor = v.FgColor
		gv.SelBgColor = v.SelBgColor
		gv.SelFgColor = v.SelFgColor
		return gv, nil
	}
	return nil, err
}

func (v *View) Detach(g *gocui.Gui) error {
	return g.DeleteView(v.Name())
}

func (v *View) UIView(g *gocui.Gui) (*gocui.View, error) {
	return g.View(v.Name())
}

func (v *View) Focus(g *gocui.Gui) (*gocui.View, error) {
	gv, err := g.SetCurrentView(v.Name())
	return gv, err
}

func (v *View) SetTitle(g *gocui.Gui, title string) error {
	v.Title = title
	gv, err := v.UIView(g)
	if err != nil {
		return err
	}
	gv.Title = v.Title
	return nil
}

func (*View) MoveDown() CommandFunc {
	return func(g *gocui.Gui, v *gocui.View) error {
		v.MoveCursor(0, 1, false)
		return nil
	}
}

func (*View) MoveUp() CommandFunc {
	return func(g *gocui.Gui, v *gocui.View) error {
		v.MoveCursor(0, -1, false)
		return nil
	}
}

func (*View) PageDown(maxY int) CommandFunc {
	return func(g *gocui.Gui, v *gocui.View) error {
		_, sy := v.Size()
		_, oy := v.Origin()
		_, cy := v.Cursor()
		ny := oy + sy
		switch {
		case maxY < sy:
			cy = maxY - 1
		case ny == maxY:
			cy = sy - 1
		case maxY-ny < sy:
			oy = maxY - sy
		default:
			oy = ny
		}
		v.SetOrigin(0, oy)
		v.SetCursor(0, cy)
		return nil
	}
}

func (*View) PageUp(_ int) CommandFunc {
	return func(g *gocui.Gui, v *gocui.View) error {
		_, sy := v.Size()
		_, oy := v.Origin()
		_, cy := v.Cursor()
		ny := oy - sy
		switch {
		case oy == 0 && cy > 0:
			cy = 0
		case ny < 0:
			oy = 0
		default:
			oy = ny
		}
		v.SetOrigin(0, oy)
		v.SetCursor(0, cy)
		return nil
	}
}

func (*View) UICursor(v *gocui.View) (int, int) {
	ox, oy := v.Origin()
	cx, cy := v.Cursor()
	return ox + cx, oy + cy
}

func (v *View) editMode() bool {
	return v.editIndex > -1
}

func (v *View) resetEdit() error {
	v.editIndex = -1
	v.editInsert = false
	return nil
}
