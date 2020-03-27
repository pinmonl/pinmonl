package tui

import (
	"fmt"
	"time"

	"github.com/jroimartin/gocui"
	"github.com/pinmonl/pinmonl/logx"
)

type App struct {
	Store *Store

	pinlList   *PinlList
	tagList    *TagList
	statusLine *StatusLine
	help       *Help
	command    *Field
}

func NewApp(endpoint string, debug bool) *App {
	store := NewStore(endpoint)
	app := &App{
		Store: store,
	}
	app.tagList = NewTagList(app)
	app.pinlList = NewPinlList(app)
	app.statusLine = NewStatusLine(app)
	app.help = NewHelp(app)
	app.command = NewField(app, SetFrame(false))
	app.command.SingleLine = true

	initLogger(debug)
	return app
}

func (a *App) Run() error {
	g, err := gocui.NewGui(gocui.Output256)
	if err != nil {
		logx.Panicln(err)
	}
	defer g.Close()

	g.InputEsc = false
	g.Mouse = false
	g.Cursor = true

	g.SetManagerFunc(a.layout)

	a.BindKeys(g)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, a.quit); err != nil {
		return err
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		logx.Println(err)
		return err
	}
	return nil
}

func (a *App) layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	tagW := 36
	statusH := 2

	a.tagList.SetPosition(
		0, 0,
		tagW, maxY-statusH)
	a.pinlList.SetPosition(
		tagW+1, 0,
		maxX-1, maxY-statusH)
	a.statusLine.SetPosition(
		-1, maxY-statusH-1,
		maxX+1, maxY-1)
	a.command.SetPosition(
		-1, maxY-statusH,
		maxX+1, maxY)

	a.tagList.Attach(g)
	a.pinlList.Attach(g)
	a.statusLine.Attach(g)

	if cv := g.CurrentView(); cv == nil {
		a.pinlList.Focus(g)
	}
	return nil
}

func (a *App) redrawPinl(g *gocui.Gui, vFunc func(*gocui.View) error) error {
	pv, err := a.pinlList.UIView(g)
	if err != nil {
		return err
	}
	vFunc(pv)
	a.pinlList.render(g)
	return nil
}

func (a *App) redrawTag(g *gocui.Gui, vFunc func(*gocui.View) error) error {
	tv, err := a.tagList.UIView(g)
	if err != nil {
		return err
	}
	if vFunc != nil {
		vFunc(tv)
	}
	a.tagList.render(g)
	return nil
}

func (a *App) BindKeys(g *gocui.Gui) error {
	g.SetKeybinding("", '1', gocui.ModAlt, a.FocusTagList())
	g.SetKeybinding("", '2', gocui.ModAlt, a.FocusPinlList())
	g.SetKeybinding("", '?', gocui.ModNone, a.ShowHelp())
	return nil
}

func (a *App) UnbindKeys(g *gocui.Gui) error {
	g.DeleteKeybinding("", '1', gocui.ModAlt)
	g.DeleteKeybinding("", '2', gocui.ModAlt)
	g.DeleteKeybinding("", '?', gocui.ModNone)
	return nil
}

func (a *App) quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func (a *App) NextView() CommandFunc {
	return func(g *gocui.Gui, v *gocui.View) error {
		cv := g.CurrentView()
		switch cv.Name() {
		case a.pinlList.Name():
			a.tagList.Focus(g)
		case a.tagList.Name():
			a.pinlList.Focus(g)
		}
		return nil
	}
}

func (a *App) PrevView() CommandFunc {
	return func(g *gocui.Gui, v *gocui.View) error {
		cv := g.CurrentView()
		switch cv.Name() {
		case a.pinlList.Name():
			a.tagList.Focus(g)
		case a.tagList.Name():
			a.pinlList.Focus(g)
		}
		return nil
	}
}

func (a *App) FocusPinlList() CommandFunc {
	return func(g *gocui.Gui, v *gocui.View) error {
		a.pinlList.Focus(g)
		return nil
	}
}

func (a *App) FocusTagList() CommandFunc {
	return func(g *gocui.Gui, v *gocui.View) error {
		a.tagList.Focus(g)
		return nil
	}
}

