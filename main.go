package main

import (
	"errors"
	"github.com/getlantern/systray"
	"golang.org/x/sys/windows"
	"syscall"
)

const (
	_ConfigFile = "config.ini"
	_UserDLL    = "user32"

	_Facebook = "https://www.facebook.com"
	_IconPath = "assets/icon.ico"

	_HotKeyCtrlAlt   = "Ctrl+Alt+F"
	_HotKeyCtrlShift = "Ctrl+Shift+F"
	_HotKeyShiftAlt  = "Shift+Alt+F"
)

var (
	user32 = syscall.MustLoadDLL(_UserDLL)
)

func main() {
	mutexName := "FesnukAppMutex"
	mutex, err := windows.CreateMutex(nil, false, windows.StringToUTF16Ptr(mutexName))
	if err != nil {
		if errors.Is(err, windows.ERROR_ALREADY_EXISTS) {
			return
		}
		return
	}
	defer func(handle windows.Handle) {
		_ = windows.CloseHandle(handle)
	}(mutex)

	config, err = loadConfig()
	if err != nil {
		return
	}

	go watchConfig()

	systray.Run(run, func() {})
}
