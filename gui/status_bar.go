package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"strings"
)

type statusBar struct {
	msgLabel *widget.Label
	widget   *fyne.Container
}

func newStatusBar() *statusBar {
	var sb statusBar

	sb.msgLabel = widget.NewLabel("")
	sb.widget = container.New(layout.NewHBoxLayout(),
		sb.msgLabel,
		layout.NewSpacer())

	return &sb
}

func (sb *statusBar) setMessage(msg string) {
	// 2022/03/16 00:26:41 main_win.go:125: Log message => 2022/03/16 00:26:41 Log message
	arr := strings.Split(msg, " ")
	s := append(arr[:2], arr[3:]...)
	sb.msgLabel.SetText(strings.Trim(strings.Join(s, " "), "\n"))
}
