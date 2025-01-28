package config

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"gopkg.in/ini.v1"
	"os"
	"os/exec"
	"sync"
	"syscall"
)

const (
	ConfigFile = "config.ini"
	UserDLL    = "user32"

	Facebook = "https://www.facebook.com"
	IconPath = "assets/icon.ico"

	VSCode    = "Code.exe"
	Goland    = "goland64.exe"
	PyCharm   = "pycharm64.exe"
	WebStorm  = "webstorm64.exe"
	RustRover = "rustrover64.exe"

	HotKeyCtrlAlt   = "Ctrl+Alt+F"
	HotKeyCtrlShift = "Ctrl+Shift+F"
	HotKeyShiftAlt  = "Shift+Alt+F"
)

var (
	User32 = syscall.MustLoadDLL(UserDLL)
)

type Doc struct {
	HotKey               string
	VSCodeRedirection    int
	GolandRedirection    int
	PyCharmRedirection   int
	WebStormRedirection  int
	RustRoverRedirection int
}

var Config *Doc

func LoadConfig() (*Doc, error) {
	cfg, err := ini.Load(ConfigFile)
	if err != nil {
		if os.IsNotExist(err) {
			cfg = ini.Empty()
			cfg.Section("Settings").Key("HotKey").SetValue(HotKeyCtrlAlt)
			cfg.Section("Settings").Key("VSCodeRedirection").SetValue("0")
			cfg.Section("Settings").Key("GolandRedirection").SetValue("0")
			cfg.Section("Settings").Key("PyCharmRedirection").SetValue("0")
			cfg.Section("Settings").Key("WebStormRedirection").SetValue("0")
			cfg.Section("Settings").Key("RustRoverRedirection").SetValue("0")
			if err := cfg.SaveTo(ConfigFile); err != nil {
				return nil, fmt.Errorf("failed to create config file: %v", err)
			}
		} else {
			return nil, fmt.Errorf("failed to load config file: %v", err)
		}
	}

	configDoc := &Doc{
		HotKey:               cfg.Section("Settings").Key("HotKey").String(),
		VSCodeRedirection:    cfg.Section("Settings").Key("VSCodeRedirection").MustInt(0),
		GolandRedirection:    cfg.Section("Settings").Key("GolandRedirection").MustInt(0),
		PyCharmRedirection:   cfg.Section("Settings").Key("PyCharmRedirection").MustInt(0),
		WebStormRedirection:  cfg.Section("Settings").Key("WebStormRedirection").MustInt(0),
		RustRoverRedirection: cfg.Section("Settings").Key("RustRoverRedirection").MustInt(0),
	}
	return configDoc, nil
}

func SaveConfig() {
	cfg := ini.Empty()
	cfg.Section("Settings").Key("HotKey").SetValue(Config.HotKey)
	cfg.Section("Settings").Key("VSCodeRedirection").SetValue(fmt.Sprintf("%d", Config.VSCodeRedirection))
	cfg.Section("Settings").Key("GolandRedirection").SetValue(fmt.Sprintf("%d", Config.GolandRedirection))
	cfg.Section("Settings").Key("PyCharmRedirection").SetValue(fmt.Sprintf("%d", Config.PyCharmRedirection))
	cfg.Section("Settings").Key("WebStormRedirection").SetValue(fmt.Sprintf("%d", Config.WebStormRedirection))
	cfg.Section("Settings").Key("RustRoverRedirection").SetValue(fmt.Sprintf("%d", Config.RustRoverRedirection))
	_ = cfg.SaveTo(ConfigFile)
	return
}

func WatchConfig() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return
	}
	defer func(watcher *fsnotify.Watcher) {
		_ = watcher.Close()
	}(watcher)

	err = watcher.Add(ConfigFile)
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
