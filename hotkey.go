package main

import (
	"bytes"
	"fmt"
	"syscall"
)

const (
	ModAlt = 1 << iota
	ModCtrl
	ModShift
	ModWin
)

type Hotkey struct {
	ID        int
	Modifiers int
	KeyCode   int
}

var HOTKEYS = map[int16]*Hotkey{
	1: {4, ModCtrl + ModAlt, 'F'},
}

func (h *Hotkey) String() string {
	mod := &bytes.Buffer{}
	if h.Modifiers&ModAlt != 0 {
		mod.WriteString("Alt+")
	}
	if h.Modifiers&ModCtrl != 0 {
		mod.WriteString("Ctrl+")
	}
	if h.Modifiers&ModShift != 0 {
		mod.WriteString("Shift+")
	}
	if h.Modifiers&ModWin != 0 {
		mod.WriteString("Win+")
	}
	return fmt.Sprintf("Hotkey[Id: %d, %s%c]", h.ID, mod, h.KeyCode)
}

func registerHotkeys(user32 *syscall.DLL) {
	hotKey := user32.MustFindProc("RegisterHotKey")
	for _, v := range HOTKEYS {
		hotKey.Call(0, uintptr(v.ID), uintptr(v.Modifiers), uintptr(v.KeyCode))
	}
}
