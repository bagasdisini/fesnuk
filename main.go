package main

import (
	"github.com/getlantern/systray"
)

type MSG struct {
	HWND   uintptr
	UINT   uintptr
	WPARAM int16
	LPARAM int64
	DWORD  int32
	POINT  struct{ X, Y int64 }
}

func main() {
	systray.Run(onReady, func() {})
}
