// +build !nogui

package main

import (
	"fyne.io/fyne/widget"
)

func MsgBox(title string, message string) {
	w := guiapp.NewWindow(title)
	w.SetContent(widget.NewVBox(
		widget.NewLabel(title),
		widget.NewLabel(message),
		widget.NewButton("Ok", func() {
			w.Close()
		}),
	))
	w.Show()
}
