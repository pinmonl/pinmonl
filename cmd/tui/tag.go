package tui

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jroimartin/gocui"
	"github.com/pinmonl/pinmonl/pmapi"
)

type TagList struct {
	*View
	opens map[string]bool

	treeView  bool
	nameField *Field
}

func NewTagList(app *App) *TagList {
	v := NewView(app, "tag_list", SetTitle("Tag"))
	t := &TagList{
		View:     v,
		opens:    make(map[string]bool),
		treeView: true,
	}
	t.nameField = NewField(app, SetFrame(false))
	t.nameField.SingleLine = true
	return t
}

func (t *TagList) Attach(g *gocui.Gui) error {
	v, err := t.Draw(g)
	t.attachFields(g)
	if v == nil {
		return err
	}

	t.BindKeys(g)
	fmt.Fprint(v, "Loading...")

	go func() {
		g.Update(t.update)
	}()
	return nil
}

func (t *TagList) attachFields(g *gocui.Gui) error {
	return nil
}

func (t *TagList) BindKeys(g *gocui.Gui) error {
	g.SetKeybinding(t.Name(), 'a', gocui.ModNone, t.AddTag())
	g.SetKeybinding(t.Name(), 'e', gocui.ModNone, t.EditTag())
	g.SetKeybinding(t.Name(), 'd', gocui.ModNone, t.DeleteTag())
	g.SetKeybinding(t.Name(), 'j', gocui.ModNone, t.MoveDown())
	g.SetKeybinding(t.Name(), 'k', gocui.ModNone, t.MoveUp())
	g.SetKeybinding(t.Name(), gocui.KeyCtrlD, gocui.ModNone, t.PageDown())
	g.SetKeybinding(t.Name(), gocui.KeyCtrlU, gocui.ModNone, t.PageUp())
	g.SetKeybinding(t.Name(), 'o', gocui.ModNone, t.ToggleChildren())
	g.SetKeybinding(t.Name(), 'l', gocui.ModNone, t.ExpandChildren())
	g.SetKeybinding(t.Name(), 'h', gocui.ModNone, t.CollapseChildren())
	g.SetKeybinding(t.Name(), '/', gocui.ModNone, t.StartSearch())
	g.SetKeybinding(t.Name(), 'x', gocui.ModNone, t.ResetSearch())
	g.SetKeybinding(t.Name(), gocui.KeySpace, gocui.ModNone, t.ToggleFilter())
	return nil
}

func (t *TagList) UnbindKeys(g *gocui.Gui) error {
	g.DeleteKeybinding(t.Name(), 'a', gocui.ModNone)
	g.DeleteKeybinding(t.Name(), 'e', gocui.ModNone)
	g.DeleteKeybinding(t.Name(), 'd', gocui.ModNone)
	g.DeleteKeybinding(t.Name(), 'j', gocui.ModNone)
	g.DeleteKeybinding(t.Name(), 'k', gocui.ModNone)
	g.DeleteKeybinding(t.Name(), gocui.KeyCtrlD, gocui.ModNone)
	g.DeleteKeybinding(t.Name(), gocui.KeyCtrlU, gocui.ModNone)
	g.DeleteKeybinding(t.Name(), 'o', gocui.ModNone)
	g.DeleteKeybinding(t.Name(), 'l', gocui.ModNone)
	g.DeleteKeybinding(t.Name(), 'h', gocui.ModNone)
	g.DeleteKeybinding(t.Name(), '/', gocui.ModNone)
	g.DeleteKeybinding(t.Name(), 'x', gocui.ModNone)
	return nil
}

func (t *TagList) update(g *gocui.Gui) error {
	return t.render(g)
}

func (t *TagList) render(g *gocui.Gui) error {
	v, err := t.UIView(g)
	if err != nil {
		return err
	}
	v.Clear()
	switch {
	case t.Len() == 0:
		fmt.Fprint(v, "Empty")
	default:
		_, uy := t.UICursor(v)
		tags := t.getVisible()
		for i, tag := range tags {
			if i > 0 {
				fmt.Fprint(v, "\n")
			}
			fmt.Fprint(v, t.renderLine(&tag, !t.editMode() && uy == i))
		}
	}
	return nil
}

func (t *TagList) renderLine(row *pmapi.Tag, highlight bool) string {
	lPad := row.Level*2 + 1
	lRune := " "
	if highlight {
		lRune = ">"
	}
	str := fmt.Sprintf("%*s", lPad, lRune)
	if row != nil && t.App.Store.InTagFilter(*row) {
		str += fmt.Sprintf(" \x1b[1m%s\x1b[0m", row.Name)
	} else {
		str += fmt.Sprintf(" %s", row.Name)
	}
	if t.getData().HasChildren(row.ID) {
		expRune := "+"
		if t.opens[row.ID] {
			expRune = "-"
		}
		str += fmt.Sprintf(" [%s]", expRune)
	}
	return str
}

func (t *TagList) Len() int {
	return len(t.getVisible())
}

