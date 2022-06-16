package gui

import (
	"fmt"
	"fyne.io/fyne/v2/dialog"
	"github.com/BabySid/gobase"
	"log"
	sidTheme "sid-desktop/theme"
)

func printErr(err error) {
	log.Printf(sidTheme.InternalErrorFormat, err)
	dialog.ShowError(err, globalWin.win)
}

func assertErr(cond bool, a ...interface{}) {
	if !cond {
		printErr(fmt.Errorf(fmt.Sprint(a)))
		gobase.AssertHere()
	}
}
