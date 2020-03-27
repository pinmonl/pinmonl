package tui

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jroimartin/gocui"
	"github.com/pinmonl/pinmonl/pmapi"
)

type PinlList struct {
	*View
	fieldWidth  int
	urlField    *Field
	titleField  *Field
	descField   *Field
	readmeField *Field
	tagField    *TagField
}

func NewPinlList(app *App) *PinlList {
	v := NewView(app, "pinl_list", SetTitle("Bookmark"))
	p := &PinlList{View: v, fieldWidth: 80}

	opt := func(v *View) *View {
		v.Wrap = true
		return v
	}
	p.urlField = NewField(p.App, opt, SetTitle("URL"))
	p.descField = NewField(p.App, opt, SetTitle("Description"))
	p.readmeField = NewField(p.App, opt, SetTitle("Readme"))
	p.titleField = NewField(p.App, opt, SetTitle("Title"), SetFrame(false))
	p.tagField = NewTagField(p.App, opt, SetTitle("Tag"))
	p.urlField.SingleLine = true
	p.titleField.SingleLine = true
	return p
}

func (p *PinlList) Attach(g *gocui.Gui) error {
	v, err := p.Draw(g)
	p.attachFields(g)
	if v == nil {
		return err
	}

	p.BindKeys(g)
	fmt.Fprint(v, "Loading...")

	go func() {
		g.Update(p.update)
	}()
	return nil
}

func (p *PinlList) attachFields(g *gocui.Gui) error {
	fX0 := p.X1() - p.fieldWidth
	fX1 := p.X1()

	urlH := 4
	descH := 10
	tagH := 4

	p.urlField.SetPosition(
		fX0, 0,
		fX1, urlH)
	p.descField.SetPosition(
		fX0, urlH+1,
		fX1, urlH+1+descH)
	p.tagField.SetPosition(
		fX0, urlH+descH+2,
		fX1, urlH+descH+2+tagH)

	p.urlField.Attach(g)
	p.descField.Attach(g)
	p.tagField.Attach(g)

	return nil
}

func (p *PinlList) update(g *gocui.Gui) error {
	return p.render(g)
}

func (p *PinlList) render(g *gocui.Gui) error {
	v, err := p.UIView(g)
	if err != nil {
		return err
	}
	v.Clear()
	switch {
	case p.Len() == 0:
		fmt.Fprint(v, "Empty")
	default:
		p.clearFields(g)
		pinls := p.getData()
		_, y := p.UICursor(v)
		for i, row := range pinls {
			if i > 0 {
				fmt.Fprint(v, "\n")
			}
			hl := i == y
			fmt.Fprint(v, p.renderLine(&row, hl))
		}
		if row := p.getAt(y); row != nil && !p.editInsert {
			p.urlField.SetValue(g, row.URL)
			p.titleField.SetValue(g, row.Title)
			p.descField.SetValue(g, row.Description)
			p.readmeField.SetValue(g, row.Readme)
			p.tagField.SetValue(g, row.Tags)
		}
		if y >= p.Len() {
			fmt.Fprint(v, "\n")
			fmt.Fprint(v, p.renderLine(nil, true))
		}
	}
	return nil
}

func (p *PinlList) renderLine(row *pmapi.Pinl, highlight bool) string {
	str := ""
	if highlight {
		if p.editMode() {
			// str += "\x1b[38;5;10m>\x1b[0m"
			str += " "
		} else {
			str += ">"
		}
	} else {
		str += " "
	}
	if row != nil {
		str += fmt.Sprintf(" %s", row.Title)
	}
	return str
}

func (p *PinlList) clearFields(g *gocui.Gui) error {
	p.titleField.SetValue(g, "")
	p.urlField.SetValue(g, "")
	p.descField.SetValue(g, "")
	p.tagField.SetValue(g, nil)
	return nil
}