func (t *TagList) getVisible() tagSlice {
	tm := t.getData().byParent()
	return t.getVisibleChildren(tm, "")
}

func (t *TagList) getVisibleChildren(tm map[string]tagSlice, id string) tagSlice {
	out := tagSlice{}
	for _, item := range tm[id] {
		out = append(out, item)
		if t.opens[item.ID] {
			out = append(out, t.getVisibleChildren(tm, item.ID)...)
		}
	}
	return out
}

func (t *TagList) getVisibleAt(i int) *pmapi.Tag {
	data := t.getVisible()
	if i >= len(data) {
		return nil
	}
	item := data[i]
	return &item
}

func (t *TagList) getAt(i int) *pmapi.Tag {
	return t.getVisibleAt(i)
}

func (t *TagList) getData() tagSlice {
	data, err := t.App.Store.GetTags(context.TODO())
	if err != nil {
		return nil
	}
	return data
}

func (t *TagList) MoveUp() CommandFunc {
	return func(g *gocui.Gui, v *gocui.View) error {
		t.View.MoveUp()(g, v)
		return t.render(g)
	}
}

func (t *TagList) MoveDown() CommandFunc {
	return func(g *gocui.Gui, v *gocui.View) error {
		t.View.MoveDown()(g, v)
		return t.render(g)
	}
}

func (t *TagList) PageUp() CommandFunc {
	return func(g *gocui.Gui, v *gocui.View) error {
		t.View.PageUp(t.Len())(g, v)
		return t.render(g)
	}
}

func (t *TagList) PageDown() CommandFunc {
	return func(g *gocui.Gui, v *gocui.View) error {
		t.View.PageDown(t.Len())(g, v)
		return t.render(g)
	}
}

func (t *TagList) ToggleChildren() CommandFunc {
	return func(g *gocui.Gui, v *gocui.View) error {
		_, uy := t.UICursor(v)
		tag := t.getAt(uy)
		return t.setOpen(!t.opens[tag.ID])(g, v)
	}
}

func (t *TagList) ExpandChildren() CommandFunc {
	return t.setOpen(true)
}

func (t *TagList) CollapseChildren() CommandFunc {
	return func(g *gocui.Gui, v *gocui.View) error {
		_, uy := t.UICursor(v)
		tag := t.getAt(uy)
		if t.getData().HasChildren(tag.ID) {
			return t.setOpen(false)(g, v)
		}

		idx := t.getVisible().findID(tag.ParentID)
		if idx == -1 {
			return nil
		}
		_, cy := v.Cursor()
		_, oy := v.Origin()
		switch {
		case oy >= idx:
			oy = idx
		default:
			cy = idx - oy
		}
		v.SetCursor(0, cy)
		v.SetOrigin(0, oy)
		return t.setOpen(false)(g, v)
	}
}

func (t *TagList) setOpen(open bool) CommandFunc {
	return func(g *gocui.Gui, v *gocui.View) error {
		_, uy := t.UICursor(v)
		tag := t.getAt(uy)
		t.opens[tag.ID] = open
		return t.render(g)
	}
}

func (t *TagList) AddTag() CommandFunc {
	return func(g *gocui.Gui, v *gocui.View) error {
		t.editIndex = t.Len()
		t.editInsert = true
		t.nameField.SetValue(g, "")
		return t.StartEditing()(g, v)
	}
}

func (t *TagList) EditTag() CommandFunc {
	return func(g *gocui.Gui, v *gocui.View) error {
		_, y := t.UICursor(v)
		tag := t.getAt(y)
		if tag == nil {
			return nil
		}
		t.editIndex = y
		t.nameField.SetValue(g, tag.Name)
		return t.StartEditing()(g, v)
	}
}

func (t *TagList) StartEditing() CommandFunc {
	return func(g *gocui.Gui, v *gocui.View) error {
		_, oy := v.Origin()
		nameY0 := t.editIndex - oy
		t.nameField.SetPosition(
			t.X0(), nameY0,
			t.X1(), nameY0+2)
		t.nameField.label = "> "
		t.nameField.Attach(g)
		t.nameField.Focus(g)

		t.App.UnbindKeys(g)
		t.UnbindKeys(g)
		t.bindFormKeys(g)
		return t.render(g)
	}
}

func (t *TagList) EndEditing() CommandFunc {
	return func(g *gocui.Gui, v *gocui.View) error {
		t.resetEdit()
		t.nameField.Detach(g)
		t.unbindFormKeys(g)
		t.BindKeys(g)
		t.App.BindKeys(g)
		t.Focus(g)
		return t.render(g)
	}
}

