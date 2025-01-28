package main

import (
	"fesnuk/assets"
	"fesnuk/internal/config"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
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

	dropdownOptions := []string{config.HotKeyCtrlAlt, config.HotKeyCtrlShift, config.HotKeyShiftAlt}
	selectedOption := widget.NewSelect(dropdownOptions, func(value string) {})
	selectedOption.PlaceHolder = "Select hotkey"

	selectedOption.Selected = config.Config.HotKey

	vsCodeCb := widget.NewCheck("VSCode", func(checked bool) {})
	golandCb := widget.NewCheck("Goland", func(checked bool) {})
	pyCharmCb := widget.NewCheck("PyCharm", func(checked bool) {})
	webStormCb := widget.NewCheck("WebStorm", func(checked bool) {})
	rustRoverCb := widget.NewCheck("RustRover", func(checked bool) {})

	checkboxGrid := container.NewGridWithColumns(2, vsCodeCb, golandCb, pyCharmCb, webStormCb, rustRoverCb)

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
		config.Config.HotKey = selectedOption.Selected

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
		selectedOption,
		widget.NewLabel("Aku malas:"),
		checkboxGrid,
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
