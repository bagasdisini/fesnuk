package main

import (
	"fesnuk/assets"
	"fesnuk/internal/config"
	"fesnuk/internal/hotkey"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"strings"
)

func main() {
	var err error
	config.Config, err = config.LoadConfig()
	if err != nil {
		return
	}

	myApp := app.New()
	myApp.SetIcon(assets.ResourceIconIco)
	myWindow := myApp.NewWindow("Fesnuk")
	myWindow.SetIcon(assets.ResourceIconIco)

	shiftCb := widget.NewCheck("Shift", func(checked bool) {})
	ctrlCb := widget.NewCheck("Ctrl", func(checked bool) {})
	winCb := widget.NewCheck("Win", func(checked bool) {})
	altCb := widget.NewCheck("Alt", func(checked bool) {})
	keyCb := widget.NewSelect(hotkey.SupportedKeys, func(value string) {})

	keyCb.PlaceHolder = " "

	hotKeyGrid := container.NewGridWithColumns(5, shiftCb, ctrlCb, winCb, altCb, keyCb)

	if strings.Contains(config.Config.HotKey, "Ctrl") {
		ctrlCb.Checked = true
	}
	if strings.Contains(config.Config.HotKey, "Shift") {
		shiftCb.Checked = true
	}
	if strings.Contains(config.Config.HotKey, "Win") {
		winCb.Checked = true
	}
	if strings.Contains(config.Config.HotKey, "Alt") {
		altCb.Checked = true
	}
	keyCb.Selected = hotkey.GetSupportedKeys()

	vsCodeCb := widget.NewCheck("VSCode", func(checked bool) {})
	golandCb := widget.NewCheck("Goland", func(checked bool) {})
	pyCharmCb := widget.NewCheck("PyCharm", func(checked bool) {})
	webStormCb := widget.NewCheck("WebStorm", func(checked bool) {})
	rustRoverCb := widget.NewCheck("RustRover", func(checked bool) {})

	ideGrid := container.NewGridWithColumns(2, vsCodeCb, golandCb, pyCharmCb, webStormCb, rustRoverCb)

	if config.Config.VSCodeRedirection != 0 {
		vsCodeCb.Checked = true
	}
	if config.Config.GolandRedirection != 0 {
		golandCb.Checked = true
	}
	if config.Config.PyCharmRedirection != 0 {
		pyCharmCb.Checked = true
	}
	if config.Config.WebStormRedirection != 0 {
		webStormCb.Checked = true
	}
	if config.Config.RustRoverRedirection != 0 {
		rustRoverCb.Checked = true
	}

	saveButton := widget.NewButton("Save", func() {
		modifiers := []string{}
		if ctrlCb.Checked {
			modifiers = append(modifiers, "Ctrl")
		}
		if shiftCb.Checked {
			modifiers = append(modifiers, "Shift")
		}
		if winCb.Checked {
			modifiers = append(modifiers, "Win")
		}
		if altCb.Checked {
			modifiers = append(modifiers, "Alt")
		}

		config.Config.HotKey = hotkey.BuildHotkey(keyCb.Selected, modifiers)

		config.Config.VSCodeRedirection = boolToInt(vsCodeCb.Checked)
		config.Config.GolandRedirection = boolToInt(golandCb.Checked)
		config.Config.PyCharmRedirection = boolToInt(pyCharmCb.Checked)
		config.Config.WebStormRedirection = boolToInt(webStormCb.Checked)
		config.Config.RustRoverRedirection = boolToInt(rustRoverCb.Checked)

		config.SaveConfig()

		dialog.ShowInformation("Success", "Configuration saved successfully!", myWindow)
	})

	content := container.NewVBox(
		widget.NewLabel("Select hotkey:"),
		hotKeyGrid,
		widget.NewLabel("Aku malas:"),
		ideGrid,
		saveButton,
	)

	myWindow.SetContent(content)
	myWindow.Resize(fyne.NewSize(100, 100))
	myWindow.ShowAndRun()
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