func (p *PinlList) gatherFields(g *gocui.Gui) (*pmapi.Pinl, error) {
	b := pmapi.Pinl{Tags: []string{}}
	if val, err := p.urlField.Value(g); err == nil {
		b.URL = val
	}
	if val, err := p.titleField.Value(g); err == nil {
		b.Title = val
	}
	if val, err := p.descField.Value(g); err == nil {
		b.Description = val
	}
	if val, err := p.readmeField.Value(g); err == nil {
		b.Readme = val
	}
	if val, err := p.tagField.Value(g); err == nil {
		b.Tags = val
	}
	if v, err := p.UIView(g); err == nil && !p.editInsert {
		_, y := p.UICursor(v)
		if row := p.getAt(y); row != nil {
			b.ID = row.ID
		}
	}
	return &b, nil
}

func (p *PinlList) getData() pinlSlice {
	data, err := p.App.Store.GetPinls(context.TODO())
	if err != nil {
		return nil
	}
	return data
}

func (p *PinlList) Len() int {
	return len(p.getData())
}

func (p *PinlList) getAt(i int) *pmapi.Pinl {
	if i >= p.Len() {
		return nil
	}
	item := p.getData()[i]
	return &item
}

func (p *PinlList) BindKeys(g *gocui.Gui) error {
	g.SetKeybinding(p.Name(), 'o', gocui.ModNone, p.OpenLink())
	g.SetKeybinding(p.Name(), 'y', gocui.ModNone, p.CopyLink())
	g.SetKeybinding(p.Name(), 'j', gocui.ModNone, p.MoveDown())
	g.SetKeybinding(p.Name(), 'k', gocui.ModNone, p.MoveUp())
	g.SetKeybinding(p.Name(), gocui.KeyCtrlD, gocui.ModNone, p.PageDown())
	g.SetKeybinding(p.Name(), gocui.KeyCtrlU, gocui.ModNone, p.PageUp())
	g.SetKeybinding(p.Name(), 'a', gocui.ModNone, p.AddBookmark())
	g.SetKeybinding(p.Name(), 'e', gocui.ModNone, p.EditBookmark())
	g.SetKeybinding(p.Name(), 'd', gocui.ModNone, p.DeleteBookmark())
	g.SetKeybinding(p.Name(), '/', gocui.ModNone, p.StartSearch())
	g.SetKeybinding(p.Name(), 'x', gocui.ModNone, p.ClearFilter())
	return nil
}

func (p *PinlList) UnbindKeys(g *gocui.Gui) error {
	g.DeleteKeybinding(p.Name(), 'o', gocui.ModNone)
	g.DeleteKeybinding(p.Name(), 'y', gocui.ModNone)
	g.DeleteKeybinding(p.Name(), 'j', gocui.ModNone)
	g.DeleteKeybinding(p.Name(), 'k', gocui.ModNone)
	g.DeleteKeybinding(p.Name(), gocui.KeyCtrlD, gocui.ModNone)
	g.DeleteKeybinding(p.Name(), gocui.KeyCtrlU, gocui.ModNone)
	g.DeleteKeybinding(p.Name(), 'a', gocui.ModNone)
	g.DeleteKeybinding(p.Name(), 'e', gocui.ModNone)
	g.DeleteKeybinding(p.Name(), 'd', gocui.ModNone)
	g.DeleteKeybinding(p.Name(), '/', gocui.ModNone)
	g.DeleteKeybinding(p.Name(), 'x', gocui.ModNone)
	return nil
}

func (p *PinlList) OpenLink() CommandFunc {
	return func(g *gocui.Gui, v *gocui.View) error {
		_, y := p.UICursor(v)
		row := p.getAt(y)
		if row == nil {
			return nil
		}
		err := OpenURL(row.URL)
		debugln(err)
		return nil
	}
}

func (p *PinlList) CopyLink() CommandFunc {
	return func(g *gocui.Gui, v *gocui.View) error {
		_, y := p.UICursor(v)
		row := p.getAt(y)
		if row == nil {
			return nil
		}
		err := WriteClipboard(row.URL)
		debugln(err)
		return nil
	}
}

func (p *PinlList) MoveDown() CommandFunc {
	return func(g *gocui.Gui, v *gocui.View) error {
		p.View.MoveDown()(g, v)
		return p.render(g)
	}
}

