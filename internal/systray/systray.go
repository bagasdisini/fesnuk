package systray

import (
	"fesnuk/internal/config"
	"fesnuk/internal/hotkey"
	"github.com/getlantern/systray"
	"github.com/pkg/browser"
	"io"
	"os"
	"os/exec"
	"path/filepath"
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

func Run() {
	defer func(user32 *syscall.DLL) {
		_ = user32.Release()
	}(config.User32)

	setupSystray()
	hotkey.UpdateHotkeys()
	hotkey.RegisterHotkeys(config.User32)

	getMsg := config.User32.MustFindProc("GetMessageW")
	for {
		var msg = &msg{}
		_, _, _ = getMsg.Call(uintptr(unsafe.Pointer(msg)), 0, 0, 0, 1)

		if id := msg.WPARAM; id != 0 {
			_ = browser.OpenURL(config.Facebook)
		}
	}
}

func setupSystray() {
	systray.SetIcon(getIcon(config.IconPath))
	systray.SetTitle("Fesnuk")
	systray.SetTooltip("Mending scroll fesbuk ygy")

	go func() {
		mGuiOrig := systray.AddMenuItem("Open config", "Open config")

		systray.AddSeparator()

		menuIde := systray.AddMenuItem("Aku malas", "Aku malas")
		subMenuVsCode := menuIde.AddSubMenuItemCheckbox("VS Code", "VS Code", config.Config.VSCodeRedirection == 1)
		subMenuGoland := menuIde.AddSubMenuItemCheckbox("Goland", "Goland", config.Config.GolandRedirection == 1)
		subMenuPyCharm := menuIde.AddSubMenuItemCheckbox("PyCharm", "PyCharm", config.Config.PyCharmRedirection == 1)
		subMenuWebStorm := menuIde.AddSubMenuItemCheckbox("WebStorm", "WebStorm", config.Config.WebStormRedirection == 1)
		subMenuRustRover := menuIde.AddSubMenuItemCheckbox("Rust Rover", "Rust Rover", config.Config.RustRoverRedirection == 1)

		menuHotkey := systray.AddMenuItem("Hotkey", "Change hotkey")
		subMenuCtrlAlt := menuHotkey.AddSubMenuItemCheckbox(config.HotKeyCtrlAlt, config.HotKeyCtrlAlt, config.Config.HotKey == config.HotKeyCtrlAlt)
		subMenuCtrlShift := menuHotkey.AddSubMenuItemCheckbox(config.HotKeyCtrlShift, config.HotKeyCtrlShift, config.Config.HotKey == config.HotKeyCtrlShift)
		subMenuShiftAlt := menuHotkey.AddSubMenuItemCheckbox(config.HotKeyShiftAlt, config.HotKeyShiftAlt, config.Config.HotKey == config.HotKeyShiftAlt)
		subMenuCustom := menuHotkey.AddSubMenuItemCheckbox("(Custom)", "Custom hotkey", config.Config.HotKey != config.HotKeyCtrlAlt && config.Config.HotKey != config.HotKeyCtrlShift && config.Config.HotKey != config.HotKeyShiftAlt)

		systray.AddSeparator()

		mQuitOrig := systray.AddMenuItem("Quit", "Quit the app")

		for {
			select {
			case <-mGuiOrig.ClickedCh:
				openGui()

			case <-subMenuVsCode.ClickedCh:
				if subMenuVsCode.Checked() {
					subMenuVsCode.Uncheck()
					config.Config.VSCodeRedirection = 0
				} else {
					subMenuVsCode.Check()
					config.Config.VSCodeRedirection = 1
				}
				config.SaveConfig()

			case <-subMenuGoland.ClickedCh:
				if subMenuGoland.Checked() {
					subMenuGoland.Uncheck()
					config.Config.GolandRedirection = 0
				} else {
					subMenuGoland.Check()
					config.Config.GolandRedirection = 1
				}
				config.SaveConfig()

			case <-subMenuPyCharm.ClickedCh:
				if subMenuPyCharm.Checked() {
					subMenuPyCharm.Uncheck()
					config.Config.PyCharmRedirection = 0
				} else {
					subMenuPyCharm.Check()
					config.Config.PyCharmRedirection = 1
				}
				config.SaveConfig()

			case <-subMenuWebStorm.ClickedCh:
				if subMenuWebStorm.Checked() {
					subMenuWebStorm.Uncheck()
					config.Config.WebStormRedirection = 0
				} else {
					subMenuWebStorm.Check()
					config.Config.WebStormRedirection = 1
				}
				config.SaveConfig()

			case <-subMenuRustRover.ClickedCh:
				if subMenuRustRover.Checked() {
					subMenuRustRover.Uncheck()
					config.Config.RustRoverRedirection = 0
				} else {
					subMenuRustRover.Check()
					config.Config.RustRoverRedirection = 1
				}
				config.SaveConfig()

			case <-subMenuCtrlAlt.ClickedCh:
				subMenuCtrlAlt.Check()
				subMenuCtrlShift.Uncheck()
				subMenuShiftAlt.Uncheck()
				config.Config.HotKey = config.HotKeyCtrlAlt
				config.SaveConfig()

			case <-subMenuCtrlShift.ClickedCh:
				subMenuCtrlAlt.Uncheck()
				subMenuCtrlShift.Check()
				subMenuShiftAlt.Uncheck()
				config.Config.HotKey = config.HotKeyCtrlShift
				config.SaveConfig()

			case <-subMenuShiftAlt.ClickedCh:
				subMenuCtrlAlt.Uncheck()
				subMenuCtrlShift.Uncheck()
				subMenuShiftAlt.Check()
				config.Config.HotKey = config.HotKeyShiftAlt
				config.SaveConfig()

			case <-subMenuCustom.ClickedCh:
				openGui()

			case <-mQuitOrig.ClickedCh:
				systray.Quit()
				return
			}
		}
	}()
}

func openGui() {
	exePath, err := os.Executable()
	if err != nil {
		return
	}
	exeDir := filepath.Dir(exePath)

	executablePath := filepath.Join(exeDir, "fesnuk-gui.exe")
	cmd := exec.Command(executablePath)

	_ = cmd.Start()
	_ = cmd.Wait()
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
