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
		menuIde := systray.AddMenuItem("Aku malas", "Aku malas")
		subMenuVsCode := menuIde.AddSubMenuItemCheckbox("VS Code", "VS Code", config.VSCodeRedirection == 1)
		subMenuGoland := menuIde.AddSubMenuItemCheckbox("Goland", "Goland", config.GolandRedirection == 1)
		subMenuPyCharm := menuIde.AddSubMenuItemCheckbox("PyCharm", "PyCharm", config.PyCharmRedirection == 1)
		subMenuWebStorm := menuIde.AddSubMenuItemCheckbox("WebStorm", "WebStorm", config.WebStormRedirection == 1)
		subMenuRustRover := menuIde.AddSubMenuItemCheckbox("Rust Rover", "Rust Rover", config.RustRoverRedirection == 1)

		menuHotkey := systray.AddMenuItem("Hotkey", "Change hotkey")
		subMenuCtrlAlt := menuHotkey.AddSubMenuItemCheckbox(_HotKeyCtrlAlt, _HotKeyCtrlAlt, config.HotKey == _HotKeyCtrlAlt)
		subMenuCtrlShift := menuHotkey.AddSubMenuItemCheckbox(_HotKeyCtrlShift, _HotKeyCtrlShift, config.HotKey == _HotKeyCtrlShift)
		subMenuShiftAlt := menuHotkey.AddSubMenuItemCheckbox(_HotKeyShiftAlt, _HotKeyShiftAlt, config.HotKey == _HotKeyShiftAlt)

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

			case <-subMenuGoland.ClickedCh:
				if subMenuGoland.Checked() {
					subMenuGoland.Uncheck()
					config.GolandRedirection = 0
				} else {
					subMenuGoland.Check()
					config.GolandRedirection = 1
				}
				saveConfig()

			case <-subMenuPyCharm.ClickedCh:
				if subMenuPyCharm.Checked() {
					subMenuPyCharm.Uncheck()
					config.PyCharmRedirection = 0
				} else {
					subMenuPyCharm.Check()
					config.PyCharmRedirection = 1
				}
				saveConfig()

			case <-subMenuWebStorm.ClickedCh:
				if subMenuWebStorm.Checked() {
					subMenuWebStorm.Uncheck()
					config.WebStormRedirection = 0
				} else {
					subMenuWebStorm.Check()
					config.WebStormRedirection = 1
				}
				saveConfig()

			case <-subMenuRustRover.ClickedCh:
				if subMenuRustRover.Checked() {
					subMenuRustRover.Uncheck()
					config.RustRoverRedirection = 0
				} else {
					subMenuRustRover.Check()
					config.RustRoverRedirection = 1
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
