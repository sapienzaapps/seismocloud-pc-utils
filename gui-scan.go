// +build !nogui

package main

import (
	"fmt"
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
	"time"
)

func guiscan() {
	w := guiapp.NewWindow("Device scan")
	w.SetContent(widget.NewVBox(
		widget.NewLabel("Scan in progress"),
		widget.NewProgressBarInfinite(),
		widget.NewButton("Abort", func() {
			stopDiscovery()
		}),
	))
	w.Show()

	devices := discovery(5 * time.Second)

	w.Close()

	if devices == nil {
		MsgBox(
			"Internal error",
			"Please check if the network connection is OK and no other instances of this program are running",
		)
		return
	} else if len(devices) == 0 {
		MsgBox(
			"Scan completed",
			"No device found! Remember to connect the SeismoCloud sensor in the same network of this PC",
		)
		return
	}

	buttons := []fyne.CanvasObject{
		widget.NewLabel("Discovered devices"),
	}
	for _, dev := range devices {
		buttons = append(buttons, widget.NewButton(
			fmt.Sprintf("Device %s (%s - v%s) - %s", dev.DeviceID, dev.Model, dev.Version, dev.RemoteIP),
			func() {
				// TODO: show options
			},
		))
	}

	buttons = append(buttons, widget.NewButton("Return to main menu", func() {
		w.Close()
	}))

	w = guiapp.NewWindow("Discovered devices")
	w.SetContent(widget.NewVBox(buttons...))
	w.Show()
}
