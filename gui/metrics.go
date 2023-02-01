package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type metrics struct {
	ok      *widget.Button
	dismiss *widget.Button

	win fyne.Window
}
