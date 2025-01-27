package main

import (
	"github.com/getlantern/systray"
	"github.com/pkg/browser"
	"io"
	"os"
	"syscall"
	"unsafe"
)

type msg struct {
	HWND   uintptr
	UINT   uintptr
	WPARAM int16
	LPARAM int64
	DWORD  int32
	POINT  struct{ X, Y int64 }
}

func run() {
	defer func(user32 *syscall.DLL) {
		_ = user32.Release()
	}(user32)

	setupSystray()
	updateHotkeys()
	registerHotkeys(user32)

	getMsg := user32.MustFindProc("GetMessageW")
	for {
		var msg = &msg{}
		_, _, _ = getMsg.Call(uintptr(unsafe.Pointer(msg)), 0, 0, 0, 1)

		if id := msg.WPARAM; id != 0 {
			_ = browser.OpenURL(_Facebook)
		}
	}
}

func setupSystray() {
	systray.SetIcon(getIcon(_IconPath))
	systray.SetTitle("Fesnuk")
	systray.SetTooltip("Mending scroll fesnuk ygy")

	go func() {
		subMenuVsCode := systray.AddMenuItemCheckbox("Aku malas", "Aku malas", config.VSCodeRedirection == 1)
		subMenuTop := systray.AddMenuItem("Hotkey", "Change hotkey")
		subMenuCtrlAlt := subMenuTop.AddSubMenuItemCheckbox(_HotKeyCtrlAlt, _HotKeyCtrlAlt, config.HotKey == _HotKeyCtrlAlt)
		subMenuCtrlShift := subMenuTop.AddSubMenuItemCheckbox(_HotKeyCtrlShift, _HotKeyCtrlShift, config.HotKey == _HotKeyCtrlShift)
		subMenuShiftAlt := subMenuTop.AddSubMenuItemCheckbox(_HotKeyShiftAlt, _HotKeyShiftAlt, config.HotKey == _HotKeyShiftAlt)

		systray.AddSeparator()

		mQuitOrig := systray.AddMenuItem("Quit", "Quit the app")

		for {
			select {
			case <-subMenuVsCode.ClickedCh:
				if subMenuVsCode.Checked() {
					subMenuVsCode.Uncheck()
					config.VSCodeRedirection = 0
				} else {
					subMenuVsCode.Check()
					config.VSCodeRedirection = 1
				}
				saveConfig()

			case <-subMenuCtrlAlt.ClickedCh:
				subMenuCtrlAlt.Check()
				subMenuCtrlShift.Uncheck()
				subMenuShiftAlt.Uncheck()
				config.HotKey = _HotKeyCtrlAlt
				saveConfig()

			case <-subMenuCtrlShift.ClickedCh:
				subMenuCtrlAlt.Uncheck()
				subMenuCtrlShift.Check()
				subMenuShiftAlt.Uncheck()
				config.HotKey = _HotKeyCtrlShift
				saveConfig()

			case <-subMenuShiftAlt.ClickedCh:
				subMenuCtrlAlt.Uncheck()
				subMenuCtrlShift.Uncheck()
				subMenuShiftAlt.Check()
				config.HotKey = _HotKeyShiftAlt
				saveConfig()

			case <-mQuitOrig.ClickedCh:
				systray.Quit()
				return
			}
		}
	}()
}

func getIcon(filePath string) []byte {
	file, err := os.Open(filePath)
	if err != nil {
		return nil
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	content, err := io.ReadAll(file)
	if err != nil {
		return nil
	}
	return content
}
