package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"sid-desktop/theme"
)

type about struct {
	aboutDialog dialog.Dialog
}

func newAbout() *about {
	var a about

	logo := canvas.NewImageFromResource(theme.ResourceAppIcon)
	//logo.FillMode = canvas.ImageFillContain
	logo.SetMinSize(fyne.NewSize(200, 200))

	content := widget.NewCard("", "",
		container.NewBorder(nil, nil, logo, nil, widget.NewRichTextFromMarkdown(theme.AboutIntro)))
	a.aboutDialog = dialog.NewCustom(theme.AboutTitle, theme.ConfirmText, content, globalWin.win)
	return &a
}

func showAbout() {
	newAbout().aboutDialog.Show()
}
