package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"gopkg.in/ini.v1"
	"os"
	"os/exec"
	"sync"
)

type Config struct {
	HotKey               string
	VSCodeRedirection    int
	GolandRedirection    int
	PyCharmRedirection   int
	WebStormRedirection  int
	RustRoverRedirection int
}

var config *Config

func loadConfig() (*Config, error) {
	cfg, err := ini.Load(_ConfigFile)
	if err != nil {
		if os.IsNotExist(err) {
			cfg = ini.Empty()
			cfg.Section("Settings").Key("HotKey").SetValue(_HotKeyCtrlAlt)
			cfg.Section("Settings").Key("VSCodeRedirection").SetValue("0")
			cfg.Section("Settings").Key("GolandRedirection").SetValue("0")
			cfg.Section("Settings").Key("PyCharmRedirection").SetValue("0")
			cfg.Section("Settings").Key("WebStormRedirection").SetValue("0")
			cfg.Section("Settings").Key("RustRoverRedirection").SetValue("0")
			if err := cfg.SaveTo(_ConfigFile); err != nil {
				return nil, fmt.Errorf("failed to create config file: %v", err)
			}
		} else {
			return nil, fmt.Errorf("failed to load config file: %v", err)
		}
	}

	configDoc := &Config{
		HotKey:               cfg.Section("Settings").Key("HotKey").String(),
		VSCodeRedirection:    cfg.Section("Settings").Key("VSCodeRedirection").MustInt(0),
		GolandRedirection:    cfg.Section("Settings").Key("GolandRedirection").MustInt(0),
		PyCharmRedirection:   cfg.Section("Settings").Key("PyCharmRedirection").MustInt(0),
		WebStormRedirection:  cfg.Section("Settings").Key("WebStormRedirection").MustInt(0),
		RustRoverRedirection: cfg.Section("Settings").Key("RustRoverRedirection").MustInt(0),
	}
	return configDoc, nil
}

func saveConfig() {
	cfg := ini.Empty()
	cfg.Section("Settings").Key("HotKey").SetValue(config.HotKey)
	cfg.Section("Settings").Key("VSCodeRedirection").SetValue(fmt.Sprintf("%d", config.VSCodeRedirection))
	cfg.Section("Settings").Key("GolandRedirection").SetValue(fmt.Sprintf("%d", config.GolandRedirection))
	cfg.Section("Settings").Key("PyCharmRedirection").SetValue(fmt.Sprintf("%d", config.PyCharmRedirection))
	cfg.Section("Settings").Key("WebStormRedirection").SetValue(fmt.Sprintf("%d", config.WebStormRedirection))
	cfg.Section("Settings").Key("RustRoverRedirection").SetValue(fmt.Sprintf("%d", config.RustRoverRedirection))
	_ = cfg.SaveTo(_ConfigFile)
	return
}

func watchConfig() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return
	}
	defer func(watcher *fsnotify.Watcher) {
		_ = watcher.Close()
	}(watcher)

	err = watcher.Add(_ConfigFile)
	if err != nil {
		return
	}

	var restartLock sync.Mutex

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}

			if event.Op&(fsnotify.Write|fsnotify.Create) != 0 {
				restartLock.Lock()
				exeApp()
				restartLock.Unlock()
			}

		case _, ok := <-watcher.Errors:
			if !ok {
				return
			}
		}
	}
}

func exeApp() {
	exePath, err := os.Executable()
	if err != nil {
		return
	}

	cmd := exec.Command(exePath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err = cmd.Start(); err != nil {
		return
	}
	os.Exit(0)
}
