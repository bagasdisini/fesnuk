package main

import (
	"github.com/getlantern/systray"
	"github.com/pkg/browser"
	"io"
	"os"
	"syscall"
	"unsafe"
)

func run() {
	user32 := syscall.MustLoadDLL("user32")
	defer user32.Release()

	setupSystray(user32)
	updateHotkeys(config.HotKey)
	registerHotkeys(user32)

	getMsg := user32.MustFindProc("GetMessageW")
	for {
		var msg = &MSG{}
		getMsg.Call(uintptr(unsafe.Pointer(msg)), 0, 0, 0, 1)

		if id := msg.WPARAM; id != 0 {
			browser.OpenURL("https://www.facebook.com")
		}
	}
}

func setupSystray(user32 *syscall.DLL) {
	systray.SetIcon(getIcon("assets/icon.ico"))
	systray.SetTitle("Fesnuk")
	systray.SetTooltip("CTRL+ALT+F to open Facebook")

	go func() {
		subMenuVsCode := systray.AddMenuItemCheckbox("Malas ngoding", "Malas ngoding", config.VSCodeRedirection == 1)
		subMenuTop := systray.AddMenuItem("Hotkey", "Change hotkey")
		subMenuCtrlAlt := subMenuTop.AddSubMenuItemCheckbox("Ctrl+Alt+F", "Ctrl+Alt+F", config.HotKey == "Ctrl+Alt+F")
		subMenuCtrlShift := subMenuTop.AddSubMenuItemCheckbox("Ctrl+Shift+F", "Ctrl+Shift+F", config.HotKey == "Ctrl+Shift+F")
		subMenuShiftAlt := subMenuTop.AddSubMenuItemCheckbox("Shift+Alt+F", "Shift+Alt+F", config.HotKey == "Shift+Alt+F")

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
				saveConfig("config.ini")

			case <-subMenuCtrlAlt.ClickedCh:
				subMenuCtrlAlt.Check()
				subMenuCtrlShift.Uncheck()
				subMenuShiftAlt.Uncheck()
				config.HotKey = "Ctrl+Alt+F"
				saveConfig("config.ini")

			case <-subMenuCtrlShift.ClickedCh:
				subMenuCtrlAlt.Uncheck()
				subMenuCtrlShift.Check()
				subMenuShiftAlt.Uncheck()
				config.HotKey = "Ctrl+Shift+F"
				saveConfig("config.ini")

			case <-subMenuShiftAlt.ClickedCh:
				subMenuCtrlAlt.Uncheck()
				subMenuCtrlShift.Uncheck()
				subMenuShiftAlt.Check()
				config.HotKey = "Shift+Alt+F"
				saveConfig("config.ini")

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
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return nil
	}
	return content
}