func (p *PinlList) MoveUp() CommandFunc {
	return func(g *gocui.Gui, v *gocui.View) error {
		p.View.MoveUp()(g, v)
		return p.render(g)
	}
}

func (p *PinlList) PageDown() CommandFunc {
	return func(g *gocui.Gui, v *gocui.View) error {
		p.View.PageDown(p.Len())(g, v)
		return p.render(g)
	}
}

func (p *PinlList) PageUp() CommandFunc {
	return func(g *gocui.Gui, v *gocui.View) error {
		p.View.PageUp(p.Len())(g, v)
		return p.render(g)
	}
}

func (p *PinlList) AddBookmark() CommandFunc {
	return func(g *gocui.Gui, v *gocui.View) error {
		p.editIndex = p.Len()
		p.editInsert = true
		p.clearFields(g)
		return p.StartEditing()(g, v)
	}
}

func (p *PinlList) EditBookmark() CommandFunc {
	return func(g *gocui.Gui, v *gocui.View) error {
		_, y := p.UICursor(v)
		p.editIndex = y
		return p.StartEditing()(g, v)
	}
}

func (p *PinlList) StartEditing() CommandFunc {
	return func(g *gocui.Gui, v *gocui.View) error {
		_, oy := v.Origin()
		titleX0 := p.X0()
		titleX1 := p.X1() - p.fieldWidth - 2
		titleY0 := p.editIndex - oy
		p.titleField.SetPosition(
			titleX0, titleY0,
			titleX1, titleY0+2)
		p.titleField.label = "> "
		p.titleField.Attach(g)

		p.App.UnbindKeys(g)
		p.UnbindKeys(g)
		p.bindFormKeys(g)
		p.urlField.Focus(g)
		return p.render(g)
	}
}

func (p *PinlList) EndEditing() CommandFunc {
	return func(g *gocui.Gui, v *gocui.View) error {
		p.titleField.Detach(g)

		p.resetEdit()
		p.unbindFormKeys(g)
		p.App.BindKeys(g)
		p.BindKeys(g)
		p.Focus(g)

		if gv, err := p.UIView(g); err == nil {
			_, oy := gv.Origin()
			_, cy := gv.Cursor()
			if oy+cy >= p.Len() {
				_, maxY := gv.Size()
				switch {
				case p.Len() < maxY:
					oy = 0
					cy = p.Len() - 1
				default:
					oy = p.Len() - maxY
					cy = p.Len() - oy
				}
				gv.SetOrigin(0, oy)
				gv.SetCursor(0, cy)
			}
		}

		return p.render(g)
	}
}

func (p *PinlList) SaveEditing() CommandFunc {
	return func(g *gocui.Gui, v *gocui.View) error {
		var model pmapi.Pinl
		if data, err := p.gatherFields(g); err == nil {
			model = *data
		}

		insert := p.editInsert
		go func() {
			var err error
			ctx := context.TODO()
			p.App.ShowMessage(g, NewMessage("Loading..."), 0)
			if insert {
				err = p.App.Store.CreatePinl(ctx, &model)
			} else {
				err = p.App.Store.UpdatePinl(ctx, &model)
			}
			if err != nil {
				p.App.ShowMessage(g, NewErrorMessage("Some errors occurred"), time.Second)
				return
			} else {
				g.Update(p.update)
				if insert {
					p.App.ShowMessage(g, NewSuccessMessage("Bookmark added"), time.Second)
				} else {
					p.App.ShowMessage(g, NewSuccessMessage("Bookmark updated"), time.Second)
				}
				p.App.redrawTag(g, nil)
			}
		}()

		return p.EndEditing()(g, v)
	}
}

func (p *PinlList) bindFormKeys(g *gocui.Gui) error {
	g.SetKeybinding("", gocui.KeyCtrlN, gocui.ModNone, p.NextField())
	g.SetKeybinding("", gocui.KeyCtrlP, gocui.ModNone, p.PrevField())
	g.SetKeybinding("", gocui.KeyCtrlQ, gocui.ModNone, p.EndEditing())
	g.SetKeybinding("", gocui.KeyCtrlS, gocui.ModNone, p.SaveEditing())
	return nil
}

