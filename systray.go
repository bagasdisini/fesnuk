package main

import (
	"github.com/getlantern/systray"
	"github.com/pkg/browser"
	"io"
	"os"
	"syscall"
	"unsafe"
)

func setupSystray() {
	systray.SetIcon(getIcon("assets/icon.ico"))
	systray.SetTitle("Fesnuk")
	systray.SetTooltip("CTRL+ALT+F to open Facebook")
	mQuitOrig := systray.AddMenuItem("Quit", "Quit the app")
	go func() {
		<-mQuitOrig.ClickedCh
		systray.Quit()
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

func onReady() {
	setupSystray()
	run()
}

func run() {
	user32 := syscall.MustLoadDLL("user32")
	defer user32.Release()

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
