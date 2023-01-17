package gui

import (
	"errors"
	"fmt"
	"fyne.io/fyne/v2/dialog"
	"github.com/BabySid/gobase"
	"github.com/mitchellh/go-wordwrap"
	"log"
	sidTheme "sid-desktop/theme"
)

func printErr(err error) {
	log.Printf(sidTheme.InternalErrorFormat, err)

	str := err.Error()
	const maxLen = 64

	target := wordwrap.WrapString(str, maxLen)
	dialog.ShowError(errors.New(target), globalWin.win)
}

func assertErr(cond bool, a ...interface{}) {
	if !cond {
		printErr(fmt.Errorf(fmt.Sprint(a)))
		gobase.AssertHere()
	}
}
