package gui

import (
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	sidTheme "sid-desktop/desktop/theme"
)

type about struct {
	aboutDialog dialog.Dialog
}

func newAbout() *about {
	var a about
	content := widget.NewCard("", "", container.NewVBox(
		widget.NewLabel(sidTheme.AboutIntro),
	))
	a.aboutDialog = dialog.NewCustom(sidTheme.AboutTitle, sidTheme.ConfirmText, content, globalWin.win)
	return &a
}

func showAbout() {
	newAbout().aboutDialog.Show()
}