func (a *App) ShowHelp() CommandFunc {
	return func(g *gocui.Gui, v *gocui.View) error {
		a.UnbindKeys(g)
		g.SetKeybinding("", gocui.KeyCtrlQ, gocui.ModNone, a.HideHelp(v.Name()))
		a.help.Attach(g)
		a.help.Focus(g)
		g.Cursor = false
		return nil
	}
}

func (a *App) HideHelp(prevView string) CommandFunc {
	return func(g *gocui.Gui, v *gocui.View) error {
		g.DeleteKeybinding("", gocui.KeyCtrlQ, gocui.ModNone)
		a.help.Detach(g)
		a.BindKeys(g)
		g.SetCurrentView(prevView)
		g.Cursor = true
		return nil
	}
}

func (a *App) StartCommand(label, value string, onEnter, onCancel CommandFunc) CommandFunc {
	return func(g *gocui.Gui, v *gocui.View) error {
		a.UnbindKeys(g)
		g.SetKeybinding("", gocui.KeyCtrlQ, gocui.ModNone, a.EndCommand(onCancel))
		a.command.SetLabel(g, label)
		a.command.SetValue(g, value)
		a.command.OnEnter = a.EndCommand(onEnter)
		a.command.Attach(g)
		a.command.Focus(g)
		return nil
	}
}

func (a *App) EndCommand(handler CommandFunc) CommandFunc {
	return func(g *gocui.Gui, v *gocui.View) error {
		handler(g, v)
		g.DeleteKeybinding("", gocui.KeyCtrlQ, gocui.ModNone)
		a.BindKeys(g)
		a.command.Detach(g)
		return nil
	}
}

func (a *App) CommandValue(g *gocui.Gui) (string, error) {
	return a.command.Value(g)
}

func (a *App) ShowMessage(g *gocui.Gui, msg Message, d time.Duration) error {
	g.Update(func(g *gocui.Gui) error {
		a.statusLine.msg = msg
		return a.statusLine.render(g)
	})
	if d > 0 {
		time.Sleep(d)
		g.Update(func(g *gocui.Gui) error {
			a.statusLine.msg = nil
			return a.statusLine.render(g)
		})
	}
	return nil
}

type Help struct {
	*View
}

func NewHelp(app *App) *Help {
	v := NewView(app, "help", SetTitle("Key binding"))
	return &Help{View: v}
}

func (h *Help) Attach(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	h.SetPosition(maxX/2-30, maxY/2-18, maxX/2+30, maxY/2+18)
	v, err := h.Draw(g)
	if v == nil {
		return err
	}

	h.render(g)
	return nil
}

func (h *Help) render(g *gocui.Gui) error {
	v, err := h.UIView(g)
	if err != nil {
		return err
	}
	v.Clear()

	bindings := []struct {
		scope string
		key   string
		desc  string
	}{
		{scope: "global"},
		{key: "alt+1", desc: "focus tag pane"},
		{key: "alt+2", desc: "focus bookmark pane"},
		{key: "ctrl+q", desc: "close/cancel"},
		{scope: "control"},
		{key: "/", desc: "search on bookmark/tag"},
		{key: "x", desc: "clear search"},
		{key: "k", desc: "move up"},
		{key: "j", desc: "move down"},
		{key: "ctrl+u", desc: "page up"},
		{key: "ctrl+d", desc: "move down"},
		{scope: "modify"},
		{key: "a", desc: "add new"},
		{key: "e", desc: "edit"},
		{key: "ctrl+n", desc: "go to next field"},
		{key: "ctrl+p", desc: "go to previous field"},
		{scope: "bookmark"},
		{key: "o", desc: "open url in browser"},
		{scope: "tag"},
		{key: "<space>", desc: "toggle filter"},
	}

	fmt.Fprint(v, "\n")
	for _, binding := range bindings {
		if binding.scope != "" {
			fmt.Fprintf(v, "\n%s\n", binding.scope)
		} else {
			fmt.Fprintf(v, "  %*s%s\n", -12, binding.key, binding.desc)
		}
	}

	return nil
}