func (t *TagList) SaveEditing() CommandFunc {
	return func(g *gocui.Gui, v *gocui.View) error {
		var model pmapi.Tag
		if data, err := t.gatherFields(g); err == nil {
			model = *data
		}

		insert := t.editInsert
		go func() {
			var err error
			ctx := context.TODO()
			t.App.ShowMessage(g, NewMessage("Loading..."), 0)
			if insert {
				err = t.App.Store.CreateTag(ctx, &model)
			} else {
				err = t.App.Store.UpdateTag(ctx, &model)
			}
			if err != nil {
				t.App.ShowMessage(g, NewErrorMessage(fmt.Sprintf("Some errors occurred: %v", err)), time.Second*2)
				return
			} else {
				g.Update(t.update)
				if insert {
					t.App.ShowMessage(g, NewSuccessMessage("Tag added"), time.Second)
				} else {
					t.App.ShowMessage(g, NewSuccessMessage("Tag updated"), time.Second)
				}
			}
		}()

		return t.EndEditing()(g, v)
	}
}

func (t *TagList) gatherFields(g *gocui.Gui) (*pmapi.Tag, error) {
	out := pmapi.Tag{}
	if val, err := t.nameField.Value(g); err == nil {
		out.Name = val
	}
	if v, err := t.UIView(g); err == nil && !t.editInsert {
		_, y := t.UICursor(v)
		if row := t.getAt(y); row != nil {
			out.ID = row.ID
		}
	}
	return &out, nil
}

func (t *TagList) bindFormKeys(g *gocui.Gui) error {
	g.SetKeybinding("", gocui.KeyCtrlQ, gocui.ModNone, t.EndEditing())
	g.SetKeybinding("", gocui.KeyCtrlS, gocui.ModNone, t.SaveEditing())
	return nil
}

func (t *TagList) unbindFormKeys(g *gocui.Gui) error {
	g.DeleteKeybinding("", gocui.KeyCtrlQ, gocui.ModNone)
	g.DeleteKeybinding("", gocui.KeyCtrlS, gocui.ModNone)
	return nil
}

func (t *TagList) DeleteTag() CommandFunc {
	return func(g *gocui.Gui, v *gocui.View) error {
		onCancel := func(g *gocui.Gui, v *gocui.View) error {
			t.Focus(g)
			return nil
		}
		onEnter := func(g *gocui.Gui, v *gocui.View) error {
			val, err := t.App.CommandValue(g)
			if err != nil {
				return err
			}
			val = strings.ToLower(val)
			if val == "y" || val == "yes" {
				v, _ := t.UIView(g)
				_, y := t.UICursor(v)
				row := t.getAt(y)
				if row != nil {
					go func() {
						t.App.ShowMessage(g, NewMessage("Loading..."), 0)
						err := t.App.Store.DeleteTag(context.TODO(), row)
						if err != nil {
							t.App.ShowMessage(g, NewErrorMessage(fmt.Sprintf("Some error occurred: %s", err)), time.Second)
						} else {
							v.MoveCursor(0, -1, false)
							t.Focus(g)
							g.Update(t.update)
							t.App.redrawPinl(g, nil)
							go func() {
								t.App.ShowMessage(g, NewSuccessMessage("Bookmark deleted"), time.Second)
							}()
						}
					}()
				}
			} else {
				onCancel(g, v)
			}
			return nil
		}
		return t.App.StartCommand("Confirm to delete? (y/n)", "", onEnter, onCancel)(g, v)
	}
}

func (t *TagList) StartSearch() CommandFunc {
	return func(g *gocui.Gui, v *gocui.View) error {
		prevSearch := t.App.Store.tagSearch
		return t.App.StartCommand("search tag>", prevSearch, t.SubmitSearch(), t.EndSearch())(g, v)
	}
}

func (t *TagList) SubmitSearch() CommandFunc {
	return func(g *gocui.Gui, v *gocui.View) error {
		search, err := t.App.CommandValue(g)
		if err != nil {
			return err
		}
		t.App.Store.SetTagSearch(search)
		tv, err := t.UIView(g)
		if err != nil {
			return err
		}
		tv.SetOrigin(0, 0)
		tv.SetCursor(0, 0)
		return t.EndSearch()(g, v)
	}
}

func (t *TagList) EndSearch() CommandFunc {
	return func(g *gocui.Gui, v *gocui.View) error {
		t.Focus(g)
		return t.render(g)
	}
}

func (t *TagList) ResetSearch() CommandFunc {
	return func(g *gocui.Gui, v *gocui.View) error {
		t.App.Store.SetTagSearch("")
		v.SetOrigin(0, 0)
		v.SetCursor(0, 0)
		return t.render(g)
	}
}

func (t *TagList) ToggleFilter() CommandFunc {
	return func(g *gocui.Gui, v *gocui.View) error {
		_, y := t.UICursor(v)
		tag := t.getAt(y)
		if tag != nil {
			if t.App.Store.InTagFilter(*tag) {
				t.App.Store.DelTagFilter(*tag)
			} else {
				t.App.Store.AddTagFilter(*tag)
			}
		}
		t.render(g)
		return t.App.redrawPinl(g, func(v *gocui.View) error {
			v.SetOrigin(0, 0)
			v.SetCursor(0, 0)
			return nil
		})
	}
}
