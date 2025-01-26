package main

import (
	"errors"
	"github.com/getlantern/systray"
	"golang.org/x/sys/windows"
	"log"
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
	mutexName := "FesnukAppMutex"
	mutex, err := windows.CreateMutex(nil, false, windows.StringToUTF16Ptr(mutexName))
	if err != nil {
		if errors.Is(err, windows.ERROR_ALREADY_EXISTS) {
			return
		}
		return
	}
	defer windows.CloseHandle(mutex)

	config, err = loadConfig("config.ini")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
		return
	}
	systray.Run(run, func() {})
}