func (p *PinlList) unbindFormKeys(g *gocui.Gui) error {
	g.DeleteKeybinding("", gocui.KeyCtrlN, gocui.ModNone)
	g.DeleteKeybinding("", gocui.KeyCtrlP, gocui.ModNone)
	g.DeleteKeybinding("", gocui.KeyCtrlQ, gocui.ModNone)
	g.DeleteKeybinding("", gocui.KeyCtrlS, gocui.ModNone)
	return nil
}

func (p *PinlList) DeleteBookmark() CommandFunc {
	return func(g *gocui.Gui, v *gocui.View) error {
		onCancel := func(g *gocui.Gui, v *gocui.View) error {
			p.Focus(g)
			return nil
		}
		onEnter := func(g *gocui.Gui, v *gocui.View) error {
			val, err := p.App.CommandValue(g)
			if err != nil {
				return err
			}
			val = strings.ToLower(val)
			if val == "y" || val == "yes" {
				v, _ := p.UIView(g)
				_, y := p.UICursor(v)
				row := p.getAt(y)
				if row != nil {
					go func() {
						p.App.ShowMessage(g, NewMessage("Loading..."), 0)
						err := p.App.Store.DeletePinl(context.TODO(), row)
						if err != nil {
							p.App.ShowMessage(g, NewErrorMessage(fmt.Sprintf("Some error occurred: %s", err)), time.Second)
						} else {
							v.MoveCursor(0, -1, false)
							g.Update(p.update)
							p.Focus(g)
							go func() {
								p.App.ShowMessage(g, NewSuccessMessage("Bookmark deleted"), time.Second)
							}()
						}
					}()
				}
			} else {
				onCancel(g, v)
			}
			return nil
		}
		return p.App.StartCommand("Confirm to delete? (y/n)", "", onEnter, onCancel)(g, v)
	}
}

func (p *PinlList) NextField() CommandFunc {
	return func(g *gocui.Gui, v *gocui.View) error {
		cv := g.CurrentView()
		switch cv.Name() {
		case p.titleField.Name():
			p.urlField.Focus(g)
		case p.urlField.Name():
			p.descField.Focus(g)
		case p.descField.Name():
			p.tagField.Focus(g)
		case p.tagField.Name():
			p.titleField.Focus(g)
		}
		return nil
	}
}

func (p *PinlList) PrevField() CommandFunc {
	return func(g *gocui.Gui, v *gocui.View) error {
		cv := g.CurrentView()
		switch cv.Name() {
		case p.tagField.Name():
			p.descField.Focus(g)
		case p.descField.Name():
			p.urlField.Focus(g)
		case p.urlField.Name():
			p.titleField.Focus(g)
		case p.titleField.Name():
			p.tagField.Focus(g)
		}
		return nil
	}
}

func (p *PinlList) StartSearch() CommandFunc {
	return func(g *gocui.Gui, v *gocui.View) error {
		prevSearch := p.App.Store.pinlSearch
		return p.App.StartCommand("search bookmark>", prevSearch, p.SubmitSearch(), p.EndSearch())(g, v)
	}
}

func (p *PinlList) SubmitSearch() CommandFunc {
	return func(g *gocui.Gui, v *gocui.View) error {
		search, err := p.App.CommandValue(g)
		if err != nil {
			return err
		}
		p.App.Store.SetPinlSearch(search)
		pv, err := p.UIView(g)
		if err != nil {
			return err
		}
		pv.SetOrigin(0, 0)
		pv.SetCursor(0, 0)
		return p.EndSearch()(g, v)
	}
}

func (p *PinlList) EndSearch() CommandFunc {
	return func(g *gocui.Gui, v *gocui.View) error {
		p.Focus(g)
		return p.render(g)
	}
}

func (p *PinlList) ClearFilter() CommandFunc {
	return func(g *gocui.Gui, v *gocui.View) error {
		p.App.Store.SetPinlSearch("")
		v.SetOrigin(0, 0)
		v.SetCursor(0, 0)
		return p.render(g)
	}
}
