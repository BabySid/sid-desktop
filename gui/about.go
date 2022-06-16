package gui

import (
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	sidTheme "sid-desktop/theme"
)

type about struct {
	aboutDialog dialog.Dialog
}

func newAbout() *about {
	var a about

	content := widget.NewCard("", "", widget.NewRichTextFromMarkdown(sidTheme.AboutIntro))
	a.aboutDialog = dialog.NewCustom(sidTheme.AboutTitle, sidTheme.ConfirmText, content, globalWin.win)
	return &a
}

func showAbout() {
	newAbout().aboutDialog.Show()
}
