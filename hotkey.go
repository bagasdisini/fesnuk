package main

import (
	"bytes"
	"fmt"
	"log"
	"strings"
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

func updateHotkeys(hotkey string) {
	modifiers, keyCode, err := parseHotkey(hotkey)
	if err != nil {
		return
	}

	HOTKEYS[1] = &Hotkey{
		ID:        1,
		Modifiers: modifiers,
		KeyCode:   keyCode,
	}
}

func parseHotkey(hotkey string) (modifiers int, keyCode int, err error) {
	parts := strings.Split(hotkey, "+")
	modifiers = 0
	keyCode = 0

	for _, part := range parts {
		switch part {
		case "Ctrl":
			modifiers += ModCtrl
		case "Alt":
			modifiers += ModAlt
		case "Shift":
			modifiers += ModShift
		case "Win":
			modifiers += ModWin
		default:
			if len(part) == 1 {
				keyCode = int(part[0])
			} else {
				log.Printf("Invalid hotkey format: %s", hotkey)
				return 0, 0, fmt.Errorf("invalid hotkey format: %s", hotkey)
			}
		}
	}

	if keyCode == 0 {
		log.Printf("No key code found in hotkey: %s", hotkey)
		return 0, 0, fmt.Errorf("no key code found in hotkey: %s", hotkey)
	}
	return modifiers, keyCode, nil
}
