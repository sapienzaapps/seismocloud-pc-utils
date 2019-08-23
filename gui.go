// +build !nogui

package main

import (
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/widget"
)

var guiapp fyne.App

func guimain() {
	guiapp = app.New()
	w := guiapp.NewWindow("Hello")
	w.SetContent(widget.NewVBox(
		widget.NewLabel("Hello Fyne!"),
		widget.NewButton("Scan for devices in LAN", guiscan),
		widget.NewButton("Quit", func() {
			guiapp.Quit()
		}),
	))
	w.ShowAndRun()
}
