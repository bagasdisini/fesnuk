package main

import (
	"errors"
	"fesnuk/internal/config"
	"fesnuk/internal/monitor"
	_systray "fesnuk/internal/systray"
	"github.com/getlantern/systray"
	"golang.org/x/sys/windows"
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

	config.Config, err = config.LoadConfig()
	if err != nil {
		return
	}

	go config.WatchConfig()

	ides := []string{}

	if config.Config.VSCodeRedirection != 0 {
		ides = append(ides, config.VSCode)
	}
	if config.Config.GolandRedirection != 0 {
		ides = append(ides, config.Goland)
	}
	if config.Config.PyCharmRedirection != 0 {
		ides = append(ides, config.PyCharm)
	}
	if config.Config.WebStormRedirection != 0 {
		ides = append(ides, config.WebStorm)
	}
	if config.Config.RustRoverRedirection != 0 {
		ides = append(ides, config.RustRover)
	}

	if len(ides) != 0 {
		go monitor.DoMonitor(ides)
	}

	go _systray.OpenGui()
	systray.Run(_systray.Run, func() {})
}
