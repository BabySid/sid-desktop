package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	sidTheme "sid/desktop/theme"
)

type logViewer struct {
	logArea *widget.Entry
	refresh *widget.Button

	Win fyne.Window
}

func newLogViewer() *logViewer {
	var lv logViewer

	lv.logArea = widget.NewMultiLineEntry()
	lv.logArea.Disable()
	lv.logArea.Wrapping = fyne.TextWrapBreak

	lv.initLogContent()

	lv.refresh = widget.NewButtonWithIcon(sidTheme.LogViewerRefreshBtn, theme.ViewRefreshIcon(), func() {
		lv.initLogContent()
	})

	lv.Win = fyne.CurrentApp().NewWindow(sidTheme.LogViewerTitle)
	lv.Win.SetContent(container.NewBorder(
		container.NewHBox(layout.NewSpacer(), lv.refresh), nil, nil, nil,
		lv.logArea))
	lv.Win.Resize(fyne.NewSize(800, 600))
	lv.Win.CenterOnScreen()

	return &lv
}

func (lv *logViewer) initLogContent() {
	cont := ""
	globalLogWriter.Traversal(func(s string) {
		cont += s
	})

	lv.logArea.CursorRow = globalLogWriter.Size()
	lv.logArea.SetText(cont)
}
