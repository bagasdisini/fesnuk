package main

import (
	"fmt"
	"gopkg.in/ini.v1"
	"os"
)

type Config struct {
	HotKey            string
	VSCodeRedirection int
}

var config *Config

func loadConfig(configPath string) (*Config, error) {
	cfg, err := ini.Load(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			cfg = ini.Empty()
			cfg.Section("Settings").Key("HotKey").SetValue("Ctrl+Alt+F")
			cfg.Section("Settings").Key("VSCodeRedirection").SetValue("0")
			if err := cfg.SaveTo(configPath); err != nil {
				return nil, fmt.Errorf("failed to create config file: %v", err)
			}
		} else {
			return nil, fmt.Errorf("failed to load config file: %v", err)
		}
	}

	config := &Config{
		HotKey:            cfg.Section("Settings").Key("HotKey").String(),
		VSCodeRedirection: cfg.Section("Settings").Key("VSCodeRedirection").MustInt(0),
	}

	return config, nil
}

func saveConfig(configPath string) {
	cfg := ini.Empty()
	cfg.Section("Settings").Key("HotKey").SetValue(config.HotKey)
	cfg.Section("Settings").Key("VSCodeRedirection").SetValue(fmt.Sprintf("%d", config.VSCodeRedirection))
	cfg.SaveTo(configPath)
	return
}
