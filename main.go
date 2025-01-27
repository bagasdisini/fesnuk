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

	_VSCode    = "Code.exe"
	_Goland    = "goland64.exe"
	_PyCharm   = "pycharm64.exe"
	_WebStorm  = "webstorm64.exe"
	_RustRover = "rustrover64.exe"

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

	ides := []string{}

	if config.VSCodeRedirection != 0 {
		ides = append(ides, _VSCode)
	}
	if config.GolandRedirection != 0 {
		ides = append(ides, _Goland)
	}
	if config.PyCharmRedirection != 0 {
		ides = append(ides, _PyCharm)
	}
	if config.WebStormRedirection != 0 {
		ides = append(ides, _WebStorm)
	}
	if config.RustRoverRedirection != 0 {
		ides = append(ides, _RustRover)
	}

	if len(ides) != 0 {
		go doMonitor(ides)
	}

	systray.Run(run, func() {})
}
