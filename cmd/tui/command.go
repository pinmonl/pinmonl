package tui

import "github.com/jroimartin/gocui"

type CommandFunc func(*gocui.Gui, *gocui.View) error
