package gui

import (
	"fyne.io/fyne/v2/dialog"
	"log"
	sidTheme "sid/desktop/theme"
)

func printErr(err error) {
	log.Printf(sidTheme.InternalErrorFormat, err)
	dialog.ShowError(err, globalWin.win)
}
