package main

import (
	"fmt"
	"strings"
	"syscall"
)

const (
	_ModAlt = 1 << iota
	_ModCtrl
	_ModShift
	_ModWin
)

type hotKey struct {
	ID        int
	Modifiers int
	KeyCode   int
}

var hotKeys = map[int16]*hotKey{
	1: {4, _ModCtrl + _ModAlt, 'F'},
}

func registerHotkeys(user32 *syscall.DLL) {
	regHotKey := user32.MustFindProc("RegisterHotKey")
	for _, v := range hotKeys {
		_, _, _ = regHotKey.Call(0, uintptr(v.ID), uintptr(v.Modifiers), uintptr(v.KeyCode))
	}
}

func updateHotkeys() {
	modifiers, keyCode, err := parseHotkey()
	if err != nil {
		return
	}

	hotKeys[1] = &hotKey{
		ID:        1,
		Modifiers: modifiers,
		KeyCode:   keyCode,
	}
}

func parseHotkey() (modifiers int, keyCode int, err error) {
	parts := strings.Split(config.HotKey, "+")
	modifiers = 0
	keyCode = 0

	for _, part := range parts {
		switch part {
		case "Ctrl":
			modifiers += _ModCtrl
		case "Alt":
			modifiers += _ModAlt
		case "Shift":
			modifiers += _ModShift
		case "Win":
			modifiers += _ModWin
		default:
			if len(part) == 1 {
				keyCode = int(part[0])
			} else {
				return 0, 0, fmt.Errorf("invalid hotkey format: %s", config.HotKey)
			}
		}
	}

	if keyCode == 0 {
		return 0, 0, fmt.Errorf("no key code found in hotkey: %s", config.HotKey)
	}
	return modifiers, keyCode, nil
}
